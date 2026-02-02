package config

import (
	"log"
	"os"

	"github.com/naoina/toml"
)

type Config struct {
	Info struct {
		Port     string
		Issuer   string
		Audience string
	} `toml:"info"`

	Cassandra struct {
		Keyspace string
	} `toml:"cassandra"`

	Kafka struct {
		ProducerClientId string
	} `toml:"kafka"`
}

func NewConfig(path string) *Config {
	c := new(Config)

	f, err := os.Open(path)
	if err != nil {
		log.Panicf("Fail to open config path: %v", err)
	}
	err = toml.NewDecoder(f).Decode(c)
	if err != nil {
		log.Panicf("Fail to decode toml: %v", err)
	}
	return c
}
