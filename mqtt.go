package goscript

import mqtt "github.com/eclipse/paho.mqtt.golang"

// SubscribeMqtt subscribes to a top with a trigger
func (gs *GoScript) SubscribeMqtt(topic string, qos byte, tr *Trigger) {
	tr = setupTrigger(tr)

	gs.mqtt.Subscribe(topic, qos, func(client mqtt.Client, message mqtt.Message) {
		task := gs.newTask(tr, nil)
		task.MqttMessage = message
		gs.taskToRun.add(task)
	})
}

// TaskMQTT wraps a trigger and TaskFunc setting up and passing the task through
func (gs *GoScript) TaskMQTT(tr *Trigger) func(message mqtt.Message, client mqtt.Client) {
	// setup the trigger
	tr = setupTrigger(tr)

	return func(message mqtt.Message, client mqtt.Client) {
		task := gs.newTask(tr, nil)
		task.MqttMessage = message
		gs.taskToRun.add(task)
	}
}
