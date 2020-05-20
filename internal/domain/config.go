package domain

// Config of domain
type Config struct {
	SQL sqlConfig `koanf:"sql"`
}

type sqlConfig struct {
	Host     string   `koanf:"host"`
	Port     int      `koanf:"port"`
	User     string   `koanf:"user"`
	Password string   `koanf:"-"`
	DBName   string   `koanf:"db-name"`
	Options  []string `koanf:"options"`
}

// NewConfig of domain
func NewConfig() *Config {
	return &Config{
		SQL: sqlConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "void",
			Password: "void123",
			DBName:   "void_test",
		},
	}
}
