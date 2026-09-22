package config

type Base struct {
	ServiceName string `env:"SERVICE_NAME,required"`
	AppEnv      string `env:"APP_ENV" envDefault:"development"`
	Host        string `env:"HTTP_HOST" envDefault:"0.0.0.0"`
	Port        string `env:"HTTP_PORT" envDefault:"9000"`
	LogLevel    string `env:"LOG_LEVEL" envDefault:"info"`
}
