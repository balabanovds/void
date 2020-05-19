package server

// Config of server
type Config struct {
	Hostname string `koanf:"hostname"`
	Port     int    `koanf:"port"`
}

// NewConfig of server
func NewConfig() *Config {
	return &Config{
		Hostname: "localhost",
		Port:     50000,
	}
}
