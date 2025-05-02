## Flightpath

### To set up and run
Code will require Go to be installed. To run the application:
```
 go run main.go
```

The application will by default log at debug level and will listen on port 8000.

Both of these can be changed by setting these ENV variables:

```
FP_LOG_LEVEL
FP_SERVER_PORT
```

### Calling the service

A single POST endpoint serves the request on the /itinerary endpoint, it accepts an array of tuples

### Tests
Run the following to confirm tests are working:

```
 go test ./...
```

### Design decisions
The core logic of the program sits in the domain directory. The HTTP server is written using the echo package. Zerolog
is used as the logging package, although for such a simple program is probably not really needed. Viper is used to 
parse the ENV variables. 

The main function calls the start CMD, which sets up the service and handles dependency injection. The functions of the
HTTP server and the domain logic are split and abstracted by the server accepting an interface of the domain logic. This 
allows for simpler unit tests and mocking on the server package, should they be needed. To save time I've just 
implemented tests in the core domain logic.

The program works out the path taken given a set of tuples. The program can handle loops up to a point. 

The solution is derived by first creating a hashmap of nodes based on the tuples given, and tracking the incoming and 
outgoing edges on that node. A node that has no incoming edges is the origin node, one with no outgoing edges is the 
termination node. The service needs at least one of these two to start the process of traversing the nodes. The program
will first attempt to start at the origin, and traverse the flight paths forward. As the traversal leaves a node, it 
will delete the edge it is travelling on, by deleting the record of that outgoing edge on the node that is being left 
behind, and by deleting the record of the incoming edge on the node that is being arrived at. Nodes that have had all
edge's connected to then traversed are removed from the map. Traversal continues until all nodes have been removed
from the hashmap. If a traversal arrives at a fork it will stop, and attempt to start traversing in the opposite 
direction from the termination node, to see if it can get back to that node down one of the two outgoing paths, 
resolving the dilemma. This allows the program to handle situations where a single loop exists. Multiple loops centred
on a single node (creating an infinity symbol) can't be handled, as it is impossible without to prioritise one path over
the other. Due to the traversal method, it is possible that one of the two end nodes (origin or termination) could be 
part of a loop, as long as the other one isn't.

All nodes need to be connected.
