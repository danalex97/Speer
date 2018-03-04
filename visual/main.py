import matplotlib
matplotlib.use('Qt4Agg')

import sys
import numpy as np
import matplotlib.pyplot as plt
# plt.ion()

from PyQt4 import QtCore
from PyQt4 import QtGui as qt

from menu import MplMultiTab
from live_plot import LivePlot
from live_graph import Node
from live_graph import Graph
import threading

if __name__ == "__main__":
    plots = [LivePlot() for _ in range(3)]
    graph = Graph()

    def animate(plot):
        x = np.linspace(1, 2*np.pi, 100)
        y = np.sin(np.pi*x)+0.1*np.random.randn(100)
        plot.update_data(x, y)
        # plot.figure.canvas.draw()

    app = qt.QApplication(sys.argv)

    ui = MplMultiTab(figures = [p.figure for p in plots] + [graph.figure])
    ui.show()

    for p in plots:
        class Updater():
            def __init__(self, plot):
                self.plot   = plot
                self.canvas = plot.figure.canvas

            def update(self):
                animate(self.plot)
                self.canvas.draw()

        ui.add_updater(Updater(p).update)

    class GraphUpdater():
        def __init__(self, graph):
            self.ctr = 0
            self.last_node = None
            self.node = None
            self.graph = graph

        def update(self):
            self.ctr += 1

            self.last_node = self.node
            self.node = Node(self.ctr)

            self.graph.add_node(self.node)
            if self.ctr > 1:
                self.graph.add_edge(self.last_node, self.node)
            self.graph.draw()

    ui.add_updater(GraphUpdater(graph).update)

    app.exec_()
