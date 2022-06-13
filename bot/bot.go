package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var BotId string
var goBot *discordgo.Session

var req *http.Request
var client = &http.Client{}

type Response struct {
	Result float64 `json:"result"`
}

func init() {

}

func Start() {
	goBot, err := discordgo.New("Bot " + os.Getenv("TOKEN"))

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	BotId = u.ID

	goBot.AddHandler(messageHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Bot is running !")
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}
	fmt.Println(m.Content)
	if m.Content == "ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}
	if len(m.Content) == 7 && m.Content[3] == '2' {
		in := strings.ToUpper(m.Content[:3])
		out := strings.ToUpper(m.Content[4:])
		_, _ = s.ChannelMessageSend(m.ChannelID, convert(in, out))
	}
}

func convert(in, out string) string {

	req, _ = http.NewRequest("GET",
		fmt.Sprintf("https://api.apilayer.com/fixer/convert?amount=1&from=%s&to=%s", in, out),
		nil,
	)
	req.Header.Set("apikey", os.Getenv("APIKEY"))

	resp, _ := client.Do(req)
	bodyBytes, _ := io.ReadAll(resp.Body)
	data := Response{}
	json.Unmarshal(bodyBytes, &data)
	return fmt.Sprintf("%s/%s: %f", in, out, data.Result)

}
