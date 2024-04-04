package core

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kjbreil/goscript/pkg/trigger"
)

// SubscribeMqtt subscribes to a top with a trigger
func (gs *Core) SubscribeMqtt(topic string, qos byte, tr *trigger.Trigger) {
	tr = trigger.SetupTrigger(tr)

	gs.mqtt.Subscribe(topic, qos, func(client mqtt.Client, message mqtt.Message) {
		task := gs.TrigRunner.NewTask(tr, nil)
		task.MqttMessage = message
		gs.TrigRunner.AddTask(task)
	})
}

// TaskMQTT wraps a trigger and TaskFunc setting up and passing the task through
func (gs *Core) TaskMQTT(tr *trigger.Trigger) func(message mqtt.Message, client mqtt.Client) {
	// setup the trigger
	tr = trigger.SetupTrigger(tr)

	return func(message mqtt.Message, client mqtt.Client) {
		task := gs.TrigRunner.NewTask(tr, nil)
		task.MqttMessage = message
		gs.TrigRunner.AddTask(task)
	}
}
