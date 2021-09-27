package herbplugin

import "fmt"

type UnauthorizePathError struct {
	Path string
}

func (e *UnauthorizePathError) Error() string {
	return fmt.Sprintf("herbplugin: unauthorize path %s", e.Path)
}

func NewUnauthorizePathError(path string) error {
	return &UnauthorizePathError{
		Path: path,
	}
}
func IsUnauthorizePathError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*UnauthorizePathError)
	return ok
}

type UnauthorizeDomainError struct {
	Domain string
}

func (e *UnauthorizeDomainError) Error() string {
	return fmt.Sprintf("herbplugin: unauthorize domain %s", e.Domain)
}

func NewUnauthorizeDomainError(domain string) error {
	return &UnauthorizeDomainError{
		Domain: domain,
	}
}
func IsUnauthorizeDomainError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*UnauthorizeDomainError)
	return ok
}

type UnauthorizePermissionError struct {
	Permission string
}

func (e *UnauthorizePermissionError) Error() string {
	return fmt.Sprintf("herbplugin: unauthorize permission %s", e.Permission)
}

func IsUnauthorizePermissionError(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*UnauthorizePermissionError)
	return ok
}

func NewUnauthorizePermissionError(prmission string) error {
	return &UnauthorizePermissionError{
		Permission: prmission,
	}
}
