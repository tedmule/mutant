package main

import (
	"github.com/caarlos0/env/v9"
	"github.com/daddvted/mutant/mutant"
	log "github.com/sirupsen/logrus"
)

var (
	config mutant.MutantConfig
)

func init() {
	//Parse configuration from environment variables
	if err := env.Parse(&config); err != nil {
		log.Error(err.Error())
	}
	log.Debugf("%+v\n", config)

	// Init Logrus, default to INFO
	if config.Production {
		log.SetFormatter(&log.JSONFormatter{})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05.00000",
		})

	}
	// log.SetFormatter(&log.JSONFormatter{})
	logLvl, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		// logLvl = log.InfoLevel
		logLvl = log.DebugLevel
	}
	log.SetLevel(logLvl)
	log.SetReportCaller(true)
}

// Starts the webserver and serves the mutate function.
func main() {

	nfsMT := mutant.MutantCSI{
		StorageClass: "nfs.csi.k8s.io",
		Annotation:   "nomutant",
	}

	nfsWebhook, err := mutant.NewMutantWebhook(&nfsMT, config)
	if err != nil {
		log.Fatal(err)
	}

	nfsWebhook.Server.Logger.Fatal(nfsWebhook.Server.StartTLS(config.Listen, config.CertFile, config.KeyFile))
}
