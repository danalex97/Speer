from scheduler_api import _return, _start, _schedule, _wait
from scheduler import Scheduler

def g(x):
    _return(x + 1)

def f():
    ans = _start(g, 5)
    _schedule()
    print(_wait(ans))

if __name__ == "__main__":
    s = Scheduler()
    s.execute(f)
    s.start()
