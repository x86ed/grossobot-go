package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

//Team for trivia
type Team struct {
	Captain  string
	Members  []string
	Scores   []Score
	Name     string
	ID       string
	Inactive bool
}

var teams = map[string]Team{}

//GetAllTimeScore gets the full pointes the team has scored
func (tm *Team) GetAllTimeScore() (s int) {
	for _, v := range tm.Scores {
		s += v.Points
	}
	return
}

//GetMaxScore gets the full pointes the team has scored
func (tm *Team) GetMaxScore() (s int) {
	for _, v := range tm.Scores {
		if v.Points > s {
			s = v.Points
		}
	}
	return
}

var actions = []string{"disband", "leave", "join", "list"}
var haveATeam = "<@!%s> is already on <@&%s>. You can leave or disband and join another team."

func (c *Command) team(s *discordgo.Session, m *discordgo.MessageCreate) {
	teamMsg := "<@!%s> has created the team <@&%s>."
	cap := m.Author.ID
	fmt.Println(c.Values, len(c.Values))
	if len(c.Values) > 1 {
		if containsVal(actions, c.Values[1]) > -1 {
			handleTeamAction(c.Values[1], s, m, c)
			return
		}
		tms := isOnActiveTeam(teams, m.Author.ID)
		if len(tms) > 0 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(haveATeam, m.Author.ID, tms[0]))
			return
		}
		role, err := s.GuildRoleCreate(m.GuildID)
		if err != nil {
			return
		}
		rand.Seed(time.Now().UnixNano())
		color := rand.Intn(0xffffff + 1)
		spacey := strings.Replace(c.Values[1], "_", " ", -1)
		_, err = s.GuildRoleEdit(m.GuildID, role.ID, spacey, color, false, role.Permissions, true)
		if err != nil {
			return
		}
		err = s.GuildMemberRoleAdd(m.GuildID, cap, role.ID)
		if err != nil {
			return
		}
		nt := Team{
			Captain: m.Author.ID,
			Members: []string{},
			Scores:  []Score{},
			Name:    c.Values[1],
			ID:      role.ID,
		}
		teams[role.ID] = nt
		archiveJSON(os.Getenv("TRIVIATEAMS"), &teams)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(teamMsg, m.Author.ID, role.ID))
	}
}

func isOnActiveTeam(ts map[string]Team, id string) (tms []string) {
	for _, v := range ts {
		if v.Captain == id && !v.Inactive {
			tms = append(tms, v.ID)
		}
		for _, w := range v.Members {
			if w == id && !v.Inactive {
				tms = append(tms, v.ID)
			}
		}
	}
	return
}

func handleTeamAction(a string, s *discordgo.Session, m *discordgo.MessageCreate, c *Command) {
	disMsg := "<@!%s> has disbanded the team %s."
	leftMsg := "<@!%s> has left <@&%s>."
	joinMsg := "<@!%s> has joined <@&%s>."
	errLeftMsg := "<@!%s> couldn't leave <@&%s>."
	errJoinMsg := "<@!%s> couldn't join <@&%s>."
	printTeam := "----------\n<@&%s>\n----------\n\nCaptain:\n<@!%s>\nMembers:\n%s\n"
	switch a {
	case "disband":
		for _, v := range teams {
			if v.Captain == m.Author.ID {
				err := s.GuildRoleDelete(m.GuildID, v.ID)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, "The team couldn't be disbanded")
					return
				}
				t := teams[v.ID]
				t.Inactive = true
				teams[v.ID] = t
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(disMsg, m.Author.ID, v.Name))
			}
		}
	case "leave":
		for _, v := range teams {
			for i, w := range v.Members {
				if w == m.Author.ID && !v.Inactive {
					v.Members = remove(v.Members, i)
					teams[v.ID] = v
					err := s.GuildMemberRoleRemove(m.GuildID, m.Author.ID, v.ID)
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(errLeftMsg, m.Author.ID, v.ID))
						return
					}
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(leftMsg, m.Author.ID, v.ID))
				}
			}
		}
	case "join":
		if len(c.Values) < 2 {
			s.ChannelMessageSend(m.ChannelID, "You gotta pick a team to join.")
			return
		}
		tms := isOnActiveTeam(teams, m.Author.ID)
		if len(tms) > 0 {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(haveATeam, m.Author.ID, tms[0]))
			return
		}
		teamID := strings.Replace(strings.Replace(c.Values[2], "<@&", "", -1), ">", "", -1)
		err := s.GuildMemberRoleAdd(m.GuildID, m.Author.ID, teamID)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(errJoinMsg, m.Author.ID, teamID))
			return
		}
		t := teams[teamID]
		if containsVal(t.Members, m.Author.ID) < 0 {
			t.Members = append(t.Members, m.Author.ID)
			teams[teamID] = t
			archiveJSON(os.Getenv("TRIVIATEAMS"), &teams)
		}
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(joinMsg, m.Author.ID, teamID))
	case "list":
		for _, v := range teams {
			if !v.Inactive {
				mm := getMems(v.Members)
				fmt.Println(v)
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(printTeam, v.ID, v.Captain, mm))
			}
		}
	}
}

func getMems(m []string) (out string) {
	for i, v := range m {
		out += fmt.Sprintf("%d. <@!%s>\n", i+1, v)
	}
	return
}
