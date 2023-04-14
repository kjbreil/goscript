[![GoDoc](https://img.shields.io/badge/pkg.go.dev-doc-blue)](http://pkg.go.dev/github.com/kjbreil/goscript)

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
Three type of triggers are available. Standard domain.entity, all domain and then periodic. All triggers are automatically deduped and added to the states array.
### Domain Entity Triggers
Triggers are an array of strings. Format for each string is "domain.entity", there is no validation that the domain entity combination exists.
### Domain Triggers
DomainTriggers is an array of strings containing just the domain. All entities within that domain will cause the trigger to fire.
### Periodic Triggers
PeriodicTriggers is an array of strings containing the cron expression matching for when the trigger should run. A blank cron expression, "", will launch the trigger at program start. All cron jobs are evaluated every minute so no periodic job can be set to run quicker than 1 minute. Cron expression parsing and matching is provided by [gronx](github.com/adhocore/gronx).
### States
States is an array of other entities that you would like the states to be available within the function that is run. All triggers are automatically added to states.
### Evaluation
Evaluation is done through a list of strings that are run through [expr](github.com/antonmedv/expr) to evaluate the output. Like with PyScript type is important in the evaluation scripts. Check out [expr](github.com/antonmedv/expr) for more details on casting and converting. You cannot mix types in a single evaluation so `state == "on" || state > 10` will always return false due to failure parsing the evaluation. Attributes are available inside the evaluations,`color_temp > 100` will work as long as color_temp exists in the attributes of the entity and the data type is a float
### TriggerFunc
Func is the function to run when the criteria are met. Within the trigger function a *Task is available to give information on the trigger. Killing the TriggerFunc panics to exit. The runner recovers this panic, this also means that if your code panics the whole program will not crash but will continue. Panic will be written to the logs.

### Example
This trigger will fire at program startup, every minute and every time input.button.test_button is pressed. It will flip input_boolean.test_toggle, wait 10 seconds and flip it back
```go
&goscript.Trigger{
    Unique:        &goscript.Unique{KillMe: true},
    Triggers:      []string{"input_button.test_button"},
    Periodic:      goscript.Periodics("* * * * *", ""),
    States:        goscript.Entities("input_button.test_button", "input_boolean.test_toggle", "input_number.test_number"),
    Eval:          nil,
    Func: func(t *goscript.Task) {
        gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
        t.Sleep(10 * time.Second)
        gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
    },
}
```

## Task
Within each TriggerFunc a task is available to get information from.

`task.Message` contains the message that caused the trigger to fire.

`task.Sates` contains all the states that were requested to be available by the trigger. The states are repopulated after each method is run to keep the current states fresh.

`task.Sleep(timeout)` will sleep for the specified duration.

`task.WaitUntil(entityId, eval, timeout)` Waits until the eval is true for the entity or the timeout is reached. Timeout of 0 means no timeout.

`task.While(entityId, eval, whileFunc)` Runs the whileFunc until the eval is false. task.Sleep should be used within your whileFunc to delay otherwise whileFunc will be run very quickly.

## Services
GoScript has a channel to put service calls onto. A set of default services to call is available in the [hass-ws](https://github.com/kjbreil/hass-ws) package however this is most likely not a complete list of services available in your Home Assistant installation since the service list is dynamic based on integrations installed. Generating your own service definitions is needed to interact properly with all your specific integrations.

From your personal GoScript project directory run these commands to install the service generator and run it. You must have a config.yml with the websocket credentials defined. HassWSService will generate a folder called services and the files within, make sure you do not already have a folder named services in the root of your project.
```bash
go install github.com/kjbreil/hass-ws/helpers/HassWSService@latest
go install github.com/campoy/jsonenums@latest
go install golang.org/x/tools/cmd/stringer@latest
HassWSService
go generate ./...
```

### Calling a Service
To call a service you create a service and then add options to the service. Check the source files for available options. They are not yet commented but will have comments in the future. I recommend using Home Assistant Developer Tools -> Services page to get a better understanding of what is needed for each call and to test. There is no reporting of requirements in the service definitions so be warned, some parameters are required and others are not, it is also conditional at times. For example `ClimateSetTemperature{}` needs `TargetTempHigh(float64)` and `TargetTempLow(float64)` when the mode is Heat/Cool however if the mode is Heat or the mode is Cool then `Temperature(float64)` is needed and both `TargetTempHigh(float64)` and `TargetTempLow(float64)` are ignored.
```go
service := services.NewClimateSetTemperature(services.Targets("climate.kitchen")).
		HvacMode(services.HvacModeheat_cool).
		TargetTempHigh(75).
		TargetTempLow(65)
gs.ServiceChan <- service
```

## Logging 
goscript.New is passed the config and a logger if you want to use one using the [logr](github.com/go-logr/logr) interface. There is a default `goscript.DefaultLogger()` available which will just print the logs to the terminal.

## Example
```go
package main

import (
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/kjbreil/goscript"
	"github.com/kjbreil/hass-ws/services"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := goscript.ParseConfig("config.yml", nil)
	if err != nil {
		panic(err)
	}
	
	gs, err := goscript.New(config, goscript.DefaultLogger())
	if err != nil {
		panic(err)
	}

	gs.AddTrigger(&goscript.Trigger{
		Unique:        &goscript.Unique{KillMe: true},
		Triggers:      []string{"input_button.test_button"},
		States:        goscript.Entities("input_button.test_button", "input_boolean.test_toggle", "input_number.test_number"),
		Func: func(t *goscript.Task) {
			gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
			t.Sleep(10 * time.Second)
			gs.ServiceChan <- services.NewInputBooleanToggle(services.Targets("input_boolean.test_toggle"))
		},
	})
	
	if err != nil {
		panic(err)
	}

	err = gs.Connect()

	if err != nil {
		panic(err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	gs.GetLogger().Info("Everything is set up")
	<-done

	gs.Close()
}

```