from live_plot import LivePlot

from log_entries import UnderlayPacketEntry
from log_entries import UnderlaySendPacketEntry
from log_entries import UnderlayRecvPacketEntry

from log_processing import get_log
from log_processing import filter_entries
from log_processing import group_by
from abc import ABCMeta
from abc import abstractmethod

class PacketPlot():
    __metaclass__ = ABCMeta

    def __init__(self, log_file):
        self.log_file = log_file
        self.plot = LivePlot(scatter=True)
        self.log = []

    @property
    def figure(self):
        return self.plot.figure

    def update_log(self):
        self.log = filter_entries(get_log(self.log_file), UnderlayPacketEntry)
        self.packet_paths = group_by(self.log, lambda e: (e.src, e.dest))

    @abstractmethod
    def update(self):
        pass

class HopPlot(PacketPlot):
    def __init__(self, log_file, overlay=False):
        super(HopPlot, self).__init__(log_file)
        self.overlay = overlay

    def update(self):
        self.update_log()

        x = []
        y = []

        for key, packets in self.packet_paths.items():
            ctr = 0
            for p in packets:
                if isinstance(p, UnderlaySendPacketEntry):
                    ctr = 0
                elif isinstance(p, UnderlayRecvPacketEntry):
                    ctr += 1
                    x.append(p.timestamp)
                    y.append(ctr)
                else:
                    if self.overlay:
                        if p.recv is not None:
                            ctr += 1
                    else:
                        ctr += 1

        self.plot.update_data(x, y)
        self.figure.canvas.draw()

class LatencyPlot(PacketPlot):
    def __init__(self, log_file):
        super(LatencyPlot, self).__init__(log_file)

    def update(self):
        self.update_log()

        x = []
        y = []

        for key, packets in self.packet_paths.items():
            send = filter_entries(packets, UnderlaySendPacketEntry)
            recv = filter_entries(packets, UnderlayRecvPacketEntry)

            for i in range(min(len(send), len(recv))):
                x.append(recv[i].timestamp)
                y.append(recv[i].timestamp - send[i].timestamp)

        self.plot.update_data(x, y)
        self.figure.canvas.draw()
