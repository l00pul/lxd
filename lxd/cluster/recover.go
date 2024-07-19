package cluster

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	dqlite "github.com/canonical/go-dqlite"
	"github.com/canonical/go-dqlite/client"

	"github.com/canonical/lxd/lxd/db"
	"github.com/canonical/lxd/lxd/node"
)

// ListDatabaseNodes returns a list of database node names.
func ListDatabaseNodes(database *db.Node) ([]string, error) {
	nodes := []db.RaftNode{}
	err := database.Transaction(context.TODO(), func(ctx context.Context, tx *db.NodeTx) error {
		var err error
		nodes, err = tx.GetRaftNodes(ctx)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to list database nodes: %w", err)
	}

	addresses := make([]string, 0)
	for _, node := range nodes {
		if node.Role != db.RaftVoter {
			continue
		}

		addresses = append(addresses, node.Address)
	}

	return addresses, nil
}

// Return the entry in the raft_nodes table that corresponds to the local
// `core.https_address`.
// Returns err if no raft_node exists for the local node.
func localRaftNode(database *db.Node) (*db.RaftNode, error) {
	var info *db.RaftNode
	err := database.Transaction(context.TODO(), func(ctx context.Context, tx *db.NodeTx) error {
		var err error
		info, err = node.DetermineRaftNode(ctx, tx)

		return err
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to determine cluster member raft role: %w", err)
	}

	// If we're not a database node, return an error.
	if info == nil {
		return nil, fmt.Errorf("This cluster member has no raft role")
	}

	return info, nil
}

// Recover rebuilds the dqlite raft configuration leaving only the current
// member in the cluster. Use `Reconfigure` if more members should remain in
// the raft configuration.
func Recover(database *db.Node) error {
	info, err := localRaftNode(database)
	if err != nil {
		return err
	}

	// If this is a standalone node not exposed to the network, return an
	// error.
	if info.Address == "" {
		return fmt.Errorf("This LXD instance is not clustered")
	}

	dir := filepath.Join(database.Dir(), "global")

	cluster := []dqlite.NodeInfo{
		{
			ID:      uint64(info.ID),
			Address: info.Address,
			Role:    client.Voter,
		},
	}

	err = dqlite.ReconfigureMembershipExt(dir, cluster)
	if err != nil {
		return fmt.Errorf("Failed to recover database state: %w", err)
	}

	// Update the list of raft nodes.
	err = database.Transaction(context.TODO(), func(ctx context.Context, tx *db.NodeTx) error {
		nodes := []db.RaftNode{
			{
				NodeInfo: client.NodeInfo{
					ID:      info.ID,
					Address: info.Address,
				},
				Name: info.Name,
			},
		}

		return tx.ReplaceRaftNodes(nodes)
	})
	if err != nil {
		return fmt.Errorf("Failed to update database nodes: %w", err)
	}

	return nil
}

// updateLocalAddress updates the cluster.https_address for this node.
func updateLocalAddress(database *db.Node, address string) error {
	err := database.Transaction(context.TODO(), func(ctx context.Context, tx *db.NodeTx) error {
		var err error
		config, err := node.ConfigLoad(ctx, tx)
		if err != nil {
			return err
		}

		newConfig := map[string]any{"cluster.https_address": address}
		_, err = config.Patch(newConfig)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Failed to update node configuration: %w", err)
	}

	return nil
}

// Reconfigure replaces the entire cluster configuration.
// Addresses and node roles may be updated. Node IDs are read-only.
func Reconfigure(database *db.Node, raftNodes []db.RaftNode) error {
	info, err := localRaftNode(database)
	if err != nil {
		return err
	}

	localAddress := info.Address

	nodes := make([]client.NodeInfo, 0, len(raftNodes))
	for _, raftNode := range raftNodes {
		nodes = append(nodes, raftNode.NodeInfo)

		// Get the new address for this node.
		if raftNode.ID == info.ID {
			localAddress = raftNode.Address
		}
	}

	// Update cluster.https_address if changed.
	if localAddress != info.Address {
		err := updateLocalAddress(database, localAddress)
		if err != nil {
			return err
		}
	}

	dir := filepath.Join(database.Dir(), "global")
	// Replace cluster configuration in dqlite.
	err = dqlite.ReconfigureMembershipExt(dir, nodes)
	if err != nil {
		return fmt.Errorf("Failed to recover database state: %w", err)
	}

	// Replace cluster configuration in local raft_nodes database.
	err = database.Transaction(context.TODO(), func(ctx context.Context, tx *db.NodeTx) error {
		return tx.ReplaceRaftNodes(raftNodes)
	})
	if err != nil {
		return err
	}

	// Create patch file for global nodes database.
	content := ""
	for _, node := range nodes {
		content += fmt.Sprintf("UPDATE nodes SET address = %q WHERE id = %d;\n", node.Address, node.ID)
	}

	if len(content) > 0 {
		filePath := filepath.Join(database.Dir(), "patch.global.sql")
		file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		defer func() { _ = file.Close() }()

		_, err = file.Write([]byte(content))
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// RemoveRaftNode removes a raft node from the raft configuration.
func RemoveRaftNode(gateway *Gateway, address string) error {
	nodes, err := gateway.currentRaftNodes()
	if err != nil {
		return fmt.Errorf("Failed to get current raft nodes: %w", err)
	}

	var id uint64
	for _, node := range nodes {
		if node.Address == address {
			id = node.ID
			break
		}
	}
	if id == 0 {
		return fmt.Errorf("No raft node with address %q", address)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	client, err := client.FindLeader(
		ctx, gateway.NodeStore(),
		client.WithDialFunc(gateway.raftDial()),
		client.WithLogFunc(DqliteLog),
	)
	if err != nil {
		return fmt.Errorf("Failed to connect to cluster leader: %w", err)
	}

	defer func() { _ = client.Close() }()
	err = client.Remove(ctx, id)
	if err != nil {
		return fmt.Errorf("Failed to remove node: %w", err)
	}

	return nil
}
