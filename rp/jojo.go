package rp

import (
	"fmt"
	"strings"
)

type stand struct {
	Name string
	Type string
	Desc string
}

func (st stand) apply(user string) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s's new Stand is ", user)
	fmt.Fprintf(&b, "[u]%s[/u] ", clean(st.Name))
	fmt.Fprintf(&b, "([i]%s[/i]): ", clean(st.Type))
	fmt.Fprintf(&b, "%s", clean(st.Desc))
	return b.String()
}

// RandJojo returns a random stand message.
func RandJojo(user string) string {
	st := stands[newRand(len(stands))]
	return st.apply(user)
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
