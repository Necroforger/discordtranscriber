package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/googollee/go-socket.io"
)

var (
	port      = flag.String("port", "8100", "port to serve on")
	directory = flag.String("d", ".", "the directory of static file to host")
	username  = flag.String("u", "", "discord username")
	password  = flag.String("p", "", "discord password")
	token     = flag.String("t", "", "discord token")
)

// SendRequest represents a request to send a message
type SendRequest struct {
	ChannelID string `json:"channelID"`
	Content   string `json:"content"`
}

func main() {
	flag.Parse()

	client, err := discordgo.New(*username, *password, *token)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Open()
	if err != nil {
		log.Fatal(err)
	}

	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Println(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("Socket connected: ", so.Id())

		so.On("guild", func(msg string) *discordgo.Guild {
			log.Println("guild requested: " + msg)
			guild, _ := client.State.Guild(msg)
			return guild
		})
		so.On("channel", func(msg string) *discordgo.Channel {
			log.Println("channel requested: " + msg)
			channel, _ := client.State.Channel(msg)
			return channel
		})
		so.On("me", func() *discordgo.User {
			log.Println("me user requested")
			return client.State.User
		})
		so.On("myAvatar", func() string {
			log.Println("avatar requested")
			if client.State.User != nil {
				return client.State.User.AvatarURL("")
			}
			return ""
		})
		so.On("send-text", func(req SendRequest) {
			log.Println("Sending: ", req)
			_, err := client.ChannelMessageSend(req.ChannelID, req.Content)
			if err != nil {
				log.Println(err)
			}
		})
	})

	http.Handle("/", http.FileServer(http.Dir(*directory)))
	http.Handle("/socket.io/", server)
	log.Printf("Serving %s on HTTP port: %s\n", *directory, *port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
