package main

import (
	"flag"
	"fmt"
	"github.com/balabanovds/void/pkg/mail"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"log"
)

var (
	confFile string
	to       string
	subject  string
)

func init() {
	flag.StringVar(&confFile, "config-file", "./config/config.yml", "App config file")
	flag.StringVar(&to, "to", "denis.balabanov@nokia.com", "Send to")
	flag.StringVar(&subject, "subj", "Test email", "Email subject")
}

func main() {
	flag.Parse()

	cfg := mail.NewConfig()
	if err := unmarshalConfig(cfg); err != nil {
		log.Fatal(err)
	}

	r := mail.NewRequest(cfg, []string{to}, subject)
}

func unmarshalConfig(cfg *mail.Config) error {
	k := koanf.New(".")
	parser := yaml.Parser()
	if err := k.Load(file.Provider(confFile), parser); err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	if err := k.Unmarshal("", cfg); err != nil {
		return fmt.Errorf("error reading config: %v", err)
	}

	return nil
}
