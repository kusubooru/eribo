package rp

import "fmt"

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