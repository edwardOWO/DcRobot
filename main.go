package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gpt3 "github.com/PullRequestInc/go-gpt3"
	"github.com/bwmarrin/discordgo"
)

var (
	DcToken      string
	ChatGptToken string
)

func init() {

	flag.StringVar(&DcToken, "d", "", "Bot Token")
	flag.StringVar(&ChatGptToken, "c", "", "Bot Token")
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
	dg, err := discordgo.New("Bot " + DcToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	private, _ := dg.UserChannelCreate("508893735854145556")

	fmt.Print(private.Name)

	// dg.ChannelMessageSend("1075065614051332116", "test")

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	//*discordgo.Session.ChannelMessageSend()

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

	// userID := "508893735854145556"

	// example for specific user to send 1 v 1 channels

	/*
		if m.Author.ID == userID {
			s.ChannelMessageSend(m.ChannelID, "Teemo 大人你好")
		}
	*/

	//s.ChannelMessageSend("1075065614051332116", "test")

	// example for specific channaleID
	if m.Content == "ping" && m.ChannelID == "1075065614051332116" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
		return
	}

	if len(m.Content) > 3 {
		if m.Content[0:5] == "!gpt " {

			question := m.Content[4:]

			result := chatgpt(question)
			s.ChannelMessageSend(m.ChannelID, result)
			return
		}
	}

	// list the channel "
	if m.Content == "list" {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("當前房間 ChannelID:  %s ", m.ChannelID))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("呼叫者ID m.Author.ID:  %s", m.Author.ID))
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("呼叫者名稱 m.Author.ID:  %s", m.Author.Username))
		return
	}

	if m.Content == "sendfile" {
		Sendfile(s, m.ChannelID)
		return
	}

	if m.Content == "issue" {

		getStock(s, m)
		return
	}

	if m.Content == "1v1" {
		if channel, err := s.UserChannelCreate(m.Author.ID); err != nil {
			fmt.Print(err.Error())
		} else {
			s.ChannelMessageSend(channel.ID, "1 v 1")
			go loopSend(s, channel.ID)
		}
	}
}
func loopSend(s *discordgo.Session, channels string) {

	for {
		time.Sleep(time.Second * 5)
		s.ChannelMessageSend(channels, "1 v 1")
		Sendfile(s, channels)
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

func Sendfile(s *discordgo.Session, channels string) {
	verificationImage, _ := os.Open("./resource/BabyElephantWalk60.wav")
	s.ChannelFileSendWithMessage(channels, "Send jpg example", "BabyElephantWalk60.wav", verificationImage)
	//s.ChannelFileSendWithMessage(channels, "Send jpg example", "lisa.jpg", verificationImage)
}
func chatgpt(question string) string {

	fmt.Print("test")
	ctx := context.Background()
	//client := gpt3.NewClient("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")
	client := gpt3.NewClient(ChatGptToken)

	test := ""
	test += "```\n"
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			question,
		},
		MaxTokens:   gpt3.IntPtr(4000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		fmt.Print(resp.Choices[0].Text)
		test += resp.Choices[0].Text

	})
	test += "```\n"
	if err != nil {
		return test
	}

	return test

}
