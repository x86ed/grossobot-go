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
	"github.com/google/uuid"
)

//Score game score
type Score struct {
	Points int
	Game   string
	Team   string
}

//Round (active trivia round)
type Round struct {
	Start    time.Time
	Question *Question
	Count    int
	ID       string
}

//CurRound current round
var CurRound = Round{
	Start: time.Now(),
	Count: 3,
	ID:    "5",
}

//Game game object
type Game struct {
	Rounds    int
	Scores    []Score
	Current   int
	ID        string
	Interval  time.Duration
	Active    bool
	Questions []string
}

//Question game object
type Question struct {
	ID               string
	Text             string
	Answers          []string
	Correct          string
	Img              string
	Points           int
	Active           bool
	AlternateAnswers []string
}

//Answer game object
type Answer struct {
	ID         string
	Value      string
	QuestionID string
	TimeStamp  time.Duration
	UserID     string
}

func getQuestion(id string) Question {
	for _, v := range Questions {
		if v.ID == id {
			return v
		}
	}
	return Question{}
}

func (a *Answer) print() *discordgo.MessageEmbed {
	color := 0x000000
	difficulty := "unknown"
	q := getQuestion(a.QuestionID)
	gm := getActiveGame()
	switch q.Points {
	case 1000:
		color = 0x00FF00
		difficulty = "easy"
	case 2000:
		color = 0xFFFF00
		difficulty = "medium"
	case 3000:
		color = 0xFF0000
		difficulty = "hard"
	}
	f := []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "ID",
			Value: a.ID,
		},
		&discordgo.MessageEmbedField{
			Name:  "difficulty",
			Value: difficulty,
		},
		&discordgo.MessageEmbedField{
			Name:  "Correct",
			Value: q.Correct,
		},
		&discordgo.MessageEmbedField{
			Name:  "Actual",
			Value: a.Value,
		},
		&discordgo.MessageEmbedField{
			Name:  "Points",
			Value: string(int(gm.Interval-a.TimeStamp/gm.Interval) * q.Points),
		},
	}
	for _, v := range q.Answers {
		f = append(f, &discordgo.MessageEmbedField{
			Name:  "Incorrect",
			Value: v,
		})
	}
	e := discordgo.MessageEmbed{
		Color:       color,
		Title:       q.ID,
		Description: q.Text,
		URL:         "https://discord.gg/zq3fyV",
		// Author: &discordgo.MessageEmbedAuthor{
		// 	URL:     "https://discord.gg/zq3fyV",
		// 	Name:    "GrossoBot",
		// 	IconURL: "https://i.ibb.co/4RBtbVC/grossobot.gif",
		// },
		Fields: f,
	}
	if len(q.Img) > 0 {
		e.Image = &discordgo.MessageEmbedImage{
			URL: q.Img,
		}
	}
	return &e
}

//Games trivia games played so far
var Games = []Game{}

//Questions for trivia
var Questions = []Question{}

//NewQuestions buffer for new questions
var NewQuestions = map[string]Question{}
var howToMsg = `
So I heard you had a trivia question you'd like to add...
Heres how to do it:

Required Fields:
**+question**  Your question text goes here.
**+correct** The correct answer shown to the triviamaster.

Optional Fields:
**+incorrect** A correct answer if the question is multiple choice. (repeat as needed)
**+image** an image to display with the question. post a URL ending with .gif/.jpg/.png
**+difficulty** easy/medium/hard. (defaults to medium)
**+cancel** aborts creating the question.
**+save** saves the question.
**+help** print this menu.
`
var adminHowTo = `
**+proctor** judge a given answer (admins only)
**+approve** approve or deny a given question (admins only)
`

func (c *Command) submit(s *discordgo.Session, m *discordgo.MessageCreate) {
	if CheckRole(m) != true {
		s.ChannelMessageSend(m.ChannelID, "Nice try. Keep practicing the art of shit talking.")
		return
	}
	dm, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		return
	}
	if _, ok := NewQuestions[m.Author.ID]; !ok {
		id, err := uuid.NewRandom()
		if err != nil {
			return
		}
		NewQuestions[m.Author.ID] = Question{
			ID:      id.String(),
			Answers: []string{},
			Points:  2000,
		}
		s.ChannelMessageSend(m.ChannelID, "Sliding into your DMs...")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Finish the question you're working on with `+save` or `+cancel`")
	}
	if containsVal(quizJudge, m.Author.ID) > -1 {
		s.ChannelMessageSend(dm.ID, howToMsg+adminHowTo)
		return
	}
	s.ChannelMessageSend(dm.ID, howToMsg)
}

func (c *Command) sub(s *discordgo.Session, m *discordgo.MessageCreate, sub string) {
	dm, err := s.UserChannelCreate(m.Author.ID)
	if err != nil {
		return
	}
	if dm.ID != m.ChannelID {
		s.ChannelMessageSend(m.ChannelID, "It's been real, but keep this in the DMs.")
		return
	}
	if _, ok := NewQuestions[m.Author.ID]; !ok && sub != "app" && sub != "proc" {
		id, err := uuid.NewRandom()
		if err != nil {
			return
		}
		NewQuestions[m.Author.ID] = Question{
			ID:      id.String(),
			Answers: []string{},
			Points:  2000,
		}
		s.ChannelMessageSend(m.ChannelID, "Creating a new question...")
	}
	c.param(s, m, sub)
}

func (c *Command) param(s *discordgo.Session, m *discordgo.MessageCreate, sub string) {
	switch sub {
	case "proc":
		if containsVal(quizJudge, m.Author.ID) > -1 {

		}
	case "app":
		if containsVal(quizJudge, m.Author.ID) > -1 {
			if len(c.Values) < 2 {
				unApproved := []int{}
				for i, v := range Questions {
					if v.Active != true {
						unApproved = append(unApproved, i)
					}
				}
				for _, v := range unApproved {
					curr := Questions[v]
					s.ChannelMessageSendEmbed(m.ChannelID, curr.print())
				}
			} else {
				approvals := c.Values[1:]
				for _, v := range approvals {
					for i, w := range Questions {
						if w.ID == v {
							w.Active = true
							Questions[i] = w
							s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s approved\n", w.ID))
							if m.ChannelID != judgeChannel {
								s.ChannelMessageSendEmbed(judgeChannel, w.print())
								s.ChannelMessageSend(judgeChannel, fmt.Sprintf("%s approved\n", w.ID))
							}
							archiveJSON(os.Getenv("TRIVIAQUESTIONS"), &Questions)
						}
					}
				}
			}
		}
	case "help":
		if containsVal(quizJudge, m.Author.ID) > -1 {
			s.ChannelMessageSend(m.ChannelID, howToMsg+adminHowTo)
			return
		}
		s.ChannelMessageSend(m.ChannelID, howToMsg)
		return
	case "question":
		if len(c.Values) > 1 {
			v := NewQuestions[m.Author.ID]
			v.Text = strings.Join(c.Values[1:], " ")
			NewQuestions[m.Author.ID] = v
		}
	case "correct":
		if len(c.Values) > 1 {
			v := NewQuestions[m.Author.ID]
			v.Correct = strings.Join(c.Values[1:], " ")
			NewQuestions[m.Author.ID] = v
		}
	case "answer":
		if len(c.Values) > 1 {
			full := strings.Join(c.Values[1:], " ")
			id, err := uuid.NewRandom()
			if err != nil {
				return
			}
			An := Answer{
				ID:    id.String(),
				Value: full,
			}
			fmt.Println(An)
		}
	case "incorrect":
		if len(c.Values) > 1 {
			v := NewQuestions[m.Author.ID]
			v.Answers = append(v.Answers, strings.Join(c.Values[1:], " "))
			NewQuestions[m.Author.ID] = v
		}
	case "image":
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
		fmt.Println("p ", p)
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
		v := NewQuestions[m.Author.ID]
		v.Img = bb.Data.URL
		NewQuestions[m.Author.ID] = v
	case "difficulty":
		if len(c.Values) > 1 {
			v := NewQuestions[m.Author.ID]
			switch c.Values[1] {
			case "easy":
				v.Points = 1000
			case "medium":
				v.Points = 2000
			case "hard":
				v.Points = 3000
			}
			NewQuestions[m.Author.ID] = v
		}
	case "save":
		v := NewQuestions[m.Author.ID]
		if len(v.Correct) > 0 && len(v.Text) > 0 {
			Questions = append(Questions, v)
			s.ChannelMessageSendEmbed(judgeChannel, v.print())
			delete(NewQuestions, m.Author.ID)
			s.ChannelMessageSend(m.ChannelID, "Your Question has been saved")
			archiveJSON(os.Getenv("TRIVIAQUESTIONS"), &Questions)
		}
	case "cancel":
		delete(NewQuestions, m.Author.ID)
		s.ChannelMessageSend(m.ChannelID, "Deleted the question in progress.")
		return
	}
	if _, ok := NewQuestions[m.Author.ID]; ok {
		v := NewQuestions[m.Author.ID]
		s.ChannelMessageSendEmbed(m.ChannelID, v.print())
	}
}

func (q *Question) print() *discordgo.MessageEmbed {
	color := 0x000000
	difficulty := "unknown"
	switch q.Points {
	case 1000:
		color = 0x00FF00
		difficulty = "easy"
	case 2000:
		color = 0xFFFF00
		difficulty = "medium"
	case 3000:
		color = 0xFF0000
		difficulty = "hard"
	}
	f := []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "ID",
			Value: q.ID,
		},
		&discordgo.MessageEmbedField{
			Name:  "difficulty",
			Value: difficulty,
		},
		&discordgo.MessageEmbedField{
			Name:  "Correct",
			Value: q.Correct,
		},
	}
	for _, v := range q.Answers {
		f = append(f, &discordgo.MessageEmbedField{
			Name:  "Incorrect",
			Value: v,
		})
	}
	e := discordgo.MessageEmbed{
		Color:       color,
		Title:       q.ID,
		Description: q.Text,
		URL:         "https://discord.gg/zq3fyV",
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://discord.gg/zq3fyV",
			Name:    "GrossoBot",
			IconURL: "https://i.ibb.co/4RBtbVC/grossobot.gif",
		},
		Fields: f,
	}
	if len(q.Img) > 0 {
		e.Image = &discordgo.MessageEmbedImage{
			URL: q.Img,
		}
	}
	return &e
}

func getActiveGame() Game {
	for _, v := range Games {
		if v.Active {
			return v
		}
	}
	return Game{}
}

func shuffleAnswers(f []*discordgo.MessageEmbedField, c string, q []string) []*discordgo.MessageEmbedField {
	all := []string{}
	labels := []string{"A.", "B.", "C.", "D.", "E."}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(q), func(i, j int) { q[i], q[j] = q[j], q[i] })
	if len(q) > 5 {
		q = q[:3]
	}
	all = append(q, c)
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(all), func(i, j int) { all[i], all[j] = all[j], all[i] })
	for i, v := range all {
		f = append(f, &discordgo.MessageEmbedField{
			Name:  labels[i],
			Value: v,
		})
	}
	return f
}

func (q *Question) ask() *discordgo.MessageEmbed {
	activeGame := getActiveGame()
	qs := fmt.Sprintf("Round %d of %d", activeGame.Current, activeGame.Rounds)
	if len(activeGame.ID) < 1 {
		return &discordgo.MessageEmbed{}
	}
	color := 0x000000
	difficulty := "unknown"
	switch q.Points {
	case 1000:
		color = 0x00FF00
		difficulty = "easy"
	case 2000:
		color = 0xFFFF00
		difficulty = "medium"
	case 3000:
		color = 0xFF0000
		difficulty = "hard"
	}
	f := []*discordgo.MessageEmbedField{
		&discordgo.MessageEmbedField{
			Name:  "difficulty",
			Value: difficulty,
		},
	}
	for _, v := range q.Answers {
		f = append(f, &discordgo.MessageEmbedField{
			Name:  "Incorrect",
			Value: v,
		})
	}
	f = shuffleAnswers(f, q.Correct, q.Answers)
	e := discordgo.MessageEmbed{
		Color:       color,
		Title:       qs,
		Description: q.Text,
		URL:         "https://discord.gg/zq3fyV",
		Author: &discordgo.MessageEmbedAuthor{
			URL:     "https://discord.gg/zq3fyV",
			Name:    "GrossoBot",
			IconURL: "https://i.ibb.co/4RBtbVC/grossobot.gif",
		},
		Fields: f,
	}
	if len(q.Img) > 0 {
		e.Image = &discordgo.MessageEmbedImage{
			URL: q.Img,
		}
	}
	return &e
}

func (c *Command) trivia(s *discordgo.Session, m *discordgo.MessageCreate) {
	if containsVal(quizJudge, m.Author.ID) < 0 {
		s.ChannelMessageSend(m.ChannelID, "Nice try buddy.")
		return
	}
	if len(c.Values) > 1 && c.Values[1] == "cancel" {
		for i, v := range Games {
			if v.Active {
				Games[i] = Games[len(Games)-1]
				Games = Games[:len(Games)-1]
			}
		}
		s.ChannelMessageSend(m.ChannelID, "Game Cancelled! No points will be awarded.")
		return
	}
	if activeGame(Games) {
		s.ChannelMessageSend(m.ChannelID, "The current game hasn't ended. Play through or cancel it with `!trivia cancel`")
		return
	}
	if len(c.Values) > 1 && c.Values[1] == "start" {
		//start trivia round
		return
	}

}

func activeGame(g []Game) (is bool) {
	for _, v := range g {
		if v.Active {
			is = true
			return
		}
	}
	return
}
