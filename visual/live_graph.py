# import networkx as nx
# import matplotlib.animation as animation
# import matplotlib.pyplot as plt
# import random
# import time
#
# fig = plt.figure()
# ax = fig.add_subplot(111)
#
# # Graph initialization
# G = nx.Graph()
# G.add_nodes_from([1, 2, 3, 4, 5, 6, 7, 8, 9])
#
# all_edges = [(1,2), (3,4), (2,5), (4,5), (6,7), (8,9), (4,7), (1,7), (3,5), (2,7), (5,8), (2,9), (5,7)]
# G.add_edges_from(all_edges)
#
# # draw and show it
# ax.relim()
# ax.autoscale_view(True,True,True)
# nx.draw(G, ax=ax)
# plt.show(block=False)
#
# # time.sleep(1)
#
# edges = all_edges
# while True:
#     try:
#         G.remove_edges_from(edges)
#
#         sample = random.sample(range(len(all_edges)), 8)
#         edges = [all_edges[i] for i in sample]
#
#         G.add_edges_from(edges)
#
#         nx.draw(G, ax=ax)
#
#         time.sleep(1)
#     except KeyboardInterrupt:
#         break

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

num_plots = 50;
pylab.show()

fig = pylab.figure()

for i in range(num_plots):
    fig = get_fig(fig)
    pause(0.1)
