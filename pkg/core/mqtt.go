package core

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kjbreil/goscript/pkg/trigger"
)

// SubscribeMqtt subscribes to a top with a trigger.
func (gs *GoScript) SubscribeMqtt(topic string, qos byte, tr *trigger.Trigger) {
	tr = trigger.SetupTrigger(tr)

	gs.mqtt.Subscribe(topic, qos, func(client mqtt.Client, message mqtt.Message) {
		task := gs.Runner.NewTask(tr, nil)
		task.MqttMessage = message
		gs.Runner.AddTask(task)
	})
}

// TaskMQTT wraps a trigger and TaskFunc setting up and passing the task through.
func (gs *GoScript) TaskMQTT(tr *trigger.Trigger) func(message mqtt.Message, client mqtt.Client) {
	// setup the trigger
	tr = trigger.SetupTrigger(tr)

	return func(message mqtt.Message, client mqtt.Client) {
		task := gs.Runner.NewTask(tr, nil)
		task.MqttMessage = message
		gs.Runner.AddTask(task)
	}
}
