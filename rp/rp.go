package rp

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// Quality is the quality of an  item, folowing the World of Warcrarft system.
type Quality int

// Possible quality values.
const (
	Unknown Quality = iota
	Poor
	Common
	Uncommon
	Rare
	Epic
	Legendary
)

func (q Quality) String() string {
	switch q {
	default:
		return "unknown"
	case Poor:
		return "poor"
	case Common:
		return "common"
	case Uncommon:
		return "uncommon"
	case Rare:
		return "rare"
	case Epic:
		return "epic"
	case Legendary:
		return "legendary"
	}
}

func makeQuality(s string) Quality {
	switch s {
	default:
		return Unknown
	case "poor":
		return Poor
	case "common":
		return Common
	case "uncommon":
		return Uncommon
	case "rare":
		return Rare
	case "epic":
		return Epic
	case "legendary":
		return Legendary
	}
}

// Weight returns the weight chance factor depending on the item quality.
func (q Quality) Weight() int {
	switch q {
	default:
		return 0 // 0
	case Poor:
		return 40 // 30
	case Common:
		return 50 // 40
	case Uncommon:
		return 30 // 20
	case Rare:
		return 20 // 10
	case Epic:
		return 5 // 5
	case Legendary:
		return 1 // 1
	}
}

// Color represents the color of an item.
type Color int

// All the possible color values.
const (
	Colorless Color = iota
	Red
	Blue
	Green
	Yellow
	Orange
	Purple
	Pink
	Fuchsia
	Black
	White
	Emerald
	Brown
	Violet
	Gray
	Turquoise
)

func (c Color) String() string {
	switch c {
	default:
		return "colorless"
	case Red:
		return "red"
	case Blue:
		return "blue"
	case Green:
		return "green"
	case Yellow:
		return "yellow"
	case Orange:
		return "orange"
	case Purple:
		return "purple"
	case Pink:
		return "pink"
	case Fuchsia:
		return "fuchsia"
	case Black:
		return "black"
	case White:
		return "white"
	case Emerald:
		return "emerald"
	case Brown:
		return "brown"
	case Violet:
		return "violet"
	case Gray:
		return "gray"
	case Turquoise:
		return "turquoise"
	}
}

func newRand(n int) int {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return r.Intn(n)
}

func clean(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\t\t", " ", -1)
	s = strings.Replace(s, "\t", " ", -1)
	return s
}

// RandFeedback returns a polite response when a player gives feedback.
func RandFeedback(name string) string {
	s := feedback[newRand(len(feedback))]
	return fmt.Sprintf(clean(s), name)
}

var feedback = []string{
	`/me bows politely, "Thank you for the feedback %s".`,
	`/me bows graciously, "Your feedback is highly appreciated %s".`,
	`/me nods affirmatively, "Understood %s. Your feedback has been recorded".`,
}

// Tomato returns a message for when the player uses the tomato command.
func Tomato(name, owner string) string {
	if name == owner {
		var s = `/me humbly offers a juicy and fresh-looking tomato to %s, "A
		pleasure to serve you Ryuunosuke-sama".`
		return fmt.Sprintf(clean(s), name)
	}
	return fmt.Sprintf("/me gives a fresh-looking tomato to %s.", name)
}

func qualityColor(q Quality) string {
	switch q {
	case Poor:
		return "gray"
	case Common:
		return "white"
	case Uncommon:
		return "green"
	case Rare:
		return "blue"
	case Epic:
		return "purple"
	case Legendary:
		return "orange"
	default:
		return "white"
	}
}

func qualityColorBBCode(q Quality, s string) string {
	return fmt.Sprintf("[color=%s]%s[/color]", qualityColor(q), s)
}

type animeOp struct {
	Raw string
}

func (j animeOp) Apply() string {
	return clean(j.Raw)
}

var animeOps = []animeOp{
	{
		Raw: `/me plays: [url=https://www.youtube.com/watch?v=yF00xX-p28Y]Jojo
		Opening 1 - Sono Chi no Sadame[/url]`,
	},
	{
		Raw: `/me plays: [url=https://www.youtube.com/watch?v=i-GWFGwbEPg]Jojo
		Opening 2 - BLOODY STREAM[/url]`,
	},
	{
		Raw: `/me plays: [url=https://www.youtube.com/watch?v=RordBk3Ztk4]Jojo
		Opening 3 - STAND PROUD[/url]`,
	},
	{
		Raw: `/me plays: [url=https://www.youtube.com/watch?v=f0yK_7adSCA]Jojo
		Opening 4 - Sono Chi no Kioku[/url]`,
	},
	{
		Raw: `/me plays: [url=https://www.youtube.com/watch?v=nNQ-Qi7pBpw]Jojo
		Opening 5 - Crazy Noisy Bizarre Town[/url]`,
	},
	{
		Raw: `/me plays: [url=https://www.youtube.com/watch?v=qr73XbBFDA8]Jojo
		Opening 6 - chase[/url]`,
	},
	{
		Raw: `/me plays: [url=https://www.youtube.com/watch?v=zoqH1Rk4ANM]Jojo
		Opening 7 - Great Days[/url]`,
	},
}
