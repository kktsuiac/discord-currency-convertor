package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	BotId   string
	client  = &http.Client{}
	cache   = make(map[string]*FixerResponse)
	timeout = 10 * time.Minute
)

type FixerResponse struct {
	Success bool    `json:"success`
	Query   Query   `json:"query"`
	Info    Info    `json:"info"`
	Date    string  `json:"date"`
	Result  float64 `json:"result"`
}

type Info struct {
	Timestamp int64   `json:"timestamp"`
	Rate      float64 `json:"rate"`
}
type Query struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
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

	if len(m.Content) == 7 && m.Content[3] == '2' {
		in := strings.ToUpper(m.Content[:3])
		out := strings.ToUpper(m.Content[4:])
		resp, err := convert(in, out)
		var message string
		if err != nil {
			message = err.Error()
		} else {
			message = resp.String()
		}
		_, _ = s.ChannelMessageSend(m.ChannelID, message)
	}
}

func convert(in, out string) (*FixerResponse, error) {
	if resp, exist := cache[in+out]; exist {
		if !time.Now().After(time.Unix(resp.Info.Timestamp, 0).Add(timeout)) {
			return resp, nil
		}
	}
	req, _ := http.NewRequest("GET",
		fmt.Sprintf("https://api.apilayer.com/fixer/convert?amount=1&from=%s&to=%s", in, out),
		nil,
	)
	req.Header.Set("apikey", os.Getenv("APIKEY"))
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	data := FixerResponse{}
	json.Unmarshal(bodyBytes, &data)
	cache[in+out] = &data
	return &data, nil
}

func (r *FixerResponse) String() string {
	return fmt.Sprintf(`%s/%s: 
	Rate: %f
	Last Update: %s`,
		r.Query.From, r.Query.To,
		r.Info.Rate,
		time.Unix(r.Info.Timestamp, 0).Format(time.UnixDate))
}
