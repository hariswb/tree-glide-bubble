package treeglide

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	vertical      string = "│"
	verticalHeavy string = "┃"

	white  = lipgloss.Color("#ffffff")
	black  = lipgloss.Color("#000000")
	purple = lipgloss.Color("#bd93f9")
)

type Styles struct {
	Shapes        lipgloss.Style
	SelectedValue lipgloss.Style
	SelectedDesc  lipgloss.Style
	Unselected    lipgloss.Style
	Help          lipgloss.Style
}

func defaultStyles() Styles {
	return Styles{
		Shapes:        lipgloss.NewStyle().Margin(0, 0, 0, 0).Foreground(purple),
		SelectedValue: lipgloss.NewStyle().Margin(0, 0, 0, 0).Background(purple),
		SelectedDesc:  lipgloss.NewStyle().Margin(0, 0, 0, 0).Background(lipgloss.AdaptiveColor{Light: "#111100", Dark: "#001100"}),
		Unselected:    lipgloss.NewStyle().Margin(0, 0, 0, 0).Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}),
		Help:          lipgloss.NewStyle().Margin(0, 0, 0, 0).Foreground(lipgloss.AdaptiveColor{Light: "#000000", Dark: "#ffffff"}),
	}
}

// KeyMap holds the key bindings for the table.
type KeyMap struct {
	Down  key.Binding
	Up    key.Binding
	Right key.Binding
	Left  key.Binding
	Quit  key.Binding

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

type Node struct {
	Value    string
	Desc     string
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
	Styles Styles

	root   *Node
	cursor Cursor

	width  int
	height int

	Help     help.Model
	showHelp bool

	AdditionalShortHelpKeys func() []key.Binding
}

func NewTree(node *Node, width int, height int) Model {
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

		Help:     help.New(),
		showHelp: true,
	}
}

func NewNode(value string, desc string, parent *Node) *Node {
	node := &Node{Value: value, Desc: desc, Parent: parent}
	if parent != nil {
		parent.Children = append(parent.Children, node)
	}
	return node
}

func (node *Node) WrapDesc(width int) []string {
	if width <= 0 {
		return []string{node.Desc}
	}

	words := strings.Fields(node.Desc)
	var lines []string
	var line string

	for _, word := range words {
		if len(line)+len(word) < width {
			if line != "" {
				line += " "
			}
			line += word
		} else {
			lines = append(lines, line)
			line = word
		}
	}

	if line != "" {
		lines = append(lines, line)
	}

	return lines
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
		case key.Matches(msg, m.KeyMap.ShowFullHelp):
			fallthrough
		case key.Matches(msg, m.KeyMap.CloseFullHelp):
			m.Help.ShowAll = !m.Help.ShowAll
		}
	}

	return m, nil
}

func (m Model) View() string {
	availableHeight := m.height

	var help string
	if m.showHelp {
		help = m.helpView()
		availableHeight -= lipgloss.Height(help)
	}

	var sections []string
	sections = append(sections, lipgloss.NewStyle().Height(20).Render(m.renderTree(m.root.Children, 0)), help)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m *Model) renderTree(remainingNodes []*Node, indent int) string {
	var b strings.Builder

	for _, node := range remainingNodes {
		var str string
		var indentStr string

		shape := m.Styles.Shapes.Render(vertical)
		wrapWidth := m.width - (indent)*3
		minCharHighlight := wrapWidth

		if m.cursor.Current == node {
			shape = m.Styles.Shapes.Render(verticalHeavy)
		}

		if indent > 0 {
			indentStr = strings.Repeat(m.Styles.Shapes.Render(vertical)+"  ", (indent)) + shape
		} else {
			indentStr = shape
		}

		valueStr := m.Styles.Unselected.Render(fmt.Sprintf("%-*s", minCharHighlight, node.Value))

		// If we are at the cursor, we add the selected style to the string
		if m.cursor.Current == node {
			valueStr = m.Styles.SelectedValue.Render(valueStr)
		}

		str += indentStr + fmt.Sprintf("%s\n", valueStr)

		for _, descStrLine := range node.WrapDesc(m.width - (indent)*3) {
			descStr := m.Styles.Unselected.Render(fmt.Sprintf("%-*s", minCharHighlight, descStrLine))
			if m.cursor.Current == node {
				descStr = m.Styles.SelectedDesc.Render(descStr)
			}
			str += indentStr + fmt.Sprintf("%s\n", descStr)
		}

		b.WriteString(str)

		if node.Children != nil {
			childStr := m.renderTree(node.Children, indent+1)
			b.WriteString(childStr)
		}

	}

	return b.String()
}

func (m Model) Width() int {
	return m.width
}

func (m Model) Height() int {
	return m.height
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Model) SetWidth(newWidth int) {
	m.SetSize(newWidth, m.height)
}

func (m *Model) SetHeight(newHeight int) {
	m.SetSize(m.width, newHeight)
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

func (m *Model) SetShowHelp() bool {
	return m.showHelp
}

func (m Model) helpView() string {
	return m.Styles.Help.Render(m.Help.View(m))
}

func (m Model) ShortHelp() []key.Binding {
	kb := []key.Binding{
		m.KeyMap.Up,
		m.KeyMap.Down,
	}

	if m.AdditionalShortHelpKeys != nil {
		kb = append(kb, m.AdditionalShortHelpKeys()...)
	}

	return append(kb,
		m.KeyMap.Quit,
	)
}

func (m Model) FullHelp() [][]key.Binding {
	kb := [][]key.Binding{{
		m.KeyMap.Up,
		m.KeyMap.Down,
	}}

	return append(kb,
		[]key.Binding{
			m.KeyMap.Quit,
			m.KeyMap.CloseFullHelp,
		})
}
