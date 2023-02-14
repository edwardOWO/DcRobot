package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

//MTA3NTA1MjkyODAwOTc4MTMwOA.GkLJ1B.LDj4V0gy0aji2HpOLeBSrKhyQSgKXkAyEzfsOo
func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

type stock []struct {
	NAMING_FAILED  string `json:"發言日期"`
	NAMING_FAILED0 string `json:"發言時間"`
	NAMING_FAILED1 string `json:"出表日期"`
	NAMING_FAILED2 string `json:"公司代號"`
	NAMING_FAILED3 string `json:"公司名稱"`
	NAMING_FAILED4 string `json:"主旨 "`
	NAMING_FAILED5 string `json:"符合條款"`
	NAMING_FAILED6 string `json:"事實發生日"`
	NAMING_FAILED7 string `json:"說明"`
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.

	userID := "508893735854145556"

	// example for specific user
	if m.Author.ID == userID {
		s.ChannelMessageSend(m.ChannelID, "Teemo 大人你好")
	}

	// example for specific channaleID
	if m.Content == "ping" && m.ChannelID == "1075065614051332116" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
		return
	}

	// list the channel "
	if m.Content == "list" {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("當前房間 ChannelID:  %s ", m.ChannelID))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("呼叫者ID m.Author.ID:  %s", m.Author.ID))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("呼叫者名稱 m.Author.ID:  %s", m.Author.Username))
		return
	}

	if m.Content == "issue" {

		getStock(s, m)
		return
	}
}
func getStock(s *discordgo.Session, m *discordgo.MessageCreate) {

	url := "https://openapi.twse.com.tw/v1/opendata/t187ap04_L"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Set("accept", "application/json")
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, _ := client.Do(req)

	b, _ := io.ReadAll(res.Body)

	stock1 := stock{}
	json.Unmarshal(b, &stock1)

	sum := ""
	for _, issue := range stock1 {

		data := fmt.Sprintf("%s %s %s\n", issue.NAMING_FAILED2, issue.NAMING_FAILED3, issue.NAMING_FAILED4)
		sum += data
		//data = fmt.Sprintf("%s\n", issue.NAMING_FAILED4)
		sum += ("\n")
		s.ChannelMessageSend(m.ChannelID, sum)
		time.Sleep(time.Microsecond * 100)

	}

}
