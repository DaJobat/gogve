package tree

import (
	"fmt"

	"github.com/DaJobat/gogve/util"
)

type baseBSTNode struct {
	key         util.Comparable
	predecessor BSTNode
	left        BSTNode
	right       BSTNode
}

func NewBSTNode(key util.Comparable) *baseBSTNode {
	b := baseBSTNode{
		key: key,
	}
	return &b
}

func (b *baseBSTNode) Key() util.Comparable {
	return b.key
}

func (b *baseBSTNode) Parent() BSTNode {
	return b.predecessor
}

func (b *baseBSTNode) SetParent(node BSTNode) {
	b.predecessor = node
}

func (b *baseBSTNode) Left() BSTNode {
	return b.left
}

func (b *baseBSTNode) SetLeft(node BSTNode) {
	b.left = node
}

func (b *baseBSTNode) Right() BSTNode {
	return b.right
}

func (b *baseBSTNode) SetRight(node BSTNode) {
	b.right = node
}

func (n *baseBSTNode) String() string {
	return fmt.Sprintf("\t%v\n[%v\t%v]", n.key, n.left, n.right)
}

func (n *baseBSTNode) Value() interface{} {
	return n.key
}

func (n *baseBSTNode) Minimum() BSTNode {
	var min BSTNode = n
	for min.Left() != nil && !min.Left().Nil() {
		min = min.Left()
	}
	return min
}

func (n *baseBSTNode) Maximum() BSTNode {
	var max BSTNode = n
	for max.Right() != nil && !max.Right().Nil() {
		max = max.Right()
	}
	return max
}

func (n *baseBSTNode) Predecessor() BSTNode {
	var node BSTNode = n
	if node.Left() != nil && !node.Left().Nil() {
		return node.Left().Maximum()
	}

	predecessor := node.Parent()
	for predecessor != nil && !predecessor.Nil() && node == predecessor.Left() {
		node = predecessor
		predecessor = node.Parent()
	}

	return predecessor
}

func (n *baseBSTNode) Successor() BSTNode {
	var node BSTNode = n
	if node.Right() != nil && !node.Right().Nil() {
		return node.Right().Minimum()
	}

	successor := node.Parent()
	for successor != nil && !successor.Nil() && node == successor.Right() {
		node = successor
		successor = node.Parent()
	}
	return successor
}

func (n *baseBSTNode) Nil() bool {
	return n == nil
}

type binarySearchTree struct {
	rootNode BSTNode
}

func NewBinarySearchTree() BinarySearchTree {
	return &binarySearchTree{}
}

func (b *binarySearchTree) Search(key util.Comparable) (node BSTNode) {
	node, _ = BSTNodeSearch(b.Root(), key)
	return node
}

func (b *binarySearchTree) Root() BSTNode {
	return b.rootNode
}

func (b *binarySearchTree) Length() int {
	i := 0
	BSTNodeWalk(b.Root(), func(node BSTNode) bool {
		i++
		return false
	})
	return i
}

func (b *binarySearchTree) Insert(key util.Comparable) BSTNode {
	node := NewBSTNode(key)
	b.insert(node)
	return node
}

func (b *binarySearchTree) insert(node BSTNode) {
	n, p := BSTNodeSearch(b.Root(), node.Key())
	if n != nil && !n.Nil() {
		// we don't want duplicate keys in this basic bst
		return
	}

	node.SetParent(p)
	if p == nil || p.Nil() {
		b.rootNode = node
	} else if node.Key().Compare(p.Key()) == util.ComparableLess {
		p.SetLeft(node)
	} else {
		p.SetRight(node)
	}
}

func (b *binarySearchTree) String() string {
	return fmt.Sprint(b.rootNode)
}

func (b *binarySearchTree) Transplant(u, v BSTNode) {
	if u.Parent() == nil || u.Parent().Nil() {
		b.rootNode = v
	} else if u.Parent().Left() == u {
		u.Parent().SetLeft(v)
	} else {
		u.Parent().SetRight(v)
	}
	if v != nil && !v.Nil() {
		v.SetParent(u.Parent())
	}
}

func (b *binarySearchTree) Delete(node BSTNode) {
	if node.Left() == nil || node.Left().Nil() {
		b.Transplant(node, node.Right())
	} else if node.Right() == nil || node.Right().Nil() {
		b.Transplant(node, node.Left())
	} else {
		minRight := node.Right().Minimum()
		if minRight.Parent() != node {
			b.Transplant(minRight, minRight.Right())
			minRight.SetRight(node.Right())
			minRight.Right().SetParent(minRight)
		}
		b.Transplant(node, minRight)
		minRight.SetLeft(node.Left())
		minRight.Left().SetParent(minRight)
	}
}

type BaseBST binarySearchTree
