package herbplugin

import "strings"

type Options interface {
	GetLocation() *Location
	GetParam(name string) string
	MustAuthorizeDomain(domain string) bool
	MustAuthorizePath(path string) bool
	MustAuthorizeDangerousAPI() bool
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
func (o NopOptions) MustAuthorizeDangerousAPI() bool {
	return false
}
func (o NopOptions) MustAuthorizePermission(permission string) bool {
	return false
}

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
	for k := range o.Trusted.Domains {
		if o.Trusted.Domains[k] == domain {
			return true
		}
	}
	return false
}
func (o *PlainOptions) MustAuthorizePath(path string) bool {
	for _, v := range o.Trusted.Paths {
		if v != "" && strings.HasPrefix(path, v) {
			return true
		}
	}
	return false
}
func (o *PlainOptions) MustAuthorizeDangerousAPI() bool {
	return o.Trusted.DangerousAPI
}
func (o *PlainOptions) MustAuthorizePermission(permission string) bool {
	for k := range o.Permissions {
		if o.Permissions[k] == permission {
			return true
		}
	}
	return false
}

func NewOptions() *PlainOptions {
	return &PlainOptions{
		Location: NewLoaction(),
		Params:   map[string]string{},
		Trusted:  NewTrusted(),
	}
}
