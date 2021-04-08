# panull
Pulseaudio client library in Golang for creating null sinks and sources
(It is compatible also with Pipewire)

## Usage

```go
package main

import (
	"github.com/akosmarton/panull"
)

func main() {
	source := panull.Source{
		Properties: map[string]interface{}{
			"device.description": "Virtual Input",
		},
	}
	source.Open()
	defer source.Close()

	sink := panull.Sink{
		Properties: map[string]interface{}{
			"device.description": "Virtual Output",
		},
	}
	sink.Open()
	defer sink.Close()

	p := make([]byte, 0)

	source.Write(p)
	sink.Read(p)
}
```
