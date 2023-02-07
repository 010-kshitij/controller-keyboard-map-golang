package main

import (
	"fmt"
	"log"
	"time"
	"os"
	"os/signal"

	"github.com/google/gousb"
)

const (
  product = 0x028e
  vendor = 0x045e
)

func process(controller *Controller) {
	for {
		controller.GetReadableInput()
		time.Sleep(time.Millisecond * 60)
	}
}

func main() {
	// Creating the usb context, always has to be closed
	context := gousb.NewContext()
	defer context.Close()

	// Open controller of specified vendor and product id
	controller, err := OpenController(context, vendor, product)
	if err != nil {
		log.Fatal("Error in initializing controller: ", err)
	}
	defer controller.Close()

	fmt.Println("Controller found: ")
	fmt.Println(controller)

	fmt.Println("Starting Process")
	go process(&controller)
	
	c := make(chan os.Signal, 1)
	var done bool
	done = false
	signal.Notify(c, os.Interrupt)
	go func(){
			for sig := range c {
					fmt.Println("Signal")
					fmt.Println(sig)
					fmt.Println("End Signal")
					done = true
			}
	}()

	for {
		if done == true {
			break
		}
	}

}
