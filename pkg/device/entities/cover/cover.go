package cover

import (
	"fmt"
	"strings"

	"github.com/brutella/hap/accessory"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iancoleman/strcase"
	hassentity "github.com/kjbreil/hass-mqtt/entities"
)

type Cover struct {
	name             string
	entity           *hassentity.Cover
	hassOptions      *hassentity.CoverOptions
	homekitAccessory *accessory.WindowCovering
}

func New(name string, options ...func(*Cover)) *Cover {
	snakeName := strcase.ToSnake(name)
	readableName := strings.Title(strings.ReplaceAll(snakeName, "_", " "))

	c := &Cover{
		name: readableName,
	}

	c.hassOptions = hassentity.NewCoverOptions().Name(readableName)

	for _, option := range options {
		option(c)
	}

	var err error
	c.entity, err = hassentity.NewCover(c.hassOptions)
	if err != nil {
		return nil
	}

	return c
}

func (c *Cover) GetHassEntity() hassentity.Entity {
	return c.entity
}

func (c *Cover) GetDomainEntity() string {
	return c.entity.GetDomainEntity()
}

func (c *Cover) UpdateState() {
	c.entity.UpdateState()
}

func (c *Cover) GetHomekitAccessory() *accessory.A {
	if c.homekitAccessory != nil {
		return c.homekitAccessory.A
	}
	return nil
}

func WithHomeKit() func(*Cover) {
	return func(c *Cover) {
		c.homekitAccessory = accessory.NewWindowCovering(accessory.Info{
			Name: c.name,
		})
		// c.hassOptions.PositionFunc(func(message mqtt.Message, client mqtt.Client) {
		// 	msgPay := string(message.Payload())
		// 	fmt.Println(msgPay)
		// 	// c.entity.SetPosition("100")
		// })
		// c.hassOptions.PositionFunc(func() string {
		// 	return "100"
		// })
		c.hassOptions.SetPositionFunc(func(message mqtt.Message, client mqtt.Client) {
			msgPay := string(message.Payload())
			fmt.Println(msgPay)
			// c.entity.Position("10")
			c.entity.State("closing")
			// for i := 0; i < 100; i++ {
			// 	c.entity.Position(fmt.Sprintf("%d", i))
			// 	fmt.Println(c.entity.States.Position)
			// 	time.Sleep(time.Second)
			// }
		})
		// c.hassOptions.SetPositionFunc(func() string {
		// 	return "100"
		// })
		c.hassOptions.CommandFunc(func(message mqtt.Message, client mqtt.Client) {
			msgPay := string(message.Payload())
			fmt.Println(msgPay)

			// if string(message.Payload()) == "ON" {
			// 	c.homekitAccessory.WindowCovering.TargetPosition.SetValue()
			// } else {
			// 	c.homekitAccessory.Lightbulb.On.SetValue(false)
			// }
		})

		c.homekitAccessory.WindowCovering.TargetPosition.OnValueRemoteUpdate(func(v int) {

		})
	}
}
