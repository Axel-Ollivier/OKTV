package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	token     = os.Getenv("DISCORD_TOKEN")
	guildID   = os.Getenv("GUILD_ID")
	channelID = os.Getenv("CHANNEL_ID")

	mu       sync.RWMutex
	statuses = map[string]string{}
	updateCh = make(chan struct{}, 1)
	messageID string
)

func main() {
	if token == "" || guildID == "" || channelID == "" {
		panic("DISCORD_TOKEN, GUILD_ID et CHANNEL_ID requis")
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		panic(err)
	}

	dg.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildPresences |
		discordgo.IntentsGuildMembers

	dg.AddHandler(onReady)
	dg.AddHandler(onGuildCreate)
	dg.AddHandler(onPresenceUpdate)

	if err := dg.Open(); err != nil {
		panic(err)
	}
	defer dg.Close()

	go debouncedUpdater(dg)

	<-make(chan struct{}) // bloque le programme
}

func onReady(s *discordgo.Session, r *discordgo.Ready) {
	ensureMessage(s)
	triggerUpdate()
}

func onGuildCreate(s *discordgo.Session, gc *discordgo.GuildCreate) {
	if gc.Guild.ID != guildID {
		return
	}
	mu.Lock()
	for _, p := range gc.Presences {
		if p.User != nil {
			statuses[p.User.ID] = string(p.Status)
		}
	}
	mu.Unlock()
	triggerUpdate()
}

func onPresenceUpdate(s *discordgo.Session, pu *discordgo.PresenceUpdate) {
	if pu.GuildID != guildID || pu.User == nil {
		return
	}
	mu.Lock()
	statuses[pu.User.ID] = string(pu.Status)
	mu.Unlock()
	triggerUpdate()
}

func debouncedUpdater(s *discordgo.Session) {
	var timer *time.Timer
	for range updateCh {
		if timer != nil {
			timer.Stop()
		}
		timer = time.AfterFunc(800*time.Millisecond, func() {
			count := onlineCount()
			msg := fmt.Sprintf("🟢 Membres en ligne : **%d**", count)

			if !ensureMessage(s) {
				// Si pas de message, il est créé dans ensureMessage
				return
			}

			_, err := s.ChannelMessageEdit(channelID, messageID, msg)
			if err != nil {
				// Si le message a été supprimé, on en refait un
				messageID = ""
				ensureMessage(s)
				s.ChannelMessageEdit(channelID, messageID, msg)
			}
		})
	}
}

func onlineCount() int {
	mu.RLock()
	defer mu.RUnlock()
	n := 0
	for _, st := range statuses {
		if st != "offline" && st != "" {
			n++
		}
	}
	return n
}

func triggerUpdate() {
	select {
	case updateCh <- struct{}{}:
	default:
	}
}

func ensureMessage(s *discordgo.Session) bool {
	if messageID != "" {
		return true
	}
	msg, err := s.ChannelMessageSend(channelID, "🟢 Membres en ligne : ...")
	if err != nil {
		fmt.Println("Erreur création message:", err)
		return false
	}
	messageID = msg.ID
	return true
}
