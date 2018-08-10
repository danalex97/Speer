import sys
import struct
import json

from queues import PipeQueue

class Environ( object ):
    def __init__(self):
        pipe_out = open(sys.argv[1], 'wb')
        pipe_in  = open(sys.argv[2], 'rb')

        self.queue = PipeQueue(pipe_in, pipe_out)

    def recv( self ):
        while True:
            elem        = self.queue.pop()
            tp, marshal = elem[:4], elem[4:]

            tp,     = struct.unpack('i', tp)
            marshal = marshal.decode('utf8').replace("'", '"')
            yield tp, json.loads(marshal)

if __name__ == "__main__":
    print("Python environment started...")

    env = Environ()
    tp, son = env.recv().__next__()
    print(tp, son)
