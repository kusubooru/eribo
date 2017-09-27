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
	`/me bends %s forward into an awaiting standing pillory and shuts it on
	their neck and wrists. A spreader bar is then cuffed to their ankles,
	forcing their legs far apart. With their [u]sides[/u] and [u]legs[/u]
	rather vulnerable, they cannot kick effectively nor see behind them.`,
	`/me suddenly pushes %s backwards into an open coffin with two slots at the
	bottom. Their ankles get caught in the slots before the lid automatically
	slams shut and locks itself, leaving their [u]feet[/u] exposed on the
	outside.`,
	`/me extends seven spider-like legs from behind and uses them to lift %s
	into the air. An eighth leg is revealed to be pulling silk webbing up from
	a large spool, and it quickly rolls them up, leaving only their
	[u]stomach[/u] and nose exposed. It then sticks their back against a wall
	so that they can't wiggle away.`,
	`/me puts on a cowboy hat and swings a lasso onto %s before drawing them in
	and wrestling them to the ground. The excess rope is wound around their
	ankles while their wrists are positioned behind them, and the slack is
	tightened, forcing them into a hogtie that renders their [u]sides[/u] and
	[u]feet[/u] quite vulnerable.`,
	`/me tightens a belt around %s that has a leather cuff dangling from either
	side. Then, it forces their hands down into the awaiting cuffs and buckles
	them shut, leaving their hands trapped by their waist.`,
	`/me deems %s is getting too unruly and takes measures to protect them. A
	straitjacket is pulled onto them and buckled shut, forcing their arms
	crossed in front of themselves. Although their upperbody is secure and
	protected, their [u]legs[/u] and [u]feet[/u] remain uncovered.`,
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
