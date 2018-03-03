import matplotlib
matplotlib.use('Qt4Agg')

from PyQt4 import QtCore
from PyQt4 import QtGui as qt

from menu import MplMultiTab
import sys

import numpy as np
import matplotlib.pyplot as plt

x = np.linspace(1, 2*np.pi, 100)
figures = []
for i in range(1,3):
    fig, ax = plt.subplots()
    y = np.sin(np.pi*i*x)+0.1*np.random.randn(100)
    ax.plot(x,y)
    figures.append( fig )

app = qt.QApplication(sys.argv)
ui = MplMultiTab(figures = figures)
ui.show()
app.exec_()
