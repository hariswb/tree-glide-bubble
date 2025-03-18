package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"

	treeglide "github.com/hariswb/tree-glide-bubble"
)

var (
	styleDoc = lipgloss.NewStyle().Padding(1)
)

type model struct {
	tree treeglide.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}
	var cmd tea.Cmd
	m.tree, cmd = m.tree.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return styleDoc.Render(m.tree.View())
}

func main() {

	root := treeglide.Node{
		Value:    "root",
		Parent:   nil,
		Children: nil}

	nodeA := treeglide.Node{
		Value:    "A",
		Parent:   &root,
		Children: nil,
	}

	nodeA1 := treeglide.Node{
		Value:    "A1",
		Parent:   &nodeA,
		Children: nil,
	}

	nodeA11 := treeglide.Node{
		Value:    "A11",
		Parent:   &nodeA1,
		Children: nil,
	}

	nodeA12 := treeglide.Node{
		Value:    "A12",
		Parent:   &nodeA1,
		Children: nil,
	}

	nodeA13 := treeglide.Node{
		Value:    "A13",
		Parent:   &nodeA1,
		Children: nil,
	}

	nodeA2 := treeglide.Node{
		Value:    "A2",
		Parent:   &nodeA,
		Children: nil,
	}

	nodeB := treeglide.Node{
		Value:    "B",
		Parent:   &root,
		Children: nil,
	}

	nodeB1 := treeglide.Node{
		Value:    "B1",
		Parent:   &nodeB,
		Children: nil,
	}

	nodeC := treeglide.Node{
		Value:    "C",
		Parent:   &root,
		Children: nil,
	}

	nodeA1.Children = append(nodeA1.Children, &nodeA11, &nodeA12, &nodeA13)

	nodeA.Children = append(nodeA.Children, &nodeA1)
	nodeA.Children = append(nodeA.Children, &nodeA2)

	nodeB.Children = append(nodeB.Children, &nodeB1)

	root.Children = append(root.Children, &nodeA)
	root.Children = append(root.Children, &nodeB)
	root.Children = append(root.Children, &nodeC)

	w, h, errTerm := term.GetSize(int(os.Stdout.Fd()))
	if errTerm != nil {
		w = 80
		h = 24
	}
	top, right, bottom, left := styleDoc.GetPadding()
	w = w - left - right
	h = h - top - bottom

	err := tea.NewProgram(model{tree: treeglide.New(&root, w, h)}).Start()
	if err != nil {
		os.Exit(1)
	}
}
