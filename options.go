package herbplugin

import "strings"

type Location struct {
	Path string
}

func NewLoaction() *Location {
	return &Location{}
}

type Trusted struct {
	Paths        []string
	Domains      []string
	DangerousAPI bool
}

func NewTrusted() *Trusted {
	return &Trusted{}
}

type Options struct {
	Location    *Location
	Params      map[string]string
	Trusted     *Trusted
	Permissions []string
}

func (o *Options) HasPermission(permission string) bool {
	for k := range o.Permissions {
		if o.Permissions[k] == permission {
			return true
		}
	}
	return false
}
func (o *Options) IsDomainTrusted(domain string) bool {
	for k := range o.Trusted.Domains {
		if o.Trusted.Domains[k] == domain {
			return true
		}
	}
	return false
}
func (o *Options) IsPathTrusted(path string) bool {
	for _, v := range o.Trusted.Paths {
		if v != "" && strings.HasPrefix(path, v) {
			return true
		}
	}
	return false
}

func NewOptions() *Options {
	return &Options{
		Location: NewLoaction(),
		Params:   map[string]string{},
		Trusted:  NewTrusted(),
	}
}
