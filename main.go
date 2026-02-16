package main

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	logFile *os.File
	err     error
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if logFile, err = os.OpenFile(".log", os.O_CREATE|os.O_WRONLY, 0777); err != nil {
		panic(err)
	}
	log.Logger = log.Output(io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stdout}, logFile))
}

func main() {
	defer deinit()
}

func deinit() {
	if err = logFile.Close(); err != nil {
		panic(err)
	}
}
