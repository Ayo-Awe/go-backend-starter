package main

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Port int    `envconfig:"PORT" default:"4000"`
	Env  string `envconfig:"ENV" default:"development"`
	DB   struct {
		Dsn string `envconfig:"DATABASE_URL" required:"true"`
	}
	Limiter struct {
		Enabled        bool          `envconfig:"LIMITER_ENABLED" default:"false"`
		Rate           int           `envconfig:"LIMITER_RATE" default:"100"`
		Window         time.Duration `envconfig:"LIMITER_PERIOD" default:"1m"`
		WhitelistedIPs []string      `envconfig:"LIMITER_WHITELISTED_IPS" default:""`
	}
	Cors struct {
		AllowedOrigins []string `envconfig:"CORS_ALLOWED_ORIGINS" default:""`
	}
	SMTP struct {
		HOST     string `envconfig:"SMTP_HOST" default:""`
		PORT     int    `envconfig:"SMTP_PORT" default:"587"`
		USERNAME string `envconfig:"SMTP_USERNAME" default:""`
		PASSWORD string `envconfig:"SMTP_PASSWORD" default:""`
		SENDER   string `envconfig:"SMTP_SENDER" default:""`
	}
	// Google struct {
	// 	ClientID             string   `envconfig:"GOOGLE_CLIENT_ID" required:"true"`
	// 	ClientSecret         string   `envconfig:"GOOGLE_CLIENT_SECRET" required:"true"`
	// 	DefaultRedirectURL   string   `envconfig:"GOOGLE_CALLBACK_URL" required:"true"`
	// 	WhitelistedRedirects []string `envconfig:"GOOGLE_WHITELISTED_REDIRECTS" required:"false"`
	// }
}

func loadConfig() (config, error) {
	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
