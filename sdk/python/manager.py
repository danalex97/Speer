from node import NodeWrapper
from node import Node
from util import Util

class NodeManager( object ):
    def __init__( self, env : 'Environ',  template : type ):
        self.nodes = {}
        self.env = env
        self.template = template

    def create( self, create ):
        node_id             = create.id
        self.nodes[node_id] = NodeWrapper.new(
            self.template,
            Util(self.env, node_id)
        )
