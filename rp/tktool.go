package rp

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/kusubooru/eribo/loot"
)

type Tktool struct {
	Name    string
	Colors  []Color
	Quality Quality
	Emote   *template.Template
	weight  int
}

func (t Tktool) NameBBCode() string {
	return qualityColorBBCode(t.Quality, t.Name)
}

func (t Tktool) Apply(user string) (string, error) {
	data := struct {
		Tool  string
		Color Color
		User  string
	}{
		Tool: t.NameBBCode(),
		User: user,
	}
	if len(t.Colors) != 0 {
		data.Color = t.Colors[newRand(len(t.Colors))]
	}
	var buf bytes.Buffer
	if err := t.Emote.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("applying %v to %q: %v", data, t.Emote.Root, err)
	}
	return clean(buf.String()), nil
}

func Tktools() []Tktool {
	return tktools
}

func (t Tktool) Weight() int {
	if t.weight == 0 {
		return t.Quality.Weight()
	}
	return t.weight * t.Quality.Weight()
}

func RandTktool(name string) (string, error) {
	table := &loot.Table{}
	for _, t := range tktools {
		table.Add(t, t.Weight())
	}
	_, roll := table.Roll(time.Now().UnixNano())
	tool, ok := roll.(Tktool)
	if !ok {
		tool = tktools[newRand(len(tktools))]
	}
	return tool.Apply(name)
}

var tktools = []Tktool{
	{
		Name:    "[Feather of Sensitivity]",
		Quality: Uncommon,
		weight:  8,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. A magic red feather
		that makes its victim more sensitive as its used.`),
	},
	{
		Name:    "[Ravaged Goose Feather]",
		Quality: Poor,
		weight:  10,
		Emote: tmplMust(`/me hands {{.User}} an old {{.Tool}}. It's too wrecked
		to be usable for any meaningful purpose.`),
	},
	{
		Name:    "[Goose Feather]",
		Quality: Common,
		Colors:  []Color{Gray, White, Black},
		weight:  8,
		Emote:   tmplMust(`/me hands {{.User}} a stiff, {{.Color}} {{.Tool}}.`),
	},
	{
		Name:    "[Pristine Goose Feather]",
		Quality: Uncommon,
		Colors:  []Color{Gray, White, Black},
		weight:  10,
		Emote: tmplMust(`/me hands {{.User}} a long, {{.Color}} {{.Tool}} with
		a pointy tip.`),
	},
	{
		Name:    "[Ruined Ostrich Feather]",
		Quality: Poor,
		weight:  10,
		Emote:   tmplMust(`/me hands {{.User}} the remains of an old {{.Tool}}.`),
	},
	{
		Name:    "[Ostrich Feather]",
		Quality: Common,
		Colors:  []Color{Black, White, Red, Orange, Blue, Turquoise, Brown, Yellow, Fuchsia, Pink, Purple, Violet, Green},
		weight:  7,
		Emote:   tmplMust(`/me hands {{.User}} a large, {{.Color}} {{.Tool}}.`),
	},
	{
		Name:    "[Jaunty Ostrich Feather]",
		Quality: Uncommon,
		Colors:  []Color{Black, White, Red, Orange, Blue, Turquoise, Brown, Yellow, Fuchsia, Pink, Purple, Violet, Green},
		weight:  10,
		Emote: tmplMust(`/me hands {{.User}} an enormous, {{.Color}} {{.Tool}}
		which forms a slight curve at the top. Its shaft at the bottom, ends
		into a sharp quill.`),
	},
	{
		Name:    "[Destroyed Feather Boa]",
		Quality: Poor,
		weight:  3,
		Emote: tmplMust(`/me hands {{.User}} an old {{.Tool}}. The remains of
		the wrecked item don't look usable for any meaningful purpose.`),
	},
	{
		Name:    "[Feather Boa]",
		Quality: Common,
		Colors:  []Color{Pink, Fuchsia, Purple, Black, Emerald, Red, Yellow, Blue},
		weight:  3,
		Emote:   tmplMust(`/me hands {{.User}} a long, {{.Color}} {{.Tool}}.`),
	},
	{
		Name:    "[Chandelle Feather Boa]",
		Quality: Uncommon,
		Colors:  []Color{Pink, Fuchsia, Purple, Black, Emerald, Red, Yellow, Blue},
		weight:  3,
		Emote: tmplMust(`/me hands {{.User}} a fluffy, long, {{.Color}}
		{{.Tool}}. With the slightest movement, its plumes animate entrancingly.`),
	},
	{
		Name:    "[Inoperable Electric Flosser]",
		Quality: Poor,
		weight:  6,
		Emote: tmplMust(`/me hands {{.User}} an {{.Tool}}. It doesn't look
		functional anymore and the tip is missing.`),
	},
	{
		Name:    "[Electric Flosser]",
		Quality: Common,
		weight:  6,
		Emote:   tmplMust(`/me hands {{.User}} an {{.Tool}}.`),
	},
	{
		Name:    "[Aqua-colored Electric Flosser]",
		Quality: Uncommon,
		weight:  6,
		Emote: tmplMust(`/me hands {{.User}} an {{.Tool}}. It is equipped with
		a fully charged battery and a flexible, nylon tip.`),
	},
	{
		Name:    "[Busted Electric Toothbrush]",
		Quality: Poor,
		weight:  8,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. It looks broken
		beyond repair and the brush is destroyed.`),
	},
	{
		Name:    "[Electric Toothbrush]",
		Quality: Common,
		weight:  8,
		Emote:   tmplMust(`/me hands {{.User}} an {{.Tool}}.`),
	},
	{
		Name:    "[Happy Electric Toothbrush]",
		Quality: Uncommon,
		weight:  8,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. The brush is
		round-shaped. Its body is light-blue and contains lots of colorful
		smiley faces.`),
	},
	{
		Name:    "[Snapped Paintbrush]",
		Quality: Poor,
		weight:  10,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. Beyond its broken
		body, there seem to be no bristles left on its tip.`),
	},
	{
		Name:    "[Small Paintbrush]",
		Quality: Common,
		weight:  10,
		Emote:   tmplMust(`/me hands {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    "[Eastern Paintbrush]",
		Quality: Uncommon,
		weight:  10,
		Emote: tmplMust(`/me hands {{.User}} a small, brown {{.Tool}} with soft
		bristles and a pointy tip. On its black, wooden body, the characters
		搔癢折磨 are inscribed in crimson red.`),
	},
	{
		Name:    "[Snapped Feather Duster]",
		Quality: Poor,
		Colors:  []Color{Gray, White, Black, Brown},
		weight:  10,
		Emote: tmplMust(`/me hands {{.User}} a wrecked, {{.Color}} {{.Tool}}. The handle
		is broken and the top part barely resembles a duster`),
	},
	{
		Name:    "[Feather Duster]",
		Quality: Common,
		Colors:  []Color{Gray, White, Black, Brown},
		weight:  10,
		Emote:   tmplMust(`/me hands {{.User}} a clean, {{.Color}} {{.Tool}}.`),
	},
	{
		Name:    "[Impeccable Feather Duster]",
		Quality: Uncommon,
		Colors:  []Color{Gray, White, Black, Brown},
		weight:  10,
		Emote: tmplMust(`/me hands {{.User}} an {{.Tool}} which looks like a
		matching accessory for a maid uniform. Its long, {{.Color}}, ostrich
		feathers look very soft and delicate.`),
	},
	{
		Name:    "[Destroyed Feather Gloves]",
		Quality: Poor,
		Colors:  []Color{Brown, Black, Violet, Purple, Red},
		weight:  5,
		Emote: tmplMust(`/me hands {{.User}} a pair of {{.Color}} {{.Tool}} The
		majority of the feathers that were previously attached on each
		fingertip seem to be missing and the ones that remain are totally
		ruined.`),
	},
	{
		Name:    "[Feather Gloves]",
		Quality: Common,
		Colors:  []Color{Brown, Black, Violet, Purple, Red},
		weight:  5,
		Emote: tmplMust(`/me hands {{.User}} a pair of {{.Color}} {{.Tool}}. On
		each fingertip there's a feather attached.`),
	},
	{
		Name:    "[Unblemished Feather Gloves]",
		Quality: Uncommon,
		Colors:  []Color{Brown, Black, Violet, Purple, Red},
		weight:  5,
		Emote: tmplMust(`/me hands {{.User}} a pair of expensive, {{.Color}}
		{{.Tool}}. They are made out of high quality leather and on each
		fingertip there's a long, pristine feather attached.`),
	},
	{
		Name:    "[Destroyed Hitachi Magic Wand]",
		Quality: Poor,
		weight:  5,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. The device is missing
		its cord and it looks broken beyond repair.`),
	},
	{
		Name:    "[Hitachi Magic Wand]",
		Quality: Common,
		weight:  5,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}} electrical
		massager.`),
	},
	{
		Name:    "[Modified Hitachi Magic Wand]",
		Quality: Uncommon,
		weight:  5,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. This model seems to
		be cordless and its switch seems to be altered. Apart from the
		traditional, O, I and II power levels, this switch supports two extra
		levels indicated as III and XXX.`),
	},
	{
		Name:    "[Dented Cat Claws]",
		Quality: Poor,
		weight:  3,
		Emote: tmplMust(`/me hands {{.User}} a set of {{.Tool}}. Each piece is
		so horribly dented that is impossible to wear.`),
	},
	{
		Name:    "[Metallic Cat Claws]",
		Quality: Common,
		weight:  4,
		Emote:   tmplMust(`/me hands {{.User}} a set of wearable {{.Tool}}.`),
	},
	{
		Name:    "[Silver Cat Claws]",
		Quality: Uncommon,
		weight:  3,
		Emote: tmplMust(`/me hands {{.User}} a set of wearable, well-crafted
		{{.Tool}}. The pointy tips of the claws seem to be sharp enough for
		play but not harm.`),
	},
	{
		Name:    "[Vial of Fish Oil]",
		Quality: Poor,
		weight:  7,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. The content looks dubious at
		best. Better not open it!`),
	},
	{
		Name:    "[Bottle of Baby Oil]",
		Quality: Common,
		weight:  7,
		Emote:   tmplMust(`/me hands {{.User}} a small {{.Tool}}.`),
	},
	{
		Name:    "[Bottle of Pure Lavender Oil]",
		Quality: Uncommon,
		weight:  7,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. The label of the
		bottle reads: 'A wonderful fragrance, for silky soft skin and soothing
		massages.'`),
	},
	{
		Name:    "[Ruined Grooming Brush]",
		Quality: Poor,
		Colors:  []Color{Black, Blue, Violet, Pink, Orange, Purple, Brown},
		weight:  8,
		Emote: tmplMust(`/me hands {{.User}} an old, {{.Color}} {{.Tool}}. Most
		of the bristles are missing.`),
	},
	{
		Name:    "[Grooming Brush]",
		Quality: Common,
		Colors:  []Color{Black, Blue, Violet, Pink, Orange, Purple, Brown},
		weight:  8,
		Emote:   tmplMust(`/me hands {{.User}} a small, {{.Color}} {{.Tool}}.`),
	},
	{
		Name:    "[Porcupine Grooming Brush]",
		Quality: Uncommon,
		Colors:  []Color{Black, Blue, Violet, Pink, Orange, Purple, Brown},
		weight:  8,
		Emote: tmplMust(`/me hands {{.User}} a large, {{.Color}} {{.Tool}} with
		dense, black nylon bristles that form a slight curve.`),
	},
	{
		Name:    "[Toothless Afro Pick]",
		Quality: Poor,
		Colors:  []Color{Black, Red, Purple, Blue},
		weight:  7,
		Emote: tmplMust(`/me hands {{.User}} an old, {{.Color}} {{.Tool}}. All
		the teeth are missing.`),
	},
	{
		Name:    "[Afro Pick]",
		Quality: Common,
		Colors:  []Color{Black, Red, Purple, Blue},
		weight:  7,
		Emote: tmplMust(`/me hands {{.User}} a plastic, {{.Color}}
		{{.Tool}}.`),
	},
	{
		Name:    "[Enchanted Afro Pick]",
		Quality: Uncommon,
		Colors:  []Color{Black, Red, Purple, Blue},
		weight:  7,
		Emote: tmplMust(`/me hands {{.User}} a {{.Color}} {{.Tool}}. Its loose,
		thick teeth are endlessly twitching at seemingly random directions
		allowing it to walk if let free on the ground.`),
	},
	{
		Name:    "[Wooden Backscratcher]",
		Quality: Common,
		weight:  8,
		Emote:   tmplMust(`/me hands {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    "[Bear Claw Backscratcher]",
		Quality: Common,
		weight:  6,
		Emote:   tmplMust(`/me hands {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    "[Battery Operated Backscratcher]",
		Quality: Uncommon,
		weight:  8,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. It ends on a claw
		which resembles a feminine hand with long, red polished nails and extra
		rubber tips on its palm which rotate at a touch of the switch.`),
	},
	{
		Name:    "[Wartenberg Wheel]",
		Quality: Common,
		weight:  4,
		Emote:   tmplMust(`/me hands {{.User}} a {{.Tool}} also known as a pinwheel.`),
	},
	{

		Name:    `[Blue Ballpoint Pen]`,
		Quality: Common,
		weight:  4,
		Emote:   tmplMust(`/me hands {{.User}} a {{.Tool}}. It's filled with blue ink.`),
	},
	{

		Name:    `[Red Ballpoint Pen]`,
		Quality: Common,
		weight:  3,
		Emote:   tmplMust(`/me hands {{.User}} a {{.Tool}}. It's filled with red ink.`),
	},
	{
		Name:    `[Pink Feather Gel Pen]`,
		Quality: Uncommon,
		weight:  6,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. It's filled with
		pink, glittery ink and a fluffy feather on the end.`)},
	{
		Name:    `[Rainbow Gel Pen]`,
		Quality: Rare,
		weight:  6,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. It's filled with
		several varieties of colors, each one activated by color-coded slide
		switches on the side.`),
	},
	{
		Name:    `[Fallen Angel Feather Duster]`,
		Quality: Epic,
		weight:  2,
		Emote: tmplMust(`/me cautiously hands {{.User}} the {{.Tool}}. Its
		finely-crafted red handle ends in a brilliant display of feathers as
		black as a moonless midnight. The feathers flutter and twist on their
		own accord as if seeking out sensitive skin to touch.`),
	},
	{
		Name:    `[Scrubby Bath Gloves]`,
		Quality: Uncommon,
		weight:  6,
		Emote: tmplMust(`/me hands {{.User}} pair of {{.Tool}} that leave a
		teasing, tingling sensation wherever they scrub.`),
	},
}
