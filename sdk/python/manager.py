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
        node_id  = create.id
        util     = Util(self.env, node_id)

        self.nodes[node_id] = NodeWrapper.new(
            self.template,
            util,
        )

        # run node code
        execution = self.scheduler.execute(self.nodes[node_id].on_join)

        # allow util to explicitly block execution
        util.set_execution(execution)
