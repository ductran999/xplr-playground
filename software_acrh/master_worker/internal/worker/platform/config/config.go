package config

type Config struct {
	ServerURL string `mapstructure:"server_url" validate:"required,url"`
}
