package rp

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

type Color int

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

	`/me manipulates %s into fully bending their legs before wrapping each one
	up in saran wrap. Their arms are fully bent, hands pressed against their
	own shoulders, before being wrapped up in the same fashion, rendering their
	[u]torso[/u] completely vulnerable.`,

	`/me shoves %s down against a wooden chair and cuffs their wrists together
	behind the back of it. Their ankles are bent underneath the seat and cuffed
	to the support stretcher that links the chair legs, exposing their
	[u]feet[/u] as well as their [u]upperbody[/u] above.`,

	`/me invites %s to join them in a Yoga session. First up: Lotus position!
	Unfortunately for them, several ropes are wrapped around their calves and
	ankles, leaving their [u]soles[/u] upturned. Their arms are also crossed
	behind their back with their forearms bound in rope, parallel with each
	other, to align their chi and leave their [u]upperbody[/u] defenseless.`,

	`/me helps %s with their daily stretches. One arm is pulled over their
	shoulder while the other is pulled down and around by their lower back.
	Their two wrists become joined by a pair of leather cuffs, leaving them in
	an awkward pose that exposes [u]one side[/u] and the [u]front of their
	upperbody[/u].`,

	`/me adds a flair of fashion to %s by placing a lovely set of leather
	bondage mittens over their hands. Once locked shut, their fingers cannot
	manipulate anything through the thick leather. To add to their distress,
	leather cuffs are attached to their ankles, with only six-inches of slack
	between their ankles for hobbling around.`,

	`/me suddenly triggers a rope noose trap underneath %s. Their ankles are
	snagged and hoisted up into the air, but only high enough to flip them onto
	their back. Both [u]feet[/u] are vulnerable in the air unless they are
	flexible enough to reach all the way up to them.`,

	`/me plays the flute, causing a pile of bandages to rise up like snakes.
	They quickly slither out and wrap around %s from their ankles to their
	shoulders, leaving their [u]head[/u] and [u]feet[/u] exposed. Once
	complete, the bandages stretch up to the ceiling and pull their wrapped
	prey with them, leaving them dangling a few inches off the floor.`,

	`/me knocks %s over and sets a chair down on top of their body. Their
	wrists are wrapped up in rope against the chair legs, and their ankles are
	hoisted up and tied to the headrest of the chair. Their [u]feet[/u] are
	vulnerable to all, while their [u]head[/u] is exposed and forced to watch
	anyone who sits down above them.`,

	`/me quickly fastens two open-stocks to dangle from the ceiling, padded for
	maximum comfort. It lifts %s into the air, their neck and wrists placed in
	one and ankles in the other, before closing them with loud smacks, leaving
	their [u]entire body[/u] exposed and accessible. It whirs, performing some
	additional calculations, then it produces some precisely-measured silken
	cords and uses them to tie their toes to the top of the stocks.`,
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

type Tktool struct {
	Raw    string
	Colors []Color
}

func (t Tktool) Apply(user string, c Color) string {
	if t.Colors != nil {
		return fmt.Sprintf(clean(t.Raw), user, c)
	}
	return fmt.Sprintf(clean(t.Raw), user)
}

var tktools = []Tktool{
	{
		Colors: []Color{Gray, White, Black},
		Raw: `/me hands %s a stiff, long, %s, [u]goose feather[/u] with a
		pointy tip.`,
	},
	{
		Colors: []Color{Black, White, Red, Orange, Blue, Turquoise, Brown, Yellow, Fuchsia, Pink, Purple, Violet, Green},
		Raw: `/me hands %s a jaunty, enormous, %s, [u]ostrich feather[/u]
		forming a slight curve at the top. Its shaft at the bottom, ends into a
		sharp quill.`,
	},
	{
		Colors: []Color{Pink, Fuchsia, Purple, Black, Emerald, Red, Yellow, Blue},
		Raw: `/me hands %s a fluffy, long, %s, [u]chandelle feather boa[/u].
		With the slightest movement, its plumes animate in an almost hypnotic
		way.`,
	},
	{
		Raw: `/me hands %s an aqua-colored [u]electric flosser[/u], equipped
		with a fully charged battery and a flexible, nylon tip.`,
	},
	{
		Raw: `/me hands %s an [u]electric toothbrush[/u]. The brush is
		round-shaped. Its body is light-blue and contains lots of colorful
		smiley faces.`,
	},
	{
		Raw: `/me hands %s a small, brown [u]paintbrush[/u] with soft bristles
		and a pointy tip. On its black, wooden body, there are the characters
		搔癢折磨 inscriped in crimson red.`,
	},
	{
		Colors: []Color{Gray, White, Black, Brown},
		Raw: `/me hands %s a [u]feather duster[/u] which looks like a matching
		accessory for a maid uniform. Its long, %s, ostrich feathers look
		very soft and delicate.`,
	},
	{
		Colors: []Color{Brown, Black, Violet, Purple, Red},
		Raw: `/me hands %s a pair of %s [u]leather gloves[/u] with long
		feathers attached to each fingertip.`,
	},
	{
		Raw: `/me hands %s a modified [u]Hitachi Magic Wand[/u]. This model
		seems to be cordless and its switch seems to be altered. Apart from the
		traditional, O, I and II power levels, this switch supports two extra
		levels indicated as III and XXX.`,
	},
	{
		Raw: `/me hands %s a set of wearable, metallic, [u]cat claws[/u]. The
		pointy tips of the claws seem to be sharp enough for play but not
		harm.`,
	},
	{
		Raw: `/me hands %s a bottle of [u]baby oil[/u]. The label of the bottle
		reads: 'A wonderful fragrance, for silky soft skin and soothing
		massages.'`,
	},
	{
		Colors: []Color{Black, Blue, Violet, Pink, Orange, Purple, Brown},
		Raw: `/me hands %s a large, %s [u]porcupine grooming brush[/u] with
		dense, black nylon bristles that form a slight curve.`,
	},
}

func RandTktool(name string) string {
	t := tktools[newRand(len(tktools))]
	return t.Apply(name, t.Colors[newRand(len(t.Colors))])
}
