from random import randint

import pylab
from matplotlib.pyplot import pause

import networkx as nx

class Node():
    canvas_size = (0, 100)

    def __init__(self, _id, x = None, y = None):
        self._id = _id
        if x is None:
            x = randint(*Node.canvas_size)
        if y is None:
            y = randint(*Node.canvas_size)

        self._pos = (x, y)

    @property
    def id(self):
        return self._id

    @property
    def pos(self):
        return self._pos

class Graph():
    def __init__(self):
        self.graph = nx.Graph()
        self.nodes = set()
        self.edges = set()

        self.fig = pylab.figure()

    @property
    def figure(self):
        return self.fig

    def add_node(self, node):
        self.nodes.add(node)
        self.graph.add_node(node.id, pos = node.pos)

    def add_edge(self, node1, node2):
        if node1 not in self.nodes:
            self.add_node(node1)
        if node2 not in self.nodes:
            self.add_node(node2)
        if (node1, node2) not in self.edges:
            self.graph.add_edge(node1.id, node2.id)
            self.edges.add((node1, node2))
            self.edges.add((node2, node1))

    def draw(self):
        self.figure.clf()
        pos = nx.get_node_attributes(self.graph, 'pos')
        nx.draw(self.graph, pos, with_labels=True)
        self.figure.canvas.draw()
