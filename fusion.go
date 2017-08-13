package main

import (
	"fmt"
	"os"
	"os/signal"

	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/alde/fusion/config"
	"github.com/alde/fusion/db"
	"github.com/alde/fusion/server"
	"github.com/alde/fusion/version"

	"github.com/braintree/manners"
	"github.com/sirupsen/logrus"
)

var (
	printVersion   = kingpin.Flag("version", "Print version and exit").Short('v').Bool()
	configLocation = kingpin.Flag("config", "Location of config file").Short('c').String()
)

func setup() *config.Config {
	cfg := config.DefaultConfig()
	err := config.ReadConfigFile(cfg, config.GetConfigFilePath(*configLocation))
	if err != nil && len(*configLocation) > 0 {
		logrus.WithError(err).Fatal("Unable to parse config")
	}
	config.ReadEnvironment(cfg)

	return cfg
}

func setupLogging(logFormat, logLevel string) {
	if logFormat == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)
}

func main() {
	kingpin.Parse()
	if *printVersion {
		fmt.Println(version.Version)
		os.Exit(0)
	}
	config := setup()

	setupLogging(config.LogFormat, config.LogLevel)

	sql := db.New(config.Database)

	// Catch interrupt
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)
		s := <-c
		if s != os.Interrupt && s != os.Kill {
			return
		}

		logrus.Info("Fusion is shutting down")
		manners.Close()
	}()

	router := server.NewRouter(sql)

	bind := fmt.Sprintf("%s:%d", config.Address, config.Port)

	logrus.WithFields(logrus.Fields{
		"version": version.Version,
		"address": config.Address,
		"port":    config.Port,
	}).Infof("Launching Fusion version %s!\nListening to %s", version.Version, bind)

	err := manners.ListenAndServe(bind, router)
	if err != nil {
		logrus.WithError(err).Fatal("Unrecoverable error")
	}
}
