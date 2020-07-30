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
	"time"

	"github.com/bwmarrin/discordgo"
)

var help string = "```yaml\n" +
	`Commands:
* !grosso help/!grossohelp - show this 
* !grosso add/!grossoadd imageurl trigger anothertrigger - add the emoji to grossobot to be triggered by a space separated list of words.
* !grosso list/!grossolist - print a full list of triggers for grossobot.
* !grosso describe/!grossodesc trigger description goes here - add or rewrite the description for the field with the given trigger.
* !grosso kook/!kook @username - kook a user for being a kook.
* !grosso unkook/!unkook @username - unkook a kooked user.
Trivia:
* !grosso team/!grossoteam teamname - create a trivia team.
* !grosso join/!grossojoin teamname authkey - join a trivia team.
*
* !grosso submit/!grossosub - suggest a new question.
* !grosso board/!grossoboard (username) - get leaderboard for user or all users if blank
* 
* !grosso trivia/!trivia (options) - start a new trivia game (admins only)
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
	"!grosso team":     "team",
	"!grossoteam":      "team",
	"!grosso submit":   "sub",
	"!grossosub":       "sub",
	"!grosso board":    "board",
	"!grossoboard":     "board",
	"!grosso kook":     "kook",
	"!kook":            "kook",
	"!grosso unkook":   "unkook",
	"!unkook":          "unkook",
	"!grosso trivia":   "trivia",
	"!trivia":          "trivia",
	"+answer":          "answer",

	"+question":   "question",
	"+correct":    "correct",
	"+incorrect":  "incorrect",
	"+image":      "image",
	"+difficulty": "difficulty",
	"+cancel":     "cancel",
	"+save":       "save",
	"+help":       "thelp",

	"+proctor": "proc",
	"+approve": "app",
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

func checkCass(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if m.Message.Member != nil && containsVal(m.Message.Member.Roles, cassanova) > -1 {
		fmt.Println(m.ChannelID, m.Message.ID)
		rand.Seed(time.Now().UnixNano())
		cc := rand.Intn(len(jeremeVids))
		s.ChannelMessageDelete(m.ChannelID, m.Message.ID)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf(jeremeVids[cc], m.Author.ID))
		return true
	}
	return false
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
	case "kook":
		c.kook(s, m)
	case "unkook":
		c.unkook(s, m)
	case "team":
		c.team(s, m)
	case "sub":
		c.submit(s, m)
	case "trivia":
		c.trivia(s, m)
	case "question":
		c.sub(s, m, "question")
	case "correct":
		c.sub(s, m, "correct")
	case "answer":
		c.sub(s, m, "answer")
	case "incorrect":
		c.sub(s, m, "incorrect")
	case "image":
		c.sub(s, m, "image")
	case "difficulty":
		c.sub(s, m, "difficulty")
	case "cancel":
		c.sub(s, m, "cancel")
	case "save":
		c.sub(s, m, "save")
	case "thelp":
		c.sub(s, m, "help")
	case "app":
		c.sub(s, m, "app")
	case "proc":
		c.sub(s, m, "proc")
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
	if m.Message.Member == nil {
		return false
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
	dm, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		return
	}
	s.ChannelMessageSend(dm.ID, "processing...")
	bot := getFile("https://i.ibb.co/4RBtbVC/grossobot.gif")
	s.ChannelFileSend(dm.ID, "grossobot.gif", bot)
	if len(c.Values) < 2 {
		s.ChannelMessageSend(dm.ID, "aww that one was wack. try again")
		return
	}
	p := c.Values[1]
	t := c.Values[2:]
	for i, v := range t {
		t[i] = strings.Replace(v, "_", " ", -1)
	}
	fmt.Println("p ", p)
	fmt.Println("t ", t)
	url := "https://api.imgbb.com/1/upload?key=" + os.Getenv("IMGBBKEY")
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("image", p)
	err = writer.Close()
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
		Images:   []string{bb.Data.URL},
	}
	cases = caseMatch(new, cases)
	go demoNew(t, p, s, m)
	go archiveMemes(os.Getenv("MEMEARCH"))
}

func archiveMemes(fn string) {
	f, err := os.Create(fn)
	if err != nil {
		return
	}

	defer f.Close()

	arch, err := json.Marshal(&cases)
	if err != nil {
		return
	}
	c, err := f.Write(arch)
	if err != nil {
		return
	}
	fmt.Println("bytes: ", c)
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

func (c *Command) list(s *discordgo.Session, m *discordgo.MessageCreate) {
	dm, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		return
	}
	items := []string{}
	for _, v := range cases {
		line := "**" + strings.Join(v.Triggers, "/") + "** - " + v.Description
		items = append(items, line)
	}
	fmt.Println(items)
	for _, v := range items {
		s.ChannelMessageSend(dm.ID, v)
	}
}

func (c *Command) help(s *discordgo.Session, m *discordgo.MessageCreate) {
	dm, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		return
	}
	s.ChannelMessageSend(dm.ID, help)
}

func (c *Command) desc(s *discordgo.Session, m *discordgo.MessageCreate) {
	if CheckRole(m) != true {
		s.ChannelMessageSend(m.ChannelID, "Nice try. Keep practicing the art of shit talking.")
		return
	}
	dm, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		return
	}
	s.ChannelMessageSend(dm.ID, "processing...")
	bot := getFile("https://i.ibb.co/4RBtbVC/grossobot.gif")
	s.ChannelFileSend(dm.ID, "grossobot.gif", bot)
	if len(c.Values) < 2 {
		s.ChannelMessageSend(dm.ID, "aww that one was wack. try again")
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

var kookc = Command{
	Trigger: []string{"!grosso kook", "!kook"},
	Action:  "kook",
}

var unkookc = Command{
	Trigger: []string{"!grosso unkook", "!unkook"},
	Action:  "unkook",
}

var triviac = Command{
	Trigger: []string{"!grosso trivia", "!trivia"},
	Action:  "trivia",
}

var teamc = Command{
	Trigger: []string{"!grosso team", "!grossoteam"},
	Action:  "team",
}

var subc = Command{
	Trigger: []string{"!grosso submit", "!grossosub"},
	Action:  "sub",
}

var boardc = Command{
	Trigger: []string{"!grosso board", "!grossoboard"},
	Action:  "board",
}

var procc = Command{
	Trigger: []string{"+proctor"},
	Action:  "proc",
}

var appc = Command{
	Trigger: []string{"+approve"},
	Action:  "app",
}

var questionc = Command{
	Trigger: []string{"+question"},
	Action:  "question",
}

var correctc = Command{
	Trigger: []string{"+correct"},
	Action:  "correct",
}

var incorrectc = Command{
	Trigger: []string{"+incorrect"},
	Action:  "incorrect",
}

var answerc = Command{
	Trigger: []string{"+answer"},
	Action:  "answer",
}

var imagec = Command{
	Trigger: []string{"+image"},
	Action:  "image",
}

var difficultyc = Command{
	Trigger: []string{"+difficulty"},
	Action:  "difficulty",
}

var cancelc = Command{
	Trigger: []string{"+cancel"},
	Action:  "cancel",
}

var savec = Command{
	Trigger: []string{"+save"},
	Action:  "save",
}

var thelpc = Command{
	Trigger: []string{"+help"},
	Action:  "thelp",
}

var commands = []Command{
	helpc,
	listc,
	addc,
	descc,
	kookc,
	unkookc,
	triviac,
	teamc,
	subc,
	boardc,
	procc,
	appc,
	questionc,
	correctc,
	incorrectc,
	answerc,
	imagec,
	difficultyc,
	cancelc,
	savec,
	thelpc,
}
