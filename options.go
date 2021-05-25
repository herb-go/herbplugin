package herbplugin

import (
	"strings"
)

type Location struct {
	Path string
}

func NewLoaction() *Location {
	return &Location{}
}

type Trusted struct {
	Paths        []string
	Domains      []string
	Permissions  []string
	DangerousAPI bool
}

func (t *Trusted) hasPermission(permission string) bool {
	for k := range t.Permissions {
		if t.Permissions[k] == permission {
			return true
		}
	}
	return false
}
func (t *Trusted) IsDomainTrusted(domain string) bool {
	for k := range t.Domains {
		if t.Domains[k] == domain {
			return true
		}
	}
	return false
}
func (t *Trusted) IsPathTrusted(path string) bool {
	for _, v := range t.Paths {
		if v != "" && strings.HasPrefix(path, v) {
			return true
		}
	}
	return false
}
func NewTrusted() *Trusted {
	return &Trusted{}
}

type Options struct {
	Location *Location
	Params   map[string]string
	Trusted  *Trusted
}

func NewOptions() *Options {
	return &Options{
		Location: NewLoaction(),
		Params:   map[string]string{},
		Trusted:  NewTrusted(),
	}
}
