package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	theme "github.com/wanstu/wails-desktop-kit-theme"
)

func main() {
	output := flag.String("write", "manifest.json", "manifest output path")
	flag.Parse()

	manifest, err := theme.BuildRemoteManifest(theme.DefaultRemoteBaseURL)
	if err != nil {
		fail(err)
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		fail(err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(*output, data, 0o644); err != nil {
		fail(err)
	}
	fmt.Printf("wrote %s with %d theme packs\n", *output, len(manifest.Packs))
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
