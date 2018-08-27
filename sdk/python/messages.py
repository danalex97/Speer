class Serializable( object ):
    def to_json( self ):
        """ Serialize to json. """
        return NotImplementedError()

    @staticmethod
    def from_json( json ):
        """ Json to message object. """
        return NotImplementedError()

class Message( Serializable ):
    def __repr__( self ):
        return str(self.to_json())

    def __str__( self ):
        return str(self.to_json())

class Create( Message ):
    @staticmethod
    def from_json( json ):
        msg = Create()
        msg.id = json['id']
        return msg

    def to_json( self ):
        return {
            'id' : self.id
        }

MAP_ID_MESSAGE = {
    0 : Create,
}

class MESSAGES( object ):
    Create = Create
