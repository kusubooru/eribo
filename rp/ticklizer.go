package rp

import "fmt"

type ticklizerCase int

const (
	normal ticklizerCase = iota
	confused
	notFound
	forbidden
)

type bodyPart struct {
	name   string
	plural bool
}

var bodyParts = []bodyPart{
	{"armpits", true},
	{"feet", true},
	{"tummy", false},
	{"chest", false},
	{"thighs", true},
	{"genitals", true},
	{"butt", false},
}

func ticklizer(name, owner, botName string, tcase ticklizerCase) string {
	part := bodyParts[newRand(len(bodyParts))]
	itOrThem := "it"
	if part.plural {
		itOrThem = "them"
	}
	var format string
	switch tcase {
	case confused:
		format = `/me found more than one targets. It got confused and zapped
		%s instead with the ticklizer beam, hitting their [u]%s[/u] making, %s
		ten times more ticklish!`
	case notFound:
		format = `/me could not find its target. It got confused and zapped %s
		instead with the ticklizer beam, hitting their [u]%s[/u], making %s ten
		times more ticklish!`
	case forbidden:
		format = `/me is forbidden from hitting that target. It turns and zaps
		%s instead with the ticklizer beam, hitting their [u]%s[/u], making %s
		ten times more ticklish!`
	default:
		format = `/me fires the ticklizer beam from the palms of its hands,
		zapping %s's [u]%s[/u], making %s ten times more ticklish!`
	}
	if name == botName {
		return fmt.Sprintf(`/me refuses to hit itself and does nothing instead.`)
	}
	if name == owner {
		return fmt.Sprintf(`/me refuses to hit its creator. It kindly offers him a tomato instead.`)
	}
	return fmt.Sprintf(clean(format), name, part.name, itOrThem)
}

func Ticklizer(name, owner, botName string) string {
	return ticklizer(name, owner, botName, normal)
}

func TicklizerConfused(name, owner, botName string) string {
	return ticklizer(name, owner, botName, confused)
}

func TicklizerNotFound(name, owner, botName string) string {
	return ticklizer(name, owner, botName, notFound)
}

func TicklizerForbidden(name, owner, botName string) string {
	return ticklizer(name, owner, botName, notFound)
}
