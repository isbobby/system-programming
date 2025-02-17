package main

import (
	"fmt"
	"reflect"
)

type Bob struct {
	Name     string `lang:"en"`
	Age      int
	Height   float64
	HousePet Pet
}

type Pet struct {
	Name    string
	Species string
}

func InspectStruct(x interface{}) {
	val := reflect.ValueOf(x)
	t := reflect.TypeOf(x)

	if val.Kind() != reflect.Struct {
		fmt.Println("Not a struct!")
		return
	}

	fmt.Println("Type:", t.Name())
	fmt.Println("Fields:", t.Name())

	for i := 0; i < val.NumField(); i++ {
		field := t.Field(i)
		value := val.Field(i)

		fmt.Printf("\tName: %s, Type: %s, Value: %v, Tag: %s\n", field.Name, field.Type, value, field.Tag)
	}
}
