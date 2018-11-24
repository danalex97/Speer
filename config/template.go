package config

import (
  "io/ioutil"
  "strings"

  "path"
  "fmt"
  "os"
)

const (
  stubFile   = "config/stub/stub.go.tmp"
  stubScript = "config/stub/stub.go"
)


func RemoveTemplate() {
  RemoveStub()
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

func RemoveStub() {
  os.Remove(stubScript)
}

func CreateStub(entry string) {
  goPath := os.Getenv("GOPATH")

  idx := strings.LastIndex(entry, "/")
  entryModule := entry[:idx]
  entryStruct := entry[idx + 1:]
  entryPath   := path.Join(goPath, "src", fmt.Sprintf("%s.go", entryModule))

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
