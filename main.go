package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"

	"github.com/thethingsnetwork/server-shared"
)

var (
	config   Config
	consumer Consumer
	mqtt     PacketHandler
	database PacketHandler
)

func main() {
	log.Print("Jolie is ALIVE")

	//// go run */*.go -config=jolie_config.yml
	var configFilePath string
	flag.StringVar(&configFilePath, "config", "/non/existent/filez", "the YAML config file")
	flag.Parse()
	cfgData, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	config = Config{}
	err = yaml.Unmarshal([]byte(cfgData), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	/*
		if err = config.Parse(cfgData); err != nil {
			log.Fatal(err)
		}
	*/
	fmt.Printf("config:\n%+v\n\n", config)

	err = connectConsumer()
	if err != nil {
		log.Fatalf("Failed to connect consumer: %s", err.Error())
	}

	queues, err := consumer.Consume()
	if err != nil {
		log.Fatalf("Failed to consume queues: %s", err.Error())
	}

	database, err = connectDatabase(queues)
	if err != nil {
		log.Fatalf("Failed to connect database: %s", err.Error())
	}

	mqtt, err = connectMqtt(queues)
	if err != nil {
		log.Fatalf("Failed to connect MQTT: %s", err.Error())
	}

	go Handle(queues, []PacketHandler{database, mqtt})

	select {}
}

func connectConsumer() error {
	var err error
	consumer, err = ConnectRabbitConsumer()
	if err != nil {
		return err
	}

	err = consumer.Configure()
	if err != nil {
		return err
	}

	return nil
}

func connectDatabase(queues *shared.ConsumerQueues) (PacketHandler, error) {
	var err error
	database, err = ConnectMongoDatabase()
	if err != nil {
		return nil, err
	}

	return database, nil
}

func connectMqtt(queues *shared.ConsumerQueues) (PacketHandler, error) {
	var err error
	mqtt, err = ConnectPaho()
	if err != nil {
		return nil, err
	}

	return mqtt, nil
}

func Handle(queues *shared.ConsumerQueues, handlers []PacketHandler) {
	for {
		select {
		case status := <-queues.GatewayStatuses:
			for _, h := range handlers {
				h.HandleStatus(status)
			}
		case packet := <-queues.RxPackets:
			for _, h := range handlers {
				h.HandlePacket(packet)
			}
		}
	}
}
