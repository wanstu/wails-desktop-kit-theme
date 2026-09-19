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
