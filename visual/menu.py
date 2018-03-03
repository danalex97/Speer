from PyQt4 import QtCore
from PyQt4 import QtGui as qt

from matplotlib.backends.backend_qt4agg import FigureCanvasQTAgg as FigureCanvas
from matplotlib.backends.backend_qt4agg import NavigationToolbar2QT as NavigationToolbar
from matplotlib.figure import Figure

import itertools

# https://stackoverflow.com/questions/29681099/python-matplotlib-and-pyqt-for-multi-tab-plotting-navigation
class MultiTabNavTool(qt.QWidget):
    def __init__(self, canvases, tabs, parent=None):
        qt.QWidget.__init__(self, parent)
        self.canvases = canvases
        self.tabs = tabs
        self.toolbars = [NavigationToolbar(canvas, parent) for canvas in self.canvases]
        vbox = qt.QVBoxLayout()
        for toolbar in self.toolbars:
            vbox.addWidget(toolbar)
        self.setLayout(vbox)
        self.switch_toolbar()
        self.tabs.currentChanged.connect(self.switch_toolbar)

    def switch_toolbar(self):
        for toolbar in self.toolbars:
            toolbar.setVisible(False)
        self.toolbars[self.tabs.currentIndex()].setVisible(True)

class MplMultiTab(qt.QMainWindow):
    def __init__(self, parent=None, figures=None, labels=None):
        qt.QMainWindow.__init__(self, parent)

        self.main_frame = qt.QWidget()
        self.tabWidget = qt.QTabWidget( self.main_frame )
        self.create_tabs( figures, labels )

        # Create the navigation toolbar, tied to the canvas
        self.mpl_toolbar = MultiTabNavTool(self.canvases, self.tabWidget, self.main_frame)

        self.vbox = vbox = qt.QVBoxLayout()
        vbox.addWidget(self.mpl_toolbar)
        vbox.addWidget(self.tabWidget)

        self.main_frame.setLayout(vbox)
        self.setCentralWidget(self.main_frame)

    def create_tabs(self, figures, labels ):
        if labels is None:
            labels = []
        figures = [Figure()] if figures is None else figures
        self.canvases = [self.add_tab(fig, lbl)
            for (fig, lbl) in itertools.zip_longest(figures, labels) ]

    def add_tab(self, fig=None, name=None):
        if fig is None:
            fig = Figure()
            ax = fig.add_subplot(111)

        canvas = fig.canvas if fig.canvas else FigureCanvas(fig)
        canvas.setParent(self.tabWidget)
        canvas.setFocusPolicy( QtCore.Qt.ClickFocus )

        name = name if name is None else name
        name = "Tab {}".format(self.tabWidget.count() + 1)

        self.tabWidget.addTab(canvas, name)

        return canvas
