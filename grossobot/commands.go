package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var help string = "```yaml\n" +
	`Commands:
* !grosso help/!grossohelp - show this 
* !grosso add/!grossoadd imageurl trigger anothertrigger - add the emoji to grossobot to be triggered by a space separated list of words.
* !grosso list/!grossolist - print a full list of triggers for grossobot.
* !grosso describe/!grossodesc trigger description goes here - add or rewrite the description for the field with the given trigger.
Emotes:
Trigger these by bolding the trigger word for example "**dogtown**" below are some examples but there are many more we didn't list.
* crob/muska/mullen/natas/etc - Clips of said pro skater. 
* hammers/yeet/send it/stoked/lucky - Clip of someone ripping. 
* oof/sucky/slam/wasted/ouch - Skin dying or a bad slam. 
* front salad/hollywood high - Owen Wilson Yeah Right copypasta! 
* Do a kickflip - DO A KICKFLIP!!! 
* footy/footage/film/clips - A kind reminder to post some footage.
* kook/penny/nickel/longboard/revive/braille - A message to send to someone after they don't post footy. 
` + "uh - uh...```"

var actionMap = map[string]string{
	"!grosso help":     "help",
	"!grossohelp":      "help",
	"!grossoadd":       "add",
	"!grosso add":      "add",
	"!grossodesc":      "desc",
	"!grosso describe": "desc",
	"!grossolist":      "list",
	"!grosso list":     "list",
}

//Command struct for grosso commands
type Command struct {
	Trigger []string
	Action  string
	Values  []string
}

//Parse maps a string to a command
func (c *Command) Parse(s string) (out bool) {
	for _, v := range c.Trigger {
		if strings.HasPrefix(s, v) {
			out = true
			fmt.Println(s, v)
			c.Action = actionMap[v]
			s = strings.Replace(s, v, "", -1)
			sa := strings.Split(s, " ")
			if len(sa) >= 1 {
				c.Values = sa
				if c.Action == sa[0] {
					c.Values = sa[1:]
				}
			}
		}
	}
	return
}

//Process processes a command object
func (c *Command) Process(s *discordgo.Session, m *discordgo.MessageCreate) {
	cont := m.Content
	short := c.Parse(cont)
	if !short {
		return
	}
	switch c.Action {
	case "add":
		c.add(s, m)
	case "list":
		c.list(s, m)
	case "help":
		c.help(s, m)
	case "desc":
		c.desc(s, m)
	}
}

//BBresponse json response from imagebb
type BBresponse struct {
	Status  int    `json:"status"`
	Success bool   `json:"success"`
	Data    BBdata `json:"data"`
}

//BBdata json response subfield
type BBdata struct {
	URL string `json:"url"`
}

//CheckRole Checks to see if the user has the min required role
func CheckRole(m *discordgo.MessageCreate) bool {
	role := "327845878578675713"
	var badRoles = []string{
		"359852475181694976",
		"636246344855453696",
		"716536970309927033",
		"416424719487860736",
		"697875957964603403",
	}
	if containsVal(m.Message.Member.Roles, role) > -1 {
		for _, v := range badRoles {
			if containsVal(m.Message.Member.Roles, v) > -1 {
				return false
			}
		}
		return true
	}
	return false
}

func demoNew(t []string, url string, s *discordgo.Session, m *discordgo.MessageCreate) {
	ran := rand.Intn(len(t))
	msg := fmt.Sprintf("<@%s> is **%s**.", m.Author.ID, t[ran])
	s.ChannelMessageSend(m.ChannelID, msg)
	buf := getFile(url)
	s.ChannelFileSend(m.ChannelID, url, buf)
}

func (c *Command) add(s *discordgo.Session, m *discordgo.MessageCreate) {
	if CheckRole(m) != true {
		s.ChannelMessageSend(m.ChannelID, "Nice try. Keep practicing the art of shit talking.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "processing...")
	if len(c.Values) < 2 {
		s.ChannelMessageSend(m.ChannelID, "aww that one was wack. try again")
		return
	}
	p := c.Values[1]
	t := c.Values[2:]
	fmt.Println("p ", p)
	fmt.Println("t ", t)
	url := "https://api.imgbb.com/1/upload?key=" + os.Getenv("IMGBBKEY")
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("image", p)
	err := writer.Close()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "hehehe not that one.")
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "you gotta tweak it.")
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Upload bailed. try again.")
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	bb := &BBresponse{}
	err = json.Unmarshal(body, bb)
	if err != nil || bb.Success != true {
		s.ChannelMessageSend(m.ChannelID, "Upload Failed. Forget it. Go Skate!")
		return
	}
	fmt.Println(err, string(body))
	fmt.Println(bb)
	new := Case{
		Triggers: t,
		Type:     CONTAINS,
		Images:   []string{p},
	}
	cases = caseMatch(new, cases)
	go demoNew(t, p, s, m)
}

func caseMatch(n Case, c []Case) []Case {
	for i, v := range c {
		for _, t := range v.Triggers {
			index := containsVal(n.Triggers, t)
			if index > -1 {
				for _, vv := range n.Triggers {
					if containsVal(c[i].Triggers, vv) < 0 {
						c[i].Triggers = append(c[i].Triggers, vv)
					}
				}
				c[i].Images = append(c[i].Images, n.Images...)
				return c
			}
		}
	}
	c = append(c, n)
	return c
}

func caseDesc(s, d string, c []Case) ([]Case, Case) {
	var out Case
	for i, v := range c {
		if containsVal(v.Triggers, s) > -1 {
			c[i].Description = d
			out = c[i]
		}
	}
	return c, out
}

func containsVal(s []string, e string) int {
	for i, a := range s {
		if a == e {
			return i
		}
	}
	return -1
}

func (c *Command) list(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("list")
	items := []string{}
	for _, v := range cases {
		line := "**" + strings.Join(v.Triggers, "/") + "** - " + v.Description
		items = append(items, line)
	}
	fmt.Println(items)
	for _, v := range items {
		s.ChannelMessageSend(m.ChannelID, v)
	}
}

func (c *Command) help(s *discordgo.Session, m *discordgo.MessageCreate) {
	s.ChannelMessageSend(m.ChannelID, help)
}

func (c *Command) desc(s *discordgo.Session, m *discordgo.MessageCreate) {
	if CheckRole(m) != true {
		s.ChannelMessageSend(m.ChannelID, "Nice try. Keep practicing the art of shit talking.")
		return
	}
	s.ChannelMessageSend(m.ChannelID, "processing...")
	bot := getFile("https://i.ibb.co/4RBtbVC/grossobot.gif")
	s.ChannelFileSend(m.ChannelID, "grossobot.gif", bot)
	if len(c.Values) < 2 {
		s.ChannelMessageSend(m.ChannelID, "aww that one was wack. try again")
		return
	}
	t := c.Values[1]
	d := strings.Join(c.Values[2:], " ")
	var nc Case
	cases, nc = caseDesc(t, d, cases)
	out := fmt.Sprintf("**%s** - %s", strings.Join(nc.Triggers, "/"), nc.Description)
	s.ChannelMessageSend(m.ChannelID, out)
	f := nc.Images[rand.Intn(len(nc.Images))]
	buf := getFile(f)
	s.ChannelFileSend(m.ChannelID, f, buf)
}

var helpc = Command{
	Trigger: []string{"!grosso help", "!grossohelp"},
	Action:  "help",
}

var listc = Command{
	Trigger: []string{"!grosso list", "!grossolist"},
	Action:  "list",
}

var addc = Command{
	Trigger: []string{"!grosso add", "!grossoadd"},
	Action:  "help",
}

var descc = Command{
	Trigger: []string{"!grosso describe", "!grossodesc"},
	Action:  "desc",
}

var commands = []Command{
	helpc,
	listc,
	addc,
	descc,
}
