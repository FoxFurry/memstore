package cluster

import (
	"context"
	"errors"
	"fmt"
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"strconv"
)

type Cluster interface {
	Execute(commands []interface{}) ([]string, error)
	Initialize(ctx context.Context) error
}

type cluster struct {
	nodes   []inode
	cHasher *consistent.Consistent
}

func New() Cluster {
	return &cluster{}
}

func (c *cluster) Execute(commands []interface{}) ([]string, error) {
	var (
		results         []string
		commandsPerNode map[int][]interface{}
	)

	for _, cmd := range commands {
		nodeString := c.cHasher.LocateKey([]byte("test")).String() // Get the ID of selected page
		nodeID, err := strconv.Atoi(nodeString)                    // Convert ID to int
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to convert node string to int: %s", nodeString))
		}

		nodeSnapshot := c.nodes[nodeID].snapshot() // Take snapshot of selected page

		result, err := nodeSnapshot.execute(cmd) // Execute command on snapshot of selected page
		if err != nil {
			return nil, err
		}

		commandsPerNode[nodeID] = append(commandsPerNode[nodeID], cmd)
		results = append(results, result) // Add to result array
	}

	for nodeID, cmds := range commandsPerNode {
		c.nodes[nodeID].addToJournal(cmds) // The commands are valid, so we add them to execute on real storage
	}

	return results, nil
}

func (c *cluster) Initialize(ctx context.Context) error {
	nodeNum := 4 // TODO: Find the way to calculate best number of nodes for different situations

	hasherConfig := consistent.Config{
		Hasher:            hasher{},
		PartitionCount:    271,
		ReplicationFactor: 40,
		Load:              1.2,
	}

	c.cHasher = consistent.New(nil, hasherConfig)

	for i := 0; i < nodeNum; i++ {
		c.nodes = append(c.nodes, newNode(ctx, i)) // Add a new page to cluster
		c.cHasher.Add(c.nodes[i])                  // And add this page to consistent hasher
	}

	return nil
}

type hasher struct{}

func (h hasher) Sum64(data []byte) uint64 {
	return xxhash.Sum64(data)
}
