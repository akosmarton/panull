package panull

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Source struct {
	Name     string
	Format   string
	Rate     int
	Channels int
	// channelMap
	properties  map[string]interface{}
	moduleIndex int
}

func (s *Source) Create() error {
	var err error

	args := make([]string, 0)
	args = append(args, "load-module")
	args = append(args, "module-null-source")
	if s.Name != "" {
		args = append(args, fmt.Sprintf("source_name=\"%s\"", s.Name))
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

	var props string

	for k, v := range s.properties {
		kv := fmt.Sprintf("%s=%v", k, v)
		kv = strings.ReplaceAll(kv, "=", "_")
		props = props + strings.ReplaceAll(kv, " ", "_") + " "
	}

	props = strings.TrimSpace(props)

	args = append(args, fmt.Sprintf("source_properties=\"%s\"", props))

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

func (s *Source) Destroy() error {
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

func (s *Source) SetProperty(key string, value interface{}) *Source {
	if s.properties == nil {
		s.properties = make(map[string]interface{})
	}

	s.properties[key] = value
	return s
}

func (s *Source) GetProperty(key string) interface{} {
	if s.properties == nil {
		return nil
	}

	return s.properties[key]
}

func GetActiveSources() ([]*Source, error) {
	sources := make([]*Source, 0)
	ls, err := getModulesList()
	if err != nil {
		return nil, err
	}
	for _, l := range ls {
		ss := strings.Split(l, "\t")
		if len(ss) < 2 {
			continue
		}
		source := &Source{}
		source.moduleIndex, _ = strconv.Atoi(ss[0])
		if ss[1] != "module-null-source" {
			continue
		}
		if len(ss) > 2 {
			for k, v := range parseArguments(ss[2], '"') {
				switch k {
				case "source_name":
					source.Name = v
				case "format":
					source.Format = v
				case "rate":
					source.Rate, _ = strconv.Atoi(v)
				case "channels":
					source.Channels, _ = strconv.Atoi(v)
				case "source_properties":
					for k, v := range parseArguments(v, '\'') {
						source.SetProperty(k, v)
					}
				}
			}
		}
		sources = append(sources, source)
	}

	return sources, nil
}
