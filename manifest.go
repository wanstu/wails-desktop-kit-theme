package theme

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"regexp"
)

const DefaultRemoteBaseURL = "https://raw.githubusercontent.com/wanstu/wails-desktop-kit-theme/master/assets/"

type RemoteManifest struct {
	SchemaVersion int                  `json:"schema_version"`
	BaseURL       string               `json:"base_url"`
	Packs         []RemoteManifestPack `json:"packs"`
}

type RemoteManifestPack struct {
	Name                   string            `json:"name"`
	DisplayName            string            `json:"display_name"`
	Description            string            `json:"description"`
	File                   string            `json:"file"`
	SHA256                 string            `json:"sha256"`
	GradientLight          string            `json:"gradient_light,omitempty"`
	GradientDark           string            `json:"gradient_dark,omitempty"`
	GradientLightSecondary string            `json:"gradient_light_secondary,omitempty"`
	GradientDarkSecondary  string            `json:"gradient_dark_secondary,omitempty"`
	Signature              bool              `json:"signature,omitempty"`
	GradientSurfacesLight  *GradientSurfaces `json:"gradient_surfaces_light,omitempty"`
	GradientSurfacesDark   *GradientSurfaces `json:"gradient_surfaces_dark,omitempty"`
}

type GradientSurfaces struct {
	App           string `json:"app"`
	Elevated      string `json:"elevated"`
	Subtle        string `json:"subtle"`
	Muted         string `json:"muted"`
	Text          string `json:"text"`
	TextSecondary string `json:"text_secondary"`
	Border        string `json:"border"`
}

func canonicalRemoteBytes(data []byte) []byte {
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	return bytes.ReplaceAll(data, []byte("\r"), []byte("\n"))
}

var primaryColorPattern = regexp.MustCompile(`--dk-primary\s*:\s*(#[0-9a-fA-F]{6})\s*;`)
var secondaryColorPattern = regexp.MustCompile(`--dk-gradient-secondary\s*:\s*(#[0-9a-fA-F]{6})\s*;`)

func gradientColors(data []byte) (string, string, error) {
	matches := primaryColorPattern.FindAllSubmatch(data, -1)
	if len(matches) < 2 {
		return "", "", fmt.Errorf("theme stylesheet must define --dk-primary in light and dark variants")
	}
	return string(matches[0][1]), string(matches[1][1]), nil
}

var surfacePattern = regexp.MustCompile(`--dk-([a-z0-9-]+)\s*:\s*(#[0-9a-fA-F]{6})\s*;`)

func parseSurfaceStyles(data []byte, pack Pack) (*GradientSurfaces, *GradientSurfaces, error) {
	if !pack.Signature {
		return nil, nil, nil
	}
	parts := bytes.SplitN(data, []byte(`:root[data-dk-theme="dark"]`), 2)
	if len(parts) != 2 {
		return nil, nil, fmt.Errorf("%s missing dark variant", pack.Name)
	}
	read := func(part []byte) (*GradientSurfaces, error) {
		values := map[string]string{}
		for _, match := range surfacePattern.FindAllSubmatch(part, -1) {
			values[string(match[1])] = string(match[2])
		}
		required := []string{"bg-app", "bg-elevated", "bg-subtle", "bg-muted", "text", "text-secondary", "border"}
		for _, key := range required {
			if values[key] == "" {
				return nil, fmt.Errorf("%s missing signature surface %s", pack.Name, key)
			}
		}
		return &GradientSurfaces{
			App: values["bg-app"], Elevated: values["bg-elevated"],
			Subtle: values["bg-subtle"], Muted: values["bg-muted"],
			Text: values["text"], TextSecondary: values["text-secondary"], Border: values["border"],
		}, nil
	}
	light, err := read(parts[0])
	if err != nil {
		return nil, nil, err
	}
	dark, err := read(parts[1])
	if err != nil {
		return nil, nil, err
	}
	return light, dark, nil
}

func BuildRemoteManifest(baseURL string) (RemoteManifest, error) {
	if baseURL == "" {
		baseURL = DefaultRemoteBaseURL
	}
	result := RemoteManifest{
		SchemaVersion: 1,
		BaseURL:       baseURL,
		Packs:         make([]RemoteManifestPack, 0, len(packs)),
	}
	for _, pack := range packs {
		data, err := fs.ReadFile(embedded, "assets/"+pack.File)
		if err != nil {
			return RemoteManifest{}, fmt.Errorf("read %s: %w", pack.File, err)
		}
		light, dark, err := gradientColors(data)
		if err != nil {
			return RemoteManifest{}, fmt.Errorf("%s: %w", pack.File, err)
		}
		secondary := secondaryColorPattern.FindAllSubmatch(data, -1)
		var lightSecondary, darkSecondary string
		if len(secondary) == 2 {
			lightSecondary, darkSecondary = string(secondary[0][1]), string(secondary[1][1])
		} else if len(secondary) != 0 {
			return RemoteManifest{}, fmt.Errorf("%s: gradient secondary colors must be provided for both modes", pack.File)
		}
		surfaceLight, surfaceDark, err := parseSurfaceStyles(data, pack)
		if err != nil {
			return RemoteManifest{}, fmt.Errorf("%s: %w", pack.File, err)
		}
		sum := sha256.Sum256(canonicalRemoteBytes(data))
		result.Packs = append(result.Packs, RemoteManifestPack{
			Name:                   pack.Name,
			DisplayName:            pack.DisplayName,
			Description:            pack.Description,
			File:                   pack.File,
			SHA256:                 hex.EncodeToString(sum[:]),
			GradientLight:          light,
			GradientDark:           dark,
			GradientLightSecondary: lightSecondary,
			GradientDarkSecondary:  darkSecondary,
			Signature:              pack.Signature,
			GradientSurfacesLight:  surfaceLight,
			GradientSurfacesDark:   surfaceDark,
		})
	}
	return result, nil
}
