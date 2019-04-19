package config

import (
	"io/ioutil"
	"strings"

	"fmt"
	"os"
	"path"
)

const (
	stubFile   = "config/stub/stub.go.tmp"
	stubScript = "config/stub/stub.go"

	defaultFile   = "config/stub/default.go.tmp"
	defaultScript = "config/stub/default.go"
)

func TemplateExists() bool {
	if _, err := os.Stat(stubScript); os.IsNotExist(err) {
		return false
	}
	return true
}

func RemoveTemplate() {
	RemoveStub()
	RewriteDefault()
}

func CreateTemplate(config *Config) {
	defer func() {
		if err := recover(); err != nil {
			RemoveStub()
			panic(err)
		}
	}()

	if config.Lang == "go" {
		if config.Entry == "" {
			panic("Entry point not provided.")
		}
		CreateStub(config.Entry)
	} else {
		panic("Lanaguage " + config.Lang + " not supported")
	}
}

func RewriteDefault() {
	defaultCode, _ := ioutil.ReadFile(defaultFile)
	ioutil.WriteFile(defaultScript, defaultCode, 0700)
}

func RemoveStub() {
	os.Remove(stubScript)
}

func CreateStub(entry string) {
	os.Remove(defaultScript)
	goPath := os.Getenv("GOPATH")

	idxl := strings.LastIndex(entry, "/")
	idx := strings.LastIndex(entry[:idxl], "/")
	entryModule := entry[:idx]
	entryStruct := entry[idxl+1:]
	entryPath := path.Join(goPath, "src", fmt.Sprintf("%s.go", entry[:idxl]))

	// read the file at entryPath
	if _, err := ioutil.ReadFile(entryPath); err != nil {
		panic("Failed finding entry point.")
	}

	// generate stub script that return the interface
	bytes, err := ioutil.ReadFile(stubFile)
	if err != nil {
		panic(err)
	}
	stubCode := string(bytes[:])
	stubCode = fmt.Sprintf(stubCode,
		fmt.Sprintf("\"%s\"", entryModule),
		entryStruct)

	// save stub code to stub module
	err = ioutil.WriteFile(stubScript, []byte(stubCode), 0700)
	if err != nil {
		panic("Failed writing stub script.")
	}
}
