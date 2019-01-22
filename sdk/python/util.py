class Util( object ):
    def __init__( self, env : 'Environ', node_id : str ) -> None:
        self.env = env

        # register util in the environment
        self.env.register(self, node_id)

    def recv( self, message ):
        print("Recv {}".format(message))
