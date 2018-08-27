class NodeInterface( object ):
    """ Node interface. """

    def __init__( self, util ):
        """ Method called at node cration. """
        raise NotImplementedError()

    def on_join( self ):
        """ Method called at node network join. """
        raise NotImplementedError()

    def on_leave( self ):
        """ Method called at node network leave. """
        raise NotImplementedError()

class NodeWrapper( object ):
    """ Wrapper around a node class. """

    def __init__( self, node ):
        self.node = node

    @staticmethod
    def new( cls, util ):
        """ Method called at node creation. """
        node = cls(util)
        return NodeWrapper(node)

    def on_join( self ):
        self.node.on_join()

    def on_leave( self ):
        self.node.on_leave()

class Node( NodeInterface ):
    """ Dummy node class which should be replaced by the entry point """

    def __init__( self, util ):
        pass

    def on_join( self ):
        pass

    def on_leave( self ):
        pass 
