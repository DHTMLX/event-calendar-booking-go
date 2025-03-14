package main

import "event-calendar-booking/data"

type ConfigServer struct {
	URL            string
	Port           string
	Cors           []string
	ResetFrequence int `yaml:"resetFrequence"`
}

type AppConfig struct {
	Server ConfigServer
	DB     data.DBConfig
}
