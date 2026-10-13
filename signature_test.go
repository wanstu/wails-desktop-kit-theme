package theme

import (
	"fmt"
	"io/fs"
	"math"
	"regexp"
	"strings"
	"testing"
)

var cssTokenPattern = regexp.MustCompile(`(--dk-[a-z0-9-]+)\s*:\s*([^;]+);`)

func parseTokens(source string) map[string]string {
	result := map[string]string{}
	for _, match := range cssTokenPattern.FindAllStringSubmatch(source, -1) {
		result[match[1]] = strings.TrimSpace(match[2])
	}
	return result
}

func contrastRatio(foreground, background string) (float64, error) {
	lum := func(hex string) (float64, error) {
		if len(hex) != 7 || hex[0] != '#' {
			return 0, fmt.Errorf("invalid color %q", hex)
		}
		var r, g, b uint8
		if _, err := fmt.Sscanf(hex, "#%02x%02x%02x", &r, &g, &b); err != nil {
			return 0, err
		}
		gamma := func(n uint8) float64 {
			value := float64(n) / 255
			if value <= 0.04045 {
				return value / 12.92
			}
			return math.Pow((value+0.055)/1.055, 2.4)
		}
		return 0.2126*gamma(r) + 0.7152*gamma(g) + 0.0722*gamma(b), nil
	}
	a, err := lum(foreground)
	if err != nil {
		return 0, err
	}
	b, err := lum(background)
	if err != nil {
		return 0, err
	}
	return (math.Max(a, b) + .05) / (math.Min(a, b) + .05), nil
}

func TestSignatureThemesHaveFullPaletteAndContrast(t *testing.T) {
	manifest, err := BuildRemoteManifest("")
	if err != nil {
		t.Fatal(err)
	}
	signatureCount := 0
	distinct := make(map[string]string)
	for _, pack := range manifest.Packs {
		if !pack.Signature {
			continue
		}
		signatureCount++
		if pack.GradientLightSecondary == "" || pack.GradientDarkSecondary == "" {
			t.Errorf("%s missing secondary gradient colors", pack.Name)
		}
		key := pack.GradientLight + "|" + pack.GradientLightSecondary + "|" + pack.GradientDark + "|" + pack.GradientDarkSecondary
		if prior := distinct[key]; prior != "" {
			t.Errorf("%s duplicates %s accent signature", pack.Name, prior)
		}
		distinct[key] = pack.Name
		file, err := fs.ReadFile(embedded, "assets/"+pack.File)
		if err != nil {
			t.Fatal(err)
		}
		// All signature packs must define both light and dark surfaces rather than recoloring a button.
		css := string(file)
		parts := strings.SplitN(css, `:root[data-dk-theme="dark"]`, 2)
		if len(parts) != 2 {
			t.Errorf("%s missing dark palette", pack.Name)
			continue
		}
		for i, part := range parts {
			tokens := parseTokens(part)
			required := []string{"--dk-bg-app", "--dk-bg-elevated", "--dk-bg-subtle", "--dk-text", "--dk-text-secondary", "--dk-border", "--dk-primary", "--dk-primary-contrast", "--dk-gradient-secondary", "--dk-app-gradient", "--dk-sidebar-gradient", "--dk-radius-lg", "--dk-shadow-panel"}
			for _, token := range required {
				if tokens[token] == "" {
					t.Errorf("%s variant=%d missing %s", pack.Name, i, token)
				}
			}
			for _, item := range []struct {
				label, front, back string
				minimum            float64
			}{
				{"body", "--dk-text", "--dk-bg-app", 7},
				{"secondary", "--dk-text-secondary", "--dk-bg-elevated", 4.5},
				{"button", "--dk-primary-contrast", "--dk-primary", 4.5},
			} {
				score, err := contrastRatio(tokens[item.front], tokens[item.back])
				if err != nil {
					t.Errorf("%s variant=%d invalid %s: %v", pack.Name, i, item.label, err)
					continue
				}
				if score < item.minimum {
					t.Errorf("%s variant=%d %s contrast %.2f < %.2f", pack.Name, i, item.label, score, item.minimum)
				}
			}
			if gradient := tokens["--dk-primary-gradient"]; gradient != "" {
				colors := regexp.MustCompile(`#[0-9a-fA-F]{6}`).FindAllString(gradient, -1)
				if len(colors) < 2 {
					t.Errorf("%s variant=%d gradient has insufficient stops", pack.Name, i)
				}
				for _, color := range colors {
					contrast, err := contrastRatio(tokens["--dk-primary-contrast"], color)
					if err != nil || contrast < 4.5 {
						t.Errorf("%s variant=%d gradient stop %s button contrast %.2f (%v)", pack.Name, i, color, contrast, err)
					}
				}
			}
		}
	}
	if signatureCount != 16 {
		t.Errorf("signature count = %d, want 16", signatureCount)
	}
	if len(manifest.Packs) != 31 {
		t.Errorf("all pack IDs must remain usable: count=%d", len(manifest.Packs))
	}
}
