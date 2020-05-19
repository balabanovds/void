package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/server"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/rs/zerolog"
)

var (
	k          = koanf.New("::")
	parser     = yaml.Parser()
	confFile   string
	dbPassword string
	debug      bool
)

type app struct {
	logger  *zerolog.Logger
	server  server.Server
	storage domain.Storage
}

type appConfig struct {
	StorageCfg *domain.Config `koanf:"storage"`
	ServerCfg  *server.Config `koanf:"server"`
	LogCfg     *logCfg        `koanf:"logging"`
}

type logCfg struct {
	Level string `koanf:"level"`
	File  string `koanf:"file"`
}

func init() {
	flag.StringVar(&confFile, "config", "./config/config.yml", "App config file")
	flag.StringVar(&dbPassword, "db-password", "", "DB password")
	flag.BoolVar(&debug, "debug", false, "Enable debug")
}

func main() {
	appCfg := appConfig{
		StorageCfg: domain.NewConfig(),
		ServerCfg:  server.NewConfig(),
	}

	flag.Parse()

	if dbPassword == "" {
		log.Println("db password was not set explicitly, default used")
	} else {
		appCfg.StorageCfg.SQL.Password = dbPassword
	}

	if err := k.Load(file.Provider(confFile), parser); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	if err := k.Unmarshal("", &appCfg); err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	fmt.Printf("%+v\n%+v\n%+v", appCfg.StorageCfg, appCfg.ServerCfg, appCfg.LogCfg)
}
