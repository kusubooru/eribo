package rp

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/kusubooru/eribo/loot"
)

type Tktool struct {
	name    string
	Colors  []Color
	Quality Quality
	Emote   *template.Template
	weight  int
}

func (t Tktool) Name() string {
	return t.name
}

func (t Tktool) NameBBCode() string {
	return qualityColorBBCode(t.Quality, t.Name())
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

type TktoolsLootTable struct {
	*loot.Table
}

func NewTktoolsLootTable() *TktoolsLootTable {
	table := &loot.Table{}
	for _, t := range tktools {
		table.Add(t, t.Weight())
	}
	return &TktoolsLootTable{Table: table}
}

func (t *TktoolsLootTable) Legendaries() int {
	legos := 0
	drops := t.Drops()
	t.RLock()
	defer t.RUnlock()
	for _, d := range drops {
		if d.Item == nil {
			continue
		}
		tietool, ok := d.Item.(Tktool)
		if !ok {
			continue
		}
		if tietool.Quality == Legendary && d.Weight > 0 {
			legos++
		}
	}
	return legos
}

func (t *TktoolsLootTable) RandTktoolDecreaseWeight(user string) (string, error) {
	legos := t.Legendaries()
	if legos == 0 {
		t = NewTktoolsLootTable()
	}
	seed := time.Now().UnixNano()
	_, roll := t.RollDecreaseWeight(seed)
	if roll == nil {
		return "", fmt.Errorf("tktool loot table returned nothing")
	}
	tool, ok := roll.(Tktool)
	if !ok {
		return "", fmt.Errorf("TktoolsLootTable contains an item that is not a Tktool")
	}
	return tool.Apply(user)
}

func RandTktool(name string) (string, error) {
	table := NewTktoolsLootTable()
	_, roll := table.Roll(time.Now().UnixNano())
	tool, ok := roll.(Tktool)
	if !ok {
		tool = tktools[newRand(len(tktools))]
	}
	return tool.Apply(name)
}

var tktools = []Tktool{
	{
		name:    "[Feather of Sensitivity]",
		Quality: Uncommon,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. A magic red feather
		that makes its victim more sensitive as its used.`),
	},
	{
		name:    "[Ravaged Goose Feather]",
		Quality: Poor,
		Emote: tmplMust(`/me hands {{.User}} an old {{.Tool}}. It's too wrecked
		to be usable for any meaningful purpose.`),
	},
	{
		name:    "[Goose Feather]",
		Quality: Common,
		Colors:  []Color{Gray, White, Black},
		Emote:   tmplMust(`/me hands {{.User}} a stiff, {{.Color}} {{.Tool}}.`),
	},
	{
		name:    "[Pristine Goose Feather]",
		Quality: Uncommon,
		Colors:  []Color{Gray, White, Black},
		Emote: tmplMust(`/me hands {{.User}} a long, {{.Color}} {{.Tool}} with
		a pointy tip.`),
	},
	{
		name:    "[Ruined Ostrich Feather]",
		Quality: Poor,
		Emote:   tmplMust(`/me hands {{.User}} the remains of an old {{.Tool}}.`),
	},
	{
		name:    "[Ostrich Feather]",
		Quality: Common,
		Colors:  []Color{Black, White, Red, Orange, Blue, Turquoise, Brown, Yellow, Fuchsia, Pink, Purple, Violet, Green},
		Emote:   tmplMust(`/me hands {{.User}} a large, {{.Color}} {{.Tool}}.`),
	},
	{
		name:    "[Jaunty Ostrich Feather]",
		Quality: Uncommon,
		Colors:  []Color{Black, White, Red, Orange, Blue, Turquoise, Brown, Yellow, Fuchsia, Pink, Purple, Violet, Green},
		Emote: tmplMust(`/me hands {{.User}} an enormous, {{.Color}} {{.Tool}}
		which forms a slight curve at the top. Its shaft at the bottom, ends
		into a sharp quill.`),
	},
	{
		name:    "[Destroyed Feather Boa]",
		Quality: Poor,
		Emote: tmplMust(`/me hands {{.User}} an old {{.Tool}}. The remains of
		the wrecked item don't look usable for any meaningful purpose.`),
	},
	{
		name:    "[Feather Boa]",
		Quality: Common,
		Colors:  []Color{Pink, Fuchsia, Purple, Black, Emerald, Red, Yellow, Blue},
		Emote:   tmplMust(`/me hands {{.User}} a long, {{.Color}} {{.Tool}}.`),
	},
	{
		name:    "[Chandelle Feather Boa]",
		Quality: Uncommon,
		Colors:  []Color{Pink, Fuchsia, Purple, Black, Emerald, Red, Yellow, Blue},
		Emote: tmplMust(`/me hands {{.User}} a fluffy, long, {{.Color}}
		{{.Tool}}. With the slightest movement, its plumes animate entrancingly.`),
	},
	{
		name:    "[Inoperable Electric Flosser]",
		Quality: Poor,
		Emote: tmplMust(`/me hands {{.User}} an {{.Tool}}. It doesn't look
		functional anymore and the tip is missing.`),
	},
	{
		name:    "[Electric Flosser]",
		Quality: Common,
		Emote:   tmplMust(`/me hands {{.User}} an {{.Tool}}.`),
	},
	{
		name:    "[Aqua-colored Electric Flosser]",
		Quality: Uncommon,
		Emote: tmplMust(`/me hands {{.User}} an {{.Tool}}. It is equipped with
		a fully charged battery and a flexible, nylon tip.`),
	},
	{
		name:    "[Busted Electric Toothbrush]",
		Quality: Poor,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. It looks broken
		beyond repair and the brush is destroyed.`),
	},
	{
		name:    "[Electric Toothbrush]",
		Quality: Common,
		Emote:   tmplMust(`/me hands {{.User}} an {{.Tool}}.`),
	},
	{
		name:    "[Happy Electric Toothbrush]",
		Quality: Uncommon,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. The brush is
		round-shaped. Its body is light-blue and contains lots of colorful
		smiley faces.`),
	},
	{
		name:    "[Snapped Paintbrush]",
		Quality: Poor,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. Beyond its broken
		body, there seem to be no bristles left on its tip.`),
	},
	{
		name:    "[Small Paintbrush]",
		Quality: Common,
		Emote:   tmplMust(`/me hands {{.User}} a {{.Tool}}.`),
	},
	{
		name:    "[Eastern Paintbrush]",
		Quality: Uncommon,
		Emote: tmplMust(`/me hands {{.User}} a small, brown {{.Tool}} with soft
		bristles and a pointy tip. On its black, wooden body, the characters
		搔癢折磨 are inscribed in crimson red.`),
	},
	{
		name:    "[Snapped Feather Duster]",
		Quality: Poor,
		Colors:  []Color{Gray, White, Black, Brown},
		Emote: tmplMust(`/me hands {{.User}} a wrecked, {{.Color}} {{.Tool}}. The handle
		is broken and the top part barely resembles a duster`),
	},
	{
		name:    "[Feather Duster]",
		Quality: Common,
		Colors:  []Color{Gray, White, Black, Brown},
		Emote:   tmplMust(`/me hands {{.User}} a clean, {{.Color}} {{.Tool}}.`),
	},
	{
		name:    "[Impeccable Feather Duster]",
		Quality: Uncommon,
		Colors:  []Color{Gray, White, Black, Brown},
		Emote: tmplMust(`/me hands {{.User}} an {{.Tool}} which looks like a
		matching accessory for a maid uniform. Its long, {{.Color}}, ostrich
		feathers look very soft and delicate.`),
	},
	{
		name:    "[Destroyed Feather Gloves]",
		Quality: Poor,
		Colors:  []Color{Brown, Black, Violet, Purple, Red},
		Emote: tmplMust(`/me hands {{.User}} a pair of {{.Color}} {{.Tool}} The
		majority of the feathers that were previously attached on each
		fingertip seem to be missing and the ones that remain are totally
		ruined.`),
	},
	{
		name:    "[Feather Gloves]",
		Quality: Common,
		Colors:  []Color{Brown, Black, Violet, Purple, Red},
		Emote: tmplMust(`/me hands {{.User}} a pair of {{.Color}} {{.Tool}}. On
		each fingertip there's a feather attached.`),
	},
	{
		name:    "[Unblemished Feather Gloves]",
		Quality: Uncommon,
		Colors:  []Color{Brown, Black, Violet, Purple, Red},
		Emote: tmplMust(`/me hands {{.User}} a pair of expensive, {{.Color}}
		{{.Tool}}. They are made out of high quality leather and on each
		fingertip there's a long, pristine feather attached.`),
	},
	{
		name:    "[Destroyed Hitachi Magic Wand]",
		Quality: Poor,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. The device is missing
		its cord and it looks broken beyond repair.`),
	},
	{
		name:    "[Hitachi Magic Wand]",
		Quality: Common,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}} electrical
		massager.`),
	},
	{
		name:    "[Modified Hitachi Magic Wand]",
		Quality: Uncommon,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. This model seems to
		be cordless and its switch seems to be altered. Apart from the
		traditional, O, I and II power levels, this switch supports two extra
		levels indicated as III and XXX.`),
	},
	{
		name:    "[Dented Cat Claws]",
		Quality: Poor,
		Emote: tmplMust(`/me hands {{.User}} a set of {{.Tool}}. Each piece is
		so horribly dented that is impossible to wear.`),
	},
	{
		name:    "[Metallic Cat Claws]",
		Quality: Common,
		Emote:   tmplMust(`/me hands {{.User}} a set of wearable {{.Tool}}.`),
	},
	{
		name:    "[Silver Cat Claws]",
		Quality: Uncommon,
		Emote: tmplMust(`/me hands {{.User}} a set of wearable, well-crafted
		{{.Tool}}. The pointy tips of the claws seem to be sharp enough for
		play but not harm.`),
	},
	{
		name:    "[Vial of Fish Oil]",
		Quality: Poor,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. The content looks dubious at
		best. Better not open it!`),
	},
	{
		name:    "[Bottle of Baby Oil]",
		Quality: Common,
		Emote:   tmplMust(`/me hands {{.User}} a small {{.Tool}}.`),
	},
	{
		name:    "[Bottle of Pure Lavender Oil]",
		Quality: Uncommon,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. The label of the
		bottle reads: 'A wonderful fragrance, for silky soft skin and soothing
		massages.'`),
	},
	{
		name:    "[Ruined Grooming Brush]",
		Quality: Poor,
		Colors:  []Color{Black, Blue, Violet, Pink, Orange, Purple, Brown},
		Emote: tmplMust(`/me hands {{.User}} an old, {{.Color}} {{.Tool}}. Most
		of the bristles are missing.`),
	},
	{
		name:    "[Grooming Brush]",
		Quality: Common,
		Colors:  []Color{Black, Blue, Violet, Pink, Orange, Purple, Brown},
		Emote:   tmplMust(`/me hands {{.User}} a small, {{.Color}} {{.Tool}}.`),
	},
	{
		name:    "[Porcupine Grooming Brush]",
		Quality: Uncommon,
		Colors:  []Color{Black, Blue, Violet, Pink, Orange, Purple, Brown},
		Emote: tmplMust(`/me hands {{.User}} a large, {{.Color}} {{.Tool}} with
		dense, black nylon bristles that form a slight curve.`),
	},
	{
		name:    "[Toothless Afro Pick]",
		Quality: Poor,
		Colors:  []Color{Black, Red, Purple, Blue},
		Emote: tmplMust(`/me hands {{.User}} an old, {{.Color}} {{.Tool}}. All
		the teeth are missing.`),
	},
	{
		name:    "[Afro Pick]",
		Quality: Common,
		Colors:  []Color{Black, Red, Purple, Blue},
		Emote: tmplMust(`/me hands {{.User}} a plastic, {{.Color}}
		{{.Tool}}.`),
	},
	{
		name:    "[Enchanted Afro Pick]",
		Quality: Uncommon,
		Colors:  []Color{Black, Red, Purple, Blue},
		Emote: tmplMust(`/me hands {{.User}} a {{.Color}} {{.Tool}}. Its loose,
		thick teeth are endlessly twitching at seemingly random directions
		allowing it to walk if let free on the ground.`),
	},
	{
		name:    "[Wooden Backscratcher]",
		Quality: Common,
		Emote:   tmplMust(`/me hands {{.User}} a {{.Tool}}.`),
	},
	{
		name:    "[Bear Claw Backscratcher]",
		Quality: Common,
		Emote:   tmplMust(`/me hands {{.User}} a {{.Tool}}.`),
	},
	{
		name:    "[Battery Operated Backscratcher]",
		Quality: Uncommon,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. It ends on a claw
		which resembles a feminine hand with long, red polished nails and extra
		rubber tips on its palm which rotate at a touch of the switch.`),
	},
	{
		name:    "[Wartenberg Wheel]",
		Quality: Common,
		Emote:   tmplMust(`/me hands {{.User}} a {{.Tool}} also known as a pinwheel.`),
	},
	{

		name:    `[Blue Ballpoint Pen]`,
		Quality: Common,
		Emote:   tmplMust(`/me hands {{.User}} a {{.Tool}}. It's filled with blue ink.`),
	},
	{

		name:    `[Red Ballpoint Pen]`,
		Quality: Common,
		Emote:   tmplMust(`/me hands {{.User}} a {{.Tool}}. It's filled with red ink.`),
	},
	{
		name:    `[Pink Feather Gel Pen]`,
		Quality: Uncommon,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. It's filled with
		pink, glittery ink and a fluffy feather on the end.`)},
	{
		name:    `[Rainbow Gel Pen]`,
		Quality: Rare,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. It's filled with
		several varieties of colors, each one activated by color-coded slide
		switches on the side.`),
	},
	{
		name:    `[Fallen Angel Feather Duster]`,
		Quality: Epic,
		Emote: tmplMust(`/me cautiously hands {{.User}} the {{.Tool}}. Its
		finely-crafted red handle ends in a brilliant display of feathers as
		black as a moonless midnight. The feathers flutter and twist on their
		own accord as if seeking out sensitive skin to touch.`),
	},
	{
		name:    `[Scrubby Bath Gloves]`,
		Quality: Uncommon,
		Emote: tmplMust(`/me hands {{.User}} pair of {{.Tool}} that leave a
		teasing, tingling sensation wherever they scrub.`),
	},
	{
		name:    `[Vial of Sensitizing Lotion]`,
		Quality: Rare,
		Emote: tmplMust(`/me hands {{.User}} a {{.Tool}}. It enhances the tickling 
		experience in three ways: first by allowing the fingers of the tickler to 
		slide along the ticklee's skin, second by making the affected nerves more 
		sensitive, and third and perhaps most devastating of all, the scent 
		heightens the subjects awareness, keeping their mind fully withing the 
		moment and unable to mentally escape the tickling or grow accustomed to it in any way.`),
	},
	{
		name:    `[Hana Hana no Mi]`,
		Quality: Legendary,
		Emote: tmplMust(`/me hands {{.User}} the legendary devil fruit {{.Tool}}. It 
		allows the eater to replicate and sprout pieces of their body from the surface 
		of any object or living thing. Sprouting extra limbs near or even on the victim, 
		can render almost any foe into submission with ease.`),
	},
	{
		name:    `[Kocho Kobra]`,
		Quality: Epic,
		Emote: tmplMust(`/me hands {{.User}} the infamous {{.Tool}}. The snake is
		absolutely covered with bumps and ridges on its velvet-soft skin, its fangs
		are vibrating nubs, its tongue is akin to a electric toothbrush, and it
		excretes an oil that sensitises the skin. The snake looks to you for
		instructions.`),
	},
}
