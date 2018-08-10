package bridge

import (
  "os"
  "os/exec"
  "syscall"
  "strings"
  "encoding/json"
  "fmt"
)

type Environ struct {
  name  string
  cmd   *exec.Cmd

  queue *IPCQueue
}

func NewEnviron(name string , command ...string) *Environ {
  processed_command := []string{}
  for _, cur := range command {
    processed_command = append(processed_command, strings.Split(cur, " ")...)
  }
  proc := processed_command[0]
  args := processed_command[1:]

  e := new(Environ)

  in_name  := "/tmp/" + name + "_out_env_pipe"
  out_name := "/tmp/" + name + "_in_env_pipe"

  syscall.Mkfifo(in_name, 0600)
  syscall.Mkfifo(out_name, 0600)
  in, _  := os.OpenFile(in_name, os.O_RDWR|syscall.O_NONBLOCK, 0600)
  out, _ := os.OpenFile(out_name, os.O_RDWR|syscall.O_NONBLOCK, 0600)
  args = append(args, in_name, out_name)

  e.name  = name
  e.cmd   = exec.Command(proc, args...)
  e.queue = NewIPCQueue(NewRawPipeQueue(in, out))

  e.cmd.Stdin = os.Stdin;
  e.cmd.Stdout = os.Stdout;
  e.cmd.Stderr = os.Stderr;

  return e
}

func (e *Environ) Start() {
  e.cmd.Start()

  fmt.Println("Do stuff...")

  x, _ := json.Marshal("Valoare")

  e.queue.Push(x)
  v, _ := e.queue.Pop()
  fmt.Println(string(v[:]))

  e.cmd.Wait()
}
