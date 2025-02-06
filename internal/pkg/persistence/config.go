package persistence

import "time"

type Configuration struct {
	Address  string        `env:"CDP_PSQL_ADDRESS,notEmpty"`
	User     string        `env:"CDP_PSQL_USER,notEmpty"`
	Password string        `env:"CDP_PSQL_PASSWORD,unset"`
	Database string        `env:"CDP_PSQL_DATABASE,notEmpty"`
	Timeout  time.Duration `env:"CDP_PSQL_TIMEOUT" envDefault:"5s"`
}
