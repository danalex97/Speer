from functools import wraps

class Scheduler( object ):
    def __init__( self ):
        pass

@wraps
def coroutine(func):
    return func
