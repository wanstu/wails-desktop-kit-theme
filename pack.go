package theme

// Pack describes one reusable Desktop Kit theme pack.
type Pack struct {
	Name        string
	DisplayName string
	Description string
	File        string
}

var packs = []Pack{
	{Name: "aurora", DisplayName: "极光", Description: "蓝紫冷色，适合作为通用默认主题", File: "aurora.css"},
	{Name: "ocean", DisplayName: "海洋", Description: "蓝青色调，清爽明亮", File: "ocean.css"},
	{Name: "forest", DisplayName: "森林", Description: "绿色系，柔和自然", File: "forest.css"},
	{Name: "sunset", DisplayName: "落日", Description: "橙粉暖色，更有生活感", File: "sunset.css"},
	{Name: "nord", DisplayName: "北境", Description: "低饱和蓝灰，克制安静", File: "nord.css"},
	{Name: "mint", DisplayName: "薄荷", Description: "清透青绿，轻盈舒缓", File: "mint.css"},
	{Name: "rose", DisplayName: "蔷薇", Description: "柔和玫红，明快但不过分甜腻", File: "rose.css"},
	{Name: "amber", DisplayName: "琥珀", Description: "金黄暖色，适合工具与生活类应用", File: "amber.css"},
	{Name: "graphite", DisplayName: "石墨", Description: "中性灰阶与紧凑圆角，偏专业工具感", File: "graphite.css"},
	{Name: "mocha", DisplayName: "摩卡", Description: "咖啡棕与暖米色，沉稳柔和", File: "mocha.css"},
	{Name: "lavender", DisplayName: "薰衣草", Description: "浅紫柔和，并使用更圆润的界面轮廓", File: "lavender.css"},
	{Name: "midnight", DisplayName: "午夜", Description: "深蓝靛色，高对比且偏开发工具风格", File: "midnight.css"},
}

// Packs returns a copy of the built-in reusable pack manifest.
func Packs() []Pack {
	result := make([]Pack, len(packs))
	copy(result, packs)
	return result
}

// Lookup returns one built-in pack by stable name.
func Lookup(name string) (Pack, bool) {
	for _, pack := range packs {
		if pack.Name == name {
			return pack, true
		}
	}
	return Pack{}, false
}

// StylesheetPath returns the mounted public path for a pack.
func StylesheetPath(name string) (string, bool) {
	pack, ok := Lookup(name)
	if !ok {
		return "", false
	}
	return "/desktopkit-theme/" + pack.File, true
}
