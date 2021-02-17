package tree

import (
	"github.com/DaJobat/gogve/util"
)

type BinarySearchTree interface {
	Search(key util.Comparable) BSTNode
	Root() BSTNode
	Insert(util.Comparable) BSTNode //If this finds that the comparable is a BSTNode it'll use it
	//otherwise, it'll wrap it in a BSTNode
	Delete(BSTNode)
	Length() int
}

type BSTNode interface {
	Key() util.Comparable
	Successor() BSTNode
	Predecessor() BSTNode
	Parent() BSTNode
	SetParent(BSTNode)
	Left() BSTNode
	SetLeft(BSTNode)
	Right() BSTNode
	SetRight(BSTNode)
	Minimum() BSTNode
	Maximum() BSTNode
	Nil() bool
}

func BSTNodeSearch(node BSTNode, key util.Comparable) (cNode, preNode BSTNode) {
	cNode = node
	preNode = nil

	for cNode != nil && !cNode.Nil() && cNode.Key().Compare(key) != util.ComparableEqual {
		preNode = cNode
		if key.Compare(cNode.Key()) == util.ComparableLess {
			cNode = cNode.Left()
		} else {
			cNode = cNode.Right()
		}
	}
	return cNode, preNode
}

type BSTWalkCallback func(BSTNode) (stop bool)

func BSTNodeWalk(node BSTNode, callback BSTWalkCallback) {
	if callback == nil {
		panic("no callback")
	}

	if node != nil {
		BSTNodeWalk(node.Left(), callback)
		if callback(node) { //if we're told to stop, stop.
			return
		}
		BSTNodeWalk(node.Right(), callback)
	}
}
