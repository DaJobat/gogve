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

type RBTreeNode struct {
	color  rbTreeColor
	key    util.Comparable
	parent *RBTreeNode
	left   *RBTreeNode
	right  *RBTreeNode
}

func NewRBTreeNode(key util.Comparable) BSTNode {
	tn := RBTreeNode{
		key: key,
	}
	return &tn
}

func (r *RBTreeNode) String() string {
	if r.Nil() {
		return "Nil"
	}

	if r.Left().Nil() && r.Right().Nil() {
		return fmt.Sprintf("Node: %v, %v[LEAF]", r.key, r.color)
	}

	return fmt.Sprintf("Node: %v, %v \n\t[L{ %v },\tR{ %v }]", r.key, r.color, r.left, r.right)
}

func (r *RBTreeNode) Key() util.Comparable {
	return r.key
}

func (r *RBTreeNode) Parent() BSTNode {
	if r.parent == nil { //We don't want to return nil for the nil node, so don't use .Nil()
		return nil
	}
	return r.parent
}

func (r *RBTreeNode) SetParent(node BSTNode) {
	r.parent = node.(*RBTreeNode)
}

func (r *RBTreeNode) Left() BSTNode {
	if r.left == nil { //We don't want to return nil for the nil node, so don't use .Nil()
		return nil
	}
	return r.left
}

func (r *RBTreeNode) SetLeft(node BSTNode) {
	r.left = node.(*RBTreeNode)
}

func (r *RBTreeNode) Right() BSTNode {
	if r.right == nil { //We don't want to return nil for the nil node, so don't use .Nil()
		return nil
	}
	return r.right
}

func (r *RBTreeNode) SetRight(node BSTNode) {
	r.right = node.(*RBTreeNode)
}

func (r *RBTreeNode) Minimum() BSTNode {
	var min BSTNode = r
	for !min.Left().Nil() {
		min = min.Left()
	}
	return min

}

func (r *RBTreeNode) Maximum() BSTNode {
	var max BSTNode = r
	for !max.Right().Nil() {
		max = max.Right()
	}
	return max
}

func (r *RBTreeNode) Predecessor() BSTNode {
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

func (r *RBTreeNode) Successor() BSTNode {
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

func (r *RBTreeNode) Nil() bool {
	return r == nil || r == rbNilNode
}

type RBTree struct {
	rootNode *RBTreeNode
}

var rbNilNode = &RBTreeNode{
	key:    NilComparable{},
	color:  rbTreeColorBlack,
	left:   nil,
	right:  nil,
	parent: nil,
}

type NilComparable struct{}

func (nc NilComparable) Compare(util.Comparable) util.ComparableResult {
	return util.ComparableLess
}

func NewRBTree() BinarySearchTree {
	rbt := &RBTree{}
	rbt.rootNode = rbNilNode

	return rbt
}

/*
func (r *RBTree) String() string {
	b := &strings.Builder{}
	return b.String()
}
*/
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

func (r *RBTree) leftRotate(node *RBTreeNode) {
	oldRight := node.right

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

func (r *RBTree) rightRotate(node *RBTreeNode) {
	oldLeft := node.left

	//Make the old left node's right subtree our left subtree
	node.SetLeft(oldLeft.right)
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
	node := NewRBTreeNode(key)
	r.insert(node)
	return node
}

func (r *RBTree) insert(node BSTNode) {
	insertNode := node.(*RBTreeNode)
	y := rbNilNode
	x := r.rootNode

	for x != rbNilNode {
		y = x
		if insertNode.Key().Compare(x.Key()) == util.ComparableLess {
			x = x.left
		} else {
			x = x.right
		}
	}

	insertNode.SetParent(y)

	if y == rbNilNode { //Node is at the root of the tree
		r.rootNode = insertNode
	} else if insertNode.Key().Compare(y.Key()) == util.ComparableLess {
		y.SetLeft(insertNode)
	} else {
		y.SetRight(insertNode)
	}

	insertNode.SetLeft(rbNilNode)
	insertNode.SetRight(rbNilNode)
	insertNode.color = rbTreeColorRed
	fmt.Printf("Tree PostInsert\n\n %v \n\n", r)
	r.insertFixup(insertNode)
}

func (r *RBTree) insertFixup(node *RBTreeNode) {
	for node.parent.color == rbTreeColorRed {
		fmt.Println("MY PARENT IS RED")
		if node.parent == node.parent.parent.left {
			fmt.Println("MY PARENT IS LEFT")
			if node.parent.parent.right.color == rbTreeColorRed {
				fmt.Println("MY PARPAR RIGHT IS RED")
				node.parent.color = rbTreeColorBlack
				node.parent.parent.right.color = rbTreeColorBlack
				node.parent.parent.color = rbTreeColorRed
				node = node.parent.parent
			} else {
				fmt.Println("MY PARPAR RIGHT IS BLACK")
				if node == node.parent.right {
					fmt.Println("LEFT ROTATE MY PARENT")
					node = node.parent
					r.leftRotate(node)
				}
				node.parent.color = rbTreeColorBlack
				node.parent.parent.color = rbTreeColorRed
				fmt.Println("RIGHT ROTATE MY PAR PAR PAR")
				r.rightRotate(node.parent.parent)
			}
		} else { //do the opposite of above
			fmt.Println("MY PARENT IS RIGHT")
			if node.parent.parent.left.color == rbTreeColorRed {
				fmt.Println("MY PARPAR LEFT IS RED")
				node.parent.color = rbTreeColorBlack
				node.parent.parent.left.color = rbTreeColorRed
				node.parent.parent.color = rbTreeColorRed
				node = node.parent.parent
			} else {
				fmt.Println("MY PAR LEFT IS BLACK")
				if node == node.parent.left {
					fmt.Println("RIGHT ROTATE")
					node = node.parent
					r.rightRotate(node)
				}

				node.parent.color = rbTreeColorBlack
				node.parent.parent.color = rbTreeColorRed
				fmt.Println("LEFT ROTATE")
				r.leftRotate(node.parent.parent)
			}
		}
	}
	r.rootNode.color = rbTreeColorBlack
}

func (r *RBTree) Transplant(u, v BSTNode) {
	if u.Parent() == rbNilNode {
		r.rootNode = v.(*RBTreeNode)
	} else if u == u.Parent().Left() {
		u.Parent().SetLeft(v)
	} else {
		u.Parent().SetRight(v)
	}

	v.SetParent(u.Parent())
}

func (r *RBTree) Delete(node BSTNode) {
	nodeBeingRemoved := node.(*RBTreeNode)
	y := node.(*RBTreeNode)

	yOrigColor := y.color
	var nodeAtOriginalY *RBTreeNode

	if nodeBeingRemoved.Left() == rbNilNode {
		// if the node being removed has no left child,
		// overwrite the node with its entire right hand subtree
		nodeAtOriginalY = nodeBeingRemoved.right
		r.Transplant(nodeBeingRemoved, nodeBeingRemoved.Right())
	} else if nodeBeingRemoved.right == rbNilNode {
		// if the node being removed has no right child,
		// overwrite the node with its entire left hand subtree
		nodeAtOriginalY = nodeBeingRemoved.left
		r.Transplant(nodeBeingRemoved, nodeBeingRemoved.Left())
	} else {
		//the node to be removed has two children
		y = nodeBeingRemoved.Right().Minimum().(*RBTreeNode)
		yOrigColor = y.color
		nodeAtOriginalY = y.right

		if y.Parent() == nodeBeingRemoved {
			nodeAtOriginalY.SetParent(y)
		} else {
			r.Transplant(y, y.right)
			y.right = nodeBeingRemoved.right
			y.right.SetParent(y)
		}

		r.Transplant(nodeBeingRemoved, y)
		y.left = nodeBeingRemoved.left
		y.left.SetParent(y)
		y.color = nodeBeingRemoved.color
	}

	if yOrigColor == rbTreeColorBlack {
		r.deleteFixup(nodeAtOriginalY)
	}

	rbNilNode.parent = nil
}

func (r *RBTree) deleteFixup(node *RBTreeNode) {
	for node != r.rootNode && node.color == rbTreeColorBlack {
		if node == node.parent.Left() {
			w := node.parent.right
			if w.color == rbTreeColorRed {
				w.color = rbTreeColorBlack
				node.parent.color = rbTreeColorRed
				r.leftRotate(node.parent)
				w = node.parent.right
			}
			//FIXME: Added in this nilnode check for an edge case when the tree has a black root 2 with left child 1 red and right child 3 black, causing w here to be the nilnode and thefore no left. mirrored below
			if w == rbNilNode || (w.left.color == rbTreeColorBlack && w.right.color == rbTreeColorBlack) {
				w.color = rbTreeColorRed
				node = node.parent
			} else {
				if w.right.color == rbTreeColorBlack {
					w.left.color = rbTreeColorBlack
					w.color = rbTreeColorRed
					r.rightRotate(w)
					w = node.parent.right
				}
				w.color = node.parent.color
				node.parent.color = rbTreeColorBlack
				w.right.color = rbTreeColorBlack
				r.leftRotate(node.parent)
				node = r.rootNode
			}
		} else {
			w := node.parent.left
			if w.color == rbTreeColorRed {
				w.color = rbTreeColorBlack
				node.parent.color = rbTreeColorRed
				r.rightRotate(node.parent)
				w = node.parent.left
			}
			if w == rbNilNode || (w.right.color == rbTreeColorBlack && w.left.color == rbTreeColorBlack) {
				w.color = rbTreeColorRed
				node = node.parent
			} else {
				if w.left.color == rbTreeColorBlack {
					w.right.color = rbTreeColorBlack
					w.color = rbTreeColorRed
					r.leftRotate(w)
					w = node.parent.left
				}
				w.color = node.parent.color
				node.parent.color = rbTreeColorBlack
				w.left.color = rbTreeColorBlack
				r.rightRotate(node.parent)
				node = r.rootNode
			}
		}
	}
	node.color = rbTreeColorBlack
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
