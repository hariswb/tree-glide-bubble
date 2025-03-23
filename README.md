# Introduction

Treeglide is a fork of [tree-bubble](https://github.com/savannahostrowski/tree-bubble) module. It affords:
- Cursor navigation between tree's parent and its children.
- Tree view like reddit comment (future).
- Collapsible (future).

# Challenges

## Tree instantiation

Initially, tree initialization is verbose as it requires pointer to the `Node` for parents and childrens, not the Node's value.

```go
type Node struct {
	Value    string
	Parent   *Node
	Children []*Node
}

root := treeglide.Node{
    Value:    "root",
    Parent:   nil,
    Children: nil,
}

nodeA := treeglide.Node{
    Value:    "A",
    Parent:   &root,
    Children: nil,
}

nodeA1 := treeglide.Node{
    Value:    "A1",
    Parent:   &nodeA1,
    Children: nil,
}

nodeB := treeglide.Node{
    Value:    "B",
    Parent:   &root,
    Children: nil,
}

nodeA.Children = append(nodeA.Children, &nodeA1)

root.Children = append(root.Children,&nodeA, &nodeB, &nodeC)

tree = treeglide.New(&root, w, h)
```

 It's a hassle compared to tree-bubble's [implementation](https://github.com/savannahostrowski/tree-bubble/blob/main/example/main.go) which seems to me more intuitive.
```go
// Tree bubble's Node:
type Node struct {
	Value    string
	Desc     string
	Children []Node
}

nodes := []tree.Node{
    {
        Value: "history | grep docker",
        Desc: "Used in a Unix-like operating system to search through the " +
            "command history for any entries that contain the word 'docker.'",
        Children: []tree.Node{
            {
            Value:    "history",
            Children: nil,
        }, 
        {
            Value:    "|",
            Children: nil,
        },
    },
}};

tree:= tree.New(nodes, w, h)
```

My initial solution was to use tree-bubbles's `Node` structure for user's initialization and convert it to treeglide's `Node` structure, but it adds even more complexity.

To solve this, I add a helper for user to easily initializes their nodes and assign their parent's. It's also convenient for user to iterate their own tree structure to fit into treeglide's.
```go

func NewNode(value string, parent *Node) *Node {
	node := &Node{Value: value, Parent: parent}
	if parent != nil {
		parent.Children = append(parent.Children, node)
	}
	return node
}

root := treeglide.NewNode("root", nil)

nodeA := treeglide.NewNode("A", root)
nodeA1 := treeglide.NewNode("A1", nodeA)
nodeB := treeglide.NewNode("B", root)

tree:= treeglide.NewTree(root, w, h)

```


