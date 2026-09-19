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
	{Name: "cobalt", DisplayName: "钴蓝", Description: "清晰强烈的蓝色工具风格", File: "cobalt.css"},
	{Name: "jade", DisplayName: "翡翠", Description: "偏沉稳的玉石绿，圆润但不轻浮", File: "jade.css"},
	{Name: "cherry", DisplayName: "樱桃", Description: "鲜明红粉，保持克制的层级对比", File: "cherry.css"},
	{Name: "sand", DisplayName: "沙丘", Description: "米沙与棕褐色，低对比暖色工具感", File: "sand.css"},
	{Name: "plum", DisplayName: "梅紫", Description: "成熟紫红色，适合偏内容型桌面应用", File: "plum.css"},
	{Name: "ice", DisplayName: "冰川", Description: "冰蓝灰背景与清冷青色强调", File: "ice.css"},
	{Name: "terminal", DisplayName: "终端", Description: "低圆角绿色系，突出日志与开发工具气质", File: "terminal.css"},
	{Name: "cyber", DisplayName: "赛博", Description: "紫色主调配青色强调，暗色模式更鲜明", File: "cyber.css"},
	{Name: "sage", DisplayName: "鼠尾草", Description: "低饱和草木绿，柔和安静", File: "sage.css"},
	{Name: "peach", DisplayName: "蜜桃", Description: "柔软珊瑚橙与更圆润的界面轮廓", File: "peach.css"},
	{Name: "denim", DisplayName: "丹宁", Description: "低饱和牛仔蓝，偏沉稳工具风格", File: "denim.css"},
	{Name: "orchid", DisplayName: "兰花", Description: "紫粉色调，适合内容与创作类应用", File: "orchid.css"},
	{Name: "paper", DisplayName: "纸张", Description: "暖白纸张感、低阴影与紧凑圆角", File: "paper.css"},
	{Name: "obsidian", DisplayName: "黑曜石", Description: "近黑暗色与冷蓝强调，突出高对比工具感", File: "obsidian.css"},
	{Name: "neon", DisplayName: "霓虹", Description: "暗色模式使用紫青高亮与更强对比", File: "neon.css"},
	{Name: "retro", DisplayName: "复古", Description: "陶土橙配青绿强调，带复古工具气质", File: "retro.css"},
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
