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
//
// TODO: Create all variables before for-loop, not inside for-loop
func (c *cluster) Execute(transaction []command.Command) ([]string, error) {
	results := make([]string, len(transaction))                          // Stores the results
	commandsPerNode := make(map[int][]command.Command, len(transaction)) // Maps all commands to their nodes

	for idx, cmd := range transaction {
		nodeString := c.cHasher.LocateKey(cmd.Key()).String() // Find node for specified key

		nodeID, err := strconv.Atoi(nodeString) // Convert node string to int
		if err != nil {
			return nil, errIdConversionFailed
		}

		targetNode := c.nodes[nodeID].snapshot()

		result, err := targetNode.execute(cmd) // Execute command on snapshot of selected page
		if err != nil {
			return nil, err
		}

		results[idx] = result // Write result to results array

		if cmd.Type() == command.Write {
			commandsPerNode[nodeID] = append(commandsPerNode[nodeID], cmd) // Also, our journal needs only write commands
		}
	}

	for nodeID, commands := range commandsPerNode {
		c.nodes[nodeID].addToJournal(commands) // The commands are valid, so we add them to execute on real storage
	}

	return results, nil
}

// Initialize a cluster with default node array and default consistent hasher parameters.
// It also initialize every node and starts journals
// Right now node array size is constant and not changing over-time.
func (c *cluster) Initialize(ctx context.Context) {
	nodeNum := 4 // TODO: Find the way to calculate best number of nodes for different situations

	hasherConfig := consistent.Config{
		Hasher:            hasher{},
		PartitionCount:    271,
		ReplicationFactor: 40,
		Load:              1.2,
	}

	c.cHasher = consistent.New(nil, hasherConfig)

	for i := 0; i < nodeNum; i++ {
		newNode := newNode(i)
		go newNode.startJournal(ctx)

		c.nodes = append(c.nodes, newNode) // Add a new page to cluster
		c.cHasher.Add(c.nodes[i])          // And add this page to consistent hasher
	}
}

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}
