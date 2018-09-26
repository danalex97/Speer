from functools   import wraps
from collections import deque

# Signals
def _return(): pass
def _start(): pass
def _schedule(): pass
def _wait(): pass

class signals:
    _return   = _return
    _start    = _start
    _schedule = _schedule
    _wait     = _wait

class Future( object ):
    def __init__( self, execution ):
        self.execution = execution

        self._value = None
        self._done  = False

    def _finished( self, value ):
        self._value = value
        self._done  = True

    @property
    def value( self ):
        return self._value

    @property
    def done( self ):
        return self._done

class Execution( object ):
    """
    An execution is a method with arguments scheduled to run.
    :blocked - execution is blocked
    :future  - future associated with execution result
    :parent  - parent execution
    """
    def __init__( self, method, *args, **kwargs ):
        self.method   = method
        self.args     = args
        self.kwargs   = kwargs

        self.parent   = None
        self.future   = None
        self._blocked  = False

    @property
    def blocked( self ):
        return self._blocked

    def block( self ):
        """ Block execution. """
        self._blocked = True

    def unblock( self ):
        """ Unblock execution. """
        self._blocked = False

    def run( self ):
        while True:
            if self.args and not self.kwargs:
                yield self.method(*self.args)
            if not self.args and self.kwargs:
                yield self.method(**self.kwargs)
            if self.args and self.kwargs:
                yield self.method(*self.args, **self.kwargs)
            else:
                yield self.method()

class Scheduler( object ):
    """
    Workflow:
     * each execution:
       - enqueues itself
       - calls scheduler
     * scheduler loop:
       - gets first job
       - runs job & gets yielded value
       - signals:
         * _return, value
           - does not enqueue anything
           - since return has been called we don't need to rechedule the
             execution
         * _schedule
           - enqueues same execution
           - a _schedule yield implies the execution has not finished yet
         * _start, child_execution
           - starts a new function
           - enqueues same execution back
           - returns a future in the original function
           : yields execution
         * _wait
           - waits for a future to execute
           - returns result
         -- this API can be extended if needed
    """
    def __init__( self ):
        self.executions = deque()

    def execute( self, function, *args, **kwargs ):
        self.executions.append(Execution(function, *args, **kwargs))

    def start( self ):
        while True:
            self.step()

    def step( self ):
        if len(self.executions) == 0:
            return

        # The coroutine will yield on a method with arguments
        execution = self.executions.popleft()

        if execution.blocked:
            # If execution is blocked, skip the execution and renqueue it
            self.executions.append(execution)

        try:
            yield_tuple = next(execution.run())
            print(yield_tuple)

            # If execution has not finished
            if yield_tuple is not None:
                if isinstance(yield_tuple, tuple):
                    yield_type, rest = yield_tuple[0], yield_tuple[1:]
                else:
                    yield_type, rest = yield_tuple[0], None

                if yield_type is _return:
                    value = rest[0]

                    # Future finished execution
                    execution.future._finished(value)

                    # The parent is, thus, no more blocked
                    if execution.parent is not None:
                        execution.parent.unblock()
                elif yield_type is _schedule:
                    # The function called the scheduler, so we renqueue the
                    # execution
                    self.executions.append(execution)

                elif yield_type is _start:
                    child_execution = rest[0]

                    self.executions.append(child_execution)
                    self.executions.append(execution)

                elif yield_type is _wait:
                    self.executions.append(execution)
        except StopIteration:
            # The method has stopped so, we don't need to restart it
            pass
