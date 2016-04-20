package main

// Config manages Croft configuration
type Config struct {
	AMQPUri         string `yaml:"amqp_uri"`
	MongoDBUri      string `yaml:"mongodb_uri"`
	MongoDBDatabase string `yaml:"mongodb_database"`
	MqttBroker      string `yaml:"mqtt_broker"`
	MqttClientID    string `yaml:"mqtt_client_id"`
}

/*
func (c Config) Parsezz(data []byte) error {
	c.OrionBaseURL = "http://OrionBaseURL"
	//err := yaml.Unmarshal(data, &c)
	//fmt.Printf("- - - %s\n", data)
	fmt.Printf("- - - %+v\n", c)
	return nil
}
*/
