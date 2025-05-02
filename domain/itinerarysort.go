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
	tryForwardPath := true
	tryBackwardPath := true
	//var sortedRoutes []Route
	nodes := createNodeMap(routes)

	originNode, err := findOriginNode(nodes)
	if err != nil {
		return nil, err
	}

	terminationNode, err := findTerminationNode(nodes)
	if err != nil {
		return nil, err
	}
	if originNode == nil && terminationNode == nil {
		return nil, errors.New("no start or end node exists, all nodes belong to a circle")
	}

	for len(nodes) > 0 {

		if originNode != nil && tryForwardPath == true {
			var newlySortedRoutes []Route
			var forwardPathFailure error
			newlySortedRoutes, originNode, forwardPathFailure = forwardPath(nodes, *originNode)
			if len(newlySortedRoutes) != 0 {
				//If any routes were successfully sorted, add them to the slice and clear the backwardPathFailure error
				//so that the loop will try the backward path at least once more as node map has changed
				forwardSortedRoutes = append(forwardSortedRoutes, newlySortedRoutes...)
				tryBackwardPath = true
			}
			if forwardPathFailure != nil {
				tryForwardPath = false
			} else {
				tryBackwardPath = false
				tryForwardPath = false
			}

		}

		if terminationNode != nil && tryBackwardPath == true {
			var newlySortedRoutes []Route
			var backwardPathFailure error
			newlySortedRoutes, terminationNode, backwardPathFailure = backwardPath(nodes, *terminationNode)
			if len(newlySortedRoutes) != 0 {
				//If any routes were successfully sorted, add them to the slice and clear the forwardPathFailure error
				//so that the loop will try the forward path again at least once more as the node map has changed
				backwardsSortedRoutes = append(backwardsSortedRoutes, newlySortedRoutes...)
				tryForwardPath = true
			}
			if backwardPathFailure != nil {
				tryBackwardPath = false
			} else {
				tryBackwardPath = false
				tryForwardPath = false
			}
		}

		//Neither backward nor forward paths were able to make any progress in this loop, progress is locked
		if tryForwardPath == false && tryBackwardPath == false {
			return nil, errors.New("path cannot be determined, unresolvable node paths")
		}
	}

	return append(forwardSortedRoutes, backwardsSortedRoutes...), nil
}

func forwardPath(nodes nodeMap, originNode node) ([]Route, *node, error) {
	var sortedRoutes []Route
	currentNode := originNode

	for len(currentNode.outgoingNodes) > 0 {
		if len(currentNode.outgoingNodes) > 1 {
			return sortedRoutes, &currentNode, errors.New("fork found during forward path execution, aborting")
		}
		//Leave the current node on the only available path
		leavingNodeName := currentNode.name
		arrivingNodeName := currentNode.outgoingNodes[0]
		currentNode.outgoingNodes = remove(currentNode.outgoingNodes, arrivingNodeName)
		//If the node you are leaving has no available routes in or out, delete it. Otherwise, save it.
		if len(currentNode.incomingNodes) == 0 && len(currentNode.outgoingNodes) == 0 {
			delete(nodes, currentNode.name)
		} else {
			nodes[currentNode.name] = currentNode
		}

		//Add this route to the sorted routes
		sortedRoutes = append(sortedRoutes, Route{Source: leavingNodeName, Destination: arrivingNodeName})
		//Arrive at the next node
		currentNode = nodes[arrivingNodeName]
		currentNode.incomingNodes = remove(currentNode.incomingNodes, leavingNodeName)
		nodes[currentNode.name] = currentNode
	}
	if len(currentNode.incomingNodes) == 0 && len(currentNode.outgoingNodes) == 0 {
		delete(nodes, currentNode.name)
	} else {
		nodes[currentNode.name] = currentNode
	}
	return sortedRoutes, nil, nil
}

func backwardPath(nodes nodeMap, terminationNode node) ([]Route, *node, error) {
	var sortedRoutes []Route
	currentNode := terminationNode

	for len(currentNode.incomingNodes) > 0 {
		if len(currentNode.incomingNodes) > 1 {
			return sortedRoutes, &currentNode, errors.New("fork found during backwards path execution, aborting")
		}
		//Leave the current node on the only available path
		leavingNodeName := currentNode.name
		arrivingNodeName := currentNode.incomingNodes[0]
		currentNode.incomingNodes = remove(currentNode.incomingNodes, arrivingNodeName)
		//If the node you are leaving has no available routes in or out, delete it. Otherwise, save it.
		if len(currentNode.incomingNodes) == 0 && len(currentNode.outgoingNodes) == 0 {
			delete(nodes, currentNode.name)
		} else {
			nodes[currentNode.name] = currentNode
		}

		//Add this route to the sorted routes
		sortedRoutes = append([]Route{{Source: arrivingNodeName, Destination: leavingNodeName}}, sortedRoutes...)
		//Arrive at the next node
		currentNode = nodes[arrivingNodeName]
		currentNode.outgoingNodes = remove(currentNode.outgoingNodes, leavingNodeName)
		nodes[currentNode.name] = currentNode
	}

	if len(currentNode.incomingNodes) == 0 && len(currentNode.outgoingNodes) == 0 {
		delete(nodes, currentNode.name)
	} else {
		nodes[currentNode.name] = currentNode
	}
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

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
