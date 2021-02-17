package tree

import (
	"fmt"

	"github.com/DaJobat/gogve/util"
)

type rbTreeColor int

func (c rbTreeColor) String() string {
	switch c {
	case rbTreeColorRed:
		return "Red"
	case rbTreeColorBlack:
		return "Black"
	default:
		panic("invalid color")
	}
}

const (
	rbTreeColorRed rbTreeColor = iota
	rbTreeColorBlack
)

type RBTreeNode interface {
	BSTNode
	rbParent() RBTreeNode
	rbLeft() RBTreeNode
	rbRight() RBTreeNode
	Color() rbTreeColor
	setColor(rbTreeColor)
}

type baseRBNode struct {
	key    util.Comparable
	parent RBTreeNode
	left   RBTreeNode
	right  RBTreeNode
	color  rbTreeColor
}

func (r *baseRBNode) String() string {
	if r == nil || r.Nil() {
		return "Nil"
	}

	if (r.Left() == nil || r.Left().Nil()) && (r.Right() == nil || r.Right().Nil()) {
		return fmt.Sprintf("Node: %v[LEAF]", r.color)
	}

	return fmt.Sprintf("Node:  %v \n\t[L{ %v },\tR{ %v }]", r.color, r.left, r.right)
}

func (r *baseRBNode) Key() util.Comparable {
	return r.key
}

func (r *baseRBNode) Color() rbTreeColor {
	return r.color
}

func (r *baseRBNode) setColor(c rbTreeColor) {
	r.color = c
}

func (r *baseRBNode) Parent() BSTNode {
	if r.parent == nil { //We don't want to return nil for the nil node, so don't use .Nil()
		return nil
	}
	return r.parent
}

func (r *baseRBNode) rbParent() RBTreeNode {
	if r.parent == nil { //We don't want to return nil for the nil node, so don't use .Nil()
		return nil
	}
	return r.parent
}

func (r *baseRBNode) SetParent(node BSTNode) {
	if node == nil {
		r.parent = nil
	} else {
		r.parent = node.(RBTreeNode)
	}
}

func (r *baseRBNode) Left() BSTNode {
	if r.left == nil { //We don't want to return nil for the nil node, so don't use .Nil()
		return nil
	}
	return r.left
}

func (r *baseRBNode) rbLeft() RBTreeNode {
	if r.left == nil { //We don't want to return nil for the nil node, so don't use .Nil()
		return nil
	}
	return r.left
}

func (r *baseRBNode) SetLeft(node BSTNode) {
	r.left = node.(RBTreeNode)
}

func (r *baseRBNode) Right() BSTNode {
	if r.right == nil { //We don't want to return nil for the nil node, so don't use .Nil()
		return nil
	}
	return r.right
}

func (r *baseRBNode) SetRight(node BSTNode) {
	r.right = node.(RBTreeNode)
}

func (r *baseRBNode) rbRight() RBTreeNode {
	if r.right == nil { //We don't want to return nil for the nil node, so don't use .Nil()
		return nil
	}
	return r.right
}

func (r *baseRBNode) Minimum() BSTNode {
	var min BSTNode = r
	for !min.Left().Nil() {
		min = min.Left()
	}
	return min

}

func (r *baseRBNode) Maximum() BSTNode {
	var max BSTNode = r
	for !max.Right().Nil() {
		max = max.Right()
	}
	return max
}

func (r *baseRBNode) Predecessor() BSTNode {
	var node BSTNode = r
	if !node.Left().Nil() {
		return node.Left().Maximum()
	}

	predecessor := node.Parent()
	for !predecessor.Nil() && node == predecessor.Left() {
		node = predecessor
		predecessor = node.Parent()
	}

	return predecessor
}

func (r *baseRBNode) Successor() BSTNode {
	var node BSTNode = r
	if !node.Right().Nil() {
		return node.Right().Minimum()
	}

	successor := node.Parent()
	for !successor.Nil() && node == successor.Right() {
		node = successor
		successor = node.Parent()
	}
	return successor
}

func (r *baseRBNode) Nil() bool {
	return r == nil || r == rbNilNode
}

func NewRBTreeNode(key util.Comparable) RBTreeNode {
	tn := baseRBNode{
		key: key,
	}
	return &tn
}

var rbNilNode RBTreeNode = func() RBTreeNode {
	tn := NewRBTreeNode(&NilComparable{})
	tn.setColor(rbTreeColorBlack)
	return tn
}()

type RBTree struct {
	rootNode RBTreeNode
}

type NilComparable struct{}

func (nc *NilComparable) Compare(util.Comparable) util.ComparableResult {
	return util.ComparableLess
}

func NewRBTree() BinarySearchTree {
	rbt := RBTree{
		rootNode: rbNilNode,
	}

	return &rbt
}

func (r *RBTree) Search(key util.Comparable) BSTNode {
	n, _ := BSTNodeSearch(r.rootNode, key)
	if n.Nil() {
		return nil
	} else {
		return n
	}
}

func (r *RBTree) Root() BSTNode {
	return r.rootNode
}

func (r *RBTree) leftRotate(node RBTreeNode) {
	oldRight := node.rbRight()

	//Make the old right node's left subtree our right subtree
	node.SetRight(oldRight.Left())
	if oldRight.Left() != rbNilNode {
		oldRight.Left().SetParent(node)
	}

	//Make our parent our old right node's parent
	oldRight.SetParent(node.Parent())
	if node.Parent() == rbNilNode {
		r.rootNode = oldRight
	} else if node == node.Parent().Left() {
		node.Parent().SetLeft(oldRight)
	} else {
		node.Parent().SetRight(oldRight)
	}

	//Make us our old right node's left subtree
	oldRight.SetLeft(node)
	node.SetParent(oldRight)
}

func (r *RBTree) rightRotate(node RBTreeNode) {
	oldLeft := node.rbLeft()

	//Make the old left node's right subtree our left subtree
	node.SetLeft(oldLeft.Right())
	if !oldLeft.Right().Nil() {
		oldLeft.Right().SetParent(node)
	}
	oldLeft.SetParent(node.Parent())

	//Make our pre our old left node's pre
	if node.Parent().Nil() {
		r.rootNode = oldLeft
	} else if node == node.Parent().Right() {
		node.Parent().SetRight(oldLeft)
	} else {
		node.Parent().SetLeft(oldLeft)
	}

	// Make us our old left node's right subtree
	oldLeft.SetRight(node)
	node.SetParent(oldLeft)
}

func (r *RBTree) Insert(key util.Comparable) BSTNode {
	node, ok := key.(RBTreeNode)
	if !ok { //if the key isn't already a treenode, wrap it in one
		node = NewRBTreeNode(key)
	}
	r.insert(node)
	return node
}

func (r *RBTree) insert(node RBTreeNode) {
	insertNode := node
	var y RBTreeNode = rbNilNode
	x := r.rootNode

	for !x.Nil() {
		y = x
		if insertNode.Key().Compare(x.Key()) == util.ComparableLess {
			x = x.rbLeft()
		} else {
			x = x.rbRight()
		}
	}

	insertNode.SetParent(y)

	if y.Nil() { //Node is at the root of the tree
		r.rootNode = insertNode
	} else if insertNode.Key().Compare(y.Key()) == util.ComparableLess {
		y.SetLeft(insertNode)
	} else {
		y.SetRight(insertNode)
	}

	insertNode.SetLeft(rbNilNode)
	insertNode.SetRight(rbNilNode)
	insertNode.setColor(rbTreeColorRed)
	r.insertFixup(insertNode)
}

func (r *RBTree) insertFixup(node RBTreeNode) {

	for node.rbParent().Color() == rbTreeColorRed {

		if node.rbParent() == node.rbParent().rbParent().Left() { //parent is left

			if node.rbParent().rbParent().rbRight().Color() == rbTreeColorRed {
				node.rbParent().setColor(rbTreeColorBlack)
				node.rbParent().rbParent().rbRight().setColor(rbTreeColorBlack)
				node.rbParent().rbParent().setColor(rbTreeColorRed)
				node = node.rbParent().rbParent()
			} else { //grandparentRight is black

				if node == node.Parent().Right() { //i am right hand side
					node = node.rbParent()
					r.leftRotate(node)
				}

				node.rbParent().setColor(rbTreeColorBlack)
				node.rbParent().rbParent().setColor(rbTreeColorRed)
				r.rightRotate(node.rbParent().rbParent())
			}

		} else { //parent is right
			if node.rbParent().rbParent().rbLeft().Color() == rbTreeColorRed {
				node.rbParent().setColor(rbTreeColorBlack)
				node.rbParent().rbParent().rbLeft().setColor(rbTreeColorBlack)
				node.rbParent().rbParent().setColor(rbTreeColorRed)
				node = node.rbParent().rbParent()
			} else {
				if node == node.Parent().Left() {
					node = node.rbParent()
					r.rightRotate(node)
				}

				node.rbParent().setColor(rbTreeColorBlack)
				node.rbParent().rbParent().setColor(rbTreeColorRed)
				r.leftRotate(node.rbParent().rbParent())
			}
		}
	}
	r.rootNode.setColor(rbTreeColorBlack)
}

func (r *RBTree) Transplant(u, v BSTNode) {
	if u.Parent() == rbNilNode {
		r.rootNode = v.(*baseRBNode)
	} else if u == u.Parent().Left() {
		u.Parent().SetLeft(v)
	} else {
		u.Parent().SetRight(v)
	}

	v.SetParent(u.Parent())
}

//Replace transplants a node in the tree, but grafts the children of the
//transplanted node into the new node
func (r *RBTree) Replace(u, v BSTNode) {
	r.Transplant(u, v)
	v.SetLeft(u.Left())
	v.SetRight(u.Right())

	if v.Left() != nil && !v.Left().Nil() {
		v.Left().SetParent(v)
	}
	if v.Right() != nil && !v.Right().Nil() {
		v.Right().SetParent(v)
	}
	v.(RBTreeNode).setColor(u.(RBTreeNode).Color())
}

func (r *RBTree) Delete(node BSTNode) {
	nodeBeingRemoved := node.(RBTreeNode)
	y := node.(RBTreeNode)

	yOrigColor := y.Color()
	var nodeAtOriginalY RBTreeNode

	if nodeBeingRemoved.Left() == rbNilNode {
		// if the node being removed has no left child,
		// overwrite the node with its entire right hand subtree
		nodeAtOriginalY = nodeBeingRemoved.rbRight()
		r.Transplant(nodeBeingRemoved, nodeBeingRemoved.Right())
	} else if nodeBeingRemoved.Right() == rbNilNode {
		// if the node being removed has no right child,
		// overwrite the node with its entire left hand subtree
		nodeAtOriginalY = nodeBeingRemoved.rbLeft()
		r.Transplant(nodeBeingRemoved, nodeBeingRemoved.Left())
	} else {
		//the node to be removed has two children
		y = nodeBeingRemoved.Right().Minimum().(RBTreeNode)
		yOrigColor = y.Color()
		nodeAtOriginalY = y.rbRight()

		if y.Parent() == nodeBeingRemoved {
			nodeAtOriginalY.SetParent(y)
		} else {
			r.Transplant(y, y.Right())
			y.SetRight(nodeBeingRemoved.Right())
			y.Right().SetParent(y)
		}

		r.Transplant(nodeBeingRemoved, y)
		y.SetLeft(nodeBeingRemoved.Left())
		y.Left().SetParent(y)
		y.setColor(nodeBeingRemoved.Color())
	}

	if yOrigColor == rbTreeColorBlack {
		r.deleteFixup(nodeAtOriginalY)
	}

	rbNilNode.SetParent(nil)
}

func (r *RBTree) deleteFixup(node RBTreeNode) {
	for node != r.rootNode && node.Color() == rbTreeColorBlack {
		if node == node.Parent().Left() {
			w := node.rbParent().rbRight()
			if w.Color() == rbTreeColorRed {
				w.setColor(rbTreeColorBlack)
				node.rbParent().setColor(rbTreeColorRed)
				r.leftRotate(node.rbParent())
				w = node.rbParent().rbRight()
			}
			//FIXME: Added in this nilnode check for an edge case when the tree has a black root 2 with left child 1 red and right child 3 black, causing w here to be the nilnode and thefore no left. mirrored below
			if w == rbNilNode || (w.rbLeft().Color() == rbTreeColorBlack && w.rbRight().Color() == rbTreeColorBlack) {
				w.setColor(rbTreeColorRed)
				node = node.rbParent()
			} else {
				if w.rbRight().Color() == rbTreeColorBlack {
					w.rbLeft().setColor(rbTreeColorBlack)
					w.setColor(rbTreeColorRed)
					r.rightRotate(w)
					w = node.rbParent().rbRight()
				}
				w.setColor(node.rbParent().Color())
				node.rbParent().setColor(rbTreeColorBlack)
				w.rbRight().setColor(rbTreeColorBlack)
				r.leftRotate(node.rbParent())
				node = r.rootNode
			}
		} else {
			w := node.rbParent().rbLeft()
			if w.Color() == rbTreeColorRed {
				w.setColor(rbTreeColorBlack)
				node.rbParent().setColor(rbTreeColorRed)
				r.rightRotate(node.rbParent())
				w = node.rbParent().rbLeft()
			}
			if w == rbNilNode || (w.rbRight().Color() == rbTreeColorBlack && w.rbLeft().Color() == rbTreeColorBlack) {
				w.setColor(rbTreeColorRed)
				node = node.rbParent()
			} else {
				if w.rbLeft().Color() == rbTreeColorBlack {
					w.rbRight().setColor(rbTreeColorBlack)
					w.setColor(rbTreeColorRed)
					r.leftRotate(w)
					w = node.rbParent().rbLeft()
				}
				w.setColor(node.rbParent().Color())
				node.rbParent().setColor(rbTreeColorBlack)
				w.rbLeft().setColor(rbTreeColorBlack)
				r.rightRotate(node.rbParent())
				node = r.rootNode
			}
		}
	}
	node.setColor(rbTreeColorBlack)
}

func (r *RBTree) Length() int {
	i := 0
	BSTNodeWalk(r.Root(), func(b BSTNode) bool {
		if !b.Nil() {
			i++
		}
		return false
	})
	return i
}

func (r *RBTree) String() string {
	return fmt.Sprint(r.Root())
}
