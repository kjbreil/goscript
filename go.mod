module github.com/kjbreil/goscript

go 1.19

require (
	github.com/adhocore/gronx v1.5.0
	github.com/antonmedv/expr v1.12.5
	github.com/go-logr/logr v1.2.4
	github.com/go-logr/zapr v1.2.3
	github.com/goccy/go-yaml v1.11.0
	github.com/google/uuid v1.3.0
	github.com/iancoleman/strcase v0.2.0
	github.com/kjbreil/hass-mqtt v0.2.2
	github.com/kjbreil/hass-ws v0.2.0
	github.com/mitchellh/mapstructure v1.5.0
	go.uber.org/zap v1.24.0
)

require (
	github.com/eclipse/paho.mqtt.golang v1.4.2 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/go-playground/validator/v10 v10.10.1 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/klauspost/compress v1.16.4 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/net v0.9.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.7.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	nhooyr.io/websocket v1.8.7 // indirect
)

replace github.com/kjbreil/hass-ws => /Users/kjell/dev/hass-ws

replace github.com/kjbreil/hass-mqtt => /Users/kjell/dev/hass-mqtt
