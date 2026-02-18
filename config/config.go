package config

import "github.com/spf13/viper"

type Config struct {
	Path           string
	ApiVersion     string `mapstructure:"API_VERSION"`
	ServerPort     string `mapstructure:"SERVER_PORT"`
	DatabaseDriver string `mapstructure:"DATABASE_DRIVER"`
	DatabaseUrl    string `mapstructure:"DATABASE_URL"`
	MigrationRunning bool `mapstructure:"MIGRATION_RUN"`
}

func NewConfig(path string) *Config {
	return &Config{Path: path}

}

func (c *Config) Load() error {
	viper.AddConfigPath(c.Path)
	viper.SetConfigName("grpc-product")
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	return viper.Unmarshal(c)
}
