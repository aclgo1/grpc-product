package models

import "time"

type ParamsInsert struct {
	Id          string    `db:"product_id"`
	Name        string    `db:"name"`
	Price       float64   `db:"price"`
	Quantity    int64     `db:"quantity"`
	Description string    `db:"description"`
	Created_At  time.Time `db:"created_at"`
	Updated_At  time.Time `db:"updated_at"`
}

type ParamsInsertResponse struct {
	Id          string    `db:"product_id"`
	Name        string    `db:"name"`
	Price       float64   `db:"price"`
	Quantity    int64     `db:"quantity"`
	Description string    `db:"description"`
	Created_At  time.Time `db:"created_at"`
	Updated_At  time.Time `db:"updated_at"`
}

type ParamsFind struct {
	Id string
}

type ParamsFindResult struct {
	Id          string    `db:"product_id"`
	Name        string    `db:"name"`
	Price       float64   `db:"price"`
	Quantity    int64     `db:"quantity"`
	Description string    `db:"description"`
	Created_At  time.Time `db:"created_at"`
	Updated_At  time.Time `db:"updated_at"`
}

type ParamFindAllProduct struct {
	Id          string    `db:"product_id"`
	Name        string    `db:"name"`
	Price       float64   `db:"price"`
	Quantity    int64     `db:"quantity"`
	Description string    `db:"description"`
	Created_At  time.Time `db:"created_at"`
	Updated_At  time.Time `db:"updated_at"`
}

type ParamsUpdate struct {
	Id          string
	Name        string
	Price       float64
	Quantity    int64
	Description string
	Updated_At  time.Time
}

type ParamsUpdateResponse struct {
	Id          string    `db:"product_id"`
	Name        string    `db:"name"`
	Price       float64   `db:"price"`
	Quantity    int64     `db:"quantity"`
	Description string    `db:"description"`
	Created_At  time.Time `db:"created_at"`
	Updated_At  time.Time `db:"updated_at"`
}

type ParamsDelete struct {
	Id string
}
