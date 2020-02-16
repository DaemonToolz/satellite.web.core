package main

import (
	"encoding/json"
	"log"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/streadway/amqp"
)

var server *gosocketio.Server

func initSocketServer() {
	server = gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	//handle connected
	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		var broadcastRoom string = "general-notification"
		var mySpaceNotif string = "myspace"
		/*
			if c.RequestHeader().Get("username") != "" {
				mySpaceNotif = c.RequestHeader().Get("username")
			}
		*/
		c.Join(broadcastRoom)
		c.Join(mySpaceNotif)
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("Disconnected")
	})

}

func BroadcastTo(message amqp.Delivery) {
	log.Printf(" [x] %s", message.Exchange)
	log.Printf(" [x] %s", message.RoutingKey)
	log.Printf(" [x] %s", message.Body)

	var content RabbitMqMsg

	if err := json.Unmarshal(message.Body, &content); err != nil {
		failOnError(err, "Couldn't unmarshal the message")
	} else {
		server.BroadcastTo("myspace", "update", content)
		//server.BroadcastToRoom(content.To, "/send", content)
	}

}
