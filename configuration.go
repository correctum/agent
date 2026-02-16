package main

import (
	"encoding/json"
	"os"
	"time"
)

var config configuration

func loadConfiguration() error {
	data, err := os.ReadFile(pathConfiguration)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	return nil
}

type configuration struct {
	begin time.Time
	end   time.Time
}
