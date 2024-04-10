path_source = ["A", "B", "C", "D"]
path_dest = ["Z", "Y", "C", "G"]
reversed_dest = list(reversed(path_dest))  # Convert reversed_dest to a list
# Check if there's an intersection between path_source and path_dest
if set(path_source) & set(reversed_dest):
    intersect = list(set(path_source) & set(reversed_dest))  # Find the intersecting elements
    intersect_index = path_source.index(intersect[0])  # Find the index of the intersecting element in path_source
    
    new_path = path_source[:intersect_index] + reversed_dest[reversed_dest.index(intersect[0]):]  # Construct the new path
    
    print(new_path)
else:
    print("No intersection found.")
