from messages import MESSAGES, Message

class Util( object ):
    def __init__( self, env : 'Environ', node_id : str ) -> None:
        self.env = env

        # register util in the environment
        self.env.register(self, node_id)
        self.node_id = node_id

    @property
    def id( self ) -> str:
        return self.node_id

    def recv( self, message : Message ) -> None:
        print("Recv {}".format(message))
