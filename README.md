# GoScript
Something like PyScript for Home Assistant but in Go. Functionality is being added as needed for my automations but once I have finished what I need I will go through PyScript and backfill any missing functionality. There will be additions to what PyScript can like the ability to add new devices to Home Assistant through MQTT.


## Configuration
Configuration is stored in a Yaml file. Only websocket is required. For home assistant setup a long lived token specific to your scripts.
```yaml
websocket:
  host: <server host or ip>
  port: 8123
  token: <super secret token>
```
To allow GoScript to create Home Assistant devices MQTT is required. Node ID is presented to the MQTT server. If it is not unique within your MQTT server messages can get lost.
```yaml
mqtt:
  node_id: goscript
  mqtt:
    host: <mqtt host or ip>
    port: 1883
    ssl: false # SSL Not Supported yet
```
Use goscript.ParseConfig(path, modules) to parse configuration. The second parameter, modules, is a map[string]interface{} used to assign other configuration entries to custom structs. For example if I have a struct Lights
```go
type Lights struct {
    Name      string
    Entities  []string
}
```
And a configuration entry like 
```yaml
lights:
  name: test
  entities:
    - light.door
    - light.door2
```
To fill in my Lights struct from the config file
```go
modules := map[string]interface{
	"lights": &Lights{}
}
```
Then ParseConfig will fill in the struct properly and can get my config back from the GoScript.GetModule(key) method. Note that GetModule will return a interface, you will need to cast that back to your type.
```go
inter, err := gs.GetModule(key)
if err != nil {
    return nil
}
lights := inter.(*Lights)
```

## Triggers
