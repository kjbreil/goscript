package service

import "github.com/kjbreil/hass-ws/services"

// Chan is a channel to send services to be run to.
type Chan chan services.Service
