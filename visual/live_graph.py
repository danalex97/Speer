import random
import pylab
from matplotlib.pyplot import pause
import networkx as nx

pylab.ion()

graph = nx.Graph()
node_number = 0
graph.add_node(node_number)

def get_fig(fig):
    global node_number
    node_number += 1

    graph.add_node(node_number)
    graph.add_edge(node_number, node_number - 1)

    nx.draw(graph)
    return fig

pylab.show()

fig = pylab.figure()

while True:
    try:
        fig.clf()
        fig = get_fig(fig)
        pause(0.1)

    except KeyboardInterrupt:
        break
