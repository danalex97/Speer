import sys
import struct
import json

from queues    import PipeQueue
from scheduler import Scheduler
from messages  import MESSAGES

class Environ( object ):
    def __init__(self):
        pipe_out = open(sys.argv[1], 'wb')
        pipe_in  = open(sys.argv[2], 'rb')

        self.queue     = PipeQueue(pipe_in, pipe_out)
        self.scheduler = Scheduler()

    def recv( self ):
        while True:
            elem        = self.queue.pop()
            tp, marshal = elem[:4], elem[4:]

            tp,     = struct.unpack('i', tp)
            marshal = marshal.decode('utf8').replace("'", '"')

            son = json.loads(marshal)

            if tp not in MESSAGES:
                print('Message id not unrecognized: {}'.format(tp))
                continue

            message = MESSAGES[tp]().from_json(son)
            yield message

if __name__ == "__main__":
    print("Python environment started...")

    env = Environ()
    message = env.recv().__next__()
    print(message)
