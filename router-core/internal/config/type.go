package config

type Config struct {
	Server struct {
		Env          string `mapstructure:"env"`
		Port         string `mapstructure:"port"`
		ExternalPort string `mapstructure:"external_port"`
		Host         string `mapstructure:"host"`
		RunMode      string `mapstructure:"run_mode"`
	} `mapstructure:"server"`

	Database struct {
		Postgres struct {
			Connection string `mapstructure:"connection"`
		} `mapstructure:"postgres"`
	} `mapstructure:"database"`

	Selector struct {
		URL string `mapstructure:"url"`
	} `mapstructure:"selector"`

	AWS struct {
		Region          string `mapstructure:"region"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
	} `mapstructure:"aws"`
}
