package cluster

import (
	"context"
	"errors"
	"github.com/FoxFurry/GoKeyValueStore/internal/command"
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"strconv"
)

var (
	errIdConversionFailed = errors.New("failed to convert node string to int")
	errExecutionFailed    = errors.New("command execution failed")
)

type Cluster interface {
	Execute(cmds []command.Command) ([]string, error)
	Initialize(ctx context.Context)
}

type cluster struct {
	nodes   []inode
	cHasher *consistent.Consistent
}

func New() Cluster {
	return &cluster{}
}

func (c *cluster) Execute(cmds []command.Command) ([]string, error) {
	results := make([]string, len(cmds))                          // Stores the results
	commandsForNode := make(map[int][]command.Command, len(cmds)) // Maps all commands to their nodes

	for idx, cmd := range cmds {
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
			commandsForNode[nodeID] = append(commandsForNode[nodeID], cmd) // Also, our journal needs only write commands
		}
	}

	for nodeID, commands := range commandsForNode {
		c.nodes[nodeID].addToJournal(commands) // The commands are valid, so we add them to execute on real storage
	}

	return results, nil
}

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
