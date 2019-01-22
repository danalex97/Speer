from messages import MESSAGES
from node import NodeWrapper
from node import Node
from util import Util

class NodeManager( object ):
    def __init__( self,
            env : 'Environ',
            template : type,
            scheduler : 'Scheduler'
        ) -> None:
        """ Initialize template and keep scheduler for runnning new nodes."""
        self.nodes = {}
        self.env = env
        self.template = template
        self.scheduler = scheduler

    def create( self, create : MESSAGES.Create ):
        """ Update node set and execute new node's on_join. """
        node_id             = create.id
        self.nodes[node_id] = NodeWrapper.new(
            self.template,
            Util(self.env, node_id)
        )

        # run node code
        self.scheduler.execute(self.nodes[node_id].on_join)
