import requests
from bs4 import BeautifulSoup
import time

class AlgorithmParams:
    def __init__(self, title, link, path=None):
        self.title = title
        self.link = link
        self.path = [] if path is None else path

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

def bidirectional_bfs(source, destination, max_depth=6):
    if source.link == destination.link:
        return [[source.link]]

    queue_source = [(source.link, [source.link])]
    queue_dest = [(destination.link, [destination.link])]
    visited_source = {source.link: True}
    visited_dest = {destination.link: True}

    while queue_source and queue_dest and len(queue_source[0][1]) <= max_depth/2 and len(queue_dest[0][1]) <= max_depth/2:
        current_url_source, path_source = queue_source.pop(0)
        current_url_dest, path_dest = queue_dest.pop(0)
        
        # Check intersect visited
        intersect_source =  (current_url_source in visited_source.keys() and current_url_source in visited_dest.keys())
        intersect_dest = (current_url_dest in visited_source.keys() and current_url_dest in visited_dest.keys())
        if intersect_source:
            print(f"Intersect found: {current_url_source}")
            print(f"Path source: {path_source}")
            print(f"Destination source: {path_dest}")
            print("Intersect found")
            return [path_source  + path_dest[::-1]] 

        if intersect_dest:
            print(f"Intersect found: {current_url_dest} dest")
            print(f"Path source: {path_source}")
            print(f"Destination source: {path_dest}")
            print("Intersect found")
            return [path_source  + path_dest[::-1]]
        try:
            links_source = scraping_handler(current_url_source)
            for link in links_source:
                if link not in visited_source and len(path_source) < max_depth/2:
                    visited_source[link] = True
                    queue_source.append((link, path_source + [link]))
        except Exception as e:
            print(f"Error while processing {current_url_source}: {e}")

        try:
            links_dest = scraping_handler(current_url_dest)
            for link in links_dest:
                if link not in visited_dest and len(path_dest) < max_depth/2:
                    visited_dest[link] = True
                    queue_dest.append((link, path_dest + [link]))
        except Exception as e:
            print(f"Error while processing {current_url_dest}: {e}")

    return []

# Example usage
start = time.time()
source = AlgorithmParams("Source", "https://en.wikipedia.org/wiki/Joko_Widodo")
destination = AlgorithmParams("Destination", "https://en.wikipedia.org/wiki/Rengasdengklok_Incident")
routes = bidirectional_bfs(source, destination)
end = time.time()
print(f"Execution time: {end - start} seconds")

for route in routes:
    print(route)
