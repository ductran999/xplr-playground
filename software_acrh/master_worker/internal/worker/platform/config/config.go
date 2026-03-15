package config

type Config struct {
	ServerURL         string `mapstructure:"server_url" validate:"required,url"`
	Version           string `mapstructure:"version" validate:"required"`
	RegistrationToken string `mapstructure:"registration_token" validate:"required"`
}
