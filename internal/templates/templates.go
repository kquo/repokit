package templates

import (
	"io/fs"
	"os"
	"path/filepath"
)

// DiskFS returns a filesystem rooted at the templates directory
// within a local repokit checkout. Used by enhance mode and
// during development before embed is wired (AC19).
func DiskFS(repoRoot string) fs.FS {
	return os.DirFS(filepath.Join(repoRoot, "internal", "templates"))
}

// DirPath returns the absolute path to the templates directory
// within a local repokit checkout.
func DirPath(repoRoot string) string {
	return filepath.Join(repoRoot, "internal", "templates")
}
