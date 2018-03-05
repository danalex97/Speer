from live_graph import Graph
from live_graph import Node

from log_entries import JoinEntry
from log_entries import QueryEntry
from log_entries import UnderlayPacketEntry
from log_entries import UnderlaySendPacketEntry
from log_entries import UnderlayRecvPacketEntry

from log_processing import get_log
from log_processing import filter_entries
from log_processing import group_by

from abc import ABCMeta
from abc import abstractmethod

class Network():
    __metaclass__ = ABCMeta

    def __init__(self, log_file, window):
        self.log_file = log_file
        self.graph = Graph()
        self.internal_time = 0
        self.window = window

    @property
    def figure(self):
        return self.graph.figure

    def update_log(self):
        self._log = get_log(self.log_file)

        self.internal_time += self.window
        self.internal_time = min(self.internal_time, self._log[-1].timestamp)

        self.log = [e for e in self._log if e.timestamp <= self.internal_time]

    @abstractmethod
    def update(self):
        pass

class OverlayNetwork(Network):
    def __init__(self, log_file, window = 5):
        super(OverlayNetwork, self).__init__(log_file, window)

    def update(self):
        self.update_log()

        joins = filter_entries(self.log, JoinEntry)
        for entry in joins:
            self.graph.add_node(Node(entry.nodeId))

        # Packets go in the underlay(i.e. the Recvs are in the underlay, but
        # the source and the destiantion represent overlay connections)
        packets = filter_entries(self.log, UnderlayPacketEntry)
        for entry in packets:
            self.graph.add_edge(Node(entry.src), Node(entry.dest))

        self.graph.draw()
