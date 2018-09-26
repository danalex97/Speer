from scheduler import Execution
from scheduler import Future

from scheduler import signals

def _return(value):
    yield signals._return, value

def _start(function, *args, **kwargs):
    execution = Execution(function, *args, **kwargs)
    future    = Future(execution)

    # ???
    yield signals._start, execution

    return future

def _schedule():
    yield signals._schedule

def _wait(future):
    while not future.done:
        yield signals._wait
    return future.value
