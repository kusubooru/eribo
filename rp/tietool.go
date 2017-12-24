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

type Tietool struct {
	Name    string
	Quality Quality
	Desc    *template.Template
	Weight  int
}

func (t Tietool) MarshalText() (string, error) {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s\n", t.Name))
	buf.WriteString(fmt.Sprintf("%s\n", t.Quality))
	buf.WriteString(fmt.Sprintf("%s\n", clean(t.Desc.Tree.Root.String())))
	return buf.String(), nil
}

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
	t.Name = string(name)
	t.Desc = desc
	t.Quality = makeQuality(string(quality))
	return nil
}

func tmplMust(s string) *template.Template {
	return template.Must(template.New("").Parse(s))
}

func (t Tietool) Apply(user string) (string, error) {
	data := struct {
		Tool string
		User string
	}{
		qualityColorBBCode(t.Quality, t.Name),
		user,
	}
	var buf bytes.Buffer
	if err := t.Desc.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to apply %v to %q: %v", data, t.Desc, err)
	}
	return clean(buf.String()), nil
}

func RandTietool(user, toolType string) (string, error) {
	table := &loot.Table{}
	tools := tietools
	switch toolType {
	case "heavy", "hard":
		tools = tietoolsHard
	}
	for _, t := range tools {
		table.Add(t, t.Quality.Weight())
	}
	_, roll := table.Roll(time.Now().UnixNano())
	if tool, ok := roll.(Tietool); ok {
		return tool.Apply(user)
	}
	return "", fmt.Errorf("tietool loot table returned nothing")
}

func Tietools(toolType string) []Tietool {
	if toolType == "heavy" || toolType == "hard" {
		return tietoolsHard
	}
	return tietools
}

var tietools = []Tietool{
	{
		Name:    `[Blindfold]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a black {{.Tool}}.`),
	},
	{
		Name:    `[Silken Sleep Mask]`,
		Quality: Uncommon,
		Desc:    tmplMust(`/me generates for {{.User}} a luxurious {{.Tool}}.`),
	},
	{
		Name:    `[Bondage Rope]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a package of soft, cotton {{.Tool}}.`),
	},
	{
		Name:    `[Toe Ties]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a set of {{.Tool}}.`),
	},
	{
		Name:    `[Shibari Rope]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a hefty amount of {{.Tool}}
		made from natural hemp.`),
	},
	{
		Name:    `[Ball Gag]`,
		Quality: Common,
		Desc: tmplMust(`/me generates for {{.User}} a red {{.Tool}} with black
		leather straps.`),
	},
	{
		Name:    `[Fuzzy Handcuffs]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a pair of pink {{.Tool}}.`),
	},
	{
		Name:    `[Leg Spreader]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    `[Large Weighted Net]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    `[Heat Shrink Tube]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a {{.Tool}}. Body heat
		is enough to make it shrink and trap anything inside into a seamless
		rubber tube.`),
	},
	{
		Name:    `[Web Crawler's Web Shooter Cuffs]`,
		Quality: Rare,
		Desc: tmplMust(`/me generates for {{.User}} two fully loaded {{.Tool}}.
		Thwip! Thwip!`),
	},
	{
		Name:    `[Braided Leather Bolas]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a set of {{.Tool}}. The
		leather balls are connected together with braided leather and contain
		some kind of weight inside.`),
	},
	{
		Name:    `[Lace Collar]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    `[Magnetic Bracers]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} two {{.Tool}}. The
		magnets are strong and easily stick to any metal surface. Their
		polarity is opposite and once they stick together they are very hard to
		separate.`),
	},
	{
		Name:    `[Chinese Finger Traps]`,
		Quality: Uncommon,
		Desc:    tmplMust(`/me generates for {{.User}} a set of colorful {{.Tool}}.`),
	},
	{
		Name:    `[Engraved Leather Collar]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a black {{.Tool}}. The
		silver tag of the collar is engraved with the name "{{.User}}".`),
	},
	{
		Name:    `[Wrist to Ankle Cuffs]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a set of black, nylon
		{{.Tool}} bondage restrains. The wrist cuffs are attached to the ankle
		cuffs strap, forcing the wearer to a bend over position with their legs
		slightly spread, prominently displaying their buttocks.`),
	},
	{
		Name:    `[Collar to Wrist Restrains]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a set of leather
		{{.Tool}}. The collar extends to a leather strap which is connected by
		a chain to a pair of wrist cuffs. The chains between the cuffs and the
		collar pass through three O-rings and form a triangle.`),
	},
	{
		Name:    `[Chloroform Soaked Rag]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    `[Saran Wrap]`,
		Quality: Common,
		Desc: tmplMust(`/me generates for {{.User}} a hefty amount of
		{{.Tool}}.`),
	},
	{
		Name:    `[Chastity Belt]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    `[Lasso of Truth]`,
		Quality: Legendary,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
}

var tietoolsHard = []Tietool{
	{
		Name:    `[Steel Collar]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    `[Straightjacket]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a white {{.Tool}}.`),
	},
	{
		Name:    `[Full Body Straightjacket]`,
		Quality: Uncommon,
		Desc:    tmplMust(`/me generates for {{.User}} a white {{.Tool}}.`),
	},
	{
		Name:    `[Nipple Clamps]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a pair of {{.Tool}}.`),
	},
	{
		Name:    `[Latex Bodysuit]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    `[Latex Dog Suit]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a {{.Tool}}. It forces the
		wearer to walk with their elbows and knees while keeping their feet up.
		It is accompanied with a detachable mask with dog ears.`),
	},
	{
		Name:    `[Gas Mask]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    `[Gimp Mask]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    `[O-Ring Gag]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} an {{.Tool}}.`),
	},
	{
		Name:    `[Slave Collar With Leash]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    `[Vacuum Bed]`,
		Quality: Common,
		Desc:    tmplMust(`/me generates for {{.User}} a {{.Tool}}.`),
	},
	{
		Name:    `[Engraved Slave Collar With Leash]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} an {{.Tool}}. On the
		collar, the words "Tickle Slut" are engraved with large, silver
		letters.`),
	},
	{
		Name:    `[Armbinder]`,
		Quality: Uncommon,
		Desc:    tmplMust(`/me generates for {{.User}} a leather {{.Tool}}.`),
	},
	{
		Name:    `[Monoglove Armbinder]`,
		Quality: Uncommon,
		Desc: tmplMust(`/me generates for {{.User}} a blue {{.Tool}} with
		Y-shaped harness configuration .`),
	},
}
