package panull

import (
	"errors"
	"os/exec"
	"strings"
)

func getModulesList() ([]string, error) {
	cmd := exec.Command("pactl", "list", "short", "modules")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out))
	}
	return strings.Split(string(out), "\n"), nil
}

func parseArguments(s string, quote rune) map[string]string {
	res := make(map[string]string)
	var key, value string
	var inKey, inValue, inQuoute bool
	inKey = true
	for _, r := range s {
		switch r {
		case '=':
			if inQuoute {
				value += string(r)
			} else {
				inValue = true
				inKey = false
			}
		case quote:
			if inQuoute {
				inQuoute = false
			} else {
				inQuoute = true
			}
		case ' ':
			if inQuoute {
				value += string(r)
			} else {
				inValue = false
				inKey = true
				res[key] = value
				key = ""
				value = ""
			}
		default:
			if inKey {
				key += string(r)
			} else if inValue {
				value += string(r)
			}
		}
	}
	res[key] = value
	return res
}
