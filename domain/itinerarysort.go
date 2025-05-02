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

func (nodes nodeMap) upsertNode(newNode node) {
	oldNode, nodeExists := nodes[newNode.name]
	if !nodeExists {
		nodes[newNode.name] = newNode
		return
	}
	oldNode.outgoingNodes = append(oldNode.outgoingNodes, newNode.outgoingNodes...)
	oldNode.incomingNodes = append(oldNode.incomingNodes, newNode.incomingNodes...)
	nodes[oldNode.name] = oldNode
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

	originNode, err := nodes.findOriginNode()
	if err != nil {
		return nil, err
	}

	terminationNode, err := nodes.findTerminationNode()
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
				if len(nodes) != 0 {
					return nil, errors.New("disconnected nodes detected")
				}
				return append(forwardSortedRoutes, backwardsSortedRoutes...), nil
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
				if len(nodes) != 0 {
					return nil, errors.New("disconnected nodes detected")
				}
				return append(forwardSortedRoutes, backwardsSortedRoutes...), nil
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
		nodes[currentNode.name] = currentNode
		//If the node you are leaving has no available routes in or out, delete it
		if len(currentNode.incomingNodes) == 0 && len(currentNode.outgoingNodes) == 0 {
			delete(nodes, currentNode.name)
		}

		//Add this route to the sorted routes
		sortedRoutes = append(sortedRoutes, Route{Source: leavingNodeName, Destination: arrivingNodeName})
		//Arrive at the next node
		currentNode = nodes[arrivingNodeName]
		currentNode.incomingNodes = remove(currentNode.incomingNodes, leavingNodeName)
		nodes[currentNode.name] = currentNode
	}

	//The path has completed successfully, the current node should be the last node and is safe to delete
	delete(nodes, currentNode.name)

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
		nodes[currentNode.name] = currentNode

		//If the node you are leaving has no available routes in or out, delete it
		if len(currentNode.incomingNodes) == 0 && len(currentNode.outgoingNodes) == 0 {
			delete(nodes, currentNode.name)
		}

		//Add this route to the sorted routes
		sortedRoutes = append([]Route{{Source: arrivingNodeName, Destination: leavingNodeName}}, sortedRoutes...)
		//Arrive at the next node
		currentNode = nodes[arrivingNodeName]
		currentNode.outgoingNodes = remove(currentNode.outgoingNodes, leavingNodeName)
		nodes[currentNode.name] = currentNode
	}

	//The path has completed successfully, the current node should be the last node and is safe to delete
	delete(nodes, currentNode.name)

	return sortedRoutes, nil, nil
}

func (nodes nodeMap) findOriginNode() (*node, error) {
	var originNodes []node
	for _, n := range nodes {
		if len(n.incomingNodes) == 0 {
			originNodes = append(originNodes, n)
		}
	}
	//No origin node found, but that doesn't necessarily trigger an error
	if len(originNodes) == 0 {
		return nil, nil
	}
	if len(originNodes) > 1 {
		return &node{}, errors.New("multiple origins found")
	}
	return &originNodes[0], nil
}

func (nodes nodeMap) findTerminationNode() (*node, error) {
	var terminationNodes []node
	for _, n := range nodes {
		if len(n.outgoingNodes) == 0 {
			terminationNodes = append(terminationNodes, n)
		}
	}
	//No termination node found, but that doesn't necessarily trigger an error
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
