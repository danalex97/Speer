import matplotlib
matplotlib.use('Qt4Agg')

import sys
import numpy as np
import matplotlib.pyplot as plt
from matplotlib.backends.backend_qt4agg import FigureCanvasQTAgg as FigureCanvas

from PyQt4 import QtCore
from PyQt4 import QtGui as qt

from menu import MplMultiTab
from live_plot import LivePlot
import threading

if __name__ == "__main__":
    plots = [LivePlot() for _ in range(3)]

    def animate(plot):
        x = np.linspace(1, 2*np.pi, 100)
        y = np.sin(np.pi*x)+0.1*np.random.randn(100)
        plot.update_data(x, y)
        # plot.figure.canvas.draw()

    app = qt.QApplication(sys.argv)

    ui = MplMultiTab(figures = [p.figure for p in plots])
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

    app.exec_()
