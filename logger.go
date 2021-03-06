package main

import (
	"os"

	"github.com/Abramovic/logrus_influxdb"
	log "github.com/Sirupsen/logrus"
	client "github.com/influxdata/influxdb/client/v2"
)

func initLog() {
	log.SetOutput(os.Stdout)
}

// MUST BE CALLED AFTER CONFIG IS READ (ReadConfig())
func configureLog() {
	configureInflux()
}

func configureInflux() {
	if config.InfluxAvailable() {

		log.Info("Connecting to influxdb: " + config.GetInfluxHost())

		cnf := &logrus_influxdb.Config{
			Tags:     []string{"cbatch"}, // use the following tags
			Database: config.GetInfluxDatabase(),
		}

		c, err := client.NewHTTPClient(client.HTTPConfig{
			Addr:     config.GetInfluxHost(),
			Username: config.GetInfluxUser(),
			Password: config.GetInfluxPassword(),
		})

		if err != nil {
			log.Fatal(err)
		}

		hook, err := logrus_influxdb.NewInfluxDB(cnf, c)
		if err == nil {
			log.AddHook(hook)
		}

		log.Info("Connected to influxdb: " + config.GetInfluxHost())
	}

	log.Info("InfluxDB not configured")
}
