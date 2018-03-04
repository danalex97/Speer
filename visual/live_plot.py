import matplotlib
matplotlib.use('Qt4Agg')

from PyQt4.QtCore import QThread
from PyQt4 import QtCore as qtcore

import matplotlib.pyplot as plt
import numpy as np
import time

plt.style.use('ggplot')

from collections import defaultdict

class LivePlot():
    def __init__(self, names=defaultdict(str), scatter=False):
        self.fig, self.ax = plt.subplots()
        self.names = names

        self.scatter = scatter
        if self.scatter:
            self.path = self.ax.scatter([], [])
        else:
            self.data, = self.ax.plot([], [])

        self.ax.set_title(names["fig_name"])
        self.ax.set_xlabel(names["x_label"])
        self.ax.set_ylabel(names["y_label"])
        self.fig.canvas.draw()

    @property
    def figure(self):
        return self.fig

    def _set_data(self, x, y):
        self.x = x
        self.y = y
        if len(self.y) != len(self.x):
            self.y += [0] * (len(self.x) - len(self.y))
            self.x += [0] * (len(self.y) - len(self.x))

        if self.scatter:
            self.path.set_offsets(np.c_[x, y])
        else:
            self.data.set_xdata(self.x)
            self.data.set_ydata(self.y)

    def update_data(self, x, y, xlim=None, ylim=None):
        xlim = [min(x), max(x)] if xlim is None else xlim
        ylim = [min(y), max(y)] if ylim is None else ylim

        self.ax.set_xlim(xlim)
        self.ax.set_ylim(ylim)

        self._set_data(x, y)
