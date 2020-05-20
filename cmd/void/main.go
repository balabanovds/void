package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/balabanovds/void/internal/domain"
	"github.com/balabanovds/void/internal/domain/pgsql"
	"github.com/balabanovds/void/internal/server"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/rs/zerolog"
)

type app struct {
	logger  *zerolog.Logger
	server  server.Server
	storage domain.Storage
}

type appConfig struct {
	StorageCfg *domain.Config `koanf:"storage"`
	ServerCfg  *server.Config `koanf:"server"`
}

var (
	confFile   string
	dbPassword string
	debug      bool
	logFile    string
)

func init() {
	flag.StringVar(&confFile, "config-file", "./config/config.yml", "App config file")
	flag.StringVar(&logFile, "log-file", "", "Log file")
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

	// parse config file
	if err := unmarshalConfig(&appCfg); err != nil {
		log.Fatal(err)
	}

	// logger
	logger, err := newLogger(logFile)
	if err != nil {
		log.Fatal(err)
	}

	// Init Postgres storage
	storage := pgsql.New(appCfg.StorageCfg, logger)
	if err := storage.Open(); err != nil {
		logger.Fatal().Str("main", "DB").Msg(err.Error())
	}
	defer storage.Close()

	logger.Info().Msgf("connected to '%s' at %s:%d\n",
		appCfg.StorageCfg.SQL.DBName,
		appCfg.StorageCfg.SQL.Host,
		appCfg.StorageCfg.SQL.Port,
	)
	// Init server
}

func unmarshalConfig(cfg *appConfig) error {
	k := koanf.New("::")
	parser := yaml.Parser()
	if err := k.Load(file.Provider(confFile), parser); err != nil {
		return fmt.Errorf("error loading config: %v", err)
	}

	if err := k.Unmarshal("", cfg); err != nil {
		return fmt.Errorf("error reading config: %v", err)
	}

	return nil
}

func newLogger(file string) (zerolog.Logger, error) {
	var w io.Writer
	var err error
	noColor := true

	if file == "" {
		w = os.Stderr
		noColor = false
	} else {
		w, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return zerolog.Nop(), err
		}
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	output := zerolog.ConsoleWriter{
		Out:        w,
		TimeFormat: time.RFC3339,
		NoColor:    noColor,
	}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}

	return zerolog.New(output).With().Timestamp().Logger(), nil
}
