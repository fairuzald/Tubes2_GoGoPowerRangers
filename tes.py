import random
import timeit

class AlgorithmParams:
    def __init__(self, title, link):
        self.title = title
        self.link = link

class Graph:
    def __init__(self):
        self.graph = {}

    def add_node(self, node):
        if node not in self.graph:
            self.graph[node] = []

    def add_edge(self, node1, node2):
        self.graph[node1].append(node2)
        self.graph[node2].append(node1)

    def generate_random_graph(self, num_nodes, max_edges_per_node):
        nodes = [f"Node_{i}" for i in range(num_nodes)]
        for node in nodes:
            self.add_node(node)

        for node in nodes:
            # Generate random edges up to maximum number of edges per node
            num_edges = random.randint(1, max_edges_per_node)
            edges = random.sample(nodes, num_edges)
            for edge in edges:
                if edge != node and edge not in self.graph[node]:
                    self.add_edge(node, edge)

    def get(self, node):
        return self.graph.get(node, [])

def bfs(source, destination, graph):
    if source == destination:
        return [[source]]

    queue = [(source, [source])]
    visited = {source: True}  # Object of keys with boolean values
    found = False
    path = []

    while queue:
        current_url, paths = queue.pop(0)
        if current_url == destination:
            path.append(paths)
        else:
            try:
                links = graph.get(current_url)
                for link in links:
                    if link not in visited:
                        visited[link] = True
                        queue.append((link, paths + [link]))
            except Exception as e:
                print(f"Error while processing {current_url}: {e}")

    return path if path else []

def bidirectional_bfs(source, destination, graph):
    if source == destination:
        return [[source]]

    forward_queue = [(source, [source])]
    forward_visited = {source: True}
    backward_queue = [(destination, [destination])]
    backward_visited = {destination: True}
    path = []

    while forward_queue and backward_queue:
        # Pencarian maju (dari sumber ke tujuan)
        current_url, paths = forward_queue.pop(0)
        if current_url in backward_visited and (current_url, []) in backward_queue:  # Perbaikan di sini
            intersect_path = paths + backward_queue[backward_queue.index((current_url, []))][1][::-1]
            path.append(intersect_path)
        try:
            links = graph.get(current_url)
            for link in links:
                if link not in forward_visited:
                    forward_visited[link] = True
                    forward_queue.append((link, paths + [link]))
        except Exception as e:
            print(f"Error while processing {current_url}: {e}")

        # Pencarian mundur (dari tujuan ke sumber)
        current_url, paths = backward_queue.pop(0)
        if current_url in forward_visited and (current_url, []) in forward_queue:  # Perbaikan di sini
            intersect_path = paths + forward_queue[forward_queue.index((current_url, []))][1][::-1]
            path.append(intersect_path)
        try:
            links = graph.get(current_url)
            for link in links:
                if link not in backward_visited:
                    backward_visited[link] = True
                    backward_queue.append((link, paths + [link]))
        except Exception as e:
            print(f"Error while processing {current_url}: {e}")

    return path if path else []


# Generate a random graph with 1500 nodes and maximum 26 edges per node
graph = Graph()
graph.generate_random_graph(100000, 26)

# Example usage with timeit
def run_bfs():
    source = AlgorithmParams("Source", "Node_0")
    destination = AlgorithmParams("Destination", "Node_499")
    routes = bidirectional_bfs("Node_10", "Node_10000", graph)
    print(f"Number of routes: {(routes)}")

# Run BFS and measure execution time
execution_time = timeit.timeit(run_bfs, number=1)
print(f"Execution time: {execution_time} seconds")
