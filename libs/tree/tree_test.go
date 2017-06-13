package tree

import "testing"

func TestPath(t *testing.T) {
	tree := NewTree()
	node := tree.Add("/a/b/c", "asd")

	if node.Path() != "/a/b/c" {
		t.Errorf("Expected path /a/b/c but found %q instead", node.Path())
	}
	a := tree.Root.Children[0]
	b := a.Children[0]
	c := b.Children[0]
	if a.Name != "a" || b.Name != "b" || c.Name != "c" {
		t.Errorf("%+v, %+v, %+v", a, b, c)
	}
}

func TestTree(t *testing.T) {
	tree := NewTree()
	node := tree.Add("/foo", "bar")
	if node.Name != "foo" {
		t.Errorf("Expected node with name foo, got %q instead", node.Name)
	}
	if node.Value != "bar" {
		t.Errorf("Expected node with value bar, got %q instead", node.Value)
	}

	node = tree.Add("bar", "")
	if node.Name != "bar" {
		t.Errorf("Expected node with name bar, got %q instead", node.Name)
	}
}

func TestTree2(t *testing.T) {
	tree := NewTree()
	tree.Add("foo/bar/baz", "")

	child := tree.Root.Children[0]
	if child.Name != "foo" {
		t.Errorf("Expected first child to be called foo instead of %q", child.Name)
	}

	child2 := child.Children[0]
	if child2.Name != "bar" {
		t.Errorf("Expected second child to be called bar instead of %q", child2.Name)
	}
	child3 := child2.Children[0]
	if child3.Name != "baz" {
		t.Errorf("Expected third child to be called baz instead of %q", child3.Name)
	}
}

func TestGet(t *testing.T) {
	tree := NewTree()

	_ = tree.Add("foo/bar", "")
	_ = tree.Add("foo/baz", "")

	n := tree.Get("/foo/bar")
	if n.Name != "bar" {
		t.Errorf("Expected to Get bar instead of %q", n.Name)
	}

}
