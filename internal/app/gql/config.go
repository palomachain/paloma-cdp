package gql

type Configuration struct {
	Port string `env:"CDP_GQL_PORT,notEmpty"`
}
