package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"

	"github.com/soliel/BulliBot/command"
	"github.com/soliel/BulliBot/configuration"
)

var (
	randomGen   *rand.Rand
	bulliMode   bool
	bulliChance int = 70

	conf    *configuration.BotConfig
	handler *command.Handler
)

func main() {
	loadBotConfBytes, err := ioutil.ReadFile("../ConfigurationFiles/BulliConf.json")
	if err != nil {
		fmt.Println(err)
		return
	}

	conf = new(configuration.BotConfig)
	err = conf.LoadConfig(loadBotConfBytes)
	if err != nil {
		fmt.Println("Error getting bot congifuration: ", err)
	}

	dg, err := discordgo.New("Bot " + conf.BotToken)
	if err != nil {
		fmt.Println("Error starting discord session: ", err)
		return
	}

	dg.AddHandler(onMessageReceived)

	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening communication with discord: ", err)
		return
	}
	defer dg.Close()

	randomGen = rand.New(rand.NewSource(time.Now().UnixNano()))
	handler = command.CreateHandler()
	registerCommands()

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

func onMessageReceived(s *discordgo.Session, m *discordgo.MessageCreate) {
	command := filterMessages(s, m)

	if command.Command == "" {
		return
	}

	handler.HandleCommand(m, s, command)
}

func registerCommands() {
	handler.Register("ping", ping)
}

func oldboi(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != "168388789591343104" && m.Content == "toggle bulli" {
		bulliMode = !bulliMode
		bulliStr := "Online"

		if !bulliMode {
			bulliStr = "Offline"
		}

		s.ChannelMessageSend(m.ChannelID, "Bulli Mode: "+bulliStr)
	}

	if m.Author.ID != "168388789591343104" && strings.HasPrefix(m.Content, "set bulli chance ") {
		contentStr := strings.Split(m.Content, "chance ")
		intChance, err := strconv.Atoi(contentStr[1])
		if err != nil {
			return
		}

		bulliChance = 100 - intChance
	}

	if !bulliMode {
		return
	}

	if m.Author.ID == "168388789591343104" {
		randomInt := randomGen.Intn(100)
		fmt.Println(randomInt)
		if randomInt > bulliChance {
			s.ChannelMessageSend(m.ChannelID, "Zach is gay and a Sub.")
		}
	}
}
