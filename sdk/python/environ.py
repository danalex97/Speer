import sys
import struct
import json

from queues    import PipeQueue
from manager   import NodeManager
from messages  import MAP_ID_MESSAGE, MESSAGES

from scheduler import Scheduler
from scheduler_api import _schedule

class Environ( object ):
    def __init__(self):
        pipe_out = open(sys.argv[1], 'wb')
        pipe_in  = open(sys.argv[2], 'rb')

        self.queue     = PipeQueue(pipe_in, pipe_out)
        self.manager   = NodeManager()

    def recv( self ):
        while True:
            elem        = yield from self.queue.pop()
            tp, marshal = elem[:4], elem[4:]

            tp,     = struct.unpack('i', tp)
            marshal = marshal.decode('utf8').replace("'", '"')

            son = json.loads(marshal)

            if tp not in MAP_ID_MESSAGE:
                print('Message id not unrecognized: {}'.format(tp))
                continue

            message = MAP_ID_MESSAGE[tp]().from_json(son)
            yield message

    def run( self ):
        for message in self.recv():
            if isinstance(message, MESSAGES.Create):
                self.manager.create(message)
            yield from _schedule()

if __name__ == "__main__":
    print("Python environment started...")

    scheduler = Scheduler()
    env = Environ()

    scheduler.execute(env.run)

    scheduler.run()
