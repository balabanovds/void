package mail

type Config struct {
	Server       string `koanf:"server"`
	Sender       string `koanf:"sender"`
	TemplatesDir string `koanf:"mail-templates-dir"`
}

func NewConfig() *Config {
	return &Config{
		Server:       "localhost",
		Sender:       "foo@bar.baz",
		TemplatesDir: "./templates/mail",
	}
}
