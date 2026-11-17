/*
BST complexity
Operation	Balanced BST	Unbalanced BST worst case
Add server	    O(log N)	    O(N)
Get server	    O(log N)	    O(N)
Remove server	O(log N)	    O(N)
Space	O(N)	O(N)

You can even use a simpler implementation like storing the hash values of server
in a sorted array and do a linear traversal to find the node hash which is greater
than server hash for a clockwise implementation.
*/

package main

import (
	"fmt"
	"hash/fnv"
)

type TreeNode struct {
	Hash   uint32
	server string
	left   *TreeNode
	right  *TreeNode
}

type BSTHashRing struct {
	Root *TreeNode
}

func hash(key string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func insert(root *TreeNode, hash uint32, server string) *TreeNode {
	if root == nil {
		return &TreeNode{server: server, Hash: hash}
	}
	if hash < root.Hash {
		root.left = insert(root.left, hash, server)
	} else if hash > root.Hash {
		root.right = insert(root.right, hash, server)
	}
	return root
}

func (r *BSTHashRing) addServer(server string) {
	for i := range 100 {
		// Consistent hashing does not guarantee equal distribution, especially when you have only
		// three servers and no virtual nodes. Virtual nodes helps in even distribution of load.
		virtualNode := fmt.Sprintf("%s#%d", server, i)
		r.Root = insert(r.Root, hash(virtualNode), server)
	}
}

func (r *BSTHashRing) getServer(server string) string {
	if r.Root == nil {
		return ""
	}
	h := hash(server)
	curr := r.Root
	var successor *TreeNode

	for curr != nil {
		if curr.Hash >= h {
			successor = curr
			curr = curr.left
		} else {
			curr = curr.right
		}
	}

	if successor == nil {
		curr = r.Root
		for curr.left != nil {
			curr = curr.left
		}
		successor = curr
	}
	return successor.server
}

func main() {
	ring := &BSTHashRing{}

	ring.addServer("server-A")
	ring.addServer("server-B")
	ring.addServer("server-C")
	ring.addServer("server-D")

	for _, key := range []string{"user1", "user2", "user3", "user-D", "user-E"} {
		fmt.Println(key, "->", ring.getServer(key))
	}
}
