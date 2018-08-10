import struct

class PipeQueue( object ):
    def __init__(self, pipe_in, pipe_out):
        self.pipe_in  = pipe_in
        self.pipe_out = pipe_out

    def pop(self):
        """
        Blocking pop operation.
        """
        to_read, = struct.unpack('i', self.pipe_in.read(4))
        return self.pipe_in.read(to_read)

    def push(self, to_write):
        """
        Blocking push operation.
        """
        n  = len(to_write)
        nb = struct.pack('i', n)
        return self.pipe_out.write(nb + to_write)
