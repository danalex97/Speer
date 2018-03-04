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

from log_processing import get_log
from metrics import HopPlot
from metrics import LatencyPlot
from networks import UnderlayNetwork

if __name__ == "__main__":
    tabs = [
        HopPlot("metrics.txt", overlay=False),
        HopPlot("metrics.txt", overlay=True),
        LatencyPlot("metrics.txt"),
        UnderlayNetwork("metrics.txt")
    ]

    app = qt.QApplication(sys.argv)

    ui = MplMultiTab(figures = [tab.figure
        for tab in tabs])
    ui.show()

    for tab in tabs:
        ui.add_updater(tab.update)
    app.exec_()
