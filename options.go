package herbplugin

import (
	"path/filepath"
	"strings"
)

type Options interface {
	GetLocation() *Location
	GetParam(name string) string
	MustAuthorizeDomain(domain string) bool
	MustAuthorizePath(path string) bool
	MustAuthorizePermission(permission string) bool
}

type NopOptions struct{}

func (o NopOptions) GetLocation() *Location {
	return nil
}
func (o NopOptions) GetParam(name string) string {
	return ""
}
func (o NopOptions) MustAuthorizeDomain(domain string) bool {
	return false
}
func (o NopOptions) MustAuthorizePath(path string) bool {
	return false
}
func (o NopOptions) MustAuthorizePermission(permission string) bool {
	return false
}

type Location struct {
	Path string
}

func (l *Location) MustCleanPath(p string) string {
	return MustCleanPath(l.Path, p)
}
func (l *Location) MustCleanInsidePath(p string) string {
	path, err := filepath.Abs(l.Path)
	if err != nil {
		panic(err)
	}
	cleanpath := l.MustCleanPath(p)
	if !strings.HasPrefix(cleanpath, path) {
		return ""
	}
	return cleanpath
}
func NewLoaction() *Location {
	return &Location{}
}

type Trusted struct {
	Paths   []string
	Domains []string
}

func (t *Trusted) Contains(target *Trusted) bool {
PATH:
	for _, targetpath := range target.Paths {
		for _, v := range t.Paths {
			if v == targetpath {
				continue PATH
			}
		}
		return false
	}
DOMAIN:
	for _, targetdomain := range target.Domains {
		for _, v := range t.Domains {
			if v == targetdomain {
				continue DOMAIN
			}
		}
		return false
	}
	return true
}
func (t *Trusted) MustAuthorizeDomain(domain string) bool {
	for k := range t.Domains {
		if MatchDomain(t.Domains[k], domain) {
			return true
		}
	}
	return false
}
func (t *Trusted) MustAuthorizePath(path string) bool {
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

type PlainOptions struct {
	Location    *Location
	Params      map[string]string
	Trusted     *Trusted
	Permissions []string
}

func (o *PlainOptions) GetLocation() *Location {
	return o.Location
}
func (o *PlainOptions) GetParam(name string) string {
	return o.Params[name]
}
func (o *PlainOptions) MustAuthorizeDomain(domain string) bool {
	return o.Trusted.MustAuthorizeDomain(domain)
}
func (o *PlainOptions) MustAuthorizePath(path string) bool {
	return o.Trusted.MustAuthorizePath(path)
}

func (o *PlainOptions) MustAuthorizePermission(permission string) bool {
	return MustAuthorizePermission(o.Permissions, permission)
}

func NewOptions() *PlainOptions {
	return &PlainOptions{
		Location: NewLoaction(),
		Params:   map[string]string{},
		Trusted:  NewTrusted(),
	}
}

func MustAuthorizePermission(permissions []string, permission string) bool {
	for k := range permissions {
		if permissions[k] == permission {
			return true
		}
	}
	return false
}

func ContainsPermissions(src []string, target []string) bool {
PERMISSION:
	for _, targetpermission := range target {
		for _, v := range src {
			if v == targetpermission {
				continue PERMISSION
			}
		}
		return false
	}
	return true
}
