package main

import (
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	styleDoc = lipgloss.NewStyle().Padding(1)
)

type Node struct {
	Value    string
	Parent   *Node
	Children []*Node
}

type Cursor struct {
	Current *Node
	Parent  *Node
	Index   int
}

type Model struct {
	KeyMap KeyMap
	root   *Node
	cursor Cursor
}

func New(node *Node) Model {
	return Model{
		KeyMap: DefaultKeyMap(),
		root:   node,
		cursor: Cursor{
			Current: node.Children[0],
			Parent:  node,
			Index:   0,
		},
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Up):
			m.NavUp()
		case key.Matches(msg, m.KeyMap.Down):
			m.NavDown()
		case key.Matches(msg, m.KeyMap.Left):
			m.NavLeft()
		case key.Matches(msg, m.KeyMap.Right):
			m.NavRight()
		}
	}

	return m, nil
}

func (m Model) View() string {
	sections := []string{m.cursor.Parent.Value, m.cursor.Current.Value, strconv.Itoa(m.cursor.Index)}

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) NavUp() {
	if m.cursor.Index > 0 {
		m.cursor.Index--
		m.cursor.Current = m.cursor.Parent.Children[m.cursor.Index]
	}
}

func (m *Model) NavDown() {
	if m.cursor.Index < len(m.cursor.Current.Parent.Children)-1 {
		m.cursor.Index++
		m.cursor.Current = m.cursor.Parent.Children[m.cursor.Index]
	}
}

func (m *Model) NavLeft() {
	if m.cursor.Parent.Parent != nil {
		m.cursor.Current = m.cursor.Parent
		m.cursor.Parent = m.cursor.Current.Parent

		for idx, child := range m.cursor.Parent.Children {
			if child == m.cursor.Current {
				m.cursor.Index = idx
				return
			}
		}
	}
}

func (m *Model) NavRight() {
	if len(m.cursor.Current.Children) > 0 {
		m.cursor.Parent = m.cursor.Current
		m.cursor.Current = m.cursor.Current.Children[0]
		m.cursor.Index = 0
	}
}

// KeyMap holds the key bindings for the table.
type KeyMap struct {
	Bottom      key.Binding
	Top         key.Binding
	SectionDown key.Binding
	SectionUp   key.Binding
	Down        key.Binding
	Up          key.Binding
	Right       key.Binding
	Left        key.Binding
	Quit        key.Binding

	ShowFullHelp  key.Binding
	CloseFullHelp key.Binding
}

// DefaultKeyMap is the default key bindings for the table.
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Down: key.NewBinding(
			key.WithKeys("down"),
			key.WithHelp("↓", "down"),
		),
		Up: key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("↑", "up"),
		),
		Right: key.NewBinding(
			key.WithKeys("right"),
			key.WithHelp("→", "right"),
		),
		Left: key.NewBinding(
			key.WithKeys("left"),
			key.WithHelp("←", "left"),
		),

		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),

		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
	}
}

//
// Implementation
//

type model struct {
	tree Model
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

	root := Node{
		Value:    "root",
		Parent:   nil,
		Children: nil}

	nodeA := Node{
		Value:    "A",
		Parent:   &root,
		Children: nil,
	}

	nodeA1 := Node{
		Value:    "A1",
		Parent:   &nodeA,
		Children: nil,
	}

	nodeA2 := Node{
		Value:    "A2",
		Parent:   &nodeA,
		Children: nil,
	}

	nodeB := Node{
		Value:    "B",
		Parent:   &root,
		Children: nil,
	}

	nodeB1 := Node{
		Value:    "B1",
		Parent:   &nodeB,
		Children: nil,
	}

	nodeC := Node{
		Value:    "C",
		Parent:   &root,
		Children: nil,
	}

	nodeA.Children = append(nodeA.Children, &nodeA1)
	nodeA.Children = append(nodeA.Children, &nodeA2)

	nodeB.Children = append(nodeB.Children, &nodeB1)

	root.Children = append(root.Children, &nodeA)
	root.Children = append(root.Children, &nodeB)
	root.Children = append(root.Children, &nodeC)

	err := tea.NewProgram(model{tree: New(&root)}).Start()
	if err != nil {
		os.Exit(1)
	}
}
