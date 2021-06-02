package herbplugin

import (
	"path/filepath"
	"strings"
)

func MatchDomain(pattern string, domain string) bool {
	if pattern == "" || domain == "" {
		return false
	}
	if pattern[0] == '.' {
		l := strings.SplitN(domain, ".", 2)
		return len(l) == 2 && l[1] == string(domain[1:])
	} else if pattern[0] == '*' {
		return strings.HasSuffix(domain, string(pattern[1:]))
	}
	return string(pattern) == domain
}

func MustCleanPath(base string, newpath string) string {
	var err error
	if !filepath.IsAbs(newpath) {
		base, err = filepath.Abs(base)
		if err != nil {
			panic(err)
		}
		newpath = filepath.Join(base, newpath)
	}

	return filepath.Clean(newpath)
}

const PermissionDangerousAPI = "DANGER!!!ALLOW-MALICIOUS-DAMAGE"
