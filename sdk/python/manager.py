from node import NodeWrapper
from node import Node
from util import Util

class NodeManager( object ):
    def __init__( self, env ):
        self.nodes = {}
        self.env   = env

    def create( self, create ):
        node_id             = create.id
        self.nodes[node_id] = NodeWrapper.new(Node, Util(self.env))
