package migrations

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"sort"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

type Transaction struct {
	tx *sql.Tx
}

func (s *Storage) BeginTranacton() (*Transaction, error) {
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return &Transaction{tx: tx}, nil
}

var (
	storage = &Storage{}
)

var (

	//go:embed *.up.sql
	engineMigrationsFS embed.FS

	appMigrationsFS fs.FS
)

func SetAppMigratioFs(storageDB *sqlx.DB, fsys fs.FS) {
	storage.db = storageDB

	appMigrationsFS = fsys
}

type migrationEntry struct {
	id       string
	filename string
	fsys     fs.FS
}

func parseMigrationFilename(filename string) (string, error) {
	if !strings.HasSuffix(filename, ".up.sql") {
		return "", fmt.Errorf("invalid migartion filename %q: must end with .up.sql", filename)
	}

	id := strings.TrimSuffix(filename, ".up.sql")

	if len(id) < 6 {
		return "", fmt.Errorf("invalid migation filenname %q: too short", filename)
	}

	prefix := id[:4]

	for _, c := range prefix {
		if c < '0' || c > '9' {
			return "", fmt.Errorf("invalid migration filename %q must start with 4-digit prefix", filename)
		}
	}

	if id[4] != '_' {
		return "", fmt.Errorf("invalid migration filename %q: digit prefix must be followed by underscore", filename)
	}

	return id, nil
}

func collectMigrations(fsys fs.FS) ([]migrationEntry, error) {
	if fsys == nil {
		return nil, nil
	}

	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []migrationEntry

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasSuffix(name, ".up.sql") {
			continue
		}

		id, err := parseMigrationFilename(name)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, migrationEntry{
			id:       id,
			filename: name,
			fsys:     fsys,
		})
	}

	return migrations, nil
}

func findMigrationFile(fsys fs.FS, version int) (string, error) {
	pattern := fmt.Sprintf("%04d_*.up.sql", version)
	matches, err := fs.Glob(fsys, pattern)
	if err != nil {
		return "", fmt.Errorf("failed to glob migration files using %q: %w", pattern, err)
	}

	if len(matches) == 0 {
		return "", fmt.Errorf(
			"no migration file matched pattern %q for version %04d",
			pattern,
			version)
	}

	if len(matches) > 1 {
		nonEmpty := make([]string, 0, len(matches))

		for _, m := range matches {
			b, rerr := fs.ReadFile(fsys, m)
			if rerr != nil {
				nonEmpty = append(nonEmpty, m)
				continue
			}

			content := strings.TrimSpace(string(b))
			lines := strings.Split(content, "\n")
			filtered := make([]string, 0, len(lines))

			for _, ln := range lines {
				s := strings.TrimSpace(ln)
				if s == "" {
					continue
				}

				if strings.HasSuffix(s, "--") {
					continue
				}

				filtered = append(filtered, s)

				if len(filtered) > 0 {
					nonEmpty = append(nonEmpty, m)
				}
			}
		}

		if len(nonEmpty) == 1 {
			return nonEmpty[0], nil
		}

		return "", fmt.Errorf(
			"multiple migration files matched pattern %q for version %04d: %v",
			pattern,
			version,
			matches,
		)
	}

	return matches[0], nil
}

func Run() error {
	engineMigrations, err := collectMigrations(engineMigrationsFS)
	if err != nil {
		return fmt.Errorf("failed to collect engine migations: %w", err)
	}

	appMigations, err := collectMigrations(appMigrationsFS)
	if err != nil {
		return fmt.Errorf("failed to collect application migations: %w", err)
	}

	allMigrations := append(engineMigrations, appMigations...)

	seen := make(map[string]string)

	for _, m := range allMigrations {
		if existing, ok := seen[m.id]; ok {
			return fmt.Errorf(
				"duplicate migration id %q: found in %q and %q",
				m.id,
				existing,
				m.filename,
			)
		}

		seen[m.id] = m.filename
	}

	sort.Slice(allMigrations, func(i, j int) bool {
		return allMigrations[i].id < allMigrations[j].id
	})

	tr, err := storage.BeginTranacton()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if tr.tx != nil {
			rberr := tr.tx.Rollback()
			if rberr != nil {
				log.Printf("failed to rollback transaction: %v", err)
			}
		}
	}()

	exists, err := chkTableExists(tr)
	if err != nil {
		return fmt.Errorf("failed to check if schema_migrations table exists: %w", err)
	}

	if !exists {
		err = createMigrationsTable(tr)
		if err != nil {
			return fmt.Errorf("failed to ensure schema_migrations table exists: %w", err)
		}
	}

	applied, err := getAppliedMigrations(tr)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	appliedCount := 0

	for _, m := range allMigrations {
		if applied[m.id] {
			continue
		}

		content, err := fs.ReadFile(m.fsys, m.filename)
		if err != nil {
			return fmt.Errorf("failed to read migration file %q: %w", m.filename, err)
		}

		log.Printf("applying migration: %s", m.id)
		if _, err := tr.tx.Exec(string(content)); err != nil {
			return fmt.Errorf("failed to apply migration %q: %w", m.id, err)
		}

		if err := recordMigration(tr, m.id); err != nil {
			return err
		}

		appliedCount++
	}

	if err := tr.tx.Commit(); err != nil {
		return fmt.Errorf("faile to commit transaction: %w", err)
	}

	tr.tx = nil

	if appliedCount == 0 {
		log.Printf("no new migrations to apply")
	} else {
		log.Printf("applied %d migration(s)", appliedCount)
	}

	return nil
}

func chkTableExists(tr *Transaction) (bool, error) {
	const query = `SELECT count(*) FROM information_schema.tables
	WHERE table_schema = 'public' AND table_name = 'schema_migrations'`

	var count int

	err := tr.tx.QueryRow(query).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check if schema_migrations table exists: %w", err)
	}

	return count > 0, nil
}

func createMigrationsTable(tr *Transaction) error {
	const createTableSQL = `CREATE TABLE IF NOT EXISTS schema_migrations
	(id TEXT PRIMARY KEY, applied_at TEXT DEFAULT CURRENT_TIMESTAMP)`

	_, err := tr.tx.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create schema_migrations table: %w", err)
	}

	return nil
}

// getAll
func getAppliedMigrations(tr *Transaction) (map[string]bool, error) {
	const query = `SELECT id FROM  schema_migrations`
	rows, err := tr.tx.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query applied migrations: %w", err)
	}

	defer rows.Close()

	applied := make(map[string]bool)

	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("failed to scan migration id: %w", err)
		}

		applied[id] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating applied migrations: %w", err)
	}

	return applied, nil
}

// insert
func recordMigration(tr *Transaction, id string) error {
	const query = `INSERT INTO schema_migrations (id) VALUES ($1);`
	if _, err := tr.tx.Exec(query, id); err != nil {
		return fmt.Errorf("failed to record migration %q: %w", id, err)
	}

	return nil
}

func getMigrationMaxTx(tr *Transaction) (int, error) {
	const query = "SELECT MAX(version) FROM schema_migrations"

	var max sql.NullInt64
	if err := tr.tx.QueryRow(query).Scan(&max); err != nil {
		return 0, fmt.Errorf("failed to get max migration version: %w", err)
	}

	if !max.Valid {
		return 0, nil
	}

	return int(max.Int64), nil
}
