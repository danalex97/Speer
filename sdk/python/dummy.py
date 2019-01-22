from scheduler_api import _schedule
from node import Node

class DummyNode(Node):
    def __init__( self, util ):
        self.util = util

    def on_join( self ):
        ctr = 0
        while True:
            message = self.util.recv()
            ctr += 1
            print("{} {}: {}".format(self.util.id, ctr, message))
            if not message:
                self.util.wait()
                yield from _schedule()
            yield from _schedule()

    def on_leave( self ):
        pass
