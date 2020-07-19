package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

//• Win a trivia round against me <@!698330051686432900>. To play say **kooktrivia**.
//• Get three <@&327840671811502081> members to !unkook you.
type unkook struct {
	id    string
	votes []string
}

var unkooklist = []unkook{}

var kooky string = `<@!%s> you have been **kooked**!
This means you can no longer post in channels other than #offtopic .
To unkook yourself you may either:
• Get three <@&327840671811502081> members to !unkook you.
• Win a trivia round against me <@!698330051686432900>. To play say **kooktrivia**.
• Or watch all the 411vms here https://www.youtube.com/watch?v=PJCoq-vvDi4&list=PLaZk_-qbK_1iKAtOwJeALna3IZVJsN5zN .
`

var unkooky string = `<@!%s> you have been **unkooked**! Put some pants on and go skate!`

func (c *Command) kook(s *discordgo.Session, m *discordgo.MessageCreate) {
	kooked := m.Author.ID
	kooks := []string{}
	if len(c.Values) < 2 ||
		containsVal(m.Message.Member.Roles, "324575381581463553") < 0 {
		c.Values = []string{}
		kooks = []string{m.Author.ID}
		if containsVal(m.Message.Member.Roles, "324575381581463553") > -1 {
			kooks = []string{}
		}
	}
	for _, v := range c.Values[1:] {
		kooked = strings.Replace(strings.Replace(v, "<@!", "", -1), ">", "", -1)
		member, err := s.GuildMember(m.GuildID, kooked)
		if err != nil {
			fmt.Println("mem", err)
			return
		}
		if containsVal(member.Roles, "324575381581463553") > -1 {
			kooks = []string{m.Author.ID}
			if containsVal(m.Message.Member.Roles, "324575381581463553") > -1 {
				kooks = []string{}
			}
			break
		}
		kooks = append(kooks, kooked)
	}
	fmt.Println(kooks)
	for _, v := range kooks {
		err := s.GuildMemberRoleAdd(m.GuildID, kooked, "359852475181694976")
		unkooklist = append(unkooklist, unkook{id: kooked, votes: []string{}})
		if err != nil {
			fmt.Println("addrole", err)
			return
		}
		go dekook(kooked, s, m, time.Hour*67)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(kooky, v))
	}
}

func (c *Command) unkook(s *discordgo.Session, m *discordgo.MessageCreate) {
	for _, v := range c.Values[1:] {
		uk := strings.Replace(strings.Replace(v, "<@!", "", -1), ">", "", -1)
		kooks := []string{}
		for _, k := range unkooklist {
			kooks = append(kooks, k.id)
		}
		i := containsVal(kooks, uk)
		if i < 0 {
			continue
		}
		if containsVal(unkooklist[i].votes, m.Author.ID) < 0 && m.Author.ID != unkooklist[i].id {
			unkooklist[i].votes = append(unkooklist[i].votes, m.Author.ID)
		}
		if len(unkooklist[i].votes) > 2 {
			dekook(unkooklist[i].id, s, m, 0)
		}
	}
}

func dekook(st string, s *discordgo.Session, m *discordgo.MessageCreate, d time.Duration) {
	time.Sleep(d)
	err := s.GuildMemberRoleRemove(m.GuildID, st, "359852475181694976")
	if err != nil {
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(unkooky, st))
}
