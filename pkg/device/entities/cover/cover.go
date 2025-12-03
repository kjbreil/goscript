package cover

import (
	"fmt"
	"strings"

	"github.com/brutella/hap/accessory"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iancoleman/strcase"
	hassentity "github.com/kjbreil/hass-mqtt/entities"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Cover represents a cover device entity with HASS and HomeKit integration.
type Cover struct {
	name             string
	entity           *hassentity.Cover
	hassOptions      *hassentity.CoverOptions
	homekitAccessory *accessory.WindowCovering
}

// New creates a new Cover entity with the given name and options.
func New(name string, options ...func(*Cover)) *Cover {
	snakeName := strcase.ToSnake(name)
	caser := cases.Title(language.English)
	readableName := caser.String(strings.ReplaceAll(snakeName, "_", " "))

	//nolint:exhaustruct // entity, hassOptions, homekitAccessory initialized below
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

// GetHassEntity returns the HASS entity for this cover.
func (c *Cover) GetHassEntity() hassentity.Entity {
	return c.entity
}

// GetDomainEntity returns the domain entity string for this cover.
func (c *Cover) GetDomainEntity() string {
	return c.entity.GetDomainEntity()
}

// UpdateState updates the state of the cover entity.
func (c *Cover) UpdateState() {
	c.entity.UpdateState()
}

// GetHomekitAccessory returns the HomeKit accessory for this cover.
func (c *Cover) GetHomekitAccessory() *accessory.A {
	if c.homekitAccessory != nil {
		return c.homekitAccessory.A
	}
	return nil
}

// WithHomeKit returns an option function that adds HomeKit support to a cover.
func WithHomeKit() func(*Cover) {
	return func(c *Cover) {
		//nolint:exhaustruct // Only Name required; other Info fields optional
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
		c.hassOptions.SetPositionFunc(func(message mqtt.Message, _ mqtt.Client) {
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
		c.hassOptions.CommandFunc(func(message mqtt.Message, _ mqtt.Client) {
			msgPay := string(message.Payload())
			fmt.Println(msgPay)

			// if string(message.Payload()) == "ON" {
			// 	c.homekitAccessory.WindowCovering.TargetPosition.SetValue()
			// } else {
			// 	c.homekitAccessory.Lightbulb.On.SetValue(false)
			// }
		})

		c.homekitAccessory.WindowCovering.TargetPosition.OnValueRemoteUpdate(func(_ int) {

		})
	}
}
