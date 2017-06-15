package tree

import (
	"fmt"
	"strings"
	"time"
)

type Tree struct {
	Root *Node
}

type Node struct {
	Name     string
	Children []*Node
	Parent   *Node
	Value    string
	Updated  time.Time
}

func (t *Tree) Get(path string) *Node {
	path = strings.TrimPrefix(path, "/")
	child := t.Root
	if path == "" {
		return child
	}
	for _, token := range strings.Split(path, "/") {
		found := false
		for _, grandchild := range child.Children {
			if grandchild.Name == token {
				child = grandchild
				found = true
				break
			}
		}
		if !found {
			return nil
		}
	}
	return child
}

func (n *Node) Append(child *Node) *Node {
	n.Children = append(n.Children, child)
	child.Parent = n
	return child
}

func (n *Node) Path() string {
	if n.Parent == nil {
		return "/"
	} else {
		p := n.Parent.Path()
		if p == "/" {
			return "/" + n.Name
		} else {
			return p + "/" + n.Name
		}
	}
}

func (t *Tree) Add(path, val string) *Node {

	node := t.Get(path)
	if node != nil {
		node.Value = val
		node.Updated = time.Now()
		return node
	}

	tokens := strings.Split(strings.TrimPrefix(path, "/"), "/")
	child := Node{Name: tokens[len(tokens)-1], Value: val, Updated: time.Now()}

	node = t.Root
	for i := 0; i < len(tokens); i++ {
		found := false
		for _, children := range node.Children {
			if children.Name == tokens[i] {
				node = children
				found = true
				break
			}
		}
		if !found && i < len(tokens)-1 {
			node = node.Append(&Node{Name: tokens[i]})
		}
	}
	child.Parent = node
	child.Updated = time.Now()
	node.Append(&child)
	return &child
}

func (t *Tree) Delete(path string) *Node {
	node := t.Get(path)
	for i, child := range node.Parent.Children {
		if child.Name == node.Name {
			node.Parent.Children = append(node.Parent.Children[:i], node.Parent.Children[i+1:]...)
		}
	}
	return node
}

func (n *Node) Iterate(nodes chan<- *Node) {
	nodes <- n
	for _, c := range n.Children {
		c.Iterate(nodes)
	}
}

func (t *Tree) String() string {
	s := ""

	nodes := make(chan *Node)
	go func() {
		t.Root.Iterate(nodes)
		close(nodes)
	}()
	for c := range nodes {
		if c.Value != "" {
			s += fmt.Sprintf("%s = %s\n", c.Path(), c.Value)
		} else {
			s += fmt.Sprintf("%s\n", c.Path())
		}
	}
	return s

}

func NewTree() *Tree {
	return &Tree{Root: &Node{}}
}
