package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

//Penalty role
type Penalty struct {
	ID       string
	Expires  bool
	Duration time.Time
	Debounce time.Time
}

var cassMap = map[string]Penalty{}

func checkCass(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	fmt.Println(cassMap[m.Author.ID].ID)
	if m.Message.Member != nil && containsVal(m.Message.Member.Roles, cassanova) < 0 && len(cassMap[m.Author.ID].ID) > 0 {
		s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, cassanova)
		delete(cassMap, m.Author.ID)
		s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, frogkook)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@!%s> is now <@&%s>ed.", m.Author.ID, frogkook))
		return false
	}
	if m.Message.Member != nil && containsVal(m.Message.Member.Roles, cassanova) > -1 {
		if isCassExpired(m.Author.ID) {
			s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, cassanova)
			delete(cassMap, m.Author.ID)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@!%s> is no longer <@&%s>ed.", m.Author.ID, cassanova))
			return false
		}

		v := cassMap[m.Author.ID]
		if v.Duration.Sub(time.Now()) > (time.Hour * 24) {
			s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, "636246344855453696")
		}
		v.Duration = v.Duration.Add(5 * time.Minute)
		var debounced bool
		if v.Debounce.Before(time.Now()) {
			v.Debounce = time.Now().Add(time.Minute)
			debounced = true
		}
		cassMap[m.Author.ID] = v
		rand.Seed(time.Now().UnixNano())
		cc := rand.Intn(len(jeremeVids))
		s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		dd := v.Duration.Sub(time.Now())
		if debounced {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(jeremeVids[cc], cassanova, dd.String(), m.Author.ID))
			return true
		}
		return true
	}
	return false
}

func isCassExpired(id string) bool {
	if val, ok := cassMap[id]; ok {
		if val.Expires && time.Now().After(val.Duration) {
			return true
		}
	}
	return false
}

func (c *Command) jereme(s *discordgo.Session, m *discordgo.MessageCreate) {
	if containsVal(m.Message.Member.Roles, boss) > -1 || containsVal(m.Message.Member.Roles, "324575381581463553") > -1 {
		if len(c.Values) < 2 {
			return
		}
		s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, boss)
		user := strings.Replace(strings.Replace(c.Values[1], "<@!", "", -1), ">", "", -1)
		if len(c.Values) > 2 && c.Values[2] == "cancel" {
			delete(cassMap, user)
			s.GuildMemberRoleRemove(m.GuildID, user, cassanova)
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@!%s> is no longer <@&%s>ed.", user, cassanova))
			return
		}
		s.GuildMemberRoleAdd(m.GuildID, user, cassanova)
		cassMap[user] = Penalty{
			ID:       user,
			Expires:  true,
			Duration: time.Now().Add(time.Minute * 30),
			Debounce: time.Now(),
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@!%s> has been <@&%s>ed for the next *30 minutes*", user, cassanova))

		return
	}
	cassMap[m.Author.ID] = Penalty{
		ID:       m.Author.ID,
		Expires:  true,
		Duration: time.Now().Add(time.Minute * 30),
		Debounce: time.Now(),
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@!%s> has been <@&%s>ed for the next *30 minutes*", m.Author.ID, cassanova))
	s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, cassanova)
	return
}
