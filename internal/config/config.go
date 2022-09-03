package config

type Config struct {
	Service  ServiceConfig
	Server   ServerConfig
	Log      LogConfig
	Postgres PostgresConfig
}

type ServiceConfig struct {
	Name string `split_words:"true" required:"true"`
}

type ServerConfig struct {
	HttpAddress string `split_words:"true" required:"true"`
}

type LogConfig struct {
	Level string `split_words:"true" default:"INFO"`
}

type PostgresConfig struct {
	Uri                string `split_words:"true" required:"true"`
	ReconnectTimeout   int    `split_words:"true" default:"30"`
	MaxReconnectTrials int    `split_words:"true" default:"10"`
	MaxIdleConnections int    `split_words:"true" default:"5"`
	MaxOpenConnections int    `split_words:"true" default:"10"`
}
