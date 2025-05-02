package domain

import "errors"

type ItineraryService struct {
}

type Route struct {
	Source      string
	Destination string
}

type node struct {
	name          string
	incomingNodes []string
	outgoingNodes []string
}

type nodeMap map[string]node

func (nm nodeMap) upsertNode(newNode node) {
	oldNode, nodeExists := nm[newNode.name]
	if !nodeExists {
		nm[newNode.name] = newNode
		return
	}
	oldNode.outgoingNodes = append(oldNode.outgoingNodes, newNode.outgoingNodes...)
	oldNode.incomingNodes = append(oldNode.incomingNodes, newNode.incomingNodes...)
	nm[oldNode.name] = oldNode
}

func createNodeMap(routes []Route) nodeMap {
	nodes := make(nodeMap)
	for _, e := range routes {

		nodes.upsertNode(node{
			name:          e.Source,
			outgoingNodes: []string{e.Destination},
		})
		nodes.upsertNode(node{
			name:          e.Destination,
			incomingNodes: []string{e.Source},
		})

	}

	return nodes
}

func (i ItineraryService) CalculatePath(routes []Route) ([]Route, error) {
	var forwardSortedRoutes []Route
	var backwardsSortedRoutes []Route
	var forwardPathFailure error
	var backwardPathFailure error
	var sortedRoutes []Route
	nodes := createNodeMap(routes)

	originNode, err := findOriginNode(nodes)
	if err != nil {
		return nil, err
	}

	terminationNode, err := findTerminationNode(nodes)
	if err != nil {
		return nil, err
	}

	for len(nodes) > 0 {

		if originNode == nil && terminationNode == nil {
			return nil, errors.New("no start or end node exists, all nodes belong to a circle")
		}

		if originNode != nil {
			forwardSortedRoutes, originNode, forwardPathFailure = forwardPath(nodes, *originNode)
			sortedRoutes = append(forwardSortedRoutes, sortedRoutes...)
			if forwardPathFailure == nil {
				return sortedRoutes, nil
			}
		}

		if terminationNode != nil {
			backwardsSortedRoutes, terminationNode, backwardPathFailure = backwardPath(nodes, *terminationNode)
			sortedRoutes = append(sortedRoutes, backwardsSortedRoutes...)
			if backwardPathFailure == nil {
				return sortedRoutes, nil
			}
		}

		if forwardPathFailure != nil && backwardPathFailure != nil {
			return nil, errors.New("path cannot be determined, unresolvable node paths")
		}
	}

	return sortedRoutes, nil
}

func forwardPath(nodes nodeMap, originNode node) ([]Route, *node, error) {
	var sortedRoutes []Route
	currentNode := originNode

	for len(currentNode.outgoingNodes) > 0 {
		if len(currentNode.outgoingNodes) > 1 {
			return nil, &currentNode, errors.New("fork found during forward path execution, aborting")
		}
		nextNodeName := currentNode.outgoingNodes[0]
		sortedRoutes = append(sortedRoutes, Route{Source: currentNode.name, Destination: nextNodeName})
		delete(nodes, currentNode.name)
		currentNode = nodes[nextNodeName]
	}
	//delete the last node
	delete(nodes, currentNode.name)
	return sortedRoutes, nil, nil
}

func backwardPath(nodes nodeMap, terminationNode node) ([]Route, *node, error) {
	var sortedRoutes []Route
	currentNode := terminationNode

	for len(currentNode.incomingNodes) > 0 {
		if len(currentNode.incomingNodes) > 1 {
			return nil, &currentNode, errors.New("fork found during backwards path execution, aborting")
		}
		nextNodeName := currentNode.incomingNodes[0]
		sortedRoutes = append(sortedRoutes, Route{Source: currentNode.name, Destination: nextNodeName})
		delete(nodes, currentNode.name)
		currentNode = nodes[nextNodeName]
	}
	//delete the last node
	delete(nodes, currentNode.name)
	return sortedRoutes, nil, nil
}

func findOriginNode(nm nodeMap) (*node, error) {
	var originNodes []node
	for _, n := range nm {
		if len(n.incomingNodes) == 0 {
			originNodes = append(originNodes, n)
		}
	}
	if len(originNodes) == 0 {
		return nil, nil
	}
	if len(originNodes) > 1 {
		return &node{}, errors.New("multiple origins found")
	}
	return &originNodes[0], nil
}

func findTerminationNode(nm nodeMap) (*node, error) {
	var terminationNodes []node
	for _, n := range nm {
		if len(n.outgoingNodes) == 0 {
			terminationNodes = append(terminationNodes, n)
		}
	}
	if len(terminationNodes) == 0 {
		return nil, nil
	}
	if len(terminationNodes) > 1 {
		return &node{}, errors.New("multiple terminations found")
	}
	return &terminationNodes[0], nil
}
