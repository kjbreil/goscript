module github.com/kjbreil/goscript

go 1.23.0

toolchain go1.24.2

require (
	github.com/adhocore/gronx v1.19.6
	github.com/brutella/hap v0.0.35
	github.com/dave/jennifer v1.7.1
	github.com/eclipse/paho.mqtt.golang v1.5.0
	github.com/expr-lang/expr v1.17.3
	github.com/go-logr/logr v1.4.2
	github.com/goccy/go-yaml v1.17.1
	github.com/google/uuid v1.6.0
	github.com/iancoleman/strcase v0.3.0
	github.com/kjbreil/hass-mqtt v0.2.3
	github.com/kjbreil/hass-ws v0.2.3
	github.com/mitchellh/mapstructure v1.5.0
	github.com/sixdouglas/suncalc v0.0.0-20250114185126-291b1938b70c
)

require (
	github.com/brutella/dnssd v1.2.14 // indirect
	github.com/go-chi/chi v1.5.5 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/miekg/dns v1.1.66 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	github.com/tadglines/go-pkgs v0.0.0-20210623144937-b983b20f54f9 // indirect
	github.com/vishvananda/netlink v1.3.1 // indirect
	github.com/vishvananda/netns v0.0.5 // indirect
	github.com/xiam/to v0.0.0-20200126224905-d60d31e03561 // indirect
	golang.org/x/crypto v0.38.0 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/net v0.40.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/tools v0.33.0 // indirect
	gopkg.in/Regis24GmbH/go-diacritics.v2 v2.0.3 // indirect
	nhooyr.io/websocket v1.8.17 // indirect
)

replace github.com/kjbreil/hass-ws => ../hass-ws

replace github.com/kjbreil/hass-mqtt => ../hass-mqtt
