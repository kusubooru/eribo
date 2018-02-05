package rp

import "fmt"

type muffin struct {
	Name string
	URL  string
}

func (m muffin) apply(user string) string {
	return fmt.Sprintf("/me prepares some fresh [url=%s]%s[/url] for %s.", m.URL, m.Name, user)
}

func RandMuffin(user string) string {
	m := muffins[newRand(len(muffins))]
	return m.apply(user)
}

var muffins = []muffin{
	{
		Name: "Grandma's Apple Muffins",
		URL:  "https://i.imgur.com/ixCx6QG.jpg",
	},
	{
		Name: "Banana Nut Muffins",
		URL:  "https://i.imgur.com/j9RIMyT.jpg",
	},
	{
		Name: "Blueberry Cream Cheese Muffins",
		URL:  "https://i.imgur.com/THeBXDX.jpg",
	},
	{
		Name: "Apple Pecan Muffins",
		URL:  "https://i.imgur.com/VWI1z44.jpg",
	},
	{
		Name: "Stuffin' Egg Muffins",
		URL:  "https://i.imgur.com/lOo4eG0.jpg",
	},
	{
		Name: "Spud Muffins",
		URL:  "https://i.imgur.com/Q8FHWMh.jpg",
	},
	{
		Name: "Whole-Wheat Sweet Potato Muffins",
		URL:  "https://i.imgur.com/4aUUQRk.jpg",
	},
	{
		Name: "Cinnamon Muffins",
		URL:  "https://i.imgur.com/3t2OG5H.jpg",
	},
	{
		Name: "Blueberry Buttermilk Muffins",
		URL:  "https://i.imgur.com/uFVy78T.jpg",
	},
	{
		Name: "Best-Ever Banana Muffins",
		URL:  "https://i.imgur.com/XNGxO1X.jpg",
	},
	{
		Name: "Coffee, Walnut & Chocolate Chip Muffins",
		URL:  "https://i.imgur.com/6NQE7Cb.jpg",
	},
	{
		Name: "Light Pumpkin-Chocolate Chip Muffins",
		URL:  "https://i.imgur.com/xnyQjoN.jpg",
	},
	{
		Name: "Classic Bran Muffins",
		URL:  "https://i.imgur.com/IQonWET.jpg",
	},
	{
		Name: "Chocolate Chip Muffins",
		URL:  "https://i.imgur.com/KrAjZ9l.jpg",
	},
	{
		Name: "Baileys & Chocolate Muffins",
		URL:  "https://i.imgur.com/GP9HKro.jpg",
	},
	{
		Name: "Cranberry-Orange Muffins",
		URL:  "https://i.imgur.com/iD487ES.jpg",
	},
	{
		Name: "Cheddar Cheese Muffins",
		URL:  "https://i.imgur.com/geJAmQq.jpg",
	},
	{
		Name: "Maple Bacon Muffins",
		URL:  "https://i.imgur.com/fa7RHVC.jpg",
	},
	{
		Name: "Snickerdoodle Mini Muffins",
		URL:  "https://i.imgur.com/fkFD163.jpg",
	},
	{
		Name: "Sausage Brunch Muffins",
		URL:  "https://i.imgur.com/bhqK4NZ.jpg",
	},
	{
		Name: "Morning Glory Muffins",
		URL:  "https://i.imgur.com/7ajyZ55.jpg",
	},
	{
		Name: "Lemon Poppy-Seed Muffins",
		URL:  "https://i.imgur.com/zCLh4nV.jpg",
	},
	{
		Name: "Peanut Butter Banana Muffins",
		URL:  "https://i.imgur.com/UgMvEpH.jpg",
	},
	{
		Name: "Dot's Corn Muffins",
		URL:  "https://i.imgur.com/Nrupl6O.jpg",
	},
	{
		Name: "Low-Fat Carrot Cake Muffins",
		URL:  "https://i.imgur.com/uD9KYX7.jpg",
	},
	{
		Name: "Gluten-Free Muffins",
		URL:  "https://i.imgur.com/TRtgavD.jpg",
	},
	{
		Name: "Raspberry Buttermilk Muffins",
		URL:  "https://i.imgur.com/BAI0Drp.jpg",
	},
	{
		Name: "Vegan Peanut Butter-Oatmeal Muffins",
		URL:  "https://i.imgur.com/gxr5ZY9.jpg",
	},
	{
		Name: "Cottage Cheese & Dill Muffins",
		URL:  "https://i.imgur.com/dMeUjEA.jpg",
	},
	{
		Name: "Miniature French Breakfast Muffin Puffs",
		URL:  "https://i.imgur.com/EXSem60.jpg",
	},
	{
		Name: "Cinnamon Streusel-Apple Cider Muffins",
		URL:  "https://i.imgur.com/SHrsEkd.jpg",
	},
	{
		Name: "Grape Muffins",
		URL:  "https://i.imgur.com/rCXneFy.jpg",
	},
	{
		Name: "Bran Date Muffins",
		URL:  "https://i.imgur.com/CvcRP6G.jpg",
	},
	{
		Name: "Feta, Onion & Rosemary Muffins",
		URL:  "https://i.imgur.com/rvx4773.jpg",
	},
	{
		Name: "Cherry Muffins",
		URL:  "https://i.imgur.com/qMgc9GX.jpg",
	},
	{
		Name: "Orange Marmalade Muffins",
		URL:  "https://i.imgur.com/L5Z0KY1.jpg",
	},
	{
		Name: "Whole-Wheat Honey-Banana Muffins",
		URL:  "https://i.imgur.com/SYtmbET.jpg",
	},
	{
		Name: "Low-Fat Oatmeal Pumpkin Spice Muffins",
		URL:  "https://i.imgur.com/BREx3OS.jpg",
	},
	{
		Name: "Garlic-Onion Dinner Muffins",
		URL:  "https://i.imgur.com/FsrcS8U.jpg",
	},
	{
		Name: "Cranberry Oatmeal Muffins",
		URL:  "https://i.imgur.com/Zq8rEcD.jpg",
	},
	{
		Name: "Coffee Cake Muffins",
		URL:  "https://i.imgur.com/pUqof4D.jpg",
	},
	{
		Name: "Ham & Cheese Muffins",
		URL:  "https://i.imgur.com/N8pd5qM.jpg",
	},
	{
		Name: "Pineapple & Sour Cream Muffins",
		URL:  "https://i.imgur.com/pTV0G47.jpg",
	},
	{
		Name: "Mango Muffins",
		URL:  "https://i.imgur.com/CtR4NUe.jpg",
	},
	{
		Name: "Brownie Muffins",
		URL:  "https://i.imgur.com/xXUtXfh.jpg",
	},
	{
		Name: "Pear & Ginger Muffins",
		URL:  "https://i.imgur.com/mgEdqDD.jpg",
	},
	{
		Name: "King & Prince Oatmeal Raisin Muffins",
		URL:  "https://i.imgur.com/b1xMCDb.jpg",
	},
	{
		Name: "Vidalia Onion & Shallot Cheese Muffins",
		URL:  "https://i.imgur.com/7NfZ0rt.jpg",
	},
	{
		Name: "Lemon Yogurt Muffins",
		URL:  "https://i.imgur.com/5Fy4ck0.jpg",
	},
	{
		Name: "Banana-Chocolate Chip Muffins",
		URL:  "https://i.imgur.com/54Fig2q.jpg",
	},
	{
		Name: "Kathie's Zucchini Muffins",
		URL:  "https://i.imgur.com/rnJuRU4.jpg",
	},
	{
		Name: "Cranberry Streusel Muffins",
		URL:  "https://i.imgur.com/QflYy2S.jpg",
	},
}
