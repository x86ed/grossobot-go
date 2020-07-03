package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//Case to trigger an image
type Case struct {
	Triggers    []string
	Images      []string
	Type        int
	Description string
}

func getFile(s string) io.Reader {
	if strings.HasPrefix(s, "../") {
		buf, _ := os.Open(s)
		return buf
	}
	r, _ := http.Get(s)

	return r.Body
}

//Process processes an image case
func (c *Case) Process(s *discordgo.Session, m *discordgo.MessageCreate) {
	cont := m.Content
	rand := rand.Intn(len(c.Images))
	for _, t := range c.Triggers {
		switch c.Type {
		case EQUALS:
			if strings.ToLower(cont) == fmt.Sprintf("**%s**", t) {
				f := c.Images[rand]
				buf := getFile(f)
				s.ChannelFileSend(m.ChannelID, f, buf)
			}
		case CONTAINS:
			if strings.Contains(strings.ToLower(cont), fmt.Sprintf("**%s**", t)) {
				f := c.Images[rand]
				buf := getFile(f)
				s.ChannelFileSend(m.ChannelID, f, buf)
			}
		case DOESNTEQUAL:
			if strings.ToLower(cont) == fmt.Sprintf("**%s**", t) != true {
				f := c.Images[rand]
				buf := getFile(f)
				s.ChannelFileSend(m.ChannelID, f, buf)
			}
		case DOESNTCONTAIN:
			if strings.Contains(strings.ToLower(cont), fmt.Sprintf("**%s**", t)) != true {
				f := c.Images[rand]
				buf := getFile(f)
				s.ChannelFileSend(m.ChannelID, f, buf)
			}
		case REGEX:
			re := regexp.MustCompile(cont)
			if len(re.FindStringIndex(fmt.Sprintf("**%s**", t))) > 0 {
				f := c.Images[rand]
				buf := getFile(f)
				s.ChannelFileSend(m.ChannelID, f, buf)
			}
		}
	}
}

var wray = Case{
	Triggers:    []string{"wray", "water tower"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/thPYncq/wray.gif"},
	Description: "Jeremy Wray laying it down.",
}

var grosso = Case{
	Triggers:    []string{"grosso", "jeff grosso", "the letters", "loveletters to skateboarding"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/2hCHH5H/grosso.gif", "https://i.ibb.co/jvf5BSK/vans.gif", "https://i.ibb.co/4RBtbVC/grossobot.gif", "https://i.ibb.co/f8HsqvS/curbjeff.gif", "https://i.ibb.co/TmTYDhD/combi.gif"},
	Description: "A loveletter to Jeff.",
}

var noGrass = Case{
	Triggers:    []string{"no grass", "do it rolling", "grass", "rolling", "doesn't count"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/Hn63Dtw/image.gif", "https://i.ibb.co/ynzyC8M/girl.gif", "https://i.ibb.co/12bqgv8/nuts.gif", "https://i.ibb.co/tK6Pb3z/grass.gif"},
	Description: "A demonstration of what your skating looks like.",
}

var pissDrunx = Case{
	Triggers:    []string{"pd", "piss drunx"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/zmM4K9m/boulala.gif", "https://i.ibb.co/K5hWG5S/sunset.gif", "https://i.ibb.co/1XzCkK2/grecs.gif", "https://i.ibb.co/9bzj9gZ/twuan.gif", "https://i.ibb.co/NScxfNR/carlsbad.gif", "https://i.ibb.co/Z2YxFQV/dustin.gif"},
	Description: "ᛩ▹: party like it's 1999",
}

var twuan = Case{
	Triggers:    []string{"antwuan", "twuan", "dixon", "antwuan dixon", "eat dat shrimp"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/9bzj9gZ/twuan.gif", "https://i.ibb.co/NScxfNR/carlsbad.gif", "https://i.ibb.co/DWVKpSf/datshrimp.gif"},
	Description: "Heelflippin to Biggy.",
}

var baker = Case{
	Triggers:    []string{"baker", "baker baker baker", "bakerbakerbaker", "deathwish", "shalom"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/w4HwfhT/bakerbakerbaker.gif", "https://i.ibb.co/K5hWG5S/sunset.gif", "https://i.ibb.co/1XzCkK2/grecs.gif", "https://i.ibb.co/9bzj9gZ/twuan.gif", "https://i.ibb.co/NScxfNR/carlsbad.gif", "https://i.ibb.co/Z2YxFQV/dustin.gif", "https://i.ibb.co/5MWfSFh/kennedy.gif"},
	Description: "You might say they have a Deathwᚯsh.",
}

var greco = Case{
	Triggers:    []string{"greco", "jim greco", "hammers usa"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/K5hWG5S/sunset.gif", "https://i.ibb.co/1XzCkK2/grecs.gif"},
	Description: "The John Cassavettes of skateboarding.",
}

var natas = Case{
	Triggers:    []string{"natas", "satan", "101"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/wwTT0xx/natas.gif"},
	Description: "Sit n' spin.",
}

var manramp = Case{
	Triggers:    []string{"manramp", "man ramp", "worble"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/njTFYCs/manramp.gif"},
	Description: "He's half man half ramp.",
}

var mullen = Case{
	Triggers:    []string{"mullen", "dwindle", "rodney", "freestyle", "tensor"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/YQN3gWz/chess.gif", "https://i.ibb.co/4RMr0w1/casper.gif", "https://i.ibb.co/KNZ5Ybk/delmar.gif", "https://i.ibb.co/TrLJvdg/vanityfair.gif"},
	Description: "Yeah he invented that trick.",
}

var bam = Case{
	Triggers:    []string{"cky", "jackass", "bam", "him", "margera", "bam margera", "heartagram"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/yVtYqmw/bam.gif", "https://i.ibb.co/8M1zqvF/cky.gif", "https://i.ibb.co/nfZM45x/famousstarsandstraps.gif"},
	Description: "Tony Hawk's Wario.",
}

var chin = Case{
	Triggers:    []string{"animal chin", "have you seen him", "bones brigade", "skate and create"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/XD9nWh4/chin.gif", "https://i.ibb.co/YQN3gWz/chess.gif", "https://i.ibb.co/KNZ5Ybk/delmar.gif", "https://i.ibb.co/WsqMrHf/ripper.gif", "https://i.ibb.co/stg6LcN/japan.gif"},
	Description: "Have you seen him outside a little town by Guadalupe?",
}

var hammers = Case{
	Triggers:    []string{"hammers", "yew", "yeet", "hype", "rad", "gnarly", "rips", "get it", "full send", "send it", "stoked", "gang gang gang", "lucky", "he on xgames mode"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/ss7krZ0/900.gif", "https://i.ibb.co/RjL95p8/lyon.gif", "https://i.ibb.co/K5hWG5S/sunset.gif", "https://i.ibb.co/9rb36tf/booth.gif", "https://i.ibb.co/1bR0b1w/mo.gif", "https://i.ibb.co/9r0ttwn/lakai.gif", "https://i.ibb.co/yFn8Txh/wall.gif", "https://i.ibb.co/HGwyRHH/tree.gif", "https://i.ibb.co/NScxfNR/carlsbad.gif", "https://i.ibb.co/cYSqN0K/staygold.gif", "https://i.ibb.co/dJ4pwdB/kirchart.gif"},
	Description: "He on X-Games Mode!!!",
}

var tony = Case{
	Triggers:    []string{"tony hawk", "thps", "birdhouse", "900"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/ss7krZ0/900.gif", "https://i.ibb.co/x8XGmv2/thps.gif", "https://i.ibb.co/9rb36tf/booth.gif", "https://i.ibb.co/xfnrGPm/hoverboard.gif"},
	Description: "The only skater your grandparents know.",
}

var dt = Case{
	Triggers:    []string{"zephyr", "sma", "santa monica airlines", "dogtown", "dog town", "venice", "locals only", "alva", "bertleman"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/wwTT0xx/natas.gif", "https://i.ibb.co/zmtnMBm/pool.gif", "https://i.ibb.co/chcPctH/dogtown.gif", "https://i.ibb.co/Nj1ysJd/stacey.gif", "https://i.ibb.co/2FGB5Z5/jay.gif", "https://i.ibb.co/RHP3Tw2/alva.gif", "https://i.ibb.co/Bg1zz89/oldschool.gif"},
	Description: "Memes for locals only.",
}

var kook = Case{
	Triggers:    []string{"kook", "longboard", "short board", "penny", "nickel", "boosted", "rip n dip", "ripndip", "revive", "braille"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/MGCnDmw/kyro.gif", "https://i.ibb.co/Hn63Dtw/image.gif", "https://i.ibb.co/5j6fxqS/savage.gif", "https://i.ibb.co/y844WKw/ollie.gif", "https://i.ibb.co/QfCRKDj/slalom.gif", "https://i.ibb.co/vwhSDyh/tophat.gif", "https://i.ibb.co/1mYGHgs/mall.gif", "https://i.ibb.co/Y36SFcg/boosted.gif", "https://i.ibb.co/LSjwhz5/vespakook.gif", "https://i.ibb.co/7VJBXLx/buscemi.gif", "https://i.ibb.co/C5qR4Qg/policekook.gif", "https://i.ibb.co/84dHtZR/skooterbooter.gif"},
	Description: "You know it when you see it.",
}

var oof = Case{
	Triggers:    []string{"oof", "slam", "wasted", "sucky", "ouch", "half send"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/WPqS8Jw/jarne.gif", "https://i.ibb.co/6J1KYt4/skin.png", "https://i.ibb.co/HLy99YM/duffman.gif", "https://i.ibb.co/VpK322W/wasted.gif", "https://i.ibb.co/zmM4K9m/boulala.gif", "https://i.ibb.co/YNPQPS2/sucky.gif", "https://i.ibb.co/xS3Pt94/sacked.gif", "https://i.ibb.co/qdRdypn/fence.gif", "https://i.ibb.co/b1vdd94/carded.gif"},
	Description: "Hall o' Meat.",
}

var weck = Case{
	Triggers:    []string{"weck", "weckingball", "don't @ me", "yummygod"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/hczZ7Jw/wek.png"},
	Description: "That post scraped.",
}

var wack = Case{
	Triggers:    []string{"wack", "dime", "jamal smith", "palace"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/5Fg0YKJ/wack.gif"},
	Description: "Fuck Wade Desarmo!",
}

var footy = Case{
	Triggers:    []string{"footy", "footage", "film", "clips", "pics", "post footy", "frog", "🐸"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/nL3QFwj/footy.jpg", "https://i.ibb.co/nLWhH42/vx.gif", "https://i.ibb.co/m8NdsDQ/image0-17.jpg", "https://i.ibb.co/BC84Fsp/image0.jpg"},
	Description: "Send Footage DH.",
}

var crob = Case{
	Triggers:    []string{"crob", "the nine club", "9 club", "nine club"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/h8m7243/crob.jpg", "https://i.ibb.co/3TpWdWv/switchflip.gif"},
	Description: "The GOAT.",
}

var kflip = Case{
	Triggers:    []string{"do a kickflip"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/DL3WGtr/kostonkickflip.jpg", "https://i.ibb.co/PxpyNFf/muskakickflip.jpg", "https://i.ibb.co/H25sS5g/colekickflip.jpg", "https://i.ibb.co/2k6Gy1M/kickfliphawk.jpg"},
	Description: "Do A Kickflip!",
}

var muska = Case{
	Triggers:    []string{"muska", "chad muska", "supra", "muskabeats"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/xM457Nz/adopted.gif", "https://i.ibb.co/hs1MnLN/boombox.gif", "https://i.ibb.co/yqvXR1M/elonmuska.gif", "https://i.ibb.co/NYP1r31/muska.gif"},
	Description: "Blastin' Muskabeats in The Muskalade.",
}

var uh = Case{
	Triggers:    []string{"uh", "uhm"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/h8m7243/crob.jpg"},
	Description: "Okay... if that's what you want",
}

var shecks = Case{
	Triggers:    []string{"sheckler", "bs flip", "Ryan Sheckler", "el toro"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/N1CNCpt/sheckler.png"},
	Description: "The most loved skater sinve Jereme Rogers.",
}

var skin = Case{
	Triggers:    []string{"skin", "god"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/6J1KYt4/skin.png"},
	Description: "Oi Cunt!",
}

var owen = Case{
	Triggers:    []string{"hollywood high", "front salad", "back salad", "front blunt", "sylmar"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/Vjy7K50/salad.png"},
	Description: "Basically it's all been done in 411 issue 52.",
}

var grecs = Case{
	Triggers:    []string{"security guard", "rent a cop", "tantrum"},
	Type:        CONTAINS,
	Images:      []string{"https://i.ibb.co/1XzCkK2/grecs.gif"},
	Description: "I ❤️ 🐖s.",
}

var cases = []Case{
	baker,
	bam,
	chin,
	crob,
	dt,
	footy,
	greco,
	grecs,
	grosso,
	hammers,
	kflip,
	kook,
	manramp,
	mullen,
	muska,
	natas,
	noGrass,
	oof,
	owen,
	pissDrunx,
	shecks,
	skin,
	tony,
	twuan,
	uh,
	wack,
	weck,
	wray,
}
