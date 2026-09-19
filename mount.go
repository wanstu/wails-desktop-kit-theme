package theme

import (
	"embed"
	"io"
	"io/fs"
	"sort"
	"strings"
	"time"

	kitui "github.com/wanstu/wails-desktop-kit/ui"
)

//go:embed assets/*.css
var embedded embed.FS

// Mount reserves desktopkit-theme/ for optional reusable theme packs.
func Mount(app fs.FS) fs.FS {
	return mountedFS{app: app}
}

// MountWithKit mounts Desktop Kit UI assets and optional theme packs together.
// app must already be rooted at its index.html.
func MountWithKit(app fs.FS) fs.FS {
	return Mount(kitui.Mount(app))
}

type mountedFS struct{ app fs.FS }

func (m mountedFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) || m.app == nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}
	switch {
	case name == ".":
		entries, err := fs.ReadDir(m.app, ".")
		if err != nil {
			return nil, err
		}
		result := make([]fs.DirEntry, 0, len(entries)+1)
		for _, entry := range entries {
			if entry.Name() != "desktopkit-theme" {
				result = append(result, entry)
			}
		}
		result = append(result, fs.FileInfoToDirEntry(directoryInfo{name: "desktopkit-theme"}))
		sort.Slice(result, func(i, j int) bool { return result[i].Name() < result[j].Name() })
		return &directory{info: directoryInfo{name: "."}, entries: result}, nil
	case name == "desktopkit-theme":
		entries, err := embedded.ReadDir("assets")
		if err != nil {
			return nil, err
		}
		return &directory{info: directoryInfo{name: "desktopkit-theme"}, entries: entries}, nil
	case strings.HasPrefix(name, "desktopkit-theme/"):
		return embedded.Open("assets/" + strings.TrimPrefix(name, "desktopkit-theme/"))
	default:
		return m.app.Open(name)
	}
}

type directoryInfo struct{ name string }

func (i directoryInfo) Name() string     { return i.name }
func (directoryInfo) Size() int64        { return 0 }
func (directoryInfo) Mode() fs.FileMode  { return fs.ModeDir | 0555 }
func (directoryInfo) ModTime() time.Time { return time.Time{} }
func (directoryInfo) IsDir() bool        { return true }
func (directoryInfo) Sys() any           { return nil }

type directory struct {
	info    directoryInfo
	entries []fs.DirEntry
	offset  int
	closed  bool
}

func (d *directory) Stat() (fs.FileInfo, error) {
	if d.closed {
		return nil, fs.ErrClosed
	}
	return d.info, nil
}

func (d *directory) Read([]byte) (int, error) {
	if d.closed {
		return 0, fs.ErrClosed
	}
	return 0, fs.ErrInvalid
}

func (d *directory) Close() error {
	d.closed = true
	return nil
}

func (d *directory) ReadDir(n int) ([]fs.DirEntry, error) {
	if d.closed {
		return nil, fs.ErrClosed
	}
	if d.offset >= len(d.entries) {
		if n > 0 {
			return nil, io.EOF
		}
		return []fs.DirEntry{}, nil
	}
	end := len(d.entries)
	if n > 0 && n < end-d.offset {
		end = d.offset + n
	}
	result := append([]fs.DirEntry(nil), d.entries[d.offset:end]...)
	d.offset = end
	return result, nil
}
