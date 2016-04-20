package main

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

// Config manages Croft configuration
type Config struct {
	AMQPUri         string `yaml:"amqp_uri"`
	MongoDBUri      string `yaml:"mongodb_uri"`
	MongoDBDatabase string `yaml:"mongodb_database"`
	MqttBroker      string `yaml:"mqtt_broker"`
	MqttClientID    string `yaml:"mqtt_client_id"`
}

func (c *Config) ParseYamlFile(configFilePath string) error {
	cfgData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("error reading %s: %v", configFilePath, err)
		return err
	}
	err = yaml.Unmarshal(cfgData, c)
	if err != nil {
		log.Fatalf("error parsing %s: %v", configFilePath, err)
		return err
	}
	log.Printf("Jolie config:\n%+v\n\n", c)
	return nil
}
