package main

import (
	"fmt"

	"github.com/akosmarton/panull"
)

func main() {
	sink := panull.Sink{Name: "Virtual Output"}
	sink.SetProperty("device.description", "Virtual Output")

	if err := sink.Create(); err != nil {
		panic(err)
	}
	defer sink.Destroy()

	sinks, _ := panull.GetActiveSinks()
	for _, v := range sinks {
		fmt.Printf("%#v\n", v)
	}
}
