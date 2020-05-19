package domain

// Config of domain
type Config struct {
	SQL sqlConfig `koanf:"sql"`
}

type sqlConfig struct {
	Driver   string   `koanf:"driver"`
	Host     string   `koanf:"host"`
	Port     int      `koanf:"port"`
	User     string   `koanf:"user"`
	Password string   `koanf:"-"`
	Options  []string `koanf:"options"`
}

// NewConfig of domain
func NewConfig() *Config {
	return &Config{
		SQL: sqlConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     5432,
			User:     "void",
			Password: "void123",
		},
	}
}
