import requests
from bs4 import BeautifulSoup
import time

class AlgorithmParams:
    def __init__(self, title, link):
        self.title = title
        self.link = link

def sanitize_url_to_title(url):
    return url.strip().replace('_', ' ').replace('/wiki/', '')

def scraping_handler(url):
    try:
        res = requests.get(url)
        res.raise_for_status()  # Raise exception for non-200 status codes
        soup = BeautifulSoup(res.text, 'html.parser')
        links = set()
        for link in soup.find_all('a', href=lambda href: href and href.startswith('/wiki/') and ':' not in href):
            href = link.get('href')
            full_link = "https://en.wikipedia.org" + href
            links.add(full_link)
        return list(links)
    except requests.exceptions.RequestException as e:
        print(f"Error while processing {url}: {e}")
        return []

def bfs(source, destination, max_depth=6):
    if source == destination:
        return [[source]]

    queue = [(source, [source])]
    visited = [{source:True}]
    path = []
    found = False

    while queue or len(paths) < max_depth and len(queue) >0 and found:
        current_url, paths = queue.pop(0)
        if current_url == destination:
            if not found:
                found = True
            path.append(paths)
            return path
        elif not found:
            try:
                links = scraping_handler(current_url)
                for link in links:
                    if (link not in visited or not visited[link]) and len(paths) < max_depth:
                        visited.append({link:True})
                        queue.append((link, paths + [link]))
            except Exception as e:
                print(f"Error while processing {current_url}: {e}")

    return path if path else []

# Example usage
start = time.time()
source ="https://en.wikipedia.org/wiki/Joko_Widodo"
destination ="Destination", "https://en.wikipedia.org/wiki/Rengasdengklok_Incident"
routes = bfs(source, destination)
end = time.time()
print(f"Execution time: {end - start} seconds")

for route in routes:
    print(route)
