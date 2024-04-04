package core

import (
	"github.com/kjbreil/goscript/pkg/periodic"
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-ws/services"
	"sync"
	"testing"
	"time"
)

func TestTrigger(t *testing.T) {
	var fns []TestFunc
	wg := &sync.WaitGroup{}

	fns = append(fns, func(gs *Core) error {
		var triggers []*trigger.Trigger

		// testing periodic start right away
		triggers = append(triggers, &trigger.Trigger{
			Periodic: periodic.Periodics(""),
			States:   trigger.Entities("switch.test_switch"),
			Func: func(task *trigger.Task) {
				gs.ServiceChan <- services.NewSwitchTurnOn(services.Targets("switch.test_switch"))
				task.Sleep(1 * time.Second)
				if s, ok := task.States.Get("switch.test_switch"); ok {
					if s.State != "on" {
						t.Fatal("switch.test_switch did not turn on")
					}
				} else {
					t.Fatal("state not found for switch.test_switch")
				}
				task.Sleep(1 * time.Second)
				gs.ServiceChan <- services.NewSwitchTurnOff(services.Targets("switch.test_switch"))
				task.Sleep(1 * time.Second)
				if s, ok := task.States.Get("switch.test_switch"); ok {
					if s.State != "off" {
						t.Fatal("switch.test_switch did not turn off")
					}
				} else {
					t.Fatal("state not found for switch.test_switch")
				}
				t.Logf("Periodic Ran, Turn On/Turn Off successful")
				wg.Done()
			},
		})

		// test triggering
		triggers = append(triggers, &trigger.Trigger{
			Triggers: trigger.Entities("switch.test_switch"),
			Func: func(task *trigger.Task) {
				if task.Message.State() == "on" {
					wg.Done()
				}
			},
		})

		// test default unique behavior
		defaultUniqueRuns := 0
		defaultUniqueWG := false
		triggers = append(triggers, &trigger.Trigger{
			Triggers: trigger.Entities("switch.test_switch"),
			Unique:   &trigger.Unique{},
			Func: func(task *trigger.Task) {
				defaultUniqueRuns++
				task.Sleep(4 * time.Second)
				if defaultUniqueRuns < 2 {
					t.Fatal("first unique run not killed")
				}
				if !defaultUniqueWG {
					defaultUniqueWG = true
					t.Logf("default unique behavior worked")
					wg.Done()
				}
			},
		})

		// test unique wait
		waitUniqueStarts := 0
		waitUniqueEnds := 0
		triggers = append(triggers, &trigger.Trigger{
			// Triggers: Entities("switch.test_switch"),
			Periodic: periodic.Periodics("*/1 * * * * *"),
			Unique:   &trigger.Unique{Wait: true},
			Func: func(task *trigger.Task) {
				if waitUniqueEnds > 5 {
					return
				}

				waitUniqueStarts++
				task.Sleep(2 * time.Second)
				waitUniqueEnds++
				if waitUniqueStarts != waitUniqueEnds {
					wg.Done()
					t.Fatal("first unique killed when it should not have been")
				}
				if waitUniqueEnds > 5 {
					t.Logf("unique wait successfully waited before running next task")
					wg.Done()
				}
			},
		})

		// test unique killme
		killMeStarts := 0
		var killMeStartTime time.Time
		triggers = append(triggers, &trigger.Trigger{
			// Triggers: Entities("switch.test_switch"),
			Periodic: periodic.Periodics("*/1 * * * * *"),
			Unique:   &trigger.Unique{Wait: true},
			Func: func(task *trigger.Task) {

				if time.Now().Sub(killMeStartTime) < time.Second*5 {
					t.Fatalf("kill me did not work")
				}
				if killMeStartTime.IsZero() {
					killMeStartTime = time.Now()
				}

				if killMeStarts > 1 {
					return
				}

				killMeStarts++
				task.Sleep(5 * time.Second)
				if killMeStarts == 1 {
					t.Logf("unique kill me succesfully killed subsequent runs")
					wg.Done()
				}
			},
		})

		gs.AddTriggers(triggers...)
		wg.Add(len(triggers))

		return nil
	})

	GoScriptTestRun(fns, nil, wg, t)

}
