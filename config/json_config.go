package config

import (
	"encoding/json"
	"io/ioutil"

	"bytes"
	"strings"
)

func trimComments(arr []byte) []byte {
	str := string(bytes.Trim(arr, "\x00"))

	out := ""
	for _, v := range strings.Split(str, "\n") {
		if len(strings.TrimSpace(v)) >= 2 {
			if strings.TrimSpace(v)[:2] != "//" {
				out = out + v
			}
		} else {
			out = out + v
		}
		out += "\n"
	}

	return []byte(out)
}

func JSONConfig(path string) *Config {
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}

	raw = trimComments(raw)
	var conf Config

	// default parameters
	conf.Latency = false
	conf.Parallel = false
	conf.Lang = "go"
	conf.LogFile = ""

	if err := json.Unmarshal(raw, &conf); err != nil {
		panic(err)
	}

	return &conf
}
