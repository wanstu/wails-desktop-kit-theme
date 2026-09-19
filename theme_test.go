package theme

import (
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
)

func TestManifest(t *testing.T) {
	t.Parallel()
	want := []string{"aurora", "ocean", "forest", "sunset"}
	got := Packs()
	if len(got) != len(want) {
		t.Fatalf("packs = %d, want %d", len(got), len(want))
	}
	for i, name := range want {
		if got[i].Name != name {
			t.Fatalf("pack[%d] = %q, want %q", i, got[i].Name, name)
		}
		path, ok := StylesheetPath(name)
		if !ok || path != "/desktopkit-theme/"+name+".css" {
			t.Fatalf("path(%q) = %q, %v", name, path, ok)
		}
	}
}

func TestThemePacksOnlyOverrideDesktopKitTokens(t *testing.T) {
	t.Parallel()
	for _, pack := range Packs() {
		data, err := fs.ReadFile(embedded, "assets/"+pack.File)
		if err != nil {
			t.Fatal(err)
		}
		css := string(data)
		if !strings.Contains(css, "data-dk-theme-pack=\""+pack.Name+"\"") {
			t.Fatalf("%s missing pack selector", pack.File)
		}
		if !strings.Contains(css, "data-dk-theme=\"dark\"") {
			t.Fatalf("%s missing dark variant", pack.File)
		}
		if strings.Contains(css, ".dk-") || strings.Contains(css, "!important") {
			t.Fatalf("%s must not override Desktop Kit component selectors", pack.File)
		}
		for _, line := range strings.Split(css, "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "--") && !strings.HasPrefix(line, "--dk-") {
				t.Fatalf("%s contains non-kit token declaration %q", pack.File, line)
			}
		}
	}
}

func TestMountWithKitServesBothAssetNamespaces(t *testing.T) {
	t.Parallel()
	app := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("app")},
	}
	mounted := MountWithKit(app)
	for _, name := range []string{
		"index.html",
		"desktopkit/tokens.css",
		"desktopkit/theme.js",
		"desktopkit-theme/aurora.css",
		"desktopkit-theme/ocean.css",
	} {
		data, err := fs.ReadFile(mounted, name)
		if err != nil || len(data) == 0 {
			t.Fatalf("read %s: len=%d err=%v", name, len(data), err)
		}
	}
}
