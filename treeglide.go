package treeglide

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	bottomLeft string = " └──"

	white  = lipgloss.Color("#ffffff")
	black  = lipgloss.Color("#000000")
	purple = lipgloss.Color("#bd93f9")
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

type Styles struct {
	Shapes     lipgloss.Style
	Selected   lipgloss.Style
	Unselected lipgloss.Style
	Help       lipgloss.Style
}

type Model struct {
	KeyMap KeyMap
	Styles Styles

	root   *Node
	cursor Cursor

	width  int
	height int
}

func New(node *Node, width int, height int) Model {
	return Model{
		KeyMap: DefaultKeyMap(),
		Styles: defaultStyles(),

		root: node,
		cursor: Cursor{
			Current: node.Children[0],
			Parent:  node,
			Index:   0,
		},

		width:  width,
		height: height,
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
	// sections := []string{m.cursor.Parent.Value, m.cursor.Current.Value, strconv.Itoa(m.cursor.Index)}

	var sections []string
	sections = append(sections, lipgloss.NewStyle().Height(24).Render(m.renderTree(m.root.Children, 0)), m.cursor.Current.Value)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderTree(remainingNodes []*Node, indent int) string {
	var b strings.Builder

	for _, node := range remainingNodes {
		var str string

		// If we aren't at the root, we add the arrow shape to the string
		if indent > 0 {
			shape := strings.Repeat(" ", (indent-1)*4) + m.Styles.Shapes.Render(bottomLeft) + " "
			str += shape
		}

		valueStr := fmt.Sprintf("%-*s", 20, node.Value)

		// If we are at the cursor, we add the selected style to the string
		if m.cursor.Current == node {
			str += fmt.Sprintf("%s\n", m.Styles.Selected.Render(valueStr))
		} else {
			str += fmt.Sprintf("%s\n", m.Styles.Unselected.Render(valueStr))
		}

		b.WriteString(str)

		if node.Children != nil {
			childStr := m.renderTree(node.Children, indent+1)
			b.WriteString(childStr)
		}

	}

	return b.String()
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

func defaultStyles() Styles {
	return Styles{
		Shapes:     lipgloss.NewStyle().Margin(0, 0, 0, 0).Foreground(purple),
		Selected:   lipgloss.NewStyle().Margin(0, 0, 0, 0).Background(purple),
		Unselected: lipgloss.NewStyle().Margin(0, 0, 0, 0).Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}),
		Help:       lipgloss.NewStyle().Margin(0, 0, 0, 0).Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}),
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
