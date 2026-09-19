package theme

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
)

const DefaultRemoteBaseURL = "https://raw.githubusercontent.com/wanstu/wails-desktop-kit-theme/master/assets/"

type RemoteManifest struct {
	SchemaVersion int                  `json:"schema_version"`
	BaseURL       string               `json:"base_url"`
	Packs         []RemoteManifestPack `json:"packs"`
}

type RemoteManifestPack struct {
	Name        string `json:"name"`
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	File        string `json:"file"`
	SHA256      string `json:"sha256"`
}

func canonicalRemoteBytes(data []byte) []byte {
	data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	return bytes.ReplaceAll(data, []byte("\r"), []byte("\n"))
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
		sum := sha256.Sum256(canonicalRemoteBytes(data))
		result.Packs = append(result.Packs, RemoteManifestPack{
			Name:        pack.Name,
			DisplayName: pack.DisplayName,
			Description: pack.Description,
			File:        pack.File,
			SHA256:      hex.EncodeToString(sum[:]),
		})
	}
	return result, nil
}
