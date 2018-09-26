package bridge

import (
  "encoding/json"
  "encoding/binary"

  "os"
  "os/exec"
  "syscall"
  "strings"
  "fmt"
  "log"
)

func GetEnviron() *Environ {
  return nil
}

type Environ struct {
  name  string
  cmd   *exec.Cmd

  nodes map[string]*BridgedTorrent

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
  e.nodes = make(map[string]*BridgedTorrent)

  e.cmd.Stdin = os.Stdin;
  e.cmd.Stdout = os.Stdout;
  e.cmd.Stderr = os.Stderr;

  return e
}

func (e *Environ) SendMessage(msg Message) {
  tp := MessageToType(msg)

  typeBytes := make([]byte, 4)
  binary.LittleEndian.PutUint32(typeBytes, uint32(tp))

  data, _ := json.Marshal(msg)
  fmt.Println(data)
  data     = append(typeBytes, data...)

  e.queue.Push(data)
}

func (e *Environ) ListenIncoming() {
  for {
    data, _ := e.queue.Pop()
    tp      := binary.LittleEndian.Uint32(data[0:4])

    msg := TypeToMessage(int(tp))
    json.Unmarshal(data, msg)

    e.nodes[msg.GetId()].envChannel <- msg
  }
}

func (e *Environ) Start() {
  e.cmd.Start()

  e.SendMessage(&Create{Id : "hello"})
  v, _ := e.queue.Pop()
  fmt.Println(string(v[:]))

  e.cmd.Wait()
}

func (e *Environ) Stop() {
  if err := e.cmd.Process.Kill(); err != nil {
    log.Fatal("Failed to kill process: ", err)
  }
}
