module github.com/kjbreil/goscript

go 1.23.0

toolchain go1.24.2

require (
	github.com/adhocore/gronx v1.6.6
	github.com/antonmedv/expr v1.15.3
	github.com/brutella/hap v0.0.33
	github.com/dave/jennifer v1.7.1
	github.com/eclipse/paho.mqtt.golang v1.5.0
	github.com/go-logr/logr v1.3.0
	github.com/goccy/go-yaml v1.17.1
	github.com/google/uuid v1.4.0
	github.com/iancoleman/strcase v0.3.0
	github.com/kjbreil/hass-mqtt v0.2.3
	github.com/kjbreil/hass-ws v0.2.2
	github.com/mitchellh/mapstructure v1.5.0
	github.com/sixdouglas/suncalc v0.0.0-20230303054245-f8bc8c69d09e
)

require (
	github.com/brutella/dnssd v1.2.10 // indirect
	github.com/go-chi/chi v1.5.4 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/miekg/dns v1.1.54 // indirect
	github.com/tadglines/go-pkgs v0.0.0-20210623144937-b983b20f54f9 // indirect
	github.com/xiam/to v0.0.0-20200126224905-d60d31e03561 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.39.0 // indirect
	golang.org/x/sync v0.13.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	gopkg.in/Regis24GmbH/go-diacritics.v2 v2.0.3 // indirect
	nhooyr.io/websocket v1.8.17 // indirect
)

replace github.com/kjbreil/hass-ws => /Users/kjell/dev/hass-ws

replace github.com/kjbreil/hass-mqtt => /Users/kjell/dev/hass-mqtt
