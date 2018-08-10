package bridge

import (
  "encoding/binary"
  "unsafe"
  "os"
)

type RawQueue interface {
  Push(data []byte)  bool
  Pop(size int)      ([]byte, bool)
}

type RawShmQueue struct {
  base uintptr
  pos  uintptr

  length   int
  size     int
}

// Shared memory queue
func NewRawShmQueue(base uintptr, size int) RawQueue {
  return &RawShmQueue{
    base     : base,
    pos      : base,
    length   : 0,
    size     : size,
  }
}

func (q *RawShmQueue) idx(base uintptr, pls int) uintptr {
  p := base + uintptr(pls)
  if p > q.base + uintptr(q.size) {
    return p - uintptr(q.size)
  }
  return p
}

func (q *RawShmQueue) ptr(base uintptr, pls int) unsafe.Pointer {
  return unsafe.Pointer(q.idx(base, pls))
}

func (q *RawShmQueue) Push(elem []byte) bool {
  if q.length + len(elem) > q.size {
    return false
  }

  for i := 0; i < len(elem); i++ {
    *(*byte)(q.ptr(q.pos, q.length + i)) = elem[i]
  }
  q.length += len(elem)

  return true
}

func (q *RawShmQueue) Pop(size int) ([]byte, bool){
  if q.size - q.length < size {
    return nil, false
  }

  ans := []byte{}
  for i := 0; i < size; i++ {
    ans = append(ans, *(*byte)(q.ptr(q.pos, i)))
  }
  q.pos     = q.idx(q.pos, size)
  q.length -= 1

  return ans, true
}

// Pipe-based queue
type RawPipeQueue struct {
  in  *os.File
  out *os.File
}

func NewRawPipeQueue(in, out *os.File) RawQueue {
  return &RawPipeQueue{
    in  : in,
    out : out,
  }
}

func (q *RawPipeQueue) Push(elem []byte) bool {
  q.out.Write(elem)
  return true
}

func (q *RawPipeQueue) Pop(size int) ([]byte, bool){
  ans := make([]byte, size)
  q.in.Read(ans)
  return ans, true
}

type IPCQueue struct {
  RawQueue
}

func NewIPCQueue(queue RawQueue) *IPCQueue {
  return &IPCQueue{
    RawQueue : queue,
  }
}

func (q *IPCQueue) Push(data []byte) bool {
  sizeBytes := make([]byte, 4)
  binary.LittleEndian.PutUint32(sizeBytes, uint32(len(data)))

  data = append(sizeBytes, data...)
  return q.RawQueue.Push(data)
}

func (q *IPCQueue) Pop() ([]byte, bool) {
  raw_size, ok := q.RawQueue.Pop(4)
  if !ok {
    return nil, false
  }
  size := int(binary.LittleEndian.Uint32(raw_size))

  return q.RawQueue.Pop(size)
}
