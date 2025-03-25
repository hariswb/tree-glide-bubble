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
	root := treeglide.NewNode("admin", "Welcome to the thread!", nil)

	// First-level comments
	user1 := treeglide.NewNode("user1", "I totally agree with this post!", root)
	user2 := treeglide.NewNode("user2", "I think there’s another perspective to consider.", root)
	user3 := treeglide.NewNode("user3", "This is hilarious! 😂", root)

	// Replies to user1
	user4 := treeglide.NewNode("user4", "Yeah, I was thinking the same thing!", user1)
	treeglide.NewNode("user5", "Not sure if I agree, but interesting take.", user4)
	treeglide.NewNode("user6", "I see your point, but have you considered XYZ?", user4)

	treeglide.NewNode("user10", "Can you please elaborate?", user1)

	// Replies to user2
	user7 := treeglide.NewNode("user7", "What do you mean by that?", user2)
	treeglide.NewNode("user8", "I think user2 has a good argument.", user7)

	// Replies to user3
	treeglide.NewNode("user9", "LOL, right? This made my day. 😂", user3)

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
