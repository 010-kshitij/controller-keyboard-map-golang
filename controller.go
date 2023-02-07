package main

import (
	"fmt"
	"log"
	"github.com/google/gousb"
	mapset "github.com/deckarep/golang-set"
)

const (
	CONTROLLER_DEFAULT = "DEFAULT"
)

//controller references on byte index and button
var buttonByte2 map[int]string = map[int]string{
	0 : "UP",
	1 : "DOWN",
	2 : "LEFT",
	3 : "RIGHT",
	4 : "SELECT",
	5 : "BACK",
	6 : "L3",
	7 : "R3",
}

var buttonByte3 map[int]string = map[int]string{
	0 : "LB",
	1 : "RB",
	2 : "HOME",
	3 : "",
	4 : "A",
	5 : "B",
	6 : "X",
	7 : "Y",
}

// end controller references on byte index and button

type Controller struct {
	Device *gousb.Device
	Config *gousb.Config
	Interface *gousb.Interface
	InEndpoint *gousb.InEndpoint
	InBuffer []byte
	CurrentKeySet mapset.Set
	OrderedKeysPressed []string
}

func OpenController(context *gousb.Context, vendorId gousb.ID, productId gousb.ID) (Controller, error) {
	device, err := context.OpenDeviceWithVIDPID(vendorId, productId)
	defer device.Close()
	var c Controller
	c.Device = device

	// Setting up the config
	config, err := device.Config(1)
	c.Config = config

	// Setting up the Interface
	intf, err := config.Interface(0, 0)
	c.Interface = intf

	// Open IN endpoint
	epIn, err := intf.InEndpoint(0x81)
	if err != nil {
		log.Fatalf("%s.InEndpoint(1): %v", intf, err)
	}
	c.InEndpoint = epIn

	buf := make([]byte, 10*epIn.Desc.MaxPacketSize)
	c.InBuffer = buf

	orderedKeys := make([]string, 0)
	c.OrderedKeysPressed = orderedKeys
	c.CurrentKeySet = mapset.NewSet()

	return c, err
}

func (c *Controller) Close() error {
	fmt.Println("Closing Controller")
	// TODO: fire keyup event in remaining elements and empty the key set
	err := c.Config.Close()
	err = c.Device.Close()
	return err
}

func (c *Controller) CheckButtons(byteIndex int, mapped map[int]string) []string {
	var byteAtIndex byte = c.InBuffer[byteIndex]
	var binaryByte string
	binaryByte = fmt.Sprintf("%8b", int(byteAtIndex) )
	pressedKeys := getKeysOnIndexes(binaryByte, mapped)

	return pressedKeys
}

func (c *Controller) GetReadableInput() string {
	c.OrderedKeysPressed = make([]string, 0)
	// Reading the data from the device
	readBytes, err := c.InEndpoint.Read(c.InBuffer)
	if err != nil {
		fmt.Println("Read returned an error:", err)
	}
	if readBytes == 0 {
		log.Fatalf("IN endpoint 6 returned 0 bytes of data.")
	}

	//fmt.Println(c.InBuffer)

	c.OrderedKeysPressed = append(c.OrderedKeysPressed, c.CheckButtons(2, buttonByte2)...)
	c.OrderedKeysPressed = append(c.OrderedKeysPressed, c.CheckButtons(3, buttonByte3)...)

	//NOTE: keypress events here
	fmt.Println("Current Key Set: ", c.CurrentKeySet)
	if(c.CurrentKeySet.Cardinality() == 0) {
		for _, button := range c.OrderedKeysPressed {
			c.CurrentKeySet.Add(button)
			SendKeyDownEvent(button)
		}
		// keydownevent here only
	} else {
		in_set := mapset.NewSet()
		for _, button := range c.OrderedKeysPressed {
			in_set.Add(button)
		}
		fmt.Println("in set: ", in_set)

		diff_key_set_in_set := c.CurrentKeySet.Difference(in_set)
		fmt.Println("diff Key Set in set: ", diff_key_set_in_set)
		// keyup event here
		it := diff_key_set_in_set.Iterator()

		for button := range it.C {
			SendKeyUpEvent(button.(string))
		}

		inter_set := c.CurrentKeySet.Intersect(in_set)
		fmt.Println("inter set: ", inter_set)
		diff_in_set_inter_set := in_set.Difference(inter_set)
		fmt.Println("diff in set inter set: ", diff_in_set_inter_set)
		// keydown event here
		it = diff_key_set_in_set.Iterator()

		for button := range it.C {
			SendKeyDownEvent(button.(string))
		}

		c.CurrentKeySet = in_set

	}
	// end keypress events here
	
	if c.InBuffer[3] == 4 {
		return "QUIT"
	}

	return CONTROLLER_DEFAULT
}
