from typing import Optional

from messages import MESSAGES, Message
from scheduler import Execution
from collections import deque
from scheduler_api import _schedule

class Util( object ):
    def __init__( self, env : 'Environ', node_id : str ) -> None:
        self.env = env

        # register util in the environment
        self.env.register(self, node_id)
        self.node_id = node_id

        self.message_queue = deque()
        self.exectuion = None

    def _handle( self, message : Message ) -> None:
        """ Handle incoming message from simulator. """
        if (len(self.message_queue) == 0 and
                self.execution and
                self.execution.blocked):
            self.execution.unblock()
        self.message_queue.append(message)

    def set_execution( self, execution : Execution ) -> None:
        self.execution = execution

    @property
    def id( self ) -> str:
        return self.node_id

    def wait( self ) -> None:
        if self.execution and not self.execution.blocked:
            self.execution.block()

    def recv( self ) -> Optional[Message]:
        if len(self.message_queue) == 0:
            return None
        return self.message_queue.popleft()

    def send( self, message : Message ) -> None:
        pass
