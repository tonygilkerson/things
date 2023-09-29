package main

import (
	"errors"
	"fmt"
	"reflect"
)

type MsgType string

type FooMsg struct {
	Kind     MsgType
	SenderID string
	Name     string
	FooThing string
}
type BarMsg struct {
	Kind     MsgType
	SenderID string
	Name     string
	BarThing string
}
type Bad1Msg struct {
	// Kind     MsgType
	SenderID string
	Name     string
}
type Bad2Msg struct {
	Kind     MsgType
	// SenderID string
	Name     string
}

func main() {
	foo := FooMsg{"Foo", "senderid", "the-name", "foo-thing"}

	generateMsgString(&foo)

	// fmt.Printf("type: %v\n", reflect.TypeOf(foo))
	// fmt.Printf("value: %v\n", reflect.ValueOf(foo).String())

	// v := reflect.ValueOf(&foo)
	// fmt.Println("type:", v.Type())
	// fmt.Println("kind is main.Foo:", v.Kind() == reflect.TypeOf(foo).Kind())
	// fmt.Println("value: ", v.Interface())

	// s := reflect.ValueOf(&foo).Elem()
	// typeOfFoo := s.Type()
	// for i := 0; i < s.NumField(); i++ {
	// 	f := s.Field(i)
	// 	fmt.Printf("%d: %s %s = %v\n", i, typeOfFoo.Field(i).Name, f.Type(), f.Interface())
	// }
}

func generateMsgString(msgStruct any) (string, error) {

	// Well known message fields
	var msgKind string
	var msgSenderID string
	
	// other message fields
	msgFields := make(map[string]string)

	//
	// Loop through fields
	//
	s := reflect.ValueOf(msgStruct).Elem()
	typeOf := s.Type()

	for i := 0; i < s.NumField(); i++ {

		f := s.Field(i)

		switch name := typeOf.Field(i).Name; name {
		case "Kind":
			msgKind = name
		case "SenderID":
			msgSenderID = name
		default:
			msgFields[name] = f.Interface().(string)
		}

		// // Look for Kind 
		// if typeOfFoo.Field(i).Name == "Kind" {
		// 	msgKind = typeOfFoo.Field(i).Name
		// }
		// // Look for senderID 
		// if typeOfFoo.Field(i).Name == "SenderID" {
		// 	msgSenderID = typeOfFoo.Field(i).Name
		// }
		// fmt.Printf("%d: %s %s = %v\n", i, typeOfFoo.Field(i).Name, f.Type(), f.Interface())


	}

	//
	// Check for well known types
	//
	if len(msgKind) == 0 {
		return "", errors.New("invalid message, 'Kind' field not found")
	}
	if len(msgSenderID) == 0 {
		return "", errors.New("invalid message, 'SenderID' field not found")
	}

	fmt.Printf("fields: %v", msgFields)

	return "test",nil
}