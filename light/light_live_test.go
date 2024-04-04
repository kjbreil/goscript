package light

import (
	"github.com/kjbreil/goscript/pkg/core"
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-mqtt/entities"
	"sync"
	"testing"
	"time"
)

func TestLightsInDevice(t *testing.T) {

	var preFns []core.TestFunc
	var postFns []core.TestFunc
	wg := &sync.WaitGroup{}

	preFns = append(preFns, func(gs *core.Core) error {
		device, err := gs.GetDevice("test_devices")
		if err != nil {
			return err
		}

		var entitiesToAdd []entities.Entity

		esp := []string{
			"light.light_1",
			"light.light_2",
		}
		// allUuid := uuid.New()

		lightOneOptions := entities.NewLightOptions()
		lightOneOptions.Name("Light 1").
			CommandFunc(gs.TaskMQTT(&trigger.Trigger{
				States: esp,
				Unique: &trigger.Unique{
					Wait: true,
					// UUID: &allUuid,
				},
				Func: func(task *trigger.Task) {
					t.Logf("LightOne Triggered")
					switch lightOneOptions.States().State {
					case "ON":
						t.Logf("LightOne Turned On")
						New().TurnOn(task, []string{"light.light_2"})
					case "OFF":
						t.Logf("LightOne Turned Off")
						New().TurnOff(task, []string{"light.light_2"})
					}
					t.Logf("LightOne Exited")
				},
			}))

		lightOneDevice, err := entities.NewLight(lightOneOptions)
		if err != nil {
			return err
		}
		entitiesToAdd = append(entitiesToAdd, lightOneDevice)

		lightTwoOptions := entities.NewLightOptions()
		lightTwoOptions.Name("Light 2").
			CommandFunc(gs.TaskMQTT(&trigger.Trigger{
				States: esp,
				Unique: &trigger.Unique{
					Wait: true,
					// UUID: &allUuid,
				},
				Func: func(task *trigger.Task) {
					t.Logf("LightTwo Triggered")
					switch lightTwoOptions.States().State {
					case "ON":
						t.Logf("LightTwo Turned On")
						New().TurnOn(task, []string{"light.light_3"})
					case "OFF":
						t.Logf("LightTwo Turned Off")
						New().TurnOff(task, []string{"light.light_3"})
					}

					t.Logf("LightTwo Exited")
				},
			}))

		lightTwoDevice, err := entities.NewLight(lightTwoOptions)
		if err != nil {
			return err
		}
		entitiesToAdd = append(entitiesToAdd, lightTwoDevice)

		lightThreeOptions := entities.NewLightOptions()
		lightThreeOptions.Name("Light 3").
			CommandFunc(gs.TaskMQTT(&trigger.Trigger{
				States: esp,
				Unique: &trigger.Unique{
					Wait: true,
					// UUID: &allUuid,
				},
				Func: func(task *trigger.Task) {
					t.Logf("LightThree Triggered")
					switch lightThreeOptions.States().State {
					case "ON":
						t.Logf("LightThree Turned On")
					case "OFF":
						t.Logf("LightThree Turned Off")
					}
					t.Logf("LightThree Exited")
				},
			}))

		lightThreeDevice, err := entities.NewLight(lightThreeOptions)
		if err != nil {
			return err
		}
		entitiesToAdd = append(entitiesToAdd, lightThreeDevice)

		device.AddEntities(entitiesToAdd)

		return nil
	})

	postFns = append(postFns, func(gs *core.Core) error {
		time.Sleep(5 * time.Minute)

		return nil
	})

	core.GoScriptTestRun(preFns, postFns, wg, t)
}
