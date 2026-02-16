package main

import (
	"flag"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
)

var (
	logFile *os.File
	err     error
	isDebug bool
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if logFile, err = os.OpenFile(".log", os.O_CREATE|os.O_WRONLY, 0777); err != nil {
		panic(err)
	}
	log.Logger = log.Output(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout}, logFile))
	flag.BoolVar(&isDebug, "debug", false, "debug mode")
}

func main() {
	defer deinit()
	flag.Parse()
	if isDebug {
		err := debug.Run("correctum-agent", &service{})
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
	} else {
		err := svc.Run("correctum-agent", &service{})
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
	}
}

func deinit() {
	if err = logFile.Close(); err != nil {
		panic(err)
	}
}
