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

	white    = lipgloss.Color("#ffffff")
	black    = lipgloss.Color("#000000")
	blue     = lipgloss.Color("#8ea1f5")
	darkBlue = lipgloss.Color("#5a6ec4")
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
		Shapes:        lipgloss.NewStyle().Margin(0, 0, 0, 0).Foreground(blue),
		SelectedValue: lipgloss.NewStyle().Margin(0, 0, 0, 0).Background(blue),
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
	Current      *Node
	Parent       *Node
	Index        int
	WindowStart  int
	WindowEnd    int
	WindowHeight int
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
	m := Model{
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

	availableHeight := m.height

	if m.showHelp {
		availableHeight -= lipgloss.Height(m.helpView())
	}

	m.cursor.WindowHeight = availableHeight
	m.cursor.WindowStart = 0
	m.cursor.WindowEnd = availableHeight

	return m
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
	var sections []string

	// Render tree
	treeStr := m.renderTree(m.root.Children, 0)
	totalLines := strings.Count(treeStr, "\n")

	visibleTree := strings.Split(treeStr, "\n")

	curStart := 0
	curEnd := 0
	curStart, curEnd = m.CurLinePos(m.root, &curStart, &curEnd, -1)

	// Offset means the cursor went beyond the window height
	isOffset := !(curStart >= m.cursor.WindowStart && curEnd < m.cursor.WindowEnd)

	// If the whole text less than the available height,
	// Set the slice target to its minimum
	m.cursor.WindowEnd = min(len(visibleTree), m.cursor.WindowEnd)

	visibleTree = visibleTree[m.cursor.WindowStart:m.cursor.WindowEnd]

	if isOffset {
		// Readjust the tree slice, refocus the cursor
		m.cursor.WindowEnd = min(curStart+m.cursor.WindowHeight, totalLines)
		visibleTree = visibleTree[m.cursor.WindowStart : m.cursor.WindowEnd+1]
	}

	sections = append(
		sections,
		lipgloss.NewStyle().Height(m.cursor.WindowHeight).Render(strings.Join(visibleTree, "\n")),
		m.helpView(),
	)

	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

func (m Model) CurLinePos(node *Node, start *int, end *int, indent int) (int, int) {
	if node == nil {
		return -1, -1
	}

	valueLen := 1
	descLen := len(node.WrapDesc(m.width - (indent)*3))

	if *start == 0 && *end == 0 {
		*end = descLen
	} else {
		*start = *end + valueLen
		*end = *start + descLen
	}

	if m.cursor.Current == node {
		return *start, *end
	}

	if node == m.root {
		*start = 0
		*end = 0
	}

	for _, child := range node.Children {
		resultStart, resultEnd := m.CurLinePos(child, start, end, indent+1)
		if resultStart != -1 { // If found, return immediately
			return resultStart, resultEnd
		}
	}

	return -1, -1
}

func (c *Cursor) IsCurrent(node *Node) bool {
	return c.Current == node
}

type Renderer struct {
	result    string
	indentStr string
	model     *Model
	node      *Node
	indent    int
	boxWidth  int
}

func (r *Renderer) CreateIndent() {
	shape := r.model.Styles.Shapes.Render(vertical)
	if r.model.cursor.IsCurrent(r.node) {
		shape = r.model.Styles.Shapes.Render(verticalHeavy)
	}

	if r.indent > 0 {
		r.indentStr = strings.Repeat(r.model.Styles.Shapes.Render(vertical)+"  ", (r.indent)) + shape
	} else {
		r.indentStr = shape
	}
}

// Define the rendered texts maximum width
func (r *Renderer) BoxWidth() {
	r.boxWidth = r.model.width - (r.indent)*4
}

func (r *Renderer) CreateValueStr() {
	valueStr := r.model.Styles.Unselected.Render(fmt.Sprintf("%-*s", r.boxWidth, r.node.Value))

	// If we are at the cursor, we add the selected style to the string
	if r.model.cursor.IsCurrent(r.node) {
		valueStr = r.model.Styles.SelectedValue.Render(valueStr)
	}

	r.result += r.indentStr + fmt.Sprintf("%s\n", valueStr)
}

func (r *Renderer) CreateDescStrs() {
	for _, descStrLine := range r.node.WrapDesc(r.boxWidth) {
		descStr := r.model.Styles.Unselected.Render(fmt.Sprintf("%-*s", r.boxWidth, descStrLine))
		if r.model.cursor.IsCurrent(r.node) {
			descStr = r.model.Styles.SelectedDesc.Render(descStr)
		}
		r.result += r.indentStr + fmt.Sprintf("%s\n", descStr)
	}
}

func (m *Model) renderTree(remainingNodes []*Node, indent int) string {
	var b strings.Builder

	for _, node := range remainingNodes {

		renderer := Renderer{
			model:  m,
			node:   node,
			indent: indent,
		}

		renderer.BoxWidth()
		renderer.CreateIndent()
		renderer.CreateValueStr()
		renderer.CreateDescStrs()

		b.WriteString(renderer.result)

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
