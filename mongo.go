package main

import (
	"log"
	"time"

	"github.com/thethingsnetwork/server-shared"
	"gopkg.in/mgo.v2"
)

const (
	MONGODB_ATTEMPTS = 20
)

type MongoDatabase struct {
	session *mgo.Session
}

func ConnectMongoDatabase() (PacketHandler, error) {
	var err error
	//uri := os.Getenv("MONGODB_URI")
	uri := config.MongoDBUri
	for i := 0; i < MONGODB_ATTEMPTS; i++ {
		var session *mgo.Session
		//session, err = mgo.Dial(fmt.Sprintf("%s:27017", uri))
		session, err = mgo.Dial(uri) // config.MongoDBUri sets both server and port
		if err != nil {
			log.Printf("Failed to connect to %s: %s", uri, err.Error())
			time.Sleep(time.Duration(2) * time.Second)
		} else {
			log.Printf("Connected to %s", uri)
			session.SetMode(mgo.Monotonic, true)
			session.SetSocketTimeout(time.Duration(6) * time.Second)
			session.SetSyncTimeout(time.Duration(6) * time.Second)
			return &MongoDatabase{session}, nil
		}
	}
	return nil, err
}

func (db *MongoDatabase) Configure() error {
	return nil
}

func (db *MongoDatabase) HandleStatus(status *shared.GatewayStatus) {
	//err := db.session.DB("jolie").C("gateway_statuses").Insert(status)
	err := db.session.DB(config.MongoDBDatabase).C("gateway_statuses").Insert(status)
	if err != nil {
		log.Printf("Failed to save status: %s", err.Error())
	}
}

func (db *MongoDatabase) HandlePacket(packet *shared.RxPacket) {
	//err := db.session.DB("jolie").C("rx_packets").Insert(packet)
	err := db.session.DB(config.MongoDBDatabase).C("rx_packets").Insert(packet)
	if err != nil {
		log.Printf("Failed to save packet: %s", err.Error())
	}
}

func (db *MongoDatabase) RecordGatewayStatus(status *shared.GatewayStatus) error {
	return nil
}

func (db *MongoDatabase) RecordRxPacket(packet *shared.RxPacket) error {
	return nil
}
