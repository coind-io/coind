package hash

import (
	"errors"
)

type MerkleTree struct {
	Depth uint
	Root  *MerkleTreeNode
}

type MerkleTreeNode struct {
	Hash  *Hash
	Left  *MerkleTreeNode
	Right *MerkleTreeNode
}

func NewMerkleTree(hashes []*Hash) (*MerkleTree, error) {
	if len(hashes) == 0 {
		return nil, errors.New("NewMerkleTree input no item error.")
	}
	mt := &MerkleTree{
		Depth: 1,
	}
	nodes := mt.generateLeaves(hashes)
	for len(nodes) > 1 {
		nodes = mt.levelUp(nodes)
		mt.Depth = mt.Depth + 1
	}
	mt.Root = nodes[0]
	return mt, nil
}

func (mt *MerkleTree) generateLeaves(hashes []*Hash) []*MerkleTreeNode {
	var leaves []*MerkleTreeNode
	for _, hash := range hashes {
		node := &MerkleTreeNode{
			Hash: hash,
		}
		leaves = append(leaves, node)
	}
	return leaves
}

func (mt *MerkleTree) levelUp(nodes []*MerkleTreeNode) []*MerkleTreeNode {
	var nextLevel []*MerkleTreeNode
	for i := 0; i < len(nodes)/2; i++ {
		data := []byte{}
		data = append(data, nodes[i*2].Hash.bytes...)
		data = append(data, nodes[i*2+1].Hash.bytes...)
		hash := SumDoubleHash256(data)
		node := &MerkleTreeNode{
			Hash:  hash,
			Left:  nodes[i*2],
			Right: nodes[i*2+1],
		}
		nextLevel = append(nextLevel, node)
	}
	if len(nodes)%2 == 1 {
		data := []byte{}
		data = append(data, nodes[len(nodes)-1].Hash.bytes...)
		data = append(data, nodes[len(nodes)-1].Hash.bytes...)
		hash := SumDoubleHash256(data)
		node := &MerkleTreeNode{
			Hash:  hash,
			Left:  nodes[len(nodes)-1],
			Right: nodes[len(nodes)-1],
		}
		nextLevel = append(nextLevel, node)
	}
	return nextLevel
}
