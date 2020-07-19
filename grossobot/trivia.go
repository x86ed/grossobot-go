package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

var quizJudge = []string{
	"329451587422519297",
	"202218126987755523", //lsd
	"313719596886523904",
}

var judgeChannel = "734466604393431170"

//Score game score
type Score struct {
	Points int
	Game   string
	Team   string
}

//Game game object
type Game struct {
	Rounds   int
	Scores   []Score
	Current  int
	ID       string
	Interval time.Duration
	Active   bool
}

//Question game object
type Question struct {
	ID      string
	Text    string
	Answers []string
	Correct string
	Img     string
	Points  int
	Active  bool
}

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
**+answer** A correct answer if the question is multiple choice. (repeat as needed)
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
					s.ChannelMessageSend(m.ChannelID, curr.print())
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
								s.ChannelMessageSend(judgeChannel, w.print())
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
			s.ChannelMessageSend(judgeChannel, v.print())
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
		s.ChannelMessageSend(m.ChannelID, v.print())
	}
}

func (q *Question) print() (out string) {
	out = "----------\n"
	out += fmt.Sprintf("ID: %s\n", q.ID)
	out += fmt.Sprintf("Q: %s\n", q.Text)
	out += (q.Img + "\n")
	out += fmt.Sprintf("A :%s\n", q.Correct)
	if len(q.Answers) > 0 {
		out += fmt.Sprintf("Incorrect Answers:\n%s\n", strings.Join(q.Answers, "\n"))
	}
	switch q.Points {
	case 1000:
		out += "Difficulty: easy\n"
	case 2000:
		out += "Difficulty: medium\n"
	case 3000:
		out += "Difficulty: hard\n"
	}
	out += "----------\n\n"
	return out
}
