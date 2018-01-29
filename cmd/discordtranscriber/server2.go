package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/Necroforger/discordtranscriber"
	"github.com/bwmarrin/discordgo"
)

// Flags
var (
	port     = flag.String("port", "8100", "server port")
	username = flag.String("u", "", "discord username")
	password = flag.String("p", "", "discord password")
	token    = flag.String("t", "", "discord token")
	dir      = flag.String("d", "assets", "asset directory")
)

func main() {
	flag.Parse()

	// Login to discord
	c, err := discordgo.New(*username, *password, *token)
	if err != nil {
		log.Fatal(err)
	}

	waitForReady := make(chan struct{})
	c.AddHandlerOnce(func(s *discordgo.Session, ready *discordgo.Ready) {
		waitForReady <- struct{}{}
	})

	// Connect to discord Gateway
	err = c.Open()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Succesfully connected to discord")

	<-waitForReady
	log.Println("Recieved ready data")

	log.Println("Server listening on port [", *port, "]")
	server := discordtranscriber.NewServer(c, *port, http.Dir(*dir))
	server.Log = true // Enable logging
	log.Fatal(server.ListenAndServe())
}
