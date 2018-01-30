package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Necroforger/discordtranscriber"
	"github.com/bwmarrin/discordgo"
)

// Generate bindata
//go:generate go-bindata-assetfs assets/...

// Flags
var (
	port     = flag.String("port", "8100", "server port")
	username = flag.String("u", "", "discord username")
	password = flag.String("p", "", "discord password")
	token    = flag.String("t", "", "discord token")
	dir      = flag.String("d", "", "asset directory")
)

// interactive input mode
func interactive() {
	rd := bufio.NewReader(os.Stdin)
	option := strings.TrimSpace(queryLine(rd, "[1] Username and password\n[2] Token\nenter an option: "))
	switch option {
	case "1":
		*username = strings.TrimSpace(queryLine(rd, "Username: "))
		*password = strings.TrimSpace(queryLine(rd, "Password: "))
	case "2":
		*token = strings.TrimSpace(queryLine(rd, "Token: "))
	}
}

func queryLine(rd *bufio.Reader, query string) string {
	fmt.Printf(query)
	line, err := rd.ReadString('\n')
	if err != nil {
		log.Println(err)
	}
	return line
}

func main() {
	flag.Parse()

	// Alternative flags
	if *username == "" {
		*username = flag.Arg(0)
	}
	if *password == "" {
		*password = flag.Arg(1)
	}

	// Enter interactive mode if login credentials are not set
	if *username == "" && *password == "" && *token == "" {
		interactive()
	}

	// Set the fileSystem
	var fileSystem http.FileSystem
	if *dir != "" {
		fileSystem = http.Dir(*dir)
	} else {
		fileSystem = assetFS()
	}

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
	log.Println("Visit http://localhost:" + *port + "/ in a WebSpeech API supporting browser (chrome)")
	server := discordtranscriber.NewServer(c, *port, fileSystem)
	server.Log = true // Enable logging
	log.Fatal(server.ListenAndServe())
}
