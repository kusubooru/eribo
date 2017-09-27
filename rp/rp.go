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

func clean(s string) string {
	s = strings.Replace(s, "\n", "", -1)
	s = strings.Replace(s, "\t", " ", -1)
	return s
}

func RandTieUp(victim string) string {
	s := tieUps[newRand(len(tieUps))]
	return fmt.Sprintf(clean(s), victim)
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
	return fmt.Sprintf(clean(s), name)
}

var feedback = []string{
	`/me bows politely, "Thank you for the feedback %s-sama".`,
	`/me bows graciously, "Your feedback is highly appreciated %s-sama".`,
	`/me nods affirmatively, "Understood %s-sama. Your feedback has been recorded".`,
}

func Tomato(name string) string {
	if name == "Ryuunosuke Akashaka" {
		var s = `/me humbly offers a juicy and fresh-looking tomato to %s, "A
		pleasure to serve you Dragon-sama".`
		return fmt.Sprintf(clean(s), name)
	}
	return fmt.Sprintf("/me gives a fresh-looking tomato to %s.", name)
}

var tktools = []string{
	`/me hands %s a stiff, long, gray, [u]goose feather[/u] with a pointy
	tip.`,
	`/me hands %s a jaunty, enormous, white, [u]ostrich feather[/u] forming a
	slight curve at the top. Its shaft at the bottom, ends into a sharp
	quill.`,
	`/me hands %s a fluffy, long, pink, [u]chandelle feather boa[/u]. With the
	slightest movement, its plumes animate in an almost hypnotic way.`,
	`/me hands %s an aqua-colored [u]electric flosser[/u], equipped with a
	fully charged battery and a flexible, nylon tip.`,
	`/me hands %s an [u]electric toothbrush[/u]. The brush is round-shaped. Its
	body is light-blue and contains lots of colorful smiley faces.`,
	`/me hands %s a small, brown [u]paintbrush[/u] with soft bristles and a
	pointy tip. On its black, wooden body, there are the characters 搔癢折磨
	inscriped in crimson red.`,
	`/me hands %s a [u]feather duster[/u] which looks like a matching accessory
	for a maid uniform. Its long, gray, ostrich feathers look very soft and
	delicate.`,
	`/me hands %s a pair of black [u]leather gloves[/u] with long feathers
	attached to each fingertip.`,
	`/me hands %s a modified [u]Hitachi Magic Wand[/u]. This model seems to be
	cordless and its switch seems to be altered. Apart from the traditional, O,
	I and II power levels, this switch supports two extra levels indicated as
	III and XXX.`,
}

func Tktool(name string) string {
	s := tktools[newRand(len(tktools))]
	return fmt.Sprintf(clean(s), name)
}
