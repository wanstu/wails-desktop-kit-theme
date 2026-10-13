package theme

// Pack describes one reusable Desktop Kit theme pack.
type Pack struct {
	Name        string
	DisplayName string
	Description string
	File        string
	Signature   bool
}

var packs = []Pack{
	{Name: "aurora", DisplayName: "极光", Description: "极夜深蓝底色，青绿与紫罗兰多层光晕", File: "aurora.css", Signature: true},
	{Name: "sakura", DisplayName: "樱花", Description: "瓷白和浅樱粉，柔和花瓣色与大圆角", File: "sakura.css", Signature: true},
	{Name: "hacker", DisplayName: "黑客", Description: "近黑终端底色、荧光绿与硬边框", File: "hacker.css", Signature: true},
	{Name: "desert", DisplayName: "沙漠", Description: "沙金与赭石，沙漠夕照的暖色层次", File: "desert.css", Signature: true},
	{Name: "ocean", DisplayName: "海洋", Description: "海沟深蓝与浅海青绿，通透水色", File: "ocean.css", Signature: true},
	{Name: "forest", DisplayName: "森林", Description: "苔绿、林木灰绿和柔和自然边界", File: "forest.css", Signature: true},
	{Name: "sunset", DisplayName: "落日", Description: "落日橙与霞光粉的斜向暖色过渡", File: "sunset.css", Signature: true},
	{Name: "nord", DisplayName: "北境", Description: "雪蓝、冰灰和低饱和极简界面", File: "nord.css", Signature: true},
	{Name: "mint", DisplayName: "薄荷", Description: "清透青绿，轻盈舒缓", File: "mint.css"},
	{Name: "rose", DisplayName: "蔷薇", Description: "柔和玫红，明快但不过分甜腻", File: "rose.css"},
	{Name: "amber", DisplayName: "琥珀", Description: "金黄暖色，适合工具与生活类应用", File: "amber.css"},
	{Name: "graphite", DisplayName: "石墨", Description: "黑白中性灰阶、硬朗小圆角与扁平层次", File: "graphite.css", Signature: true},
	{Name: "mocha", DisplayName: "摩卡", Description: "咖啡棕与暖米色，沉稳柔和", File: "mocha.css"},
	{Name: "lavender", DisplayName: "薰衣草", Description: "奶油薰衣草紫、柔和光晕与大圆角", File: "lavender.css", Signature: true},
	{Name: "midnight", DisplayName: "午夜", Description: "午夜靛蓝、深海军蓝和克制的工具界面", File: "midnight.css", Signature: true},
	{Name: "cobalt", DisplayName: "钴蓝", Description: "清晰强烈的蓝色工具风格", File: "cobalt.css"},
	{Name: "jade", DisplayName: "翡翠", Description: "矿物翡翠绿、温润玉色和玻璃质感", File: "jade.css", Signature: true},
	{Name: "cherry", DisplayName: "樱桃", Description: "酒红、樱桃红和浓郁莓果层次", File: "cherry.css", Signature: true},
	{Name: "sand", DisplayName: "沙丘", Description: "米沙与棕褐色，低对比暖色工具感", File: "sand.css"},
	{Name: "plum", DisplayName: "梅紫", Description: "成熟紫红色，适合偏内容型桌面应用", File: "plum.css"},
	{Name: "ice", DisplayName: "冰川", Description: "冰蓝灰背景与清冷青色强调", File: "ice.css"},
	{Name: "terminal", DisplayName: "终端", Description: "低圆角绿色系，突出日志与开发工具气质", File: "terminal.css"},
	{Name: "cyber", DisplayName: "赛博", Description: "深紫电子夜景、电光青与锐利高对比", File: "cyber.css", Signature: true},
	{Name: "sage", DisplayName: "鼠尾草", Description: "低饱和草木绿，柔和安静", File: "sage.css"},
	{Name: "peach", DisplayName: "蜜桃", Description: "柔软珊瑚橙与更圆润的界面轮廓", File: "peach.css"},
	{Name: "denim", DisplayName: "丹宁", Description: "低饱和牛仔蓝，偏沉稳工具风格", File: "denim.css"},
	{Name: "orchid", DisplayName: "兰花", Description: "紫粉色调，适合内容与创作类应用", File: "orchid.css"},
	{Name: "paper", DisplayName: "纸张", Description: "羊皮纸、深褐墨色和扁平印刷质感", File: "paper.css", Signature: true},
	{Name: "obsidian", DisplayName: "黑曜石", Description: "近黑暗色与冷蓝强调，突出高对比工具感", File: "obsidian.css"},
	{Name: "neon", DisplayName: "霓虹", Description: "暗色模式使用紫青高亮与更强对比", File: "neon.css"},
	{Name: "retro", DisplayName: "复古", Description: "复古奶油黄、砖红与青绿色撞色", File: "retro.css", Signature: true},
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
