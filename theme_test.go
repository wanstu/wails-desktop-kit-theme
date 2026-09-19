package theme

import (
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
)

func TestManifestMatchesEmbeddedThemeFiles(t *testing.T) {
	t.Parallel()

	packs := Packs()
	if len(packs) == 0 {
		t.Fatal("expected at least one theme pack")
	}

	registered := make(map[string]bool, len(packs))
	for _, pack := range packs {
		if pack.Name == "" || pack.DisplayName == "" || pack.Description == "" {
			t.Fatalf("incomplete pack metadata: %+v", pack)
		}
		if pack.File != pack.Name+".css" {
			t.Fatalf("pack %q file = %q, want %q", pack.Name, pack.File, pack.Name+".css")
		}
		if registered[pack.Name] {
			t.Fatalf("duplicate pack name %q", pack.Name)
		}
		registered[pack.Name] = true

		path, ok := StylesheetPath(pack.Name)
		if !ok || path != "/desktopkit-theme/"+pack.File {
			t.Fatalf("path(%q) = %q, %v", pack.Name, path, ok)
		}
		if _, err := fs.Stat(embedded, "assets/"+pack.File); err != nil {
			t.Fatalf("registered pack %q missing embedded file: %v", pack.Name, err)
		}
	}

	entries, err := fs.ReadDir(embedded, "assets")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".css") {
			continue
		}
		name := strings.TrimSuffix(entry.Name(), ".css")
		if !registered[name] {
			t.Fatalf("embedded theme %q is not registered in Packs()", entry.Name())
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
		"desktopkit-theme/midnight.css",
	} {
		data, err := fs.ReadFile(mounted, name)
		if err != nil || len(data) == 0 {
			t.Fatalf("read %s: len=%d err=%v", name, len(data), err)
		}
	}
}
