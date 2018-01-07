package rp

import (
	"fmt"
	"time"
)

type Vonprove struct {
	Raw         string
	HasDate     bool
	HasDuration bool
	HasUser     bool
}

func (v Vonprove) Apply(user string) string {
	var vonproved = time.Date(2017, 9, 26, 0, 0, 0, 0, time.UTC)
	if v.HasDate {
		return fmt.Sprintf(clean(v.Raw), vonproved.Format("Monday, 02 Jan 2006"))
	}
	if v.HasDuration {
		return fmt.Sprintf(clean(v.Raw), time.Since(vonproved))
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
