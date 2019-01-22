import sys
import struct
import json

from typing    import Generator

from util      import Util
from dummy     import DummyNode
from queues    import PipeQueue
from manager   import NodeManager
from messages  import MAP_ID_MESSAGE, MAP_MESSAGE_ID, MESSAGES

from scheduler import Scheduler
from scheduler_api import _schedule

class Environ( object ):
    def __init__(self, template : type) -> None:
        """ The environ needs to receive a node template. """
        pipe_out = open(sys.argv[1], 'wb')
        pipe_in  = open(sys.argv[2], 'rb')

        self.queue     = PipeQueue(pipe_in, pipe_out)
        self.manager   = NodeManager(self, template)
        self.utils     = {}

    def recv( self ) -> Generator:
        """ Yields messages received from the simulator. """
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

    def send( self, message ) -> None:
        """ Utility for sending messages to the simulator. """
        tp = type(message)
        tp = MAP_MESSAGE_ID[tp]

        msg = str(message.to_json())
        msg = msg.replace("'", '"')
        msg = msg.encode('utf8')

        self.queue.push(struct.pack('i', tp) + msg)

    def register( self, util : Util, node_id : str ) -> None:
        """ Register the util to allow message passing. """
        self.utils[node_id] = util

    def run( self ):
        for message in self.recv():
            if isinstance(message, MESSAGES.Create):
                self.manager.create(message)
            if hasattr(message, "id"):
                self.utils[message.id].recv(message)
            yield from _schedule()

if __name__ == "__main__":
    print("Python environment started...")

    scheduler = Scheduler()
    env = Environ(DummyNode)

    scheduler.execute(env.run)
    scheduler.run()
