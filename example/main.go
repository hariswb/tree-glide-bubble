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
	root := treeglide.NewNode("root", nil)

	nodeA := treeglide.NewNode("A", root)
	nodeA1 := treeglide.NewNode("A1", nodeA)
	treeglide.NewNode("A11", nodeA1)
	treeglide.NewNode("A12", nodeA1)
	treeglide.NewNode("A13", nodeA1)
	treeglide.NewNode("A2", nodeA)

	nodeB := treeglide.NewNode("B", root)
	treeglide.NewNode("B1", nodeB)

	treeglide.NewNode("C", root)

	w, h, errTerm := term.GetSize(int(os.Stdout.Fd()))
	if errTerm != nil {
		w = 80
		h = 24
	}
	top, right, bottom, left := styleDoc.GetPadding()
	w = w - left - right
	h = h - top - bottom

	err := tea.NewProgram(model{tree: treeglide.NewTree(root, w, h)}).Start()
	if err != nil {
		os.Exit(1)
	}
}
