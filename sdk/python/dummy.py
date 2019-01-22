from scheduler_api import _schedule
from node import Node

class DummyNode(Node):
    def __init__( self, util ):
        self.util = util

    def on_join( self ):
        while True:
            # print("Alive: {}".format(self.util.node_id))
            yield from _schedule()

    def on_leave( self ):
        pass
