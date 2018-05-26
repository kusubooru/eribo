package rp

import "fmt"

type tieupCase int

const (
	tieUpNormal tieupCase = iota
	tieUpConfused
	tieUpNotFound
)

func randTieUp(victim, owner, botName string, tieupCase tieupCase, filter string) string {

	switch tieupCase {
	case tieUpConfused:
		return fmt.Sprintf(`/me was unable to identify the correct target and takes no action.`)
	case tieUpNotFound:
		return fmt.Sprintf(`/me was unable to find the target and stays idle.`)
	default:
		if victim == botName {
			return fmt.Sprintf(`/me refuses to tie itself up and does nothing instead.`)
		}
		if victim == owner {
			return fmt.Sprintf(`/me refuses to tie its creator. It kindly offers him a tomato instead.`)
		}

		ties := filterTieUps(filter)
		tie := ties[newRand(len(ties))]
		return fmt.Sprintf(clean(tie.msg), victim)
	}
}

func InTieUpTags(filter string) bool {
	tags := tieUpTags()
	for _, t := range tags {
		if filter == t {
			return true
		}
	}
	return false
}

func tieUpTags() []string {
	m := make(map[string]struct{})
	for _, tie := range tieUps {
		for _, tag := range tie.tags {
			m[tag] = struct{}{}
		}
	}
	tags := make([]string, len(m))
	for k := range m {
		tags = append(tags, k)
	}
	return tags
}

func RandTieUp(victim, owner, botName, filter string) string {
	return randTieUp(victim, owner, botName, tieUpNormal, filter)
}

func RandTieUpConfused(name, owner, botName, filter string) string {
	return randTieUp(name, owner, botName, tieUpConfused, filter)
}

func RandTieUpNotFound(name, owner, botName, filter string) string {
	return randTieUp(name, owner, botName, tieUpNotFound, filter)
}

type tieUp struct {
	tags []string
	msg  string
}

func filterTieUps(tag string) []tieUp {
	if tag == "" {
		return tieUps
	}

	f := make([]tieUp, 0)
	for _, tie := range tieUps {
		for _, t := range tie.tags {
			if tag == t {
				f = append(f, tie)
			}
		}
	}
	return f
}

var tieUps = []tieUp{
	{
		tags: []string{"feet", "wrap"},
		msg: `/me grabs %s and wraps their body tightly using saran wrap,
		leaving only their [u]head[/u] and [u]feet[/u] exposed. Then places the
		wrapped body on the table and starts strapping it. It applies tight
		straps above and under their chest, on their waist, thighs, knees and
		ankles rendering the victim immobile.`,
	},
	{
		tags: []string{"feet", "underarms", "wrap"},
		msg: `/me lifts %s up by their arms holding them above their head and
		swiftly ties them together. Then it proceeds to wrap the rest of the
		victim's body up with saran wrap, leaving their [u]head[/u],
		[u]underarms[/u] and [u]feet[/u] vulnerable. Lastly, it places the
		victim's body on a rack and applies straps on their wrists, elbows,
		waist, thighs, knees and ankles rendering them immobile.`,
	},
	{
		tags: []string{"feet", "stocks"},
		msg: `/me grabs %s and forces their arms behind their back, locking
		their wrists in a pair of leather cuffs. Then sits them down, placing
		their ankles in the stocks and finally locking them up, leaving their
		[u]feet[/u] vulnerable.`,
	},
	{
		tags: []string{"ub", "sides", "legs"},
		msg: `/me bends %s forward into an awaiting standing pillory and shuts
		it on their neck and wrists. A spreader bar is then cuffed to their
		ankles, forcing their legs far apart. With their [u]sides[/u] and
		[u]legs[/u] rather vulnerable, they cannot kick effectively nor see
		behind them.`,
	},
	{
		tags: []string{"feet"},
		msg: `/me suddenly pushes %s backwards into an open coffin with two
		slots at the bottom. Their ankles get caught in the slots before the
		lid automatically slams shut and locks itself, leaving their
		[u]feet[/u] exposed on the outside.`,
	},
	{
		tags: []string{"ub", "tummy"},
		msg: `/me extends seven spider-like legs from behind and uses them to
		lift %s into the air. An eighth leg is revealed to be pulling silk
		webbing up from a large spool, and it quickly rolls them up, leaving
		only their [u]stomach[/u] and nose exposed. It then sticks their back
		against a wall so that they can't wiggle away.`,
	},
	{
		tags: []string{"sides", "feet"},
		msg: `/me puts on a cowboy hat and swings a lasso onto %s before
		drawing them in and wrestling them to the ground. The excess rope is
		wound around their ankles while their wrists are positioned behind
		them, and the slack is tightened, forcing them into a hogtie that
		renders their [u]sides[/u] and [u]feet[/u] quite vulnerable.`,
	},
	{
		tags: []string{"ub"},
		msg: `/me tightens a belt around %s that has a leather cuff dangling
		from either side. Then, it forces their hands down into the awaiting
		cuffs and buckles them shut, leaving their hands trapped by their
		waist.`,
	},
	{
		tags: []string{"legs", "feet"},
		msg: `/me deems %s is getting too unruly and takes measures to protect
		them. A straitjacket is pulled onto them and buckled shut, forcing
		their arms crossed in front of themselves. Although their upperbody is
		secure and protected, their [u]legs[/u] and [u]feet[/u] remain
		uncovered.`,
	},
	{
		tags: []string{"ub", "torso", "wrap"},
		msg: `/me manipulates %s into fully bending their legs before wrapping
		each one up in saran wrap. Their arms are fully bent, hands pressed
		against their own shoulders, before being wrapped up in the same
		fashion, rendering their [u]torso[/u] completely vulnerable.`,
	},
	{
		tags: []string{"feet"},
		msg: `/me shoves %s down against a wooden chair and cuffs their wrists
		together behind the back of it. Their ankles are bent underneath the
		seat and cuffed to the support stretcher that links the chair legs,
		exposing their [u]feet[/u] as well as their [u]upperbody[/u] above.`,
	},
	{
		tags: []string{"feet", "yoga"},
		msg: `/me invites %s to join them in a Yoga session. First up: Lotus
		position! Unfortunately for them, several ropes are wrapped around
		their calves and ankles, leaving their [u]soles[/u] upturned. Their
		arms are also crossed behind their back with their forearms bound in
		rope, parallel with each other, to align their chi and leave their
		[u]upperbody[/u] defenseless.`,
	},
	{
		tags: []string{"ub", "sides"},
		msg: `/me helps %s with their daily stretches. One arm is pulled over
		their shoulder while the other is pulled down and around by their lower
		back.  Their two wrists become joined by a pair of leather cuffs,
		leaving them in an awkward pose that exposes [u]one side[/u] and the
		[u]front of their upperbody[/u].`,
	},
	{
		tags: []string{"ub"},
		msg: `/me adds a flair of fashion to %s by placing a lovely set of
		leather bondage mittens over their hands. Once locked shut, their
		fingers cannot manipulate anything through the thick leather. To add to
		their distress, leather cuffs are attached to their ankles, with only
		six-inches of slack between their ankles for hobbling around.`,
	},
	{
		tags: []string{"feet"},
		msg: `/me suddenly triggers a rope noose trap underneath %s. Their
		ankles are snagged and hoisted up into the air, but only high enough to
		flip them onto their back. Both [u]feet[/u] are vulnerable in the air
		unless they are flexible enough to reach all the way up to them.`,
	},
	{
		tags: []string{"feet", "wrap", "flute"},
		msg: `/me plays the flute, causing a pile of bandages to rise up like
		snakes.  They quickly slither out and wrap around %s from their ankles
		to their shoulders, leaving their [u]head[/u] and [u]feet[/u] exposed.
		Once complete, the bandages stretch up to the ceiling and pull their
		wrapped prey with them, leaving them dangling a few inches off the
		floor.`,
	},
	{
		tags: []string{"feet"},
		msg: `/me knocks %s over and sets a chair down on top of their body.
		Their wrists are wrapped up in rope against the chair legs, and their
		ankles are hoisted up and tied to the headrest of the chair. Their
		[u]feet[/u] are vulnerable to all, while their [u]head[/u] is exposed
		and forced to watch anyone who sits down above them.`,
	},
	{
		tags: []string{"feet", "stocks"},
		msg: `/me quickly fastens two open-stocks to dangle from the ceiling,
		padded for maximum comfort. It lifts %s into the air, their neck and
		wrists placed in one and ankles in the other, before closing them with
		loud smacks, leaving their [u]entire body[/u] exposed and accessible.
		It whirs, performing some additional calculations, then it produces
		some precisely-measured silken cords and uses them to tie their toes to
		the top of the stocks.`,
	},
	{
		tags: []string{"feet", "slinky"},
		msg: `/me starts producing very rapidly a colorful, spiral-shaped,
		plastic material which forms a giant slinky and then launches it
		upwards, in an arc, perfectly calculated to fall right on %s, trapping
		them in the middle and leaving only their [u]head[/u] and [u]feet[/u]
		exposed. The victim quickly finds out that while they are able to twist
		and wriggle, they are totally unable to escape on their own.`,
	},
	{
		tags: []string{"thighs", "ub", "matrix", "goo"},
		msg: `/me unveils two goo guns and does a Matrix dive in slow motion,
		firing at %s. Goo splatters against their wrists, pinning them against
		the wall with their arms outstretched to either-side. A series of
		rapid-fire shots coats the victim's calves against the wall as well.
		Their upperbody is extremely vulnerable along with their thighs.`,
	},
	{
		tags: []string{"feet", "scarecrow"},
		msg: `/me is ready to start a farm, but the most important part is
		missing: a Scarecrow! It grabs %s and hoists them up to a cross-shaped
		piece of wood, tying them against it with rope around their arms and
		waist. It also pulls their feet back and crosses them behind the wooden
		post before tying them in place with rope too.`,
	},
	{
		tags: []string{"ub", "cannon"},
		msg: `/me pushes in a large cannon and lights the fuse. KA-BOOM! It
		fires a wide net at %s which wraps around them several times,
		covering them from their shoulders to their knees and pinning their
		arms to their sides. The holes in the net are large enough for anyone
		to stick their hands through.`,
	},
	{
		tags: []string{"feet", "surgeon"},
		msg: `/me rushes over in surgeon scrubs. It slaps plaster all over %s's
		arms and legs before going over the new casts with a blowdryer until
		they harden. Their upperbody is still quite exposed as well as their
		feet. They also gets a $35,500 medical bill.`,
	},
	{
		tags: []string{"feet", "wizard"},
		msg: `/me puts on their robe and wizard hat. For the first trick, they
		pull a rabbit out of their hat! For the second trick, they push %s into
		a box and lock it shut, only their head and feet sticking out. A
		chainsaw is produced and used to saw down the middle of the box. Oh my
		god, there's blood everywhere! Just kidding. It passes through cleanly
		and the box is separated. The victim's body is now split in two. Their
		feet are utterly helpless.`,
	},
	{
		tags: []string{"ub", "hamster"},
		msg: `/me squeezes %s into a hamster ball that's far too small for
		them, forcing them to curl up into a ball. Worse yet, it latches shut
		from the outside! Don't worry, the ball has plenty of air holes that
		are big enough for other people to reach into.`,
	},
}
