package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Create a new Discord session using the provided bot token.
	if fileExists(os.Getenv("MEMEARCH")) {
		dat, err := ioutil.ReadFile(os.Getenv("MEMEARCH"))
		if err != nil {
			return
		}
		json.Unmarshal(dat, &cases)
	}
	if fileExists(os.Getenv("TRIVIAQUESTIONS")) {
		unarchiveJSON(os.Getenv("TRIVIAQUESTIONS"), &Questions)
	}
	if fileExists(os.Getenv("TRIVIATEAMS")) {
		unarchiveJSON(os.Getenv("TRIVIATEAMS"), &teams)
	}
	dg, err := discordgo.New("Bot " + os.Getenv("BOTID"))
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(MessageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
