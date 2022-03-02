package cluster

import (
	"KeyValueHTTPStore/internal/command"
	"context"
	"errors"
	"fmt"
	"github.com/buraksezer/consistent"
	"github.com/cespare/xxhash"
	"strconv"
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
	results := make([]string, len(cmds))
	commandsPerNode := make(map[inode][]command.Command, len(cmds))
	_ = commandsPerNode
	for idx, cmd := range cmds {
		nodeString := c.cHasher.LocateKey([]byte(cmd.Key())).String() // Get the ID of selected page

		nodeID, err := strconv.Atoi(nodeString) // Convert ID to int
		if err != nil {
			return nil, errors.New(fmt.Sprintf("failed to convert node string to int: %s", nodeString))
		}

		targetNode := c.nodes[nodeID]

		result, err := targetNode.execute(cmd) // Execute command on snapshot of selected page
		if err != nil {
			return nil, err
		}

		results[idx] = result

		if cmd.Type() == command.Write {
			commandsPerNode[targetNode] = append(commandsPerNode[targetNode], cmd) // Also, our journal needs only write commands
		}
	}

	for targetNode, cmds := range commandsPerNode {
		targetNode.addToJournal(cmds) // The commands are valid, so we add them to execute on real storage
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
