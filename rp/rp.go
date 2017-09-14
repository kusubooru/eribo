package rp

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func newRand(n int) int {
	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))
	return r.Intn(n)
}

func RandTieUp(victim string) string {
	s := tieUps[newRand(len(tieUps))]
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\t", " ", -1)
	return fmt.Sprintf(s, victim)
}

var tieUps = []string{
	`/me grabs %s and wraps their body tightly using saran wrap, leaving only
	their [u]head[/u] and [u]feet[/u] exposed. Then places the wrapped body on
	the table and starts strapping it. It applies tight straps above and under
	their chest, on their waist, thighs, knees and ankles rendering the victim
	immobile.`,
	`/me lifts %s up by their arms holding them above their head and swiftly
	ties them together. Then it proceeds to wrap the rest of the victim's body
	up with saran wrap, leaving their [u]head[/u], [u]underarms[/u] and
	[u]feet[/u] vulnerable. Lastly, it places the victim's body on a rack and
	applies straps on their wrists, elbows, waist, thighs, knees and ankles
	rendering them immobile.`,
	`/me grabs %s and forces their arms behind their back, locking their wrists
	in a pair of leather cuffs. Then sits them down, placing their ankles in
	the stocks and finally locking them up, leaving their [u]feet[/u]
	vulnerable.`,
}

func RandFeedback(name string) string {
	s := feedback[newRand(len(feedback))]
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\t", " ", -1)
	return fmt.Sprintf(s, name)
}

var feedback = []string{
	`/me bows politely, "Thank you for the feedback %s-sama".`,
	`/me bows graciously, "Your feedback is highly appreciated %s-sama".`,
	`/me nods affirmatively, "Understood %-sama. Your feedback has been recorded".`,
}
