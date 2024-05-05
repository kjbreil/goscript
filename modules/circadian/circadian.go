package circadian

import (
	"fmt"
	"github.com/kjbreil/goscript/light"
	"github.com/kjbreil/goscript/pkg/device"
	"github.com/kjbreil/goscript/pkg/module"
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/sixdouglas/suncalc"
	"math"
	"time"
)

var key = "circadian"

type Circadian struct {
	Lat              float64
	Long             float64
	MinTemperature   float64
	MaxTemperature   float64
	MinBrightnessPct float64
	MaxBrightnessPct float64

	Transition float64

	currentTemperature   float64
	currentBrightnessPct float64

	module.Base
}

func (c *Circadian) Run() error {
	if c.Lat == 0 && c.Long == 0 {
		return fmt.Errorf("circadian requires lat and long")
	}
	return nil
}

func (c *Circadian) Triggers() trigger.Triggers {
	return nil

}

func (c *Circadian) Devices() device.Devices {
	return nil
}

func (c *Circadian) Update() error {
	return nil
}

func (c *Circadian) Name() string {
	return key
}

func (c *Circadian) Close() error {
	return nil

}

func (c *Circadian) ChangeEnough(temperature, brightness float64) bool {
	brightnessStepSize := (c.MaxBrightnessPct - c.MinBrightnessPct) / 100
	temperatureStepSize := (c.MaxTemperature - c.MinTemperature) / 100

	return math.Abs(temperature-c.currentTemperature) <= temperatureStepSize &&
		math.Abs(brightness-c.currentBrightnessPct) <= brightnessStepSize

}

func (c *Circadian) Calculate() (float64, float64) {
	now := time.Now()
	times := suncalc.GetTimes(now, c.Lat, c.Long)

	currentPos := suncalc.GetPosition(now, c.Lat, c.Long)
	noonPos := suncalc.GetPosition(times[suncalc.SolarNoon].Value, c.Lat, c.Long)

	// l.gs.CallService(services.NewInputBooleanToggle([]string{"input_boolean.test_toggle"}))

	// azPct := (currentPos.Azimuth * 180 / math.Pi) / (noonPos.Azimuth * 180 / math.Pi)
	// altPct := (currentPos.Altitude / noonPos.Altitude) * 100
	altPct := math.Round((currentPos.Altitude*180/math.Pi)/(noonPos.Altitude*180/math.Pi)*10000) / 100

	c.currentTemperature = mapRange(-80, 100, c.MaxTemperature, c.MinTemperature, altPct)
	c.currentBrightnessPct = mapRange(-50, 100, c.MinBrightnessPct, c.MaxBrightnessPct, altPct)
	c.within()

	return c.currentTemperature, c.currentBrightnessPct
}

func (c *Circadian) Temperature() float64 {
	return c.currentTemperature
}

func (c *Circadian) BrightnessPct() float64 {
	return c.currentBrightnessPct
}

func (c *Circadian) TurnOn(t *trigger.Task, entities ...string) {
	if len(entities) == 0 {
		return
	}
	c.Calculate()

	light.New().
		ColorTemp(c.Temperature()).
		BrightnessPct(c.BrightnessPct()).
		Transition(c.Transition).
		TurnOn(t, entities)

}

func (c *Circadian) TurnOnTemperature(t *trigger.Task, entities ...string) {
	if len(entities) == 0 {
		return
	}
	c.Calculate()

	light.New().
		ColorTemp(c.Temperature()).
		Transition(c.Transition).
		TurnOn(t, entities)
}

func (c *Circadian) TurnOnTemperatureManualBrightness(t *trigger.Task, brightness float64, entities ...string) {
	if len(entities) == 0 {
		return
	}
	c.Calculate()

	light.New().
		ColorTemp(c.Temperature()).
		BrightnessPct(brightness).
		Transition(c.Transition).
		TurnOn(t, entities)
}

func (c *Circadian) TurnOff(t *trigger.Task, entities ...string) {
	light.New().Transition(c.Transition).TurnOff(t, entities)
}

func mapRange(rangeLow, rangeHigh, mapLow, mapHigh, value float64) float64 {
	return mapLow + ((value - rangeLow) / (rangeHigh - (rangeLow)) * (mapHigh - (mapLow)))
}

func (c *Circadian) within() {
	if c.currentTemperature > c.MaxTemperature {
		c.currentTemperature = c.MaxTemperature
	}
	if c.currentTemperature < c.MinTemperature {
		c.currentTemperature = c.MinTemperature
	}
	if c.currentBrightnessPct > c.MaxBrightnessPct {
		c.currentBrightnessPct = c.MaxBrightnessPct
	}
	if c.currentBrightnessPct < c.MinBrightnessPct {
		c.currentBrightnessPct = c.MinBrightnessPct
	}
}
