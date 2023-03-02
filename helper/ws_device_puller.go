package main

import (
	"fmt"
	"github.com/kjbreil/hass-ws/entities"
	"reflect"
)

var Devices = []entities.Entity{
	entities.Climate{},
}

type Device struct {
	Name       string
	Attributes map[string]attributes

	entity entities.Entity
}

func DevicesInit() (retval []Device) {

	for _, dev := range Devices {
		name := reflect.TypeOf(dev).Name()
		d := Device{
			Name:   name,
			entity: dev,
		}
		d.init()
		retval = append(retval, d)
	}
	return retval
}

func (dev *Device) init() {
	dev.Attributes = getAttributes(dev.entity)
}

func getAttributes(entity entities.Entity) map[string]attributes {

	t := reflect.TypeOf(entity)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		name := field.Name
		ty := field.Type.Name()
		if field.Type.Kind() == reflect.Ptr {
			ty = field.Type.Elem().Name()
		}
		fmt.Println(name, ty)
	}

	return nil
}

type attributes struct {
	Name     string
	DataType string
	Required bool
}
