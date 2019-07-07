package rp

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/kusubooru/eribo/loot"
)

// Tietool is a tool meant to tie its victim.
type Tietool struct {
	name    string
	Quality Quality
	Desc    *template.Template
	Weight  int
}

// Name returns the name of the tool.
func (t Tietool) Name() string {
	return t.name
}

// NameBBCode returns the name of the tool in BBCode.
func (t Tietool) NameBBCode() string {
	return qualityColorBBCode(t.Quality, t.Name())
}

// MarshalText indicates how a tool will appear as text.
func (t Tietool) MarshalText() (string, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s\n", t.Name()))
	buf.WriteString(fmt.Sprintf("%s\n", t.Quality))
	buf.WriteString(fmt.Sprintf("%s\n", clean(t.Desc.Tree.Root.String())))
	return buf.String(), nil
}

// UnmarshalText can scan a tool from text.
func (t *Tietool) UnmarshalText(s string) error {
	r := strings.NewReader(s)
	br := bufio.NewReader(r)
	name, _, err := br.ReadLine()
	if err != nil {
		return err
	}
	quality, _, err := br.ReadLine()
	if err != nil {
		return err
	}
	tmpl, _, err := br.ReadLine()
	if err != nil {
		return err
	}
	desc, err := template.New("").Parse(string(tmpl))
	if err != nil {
		return err
	}
	t.name = string(name)
	t.Desc = desc
	t.Quality = makeQuality(string(quality))
	return nil
}

func tmplMust(s string) *template.Template {
	return template.Must(template.New("").Parse(s))
}

// Apply applies the user name to the tool template.
func (t Tietool) Apply(user string) (string, error) {
	data := struct {
		Tool string
		User string
	}{
		qualityColorBBCode(t.Quality, t.Name()),
		user,
	}
	var buf bytes.Buffer
	if err := t.Desc.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to apply %v to %v: %v", data, t.Desc, err)
	}
	return clean(buf.String()), nil
}

// TietoolsLootTable is the loot table of the tietools.
type TietoolsLootTable struct {
	*loot.Table
	ToolType string
}

// NewTietoolsLootTable creates a new loot table for the tietools.
func NewTietoolsLootTable(toolType string) *TietoolsLootTable {
	table := &loot.Table{}
	tools := tietools
	switch toolType {
	case "heavy", "hard":
		tools = tietoolsHard
	}
	for _, t := range tools {
		table.Add(t, t.Quality.Weight())
	}
	return &TietoolsLootTable{Table: table, ToolType: toolType}
}

// Legendaries returns how many legendaries are left on the loot table.
func (t *TietoolsLootTable) Legendaries() int {
	legos := 0
	drops := t.Drops()
	t.RLock()
	defer t.RUnlock()
	for _, d := range drops {
		if d.Item == nil {
			continue
		}
		tietool, ok := d.Item.(Tietool)
		if !ok {
			continue
		}
		if tietool.Quality == Legendary && d.Weight > 0 {
			legos++
		}
	}
	return legos
}

// RandTietoolDecreaseWeight returns a random tietool but also decreases its weight.
func (t *TietoolsLootTable) RandTietoolDecreaseWeight(user string) (string, error) {
	legos := t.Legendaries()
	if legos == 0 {
		t = NewTietoolsLootTable(t.ToolType)
	}
	seed := time.Now().UnixNano()
	_, roll := t.RollDecreaseWeight(seed)
	if roll == nil {
		return "", fmt.Errorf("tietool loot table returned nothing")
	}
	tool, ok := roll.(Tietool)
	if !ok {
		return "", fmt.Errorf("TietoolsLootTable contains an item that is not a Tietool")
	}
	return tool.Apply(user)
}

// RandTietool returns a random tietool.
func RandTietool(user, toolType string) (string, error) {
	table := NewTietoolsLootTable(toolType)
	_, roll := table.Roll(time.Now().UnixNano())
	if tool, ok := roll.(Tietool); ok {
		return tool.Apply(user)
	}
	return "", fmt.Errorf("tietool loot table returned nothing")
}

// Tietools returns all the tietools.
func Tietools(toolType string) []Tietool {
	if toolType == "heavy" || toolType == "hard" {
		return tietoolsHard
	}
	return tietools
}

var tietools = []Tietool{
	{
		name:    `[Blindfold]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a black {{.Tool}}.`),
	},
	{
		name:    `[Silken Sleep Mask]`,
		Quality: Uncommon,
		Desc:    tmplMust(`/me generates for {{.User}} a luxurious {{.Tool}}.`),
	},
	{
		name:    `[Bondage Rope]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a package of soft, cotton {{.Tool}}.`),
	},
	{
		name:    `[Toe Ties]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a set of {{.Tool}}.`),
	},
	{
		name:    `[Shibari Rope]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a hefty amount of {{.Tool}}
		made from natural hemp.`),
	},
	{
		name:    `[Ball Gag]`,
		Quality: Common,
		Desc: tmplMust(`/me generates for {{.User}} a red {{.Tool}} with black
		leather straps.`),
	},
	{
		name:    `[Fuzzy Handcuffs]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a pair of pink {{.Tool}}.`),
	},
	{
		name:    `[Leg Spreader]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		name:    `[Large Weighted Net]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		name:    `[Heat Shrink Tube]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a {{.Tool}}. Body heat
		is enough to make it shrink and trap anything inside into a seamless
		rubber tube.`),
	},
	{
		name:    `[Web Crawler's Web Shooter Cuffs]`,
		Quality: Rare,
		Desc: tmplMust(`/me generates for {{.User}} two fully loaded {{.Tool}}.
		Thwip! Thwip!`),
	},
	{
		name:    `[Braided Leather Bolas]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a set of {{.Tool}}. The
		leather balls are connected together with braided leather and contain
		some kind of weight inside.`),
	},
	{
		name:    `[Lace Collar]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		name:    `[Magnetic Bracers]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} two {{.Tool}}. The
		magnets are strong and easily stick to any metal surface. Their
		polarity is opposite and once they stick together they are very hard to
		separate.`),
	},
	{
		name:    `[Chinese Finger Traps]`,
		Quality: Uncommon,
		Desc:    tmplMust(`/me generates for {{.User}} a set of colorful {{.Tool}}.`),
	},
	{
		name:    `[Engraved Leather Collar]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a black {{.Tool}}. The
		silver tag of the collar is engraved with the name "{{.User}}".`),
	},
	{
		name:    `[Wrist to Ankle Cuffs]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a set of black, nylon
		{{.Tool}} bondage restrains. The wrist cuffs are attached to the ankle
		cuffs strap, forcing the wearer to a bend over position with their legs
		slightly spread, prominently displaying their buttocks.`),
	},
	{
		name:    `[Collar to Wrist Restrains]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a set of leather
		{{.Tool}}. The collar extends to a leather strap which is connected by
		a chain to a pair of wrist cuffs. The chains between the cuffs and the
		collar pass through three O-rings and form a triangle.`),
	},
	{
		name:    `[Chloroform Soaked Rag]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		name:    `[Saran Wrap]`,
		Quality: Common,
		Desc: tmplMust(`/me generates for {{.User}} a hefty amount of
		{{.Tool}}.`),
	},
	{
		name:    `[Chastity Belt]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		name:    `[Lasso of Truth]`,
		Quality: Legendary,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		name:    `[Bondage Yoke]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me hands {{.User}} a {{.Tool}}. It's made of
		aluminum, with padded leather lining the neck and wrist restraints.`),
	},
	{
		name:    `[Bolas]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me hands {{.User}} a {{.Tool}}. The soft, cotton
		rope has two padded weights on either end.`)},
	{
		name:    `[Panel Gag Harness]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me hands {{.User}} a {{.Tool}}. It comes with an
		optional padlock if needed to secure in place.`)},
	{
		name:    `[Green Goo Ball]`,
		Quality: Common,
		Desc: tmplMust(`/me hands {{.User}} a {{.Tool}}, and wipes the
		excess off on the side of a couch. It's very sticky!`)},
	{
		name:    `[Red Goo Ball]`,
		Quality: Common,
		Desc: tmplMust(`/me hands {{.User}} a {{.Tool}}, and wipes the
		excess off on the back of someone's shirt. It's very sticky and
		unusually warm!`),
	},
	{
		name:    `[Silk Sashes]`,
		Quality: Common,
		Desc: tmplMust(`/me hands {{.User}} several sturdy {{.Tool}}. Say
		that five times fast.`)},
	{
		name:    `[Eyeless Balaclava]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me hands {{.User}} a {{.Tool}}. There is only a
		hole cut out for the mouth.`)},
	{
		name:    `[Upperbody Posture Brace]`,
		Quality: Rare,
		Desc: tmplMust(`/me hands {{.User}} a {{.Tool}}. It's far more rigid
		and restrictive than usual.`)},
	{
		name:    `[Knee Brace]`,
		Quality: Rare,
		Desc: tmplMust(`/me hands {{.User}} a {{.Tool}}. It's made of steel
		and doesn't bend.`)},
	{
		name:    `[Invisible Straitjacket]`,
		Quality: Epic,
		Desc: tmplMust(`/me hands {{.User}} a {{.Tool}}, maybe? It's hard to
		see if anything's actually there.`)},
	{
		name:    `[VonVitae-Seeking Bondage Mittens]`,
		Quality: Legendary,
		Desc: tmplMust(`/me hands {{.User}} a pair of {{.Tool}}. Push the
		button on the side and they'll seek out their pre-programmed target.`),
	},
	{
		name:    `[Paralysis Cattle Prod]`,
		Quality: Epic,
		Desc: tmplMust(`/me hands {{.User}} a {{.Tool}} with a full battery.
		Instead of a painful shock the afflicted area is paralyzed for a few
		minutes.`),
	},
}

var tietoolsHard = []Tietool{
	{
		name:    `[Steel Collar]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		name:    `[Straightjacket]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a white {{.Tool}}.`),
	},
	{
		name:    `[Full Body Straightjacket]`,
		Quality: Uncommon,
		Desc:    tmplMust(`/me generates for {{.User}} a white {{.Tool}}.`),
	},
	{
		name:    `[Nipple Clamps]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a pair of {{.Tool}}.`),
	},
	{
		name:    `[Latex Bodysuit]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		name:    `[Latex Dog Suit]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a {{.Tool}}. It forces the
		wearer to walk with their elbows and knees while keeping their feet up.
		It is accompanied with a detachable mask with dog ears.`),
	},
	{
		name:    `[Gas Mask]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		name:    `[Gimp Mask]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		name:    `[O-Ring Gag]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} an {{.Tool}}.`),
	},
	{
		name:    `[Slave Collar With Leash]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		name:    `[Vacuum Bed]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		name:    `[Engraved Slave Collar With Leash]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} an {{.Tool}}. On the
		collar, the words "Tickle Slut" are engraved with large, silver
		letters.`),
	},
	{
		name:    `[Armbinder]`,
		Quality: Uncommon,
		Desc:    tmplMust(`/me generates for {{.User}} a leather {{.Tool}}.`),
	},
	{
		name:    `[Monoglove Armbinder]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a blue {{.Tool}} with
		Y-shaped harness configuration .`),
	},
	{
		name:    `[Chastitease Belt]`,
		Quality: Rare,
		Desc: tmplMust(`/me slips a {{.Tool}} on {{.User}} and locks it on. The
		belt is designed to block the wearer from external stimulation while
		the malleable, almost sentient, interior holds said wearer on the edge
		of orgasm.`),
	},
	{
		name:    `[Corset of Gargalesis]`,
		Quality: Epic,
		Desc: tmplMust(`/me laces and locks the {{.Tool}} onto {{.User}}. The
		enchanted cloth, as if it has a mind of its own, it starts tightening
		and immediately begins inflicting the feeling of fingers poking,
		prodding and wiggling into the victim's ribs sides and tummy.`),
	},
}
