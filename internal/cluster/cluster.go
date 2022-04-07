/*
Package cluster

Provides ready-to-use interface for memstore code
*/
package cluster

import (
	"context"
	"errors"
	"github.com/FoxFurry/memstore/internal/command"
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"strconv"
)

var (
	errIdConversionFailed = errors.New("failed to convert node string to int")
	errExecutionFailed    = errors.New("command execution failed")
)

// hasher object for consistent hashing
type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}

// Cluster handles distribution of commands between nodes. Besides of array of nodes it contains consistent-hasher which
// maps key of command to specific node.
type Cluster interface {
	Execute(cmds []command.Command) ([]string, error)
	Initialize(ctx context.Context)
}

type cluster struct {
	nodes   []inode
	cHasher *consistent.Consistent
}

// New creates an empty cluster without initialization
func New() Cluster {
	return &cluster{}
}

// Execute implements atomic and fully isolated execution of a single transaction
// For every command it uses consistent hasher to find target node based on command key, after this it creates
// a snapshot of target node and executes command on the snapshot.
// If error occurs - return it immediately, otherwise - append result to results array.
// If all commands pass without errors - add commands to target node queue
func (c *cluster) Execute(transaction []command.Command) ([]string, error) {
	var ( // Declare all variables before loop
		results            = make([]string, 0, len(transaction))               // Stores the results
		commandsPerNode    = make(map[int][]command.Command, len(transaction)) // Maps all commands to their nodes
		existingNodeCopies = make(map[int]inode, len(c.nodes))                 // This map will help avoid getting snapshots of same node
		cmd                command.Command
		cmds               []command.Command
		nodeID             int
		resultBuffer       string
		err                error
	)

	for _, cmd = range transaction { // Find node for specified key
		nodeID, err = strconv.Atoi(c.cHasher.LocateKey(cmd.Key()).String()) // Convert node string to int
		if err != nil {
			return nil, errIdConversionFailed
		}

		var targetNode inode
		if _, ok := existingNodeCopies[nodeID]; ok { // if we already made a snapshot - use it
			targetNode = existingNodeCopies[nodeID]
		} else {
			targetNode = c.nodes[nodeID].snapshot() // Otherwise - take a snapshot and save it to the map
			existingNodeCopies[nodeID] = targetNode
		}

		resultBuffer, err = targetNode.execute(cmd) // Execute command on snapshot of selected page
		if err != nil {
			return nil, err
		}

		results = append(results, resultBuffer) // Write result to results array

		if cmd.Type() == command.Write {
			commandsPerNode[nodeID] = append(commandsPerNode[nodeID], cmd) // Also, our queue needs only write commands
		}
	}

	for nodeID, cmds = range commandsPerNode {
		c.nodes[nodeID].addToQueue(cmds) // The commands are valid, so we add them to execute on real storage
	}

	return results, nil
}

// Initialize a cluster with default node array and default consistent hasher parameters.
// It also initializes every node and starts queues
// Right now node array size is constant and not changing over-time.
func (c *cluster) Initialize(ctx context.Context) {
	nodeNum := 4 // TODO: Find the way to calculate best number of nodes for different situations

	hasherConfig := consistent.Config{
		Hasher:            hasher{},
		PartitionCount:    271, // TODO: Learn how to find perfect values for hasher config
		ReplicationFactor: 40,
		Load:              1.2,
	}

	c.cHasher = consistent.New(nil, hasherConfig)

	for i := 0; i < nodeNum; i++ {
		newNode := newNode(i)              // Create new node
		go newNode.startQueueListener(ctx) // Immediately start queue for this node
		// TODO: Awkward goroutine start
		c.nodes = append(c.nodes, newNode) // Add a new node to cluster
		c.cHasher.Add(c.nodes[i])          // And add this node to consistent hasher
	}
}
