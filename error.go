package herbplugin

import "errors"

var ErrPluginNotStarted = errors.New("herbplugin: plugin not started")

var ErrPluginStarted = errors.New("herbplugin: plugin started")
