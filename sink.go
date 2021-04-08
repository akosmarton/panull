package panull

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Sink struct {
	Name     string
	Format   string
	Rate     int
	Channels int
	// ChannelMap
	UseSystemClockForTiming bool
	properties              map[string]interface{}
	moduleIndex             int
}

func (s *Sink) Create() error {
	var err error

	args := make([]string, 0)
	args = append(args, "load-module")
	args = append(args, "module-null-sink")
	if s.Name != "" {
		args = append(args, fmt.Sprintf("sink_name=\"%s\"", s.Name))
	}
	if s.Format != "" {
		args = append(args, fmt.Sprintf("format=%s", s.Format))
	}
	if s.Rate > 0 {
		args = append(args, fmt.Sprintf("rate=%d", s.Rate))
	}
	if s.Channels > 0 {
		args = append(args, fmt.Sprintf("channels=%d", s.Channels))
	}

	if s.UseSystemClockForTiming {
		args = append(args, "use_system_clock_for_timing=yes")
	}

	var props string

	for k, v := range s.properties {
		kv := fmt.Sprintf("%s=%v", k, v)
		props = props + strings.ReplaceAll(kv, " ", "_") + " "
	}

	props = strings.TrimSpace(props)

	args = append(args, fmt.Sprintf("sink_properties=\"%s\"", props))

	cmd := exec.Command("pactl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}
	if _, err := fmt.Sscanf(string(out), "%d", &s.moduleIndex); err != nil {
		return err
	}

	return nil
}

func (s *Sink) Destroy() error {
	args := make([]string, 0)
	args = append(args, "unload-module")
	args = append(args, fmt.Sprintf("%d", s.moduleIndex))

	cmd := exec.Command("pactl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return errors.New(string(out))
	}

	return nil
}

func (s *Sink) SetProperty(key string, value interface{}) *Sink {
	if s.properties == nil {
		s.properties = make(map[string]interface{})
	}

	s.properties[key] = value
	return s
}

func (s *Sink) GetProperty(key string) interface{} {
	if s.properties == nil {
		return nil
	}

	return s.properties[key]
}

func GetActiveSinks() ([]*Sink, error) {
	sinks := make([]*Sink, 0)
	ls, err := getModulesList()
	if err != nil {
		return nil, err
	}
	for _, l := range ls {
		ss := strings.Split(l, "\t")
		if len(ss) < 2 {
			continue
		}
		sink := &Sink{}
		sink.moduleIndex, _ = strconv.Atoi(ss[0])
		if ss[1] != "module-null-sink" {
			continue
		}
		if len(ss) > 2 {
			for k, v := range parseArguments(ss[2], '"') {
				switch k {
				case "sink_name":
					sink.Name = v
				case "format":
					sink.Format = v
				case "rate":
					sink.Rate, _ = strconv.Atoi(v)
				case "channels":
					sink.Channels, _ = strconv.Atoi(v)
				case "use_system_clock_for_timing":
					if v == "yes" {
						sink.UseSystemClockForTiming = true
					}
				case "sink_properties":
					for k, v := range parseArguments(v, '\'') {
						sink.SetProperty(k, v)
					}
				}
			}
		}
		sinks = append(sinks, sink)
	}

	return sinks, nil
}
