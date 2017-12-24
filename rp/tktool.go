package rp

import (
	"fmt"
	"time"

	"github.com/kusubooru/eribo/loot"
)

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

// TODO(jin): Pinwheel, backscratcher

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
