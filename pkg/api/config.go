package api

import (
	"log"

	"github.com/jasosa/backend/models"
)

type DBConfig struct {
	DSN string
}

type ApiConfig struct {
	Port   int
	Env    string
	Db     *DBConfig
	Jwt    string
	Logger *log.Logger
	Models models.Models
}

var App *ApiConfig

func New(api *ApiConfig) error {
	App = api
	return nil
}
