package loadbalancer

var (
	nodes       []string
	connections map[string]bool
)

func init() {
	nodes = []string{}
	connections = make(map[string]bool)
}

func GetNodeRoundRobin(current []string) string {
	updateNodes(current)
	return chooseNode()
}

func updateNodes(current []string) {
	deleteRemovedNodes(current)
	addNewNodes(current)
}

func deleteRemovedNodes(current []string) {
	removed := 0
	for i, node := range nodes {
		if !contains(current, node) {
			index := i - removed
			nodes = append(nodes[:index], nodes[index+1:]...)
			delete(connections, node)
			removed += 1
		}
	}
}

func addNewNodes(current []string) {
	for _, node := range current {
		if !contains(nodes, node) {
			nodes = append(nodes, node)
			connections[node] = false
		}
	}
}

func chooseNode() string {
	for _, node := range nodes {
		if connections[node] == false {
			connections[node] = true
			return node
		}
	}

	resetConnections()
	return chooseNode()
}

func resetConnections() {
	for node, _ := range connections {
		connections[node] = false
	}
}

func contains(slice []string, elem string) bool {
	for _, item := range slice {
		if item == elem {
			return true
		}
	}

	return false
}
