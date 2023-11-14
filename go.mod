module github.com/kjbreil/goscript

go 1.19

require (
	github.com/adhocore/gronx v1.6.6
	github.com/antonmedv/expr v1.15.3
	github.com/eclipse/paho.mqtt.golang v1.4.3
	github.com/go-logr/logr v1.3.0
	github.com/goccy/go-yaml v1.11.2
	github.com/google/uuid v1.4.0
	github.com/iancoleman/strcase v0.3.0
	github.com/kjbreil/hass-mqtt v0.2.3
	github.com/kjbreil/hass-ws v0.2.2
	github.com/mitchellh/mapstructure v1.5.0
)

require (
	github.com/fatih/color v1.16.0 // indirect
	github.com/go-playground/validator/v10 v10.16.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sync v0.5.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/xerrors v0.0.0-20231012003039-104605ab7028 // indirect
	nhooyr.io/websocket v1.8.10 // indirect
)

replace github.com/kjbreil/hass-ws => /Users/kjell/dev/hass-ws

replace github.com/kjbreil/hass-mqtt => /Users/kjell/dev/hass-mqtt
