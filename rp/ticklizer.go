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
	{"underarms", true},
	{"ribs", true},
	{"sides", true},
	{"feet", true},
	{"tummy", false},
	{"chest", false},
	{"thighs", true},
	{"hips", true},
	{"genitals", true},
	{"butt", false},
}

func filterBodyParts(filter string) ([]bodyPart, bool) {
	parts := make([]bodyPart, 0)
	switch {
	case filter == "ub":
		for _, p := range bodyParts {
			if p.name != "feet" {
				parts = append(parts, p)
			}
		}
		return parts, true
	case filter != "":
		for _, p := range bodyParts {
			if filter == p.name {
				parts = append(parts, p)
			}
		}
		return parts, true
	default:
		return bodyParts, false
	}
}

func ticklizerFilters() []string {
	filters := make([]string, len(bodyParts)+1)
	for _, p := range bodyParts {
		filters = append(filters, p.name)
	}
	filters = append(filters, "ub")
	return filters
}

func InTicklizerFilters(arg string) bool {
	filters := ticklizerFilters()
	for _, f := range filters {
		if arg == f {
			return true
		}
	}
	return false
}

func ticklizer(name, owner, botName string, tcase ticklizerCase, filter string) string {
	factors := []int{2, 5, 10}
	intensity := factors[newRand(len(factors))]

	parts, hasFilter := filterBodyParts(filter)
	part := parts[newRand(len(parts))]
	itOrThem := "it"
	if part.plural {
		itOrThem = "them"
	}
	var format string
	filterMsg := ""
	switch tcase {
	case confused:
		format = `/me found more than one targets. It got confused and zapped
		%s instead with the ticklizer beam, hitting their [u]%s[/u] making, %s
		ten times more ticklish!`
		return fmt.Sprintf(clean(format), name, part.name, itOrThem)
	case notFound:
		format = `/me could not find its target. It got confused and zapped %s
		instead with the ticklizer beam, hitting their [u]%s[/u], making %s ten
		times more ticklish!`
		return fmt.Sprintf(clean(format), name, part.name, itOrThem)
	case forbidden:
		format = `/me is forbidden from hitting that target. It turns and zaps
		%s instead with the ticklizer beam, hitting their [u]%s[/u], making %s
		ten times more ticklish!`
		return fmt.Sprintf(clean(format), name, part.name, itOrThem)
	default:
		if hasFilter {
			filterMsg = "concentrates its aim and"
		}
		switch intensity {
		default:
			format = `/me %s fires two small rays of the ticklizer beam from
			the tips of its index fingers, zapping %s's [u]%s[/u], making %s
			two times more ticklish!`
		case 5:
			format = `/me %s fires the ticklizer beam from the palms of its
			hands, zapping %s's [u]%s[/u], making %s five times more ticklish!`
		case 10:
			format = `/me %s fires a large ray of the ticklizer beam from the
			center of its chest, zapping %s's [u]%s[/u], making %s ten times
			more ticklish!`
		}
	}
	if name == botName {
		return fmt.Sprintf(`/me refuses to hit itself and does nothing instead.`)
	}
	if name == owner {
		return fmt.Sprintf(`/me refuses to hit its creator. It kindly offers him a tomato instead.`)
	}
	return fmt.Sprintf(clean(format), filterMsg, name, part.name, itOrThem)
}

func Ticklizer(name, owner, botName, filter string) string {
	return ticklizer(name, owner, botName, normal, filter)
}

func TicklizerConfused(name, owner, botName, filter string) string {
	return ticklizer(name, owner, botName, confused, filter)
}

func TicklizerNotFound(name, owner, botName, filter string) string {
	return ticklizer(name, owner, botName, notFound, filter)
}

func TicklizerForbidden(name, owner, botName, filter string) string {
	return ticklizer(name, owner, botName, notFound, filter)
}
