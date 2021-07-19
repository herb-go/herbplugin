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

type UnauthorizePermissionError struct {
	Permission string
}

func (e *UnauthorizePermissionError) Error() string {
	return fmt.Sprintf("herbplugin: unauthorize prmission %s", e.Permission)
}

func NewUnauthorizePermissionError(prmission string) error {
	return &UnauthorizePermissionError{
		Permission: prmission,
	}
}
