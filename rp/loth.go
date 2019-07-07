package rp

import (
	"fmt"

	"github.com/kusubooru/eribo/eribo"
)

// LothTime returns a message indicating the remaining time for a loth.
func LothTime(loth *eribo.Loth) string {
	if loth == nil {
		return "There's no loth."
	}
	if loth.Expired() {
		return fmt.Sprintf("Time is up for %s. A new 'lee of the hour can be chosen!", loth.Name)
	}
	return fmt.Sprintf("Current 'lee of the hour is %s. Time left is %s.", loth.Name, loth.TimeLeft())
}

// LothWarning returns a warning message before the loth command proceeds.
func LothWarning() string {
	return "By using this command, you agree that you intend to play with the randomly chosen victim (assuming they are not afk). To continue, type: !loth confirm"
}

// Loth returns a different message depending on the different states of loth.
func Loth(user string, loth *eribo.Loth, isNew bool, targets []*eribo.Player) string {
	switch {
	case loth == nil:
		msg := `Unable to find eligible target.`

		return clean(msg)
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
