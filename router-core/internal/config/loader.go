package config

import (
	"fmt"

	"github.com/spf13/viper"
)

func LoadConfig() (*Config, error) {
	v := viper.New()
	v.AutomaticEnv()

	bindEnvVars(v)

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

func bindEnvVars(v *viper.Viper) {
	envMap := map[string]string{
		"server.env":           "ENV",
		"server.port":          "PORT",
		"server.external_port": "EXTERNAL_PORT",
		"server.host":          "HOST",
		"server.run_mode":      "RUN_MODE",

		"database.postgres.connection": "MAIN_DB_CONNECTION",

		"aws.region":            "AWS_REGION",
		"aws.access_key_id":     "AWS_ACCESS_KEY_ID",
		"aws.secret_access_key": "AWS_SECRET_ACCESS_KEY",

		"selector.url": "SELECTOR_SERVICE_URL",
	}

	for key, env := range envMap {
		_ = v.BindEnv(key, env)
	}
}
