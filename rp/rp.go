package rp

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/kusubooru/eribo/eribo"
	"github.com/kusubooru/eribo/loot"
)

type Quality int

const (
	Unknown Quality = iota
	Poor
	Common
	Uncommon
	Rare
	Epic
	Legendary
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
	`/me bows politely, "Thank you for the feedback %s".`,
	`/me bows graciously, "Your feedback is highly appreciated %s".`,
	`/me nods affirmatively, "Understood %s. Your feedback has been recorded".`,
}

func Tomato(name, owner string) string {
	if name == owner {
		var s = `/me humbly offers a juicy and fresh-looking tomato to %s, "A
		pleasure to serve you Ryuunosuke-sama".`
		return fmt.Sprintf(clean(s), name)
	}
	return fmt.Sprintf("/me gives a fresh-looking tomato to %s.", name)
}

type Tktool struct {
	Name     string
	Colors   []Color
	Poor     string
	Common   string
	Uncommon string
	Weight   int
}

func (t Tktool) Apply(user string, q Quality) string {
	desc := ""
	switch q {
	case Poor:
		desc = t.Poor
	case Common:
		desc = t.Common
	case Uncommon:
		desc = t.Uncommon
	default:
		desc = t.Common
	}

	if t.Colors != nil {
		c := t.Colors[newRand(len(t.Colors))]
		return fmt.Sprintf(clean(desc), user, c)
	}
	return fmt.Sprintf(clean(desc), user)
}

var tktools = []Tktool{
	{
		Name:   "[Goose Feather]",
		Colors: []Color{Gray, White, Black},
		Poor: `/me hands %s an old, %s [color=gray][Ravaged Goose
		Feather][/color]. It's too wrecked to be usable for any meaningful
		purpose.`,
		Common: `/me hands %s a stiff, %s [color=white][Goose Feather][/color].`,
		Uncommon: `/me hands %s a long, %s [color=green][Pristine Goose
		Feather][/color] with a pointy tip.`,
		Weight: 10,
	},
	{
		Name:   "[Ostrich Feather]",
		Colors: []Color{Black, White, Red, Orange, Blue, Turquoise, Brown, Yellow, Fuchsia, Pink, Purple, Violet, Green},
		Poor: `/me hands %s the remains of an old, %s [color=gray][Ruined
		Ostrich Feather][/color].`,
		Common: `/me hands %s a large, %s [color=white][Ostrich
		Feather][/color].`,
		Uncommon: `/me hands %s an enormous, %s [color=green][Jaunty Ostrich
		Feather][/color] which forms a slight curve at the top. Its shaft at
		the bottom, ends into a sharp quill.`,
		Weight: 10,
	},
	{
		Name:   "[Feather Boa]",
		Colors: []Color{Pink, Fuchsia, Purple, Black, Emerald, Red, Yellow, Blue},
		Poor: `/me hands %s an old, %s [color=gray][Destroyed Feather
		Boa][/color]. The remains of the wrecked item don't look usable for any
		meaningful purpose.`,
		Common: `/me hands %s a long, %s [color=white][Feather Boa][/color].`,
		Uncommon: `/me hands %s a fluffy, long, %s [color=green][Chandelle
		Feather Boa][/color]. With the slightest movement, its plumes animate
		in an almost hypnotic way.`,
		Weight: 3,
	},
	{
		Name: "[Electric Flosser]",
		Poor: `/me hands %s an [color=gray][Inoperable Electric Flosser][/color]. It
		doesn't look functional anymore and the tip is missing.`,
		Common: `/me hands %s an [color=white][Electric Flosser][/color].`,
		Uncommon: `/me hands %s an [color=green][Aqua-colored Electric
		Flosser][/color]. It is equipped with a fully charged battery and a
		flexible, nylon tip.`,
		Weight: 6,
	},
	{
		Name: "[Electric Toothbrush]",
		Poor: `/me hands %s a [color=gray][Busted Electric Toothbrush][/color].
		It looks broken beyond repair and the brush is destroyed.`,
		Common: `/me hands %s an [color=white][Electric Toothbrush][/color].`,
		Uncommon: `/me hands %s a [color=green][Happy Electric
		Toothbrush][/color]. The brush is round-shaped. Its body is light-blue
		and contains lots of colorful smiley faces.`,
		Weight: 8,
	},
	{
		Name: "[Small Paintbrush]",
		Poor: `/me hands %s a [color=gray][Snapped Paintbrush][/color]. Beyond
		its broken body, there seem to be no bristles left on its tip.`,
		Common: `/me hands %s a [color=white][Small Paintbrush][/color].`,
		Uncommon: `/me hands %s a small, brown [color=green][Eastern
		Paintbrush][/color] with soft bristles and a pointy tip. On its black,
		wooden body, there are the characters 搔癢折磨 inscribed in crimson
		red.`,
		Weight: 10,
	},
	{
		Name:   "[Feather Duster]",
		Colors: []Color{Gray, White, Black, Brown},
		Poor: `/me hands %s a wrecked, %s [color=gray][Snapped Feather
		Duster][/color]. The handle is broken and the top part, barely
		resembles a duster`,
		Common: `/me hands %s a clean, %s [color=white][Feather
		Duster][/color].`,
		Uncommon: `/me hands %s an [color=green][Impeccable Feather
		Duster][/color] which looks like a matching accessory for a maid
		uniform. Its long, %s, ostrich feathers look very soft and delicate.`,
		Weight: 10,
	},
	{
		Name:   "[Feather Gloves]",
		Colors: []Color{Brown, Black, Violet, Purple, Red},
		Poor: `/me hands %s a pair of %s [color=gray][Destroyed Feather
		Gloves][/color]. The majority of the feathers that were previously
		attached on each fingertip seem to be missing and the ones that remain
		are totally ruined.`,
		Common: `/me hands %s a pair of %s [color=white][Feather
		Gloves][/color]. On each fingertip there's a feather attached.`,
		Uncommon: `/me hands %s a pair of expensive, %s
		[color=green][Unblemished Feather Gloves][/color]. They are made out of
		high quality leather and on each fingertip there's a long, pristine
		feather attached.`,
		Weight: 5,
	},
	{
		Name: "[Hitachi Magic Wand]",
		Poor: `/me hands %s a [color=gray][Destroyed Hitachi Magic
		Wand][/color]. The device is missing its cord and it looks broken
		beyond repair.`,
		Common: `/me hands %s a [color=white][Hitachi Magic Wand][/color]
		electrical massager.`,
		Uncommon: `/me hands %s a [color=green][Modified Hitachi Magic
		Wand][/color]. This model seems to be cordless and its switch seems to
		be altered. Apart from the traditional, O, I and II power levels, this
		switch supports two extra levels indicated as III and XXX.`,
		Weight: 5,
	},
	{
		Name: "[Metallic Cat Claws]",
		Poor: `/me hands %s a set of [color=gray][Dented Cat Claws][/color].
		Each piece is so horribly dented that is impossible to wear.`,
		Common: `/me hands %s a set of wearable [color=white][Metallic Cat
		Claws][/color].`,
		Uncommon: `/me hands %s a set of wearable, well-crafted
		[color=green][Silver Cat Claws][/color]. The pointy tips of the claws
		seem to be sharp enough for play but not harm.`,
		Weight: 3,
	},
	{
		Name: "[Bottle of Baby Oil]",
		Poor: `/me hands %s a [color=gray][Vial of Fish Oil][/color]. The
		content looks dubious at best. Better not open it!`,
		Common: `/me hands %s a small [color=white][Bottle of Baby Oil][/color].`,
		Uncommon: `/me hands %s a [color=green][Bottle of Pure Lavender
		Oil][/color]. The label of the bottle reads: 'A wonderful fragrance,
		for silky soft skin and soothing massages.'`,
		Weight: 7,
	},
	{
		Name:   "[Grooming Brush]",
		Colors: []Color{Black, Blue, Violet, Pink, Orange, Purple, Brown},
		Poor: `/me hands %s an old, %s [color=gray][Ruined Grooming
		Brush][/color]. Most of the bristles are missing.`,
		Common: `/me hands %s a small, %s [color=white][Grooming
		Brush][/color].`,
		Uncommon: `/me hands %s a large, %s [color=green][Porcupine Grooming
		Brush][/color] with dense, black nylon bristles that form a slight
		curve.`,
		Weight: 8,
	},
}

func Tktools() []Tktool {
	return tktools
}

func RandTktool(name string) string {
	d := []loot.Drop{
		{Item: Poor, Weight: 8},
		{Item: Common, Weight: 40},
		{Item: Uncommon, Weight: 4},
	}
	table := loot.NewTable(d)
	_, roll := table.Roll(time.Now().UnixNano())
	quality := Unknown
	if q, ok := roll.(Quality); ok {
		quality = q
	}

	table = &loot.Table{}
	for _, t := range tktools {
		table.Add(t, t.Weight)
	}
	_, roll = table.Roll(time.Now().UnixNano())
	if tool, ok := roll.(Tktool); ok {
		return tool.Apply(name, quality)
	}

	// Just in case.
	tool := tktools[newRand(len(tktools))]
	return tool.Apply(name, quality)
}

type Vonprove struct {
	Raw         string
	HasDate     bool
	HasDuration bool
	HasUser     bool
}

var vonproves = []Vonprove{
	{
		Raw: `/me turns around and points at its own butt. Upon a closer
		inspection of the curvy surface, Von's seal of approval can be seen.`,
	},
	{
		Raw: `/me opens a small drawer on its body where Seal, its animal
		companion, can be found sleeping. When Seal realizes the drawer is
		open, quickly wears a paper mask of Von Vitae and starts nodding in
		approval.`,
	},
	{
		Raw: `/me turns on the monitor on its chest, where the classic "seal of
		approval" meme appears with the word "Von" written on top in Impact
		font.`,
	},
	{
		Raw: `/me strikes a pose and proudly shows off the renowned seal of
		approval, while Von Vitae's theme song plays from its speakers in 8-bit
		chiptune.`,
	},
	{
		HasUser: true,
		Raw: `/me quickly launches itself towards %s, smacks their forehead
		with a rubber stamp bearing Von Vitae's seal of approval and shouts,
		"Vonproved™!".`,
	},
	{
		HasDate: true,
		Raw: `/me stands still, looks upwards and after a second it says with a
		monotonous, robotic voice, "Von Vitae's seal of approval has been given
		at %v".`,
	},
	{
		HasDuration: true,
		Raw: `/me starts bleeping, performing quick calculations and then
		blurts out, "I have been Vonproved™ precisely for %v".`,
	},
}

func (v Vonprove) Apply(user string) string {
	var vonproved = time.Date(2017, 9, 26, 0, 0, 0, 0, time.UTC)
	if v.HasDate {
		return fmt.Sprintf(clean(v.Raw), vonproved.Format("Monday, 02 Jan 2006"))
	}
	if v.HasDuration {
		return fmt.Sprintf(clean(v.Raw), time.Now().Sub(vonproved))
	}
	if v.HasUser {
		return fmt.Sprintf(clean(v.Raw), user)
	}
	return fmt.Sprintf(clean(v.Raw))
}

func RandVonprove(user string) string {
	v := vonproves[newRand(len(vonproves))]
	return v.Apply(user)
}

type stand struct {
	Name string
	Type string
	Desc string
}

func (st stand) Apply(user string) string {
	s := fmt.Sprintf("%s's new Stand is ", user)
	s += fmt.Sprintf("[u]%s[/u] ", clean(st.Name))
	s += fmt.Sprintf("([i]%s[/i]): ", clean(st.Type))
	s += fmt.Sprintf("%s", clean(st.Desc))
	return s
}

// Menacing BBCode:
// [color=purple][b][sub]ゴ[/sub]ゴ[sup]ゴゴ[/sup][i][sup]ゴ[/sup][/i][i]ゴ[sub]ゴゴ...[/sub][/i][/b][/color]

var stands = []stand{
	{
		Name: `Don't Stop Me Now`,
		Type: `Artificial Non-Humanoid Stand`,
		Desc: `The Stand can move at supersonic speeds for 3 seconds, but takes
		30 seconds to recharge. Moving fewer body parts prolongs the effect.`,
	},
	{
		Name: `Take A Chance On Me`,
		Type: `Long-Distance Manipulate Type`,
		Desc: `Makes any set of possible outcomes equally probable within a 10
		meter radius of the target area.`,
	},
	{
		Name: `They Might be Giants`,
		Type: `Close-Range Humanoid Stand`,
		Desc: `Physical contact on a target can shift its elemental makeup up
		or down 1 number on the Periodic table. The effect can be reversed and
		the target will no longer be affected from any future attempts.`,
	},
	{
		Name: `Cranberry`,
		Type: `Long-Distance Manipulate Type`,
		Desc: `Allows the user to be a voice inside someone's head, and in turn
		hear that person's thoughts.`,
	},
	{
		Name: `Pachelbel`,
		Type: `Close-Range Power Type`,
		Desc: `Can disable 1 sense (ie: Touch, Sight, Hearing, Taste, Smell) on
		a target, or disable 1 of the user's senses in exchange to boost the
		other senses dramatically.`,
	},
	{
		Name: `Juke Box Hero`,
		Type: `Close-Range Power Type`,
		Desc: `Produces music that can affect gravity within audible range.`,
	},
	{
		Name: `One Way or Another`,
		Type: `Range-Irrelevant Humanoid Stand`,
		Desc: `Copies the shape, power, and abilities of its target. However it
		will attack both its target and user, and will relentlessly pursue
		whoever is closest. Once its copied target is dead it returns to its
		harmless base form.`,
	},
	{
		Name: `Eurythmics`,
		Type: `Close-Range Humanoid Stand`,
		Desc: `Gets stronger for each nearby person that is asleep.`,
	},
	{
		Name: `Killing You Softly`,
		Type: `Artificial Humanoid Stand`,
		Desc: `The user can set a three-word-phrase during the day. Anyone
		within eye-sight who repeats the phrase suffers a potentially-fatal
		heart attack. The phrase resets come sunrise and cannot be used
		again.`,
	},
	{
		Name: `Everything You Know Is Wrong`,
		Type: `Artificial Non-Humanonid Stand`,
		Desc: `Reverses temperature interactions within 100 meters of the user.
		Ice burns, boiling water freezes, etc. The user is unaffected by these
		changes.`,
	},
	{
		Name: `Jimmy Buffet`,
		Type: `Range Irrelevant Artificial Stand`,
		Desc: `Turns photographs of cooked food and bottled drinks into real,
		3D objects.`,
	},
	{
		Name: `Licensed to Ill`,
		Type: `Close-Range Power Type`,
		Desc: `The Stand utilizes a different weapon each day of the week, but
		is an expert no matter which one.`,
	},
	{
		Name: `Berlin`,
		Type: `Automatic Type`,
		Desc: `Disables all Stands and associated powers within line-of-sight
		of the user. Stands return once they are out of sight.`,
	},
	{
		Name: `Hard-Boiled`,
		Type: `Sentient Stand`,
		Desc: `An immensely powerful hand-to-hand fighter, but will only follow
		commands/directions if the user narrates as if they are in a Noir
		detective film.`,
	},
	{
		Name: `Forever Your Girl`,
		Type: `Close-Range Humanoid Stand`,
		Desc: `Wields the power of fire when in darkness, and the power of ice
		when in light.`,
	},
	{
		Name: `Lady Soul`,
		Type: `Close-Range Manipulate Stand`,
		Desc: `Can transfer locks, whether physical or digital, onto adjacent
		objects.`,
	},
	{
		Name: `Colors of the Wind`,
		Type: `Phenomenon Stand`,
		Desc: `The user can climb into paintings and even pull in other people
		with them. The user cannot enter or exit a painting that is covered up.
		Any changes made within the panting can be seen by outside viewers.`,
	},
	{
		Name: `Springsteen`,
		Type: `Phenomenon Stand`,
		Desc: `The user can learn any skill instantly at the cost of forgetting
		another skill. The forgotten skill cannot be relearned for 72-hours.`,
	},
	{
		Name: `Billy Joel`,
		Type: `Close-Range Power Stand`,
		Desc: `The Stand is faster-yet-weaker if it is hot, and
		slower-yet-stronger if it is cold. It can never be so cold that it
		isn't able to move, nor so hot that it is completely harmless.`,
	},
	{
		Name: `Stone Temple`,
		Type: `Close-Range Manipulate Stand`,
		Desc: `On physical contact the Stand can disable half of anything that
		a person has a pair of (arms, eyes, lungs, etc), but only one at a
		time. Repeated contact is needed to disable more pairs. The active
		effects can be disabled by the user at any point.`,
	},
	{
		Name: `The Supremes`,
		Type: `Colony Stand`,
		Desc: `This Stand takes on the form of three individual-yet-identical
		Stands, always standing around the user. They reflect damage back
		against the attacker based on the attacker's confidence of winning. The
		more confident they are, then the damage is multiplied further.`,
	},
	{
		Name: `Neighborhood`,
		Type: `Long-Range Manipulate Stand`,
		Desc: `Changes its appearance to match its target's greatest fear. If
		there are multiple targets, it will combine all appearances into one
		form. Does not work on intangible fears, such as "being alone," unless
		there is something physical that the target associates with it.`,
	},
	{
		Name: `Simple Minds`,
		Type: `Phenomenon Stand`,
		Desc: `The user can forget about a required bodily function and survive
		without it (ie: forget about breathing and no longer need to breathe)
		so long as no one reminds them of it and causes them to remember.`,
	},
	{
		Name: `Tenacious D`,
		Type: `Close-Range Manipulate Stand`,
		Desc: `In exchange for money and a drop of blood of the requested
		target, the Stand can forge handwriting of the person whose blood was
		given. It also possesses enough knowledge to do your homework.`,
	},
	{
		Name: `Destiny's Child`,
		Type: `Range Irrelevant`,
		Desc: `Anyone within eyesight of the user, or can hear the user, cannot
		lie when asked a question by them. The Stand user may lie, but by doing
		so the Stand will then transfer to the person that they lied to. If
		they lied to a group, it transfers to the closest person.`,
	},
	{
		Name: `Def Leppard`,
		Type: `Automatic Stand`,
		Desc: `The Stand can disguise itself as food. If eaten, the victim's
		metabolism rapidly increases to a point that their body is visibly
		consuming itself for nourishment. The effects end if the Stand is
		removed by any method, or the user is killed.`,
	},
	{
		Name: `Yo-Yo Ma`,
		Type: `Close-Range Power Stand`,
		Desc: `The Stand fights with a yo-yo that is covered in spinning
		blades. Fancy tricks allow it to attack in unpredictable patterns or
		angles.`,
	},
	{
		Name: `Good Charlotte`,
		Type: `Close-Range Humanoid Stand`,
		Desc: `This Stand's right hand can heal wounds on contact. If the
		wounds are severe, the Stand must be supercharged by using its left
		hand to drain the life energy of another being. The user can also
		absorb a percentage of their target's wounds into themselves to ease
		their burden.`,
	},
	{
		Name: `Madonna`,
		Type: `Phenomenon Stand`,
		Desc: `The user's Hamon abilities are strengthened. They also no longer
		need to breathe to use their abilities, allowing their Hamon strikes to
		retain full power even when they are unable to breathe properly.`,
	},
	{
		Name: `Dark Side of the Moon`,
		Type: `Long-Distance Manipulate Type`,
		Desc: `The Stand can refract or focus a light into offensive energy
		attacks. Dispersed light can illuminate the area and singe skin, while
		focused light can burn through solid steel within seconds.`,
	},
	{
		Name: `Topsy Turvy`,
		Type: `Close-Range Humanoid Stand`,
		Desc: `Actions placed upon and imparted by the Stand or user have an
		unequal response in physics. A direct punch will feel like a light slap
		and barely move the person, but a gentle tap can send the target
		sailing as though they have been struck hard. Objects thrown by the
		Stand user retain these unequal properties until they come to a
		complete stop.`,
	},
	{
		Name: `Vitamin C`,
		Type: `Colony Stand`,
		Desc: `The user can toggle the ability to bring fruits and vegetables
		to life. The food grows two arms and two legs, and will obey any
		command given to them. They contain enough sentience to listen,
		remember, and interpret information. However, they only have the
		relative strength of whatever fruit or vegetable they're formed from,
		and they can only speak in puns relating to the food they resemble.
		They return to normal if eaten.`,
	},
	{
		Name: `Reznor`,
		Type: `Long-Distance Power Type`,
		Desc: `This Stand, resembling a large metallic bird of prey, is capable
		of diving at incredible speeds to slash and slice unsuspecting prey.
		Its feathers are actually individual daggers which can be wielded or
		thrown with rapid precision.`,
	},
	{
		Name: `Rocky Horror Picture Show`,
		Type: `Close-Range Power Type`,
		Desc: `The user is able to rewind time up to 10 seconds, however any
		physical effects caused up to that point will remain despite the
		reversal; wounds do not disappear and objects in motion will retain
		their momentum. No one else is aware of the ability when it's
		utilized.`,
	},
	{
		Name: `Once More With Feeling`,
		Type: `Phenomenon Stand`,
		Desc: `While activated, anyone within a mile of the user can only speak
		while singing. Everyone loves a musical episode.`,
	},
	{
		Name: `Antipode`,
		Type: `Close-Range Power Type`,
		Desc: `The Stand has hands of two different temperatures: one is always
		on fire and the other is always at absolute zero. If the Stand claps
		its hands together, it can generate an explosive shock wave which can
		be directed at any angle in front of itself.`,
	},
}

type animeOp struct {
	Raw string
}

func (j animeOp) Apply() string {
	return fmt.Sprintf(clean(j.Raw))
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

func RandJojo(user string) string {
	st := stands[newRand(len(stands))]
	return st.Apply(user)
}

func Loth(user string, loth *eribo.Loth, isNew bool, targets []*eribo.Player) string {
	switch {
	case loth == nil:
		msg := `Unable to find eligible target.`

		return fmt.Sprintf(clean(msg))
	case loth != nil && !isNew:
		msg := `Current 'lee of the hour is %s. Time left is %s.`

		return fmt.Sprintf(clean(msg), loth.Name, loth.TimeLeft())
	case loth != nil && isNew && user == loth.Name && len(targets) == 1:
		msg := `/me looks around the room while performing calculations and
		seeking potentials targets. After a few seconds it stops and stares at
		what seems to be the only eligible target. It grabs %s and injects them
		with a powerful serum which numbs their strength and reflexes but
		sharply increases their sensitivity. It leaves the victim half
		incapacitated on the floor then proceeds to announce to the whole room:
		"New 'lee of the hour is %s!"`

		return fmt.Sprintf(clean(msg), loth.Name, loth.Name)
	case loth != nil && isNew && user == loth.Name && len(targets) != 1:
		msg := `/me appears to be malfunctioning as it doesn't seem to be
		seeking for other targets and turns towards the person that issued the
		command. It grabs %s and injects them with the serum instead!`

		return fmt.Sprintf(clean(msg), loth.Name)
	case loth != nil && isNew && user != loth.Name:
		msg := `/me grabs %s and injects them with a powerful serum which numbs
		their strength and reflexes but sharply increases their sensitivity. It
		leaves the victim half incapacitated on the floor then proceeds to
		announce to the whole room: "New 'lee of the hour is %s!"`

		return fmt.Sprintf(clean(msg), loth.Name, loth.Name)
	default:
		return fmt.Sprintf("/me looks confused and doesn't do anything at all.")
	}
}
