package herbplugin

import "strings"

type Authorizer interface {
	MustAuthorizeDomain(domain string) bool
	MustAuthorizePath(path string) bool
	MustAuthorizeDangerousAPI() bool
	MustAuthorizePermission(permission string) bool
}

type FixedAuthorizer bool

func (a FixedAuthorizer) MustAuthorizeDomain(domain string) bool {
	return bool(a)
}
func (a FixedAuthorizer) MustAuthorizePath(path string) bool {
	return bool(a)
}
func (a FixedAuthorizer) MustAuthorizeDangerousAPI() bool {
	return bool(a)
}
func (a FixedAuthorizer) MustAuthorizePermission(permission string) bool {
	return bool(a)
}

var DenyAll = FixedAuthorizer(false)

var AllowAll = FixedAuthorizer(true)

type BasicPluginAuthorizer struct {
	Plugin *BasicPlugin
}

func NewBasicPluginAuthorizer(p *BasicPlugin) *BasicPluginAuthorizer {
	return &BasicPluginAuthorizer{
		Plugin: p,
	}
}

func (a *BasicPluginAuthorizer) MustAuthorizeDomain(domain string) bool {
	o := a.Plugin.options
	for k := range o.Trusted.Domains {
		if o.Trusted.Domains[k] == domain {
			return true
		}
	}
	return false
}
func (a *BasicPluginAuthorizer) MustAuthorizePath(path string) bool {
	o := a.Plugin.options
	for _, v := range o.Trusted.Paths {
		if v != "" && strings.HasPrefix(path, v) {
			return true
		}
	}
	return false
}
func (a *BasicPluginAuthorizer) MustAuthorizeDangerousAPI() bool {
	o := a.Plugin.options
	return o.Trusted.DangerousAPI
}
func (a *BasicPluginAuthorizer) MustAuthorizePermission(permission string) bool {
	o := a.Plugin.options
	for k := range o.Permissions {
		if o.Permissions[k] == permission {
			return true
		}
	}
	return false
}
