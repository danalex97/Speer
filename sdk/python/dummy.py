from node import Node

class DummyNode(Node):
    def __init__( self, util ):
        self.util = util

    def on_join( self ):
        print("Joined")

    def on_leave( self ):
        pass
