package ui

import (
	"strings"
	"time"

	"github.com/76creates/stickers/flexbox"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattsp1290/october-talks-2025/example/client/internal/message"
	"github.com/mattsp1290/sweetie16/pkg/color"
)

var (
	styleBackground    = lipgloss.NewStyle().Align(lipgloss.Center).Background(color.BlackColor)
	chatStyle          = lipgloss.NewStyle().Align(lipgloss.Left).Background(color.BlackColor)
	inputStyle         = lipgloss.NewStyle().Align(lipgloss.Left).Background(color.BlackColor)
	dividerStyle       = lipgloss.NewStyle().Align(lipgloss.Center).Background(color.BlackColor)
	userNameStyle      = lipgloss.NewStyle().Foreground(color.BlueColor)
	userInputStyle     = lipgloss.NewStyle().Foreground(color.LightGreyColor)
	assistantNameStyle = lipgloss.NewStyle().Foreground(color.RedColor)
	tableRowIndex      = 1
	tableCellIndex     = 1
	cellString         = "Hello World"
)

type UIMessage struct {
	User      string
	Content   string
	Timestamp string
}

func NewUIMessage(user, content string) UIMessage {
	// now
	timestamp := time.Now().Format("15:04:05")
	if strings.Contains(user, "You") {
		content = userInputStyle.Render(content)
	}
	return UIMessage{
		User:      user,
		Content:   content,
		Timestamp: timestamp,
	}
}

type Model struct {
	flexBox   *flexbox.FlexBox
	table     table.Model
	chatRow   *flexbox.Row
	helpRow   *flexbox.Row
	statusRow *flexbox.Row
	inputRow  *flexbox.Row
	messages  []UIMessage
	textarea  textarea.Model
	userInput chan string
	err       error
}

func getBox() *flexbox.FlexBox {
	box := flexbox.New(0, 0).SetStyle(styleBackground)
	return box
}

func getChatRow(box *flexbox.FlexBox) *flexbox.Row {
	return box.NewRow().AddCells(
		flexbox.NewCell(9, 10).SetStyle(chatStyle),
		flexbox.NewCell(1, 10).SetStyle(dividerStyle),
		flexbox.NewCell(3, 10).SetStyle(styleBackground),
	)
}

func getInputRow(box *flexbox.FlexBox) *flexbox.Row {
	return box.NewRow().AddCells(
		flexbox.NewCell(9, 2).SetStyle(inputStyle),
		flexbox.NewCell(1, 2).SetStyle(dividerStyle),
	)
}

func getStatusBar(box *flexbox.FlexBox) *flexbox.Row {
	return box.NewRow().AddCells(
		flexbox.NewCell(1, 1).SetStyle(dividerStyle),
		flexbox.NewCell(9, 1).SetStyle(styleBackground),
		flexbox.NewCell(1, 1).SetStyle(dividerStyle),
	)
}

func getHelpBar(box *flexbox.FlexBox) *flexbox.Row {
	return box.NewRow().AddCells(
		flexbox.NewCell(9, 1).SetStyle(styleBackground),
		flexbox.NewCell(1, 1).SetStyle(dividerStyle),
	)
}

func getTableStyle() table.Styles {
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(color.LightGreyColor).
		BorderBottom(false).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color(color.LightBlueColor)).
		Background(lipgloss.Color(color.DarkBlueColor)).
		Bold(false)
	return s
}

func getColumns() []table.Column {
	return []table.Column{
		{Title: "Time", Width: 15},
		{Title: "User", Width: 35},
		{Title: "Message", Width: 150},
	}
}

func getRows(messages []UIMessage) []table.Row {
	var rows []table.Row
	for _, msg := range messages {
		rows = append(rows, table.Row{
			msg.Timestamp,
			msg.User,
			msg.Content,
		})
	}
	return rows
}

func getTable(rows []table.Row) table.Model {
	t := table.New(
		table.WithColumns(getColumns()),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithWidth(250),
	)
	t.SetStyles(getTableStyle())

	return t
}

func getTextarea() textarea.Model {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = inputStyle
	ta.FocusedStyle.Text = inputStyle
	ta.BlurredStyle.Text = inputStyle

	ta.ShowLineNumbers = false
	return ta
}

func InitialModel(userInput chan string) *Model {
	box := getBox()
	chatRow := getChatRow(box)
	textareaRow := getInputRow(box)
	statusRow := getStatusBar(box)
	helpRow := getHelpBar(box)
	rows := []*flexbox.Row{chatRow, textareaRow, statusRow, helpRow}

	box.AddRows(rows)

	return &Model{
		flexBox:   box,
		table:     getTable(nil),
		userInput: userInput,
		chatRow:   chatRow,
		helpRow:   helpRow,
		statusRow: statusRow,
		inputRow:  textareaRow,
		textarea:  getTextarea(),
	}
}

func (m *Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		windowHeight := msg.Height
		windowWidth := msg.Width
		m.flexBox.SetWidth(windowWidth)
		m.flexBox.SetHeight(windowHeight)
		m.table.SetWidth(windowWidth)
		m.table.SetHeight(windowHeight)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			m.userInput <- m.textarea.Value()

			uiMsg := NewUIMessage(inputStyle.Render("You"), m.textarea.Value())
			m.messages = append(m.messages, uiMsg)
			//m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
			m.textarea.Reset()
			//m.viewport.GotoBottom()
		}
	case *message.Message:
		for _, currMsg := range msg.Strings() {
			uiMsg := NewUIMessage(assistantNameStyle.Render("Assistant"), currMsg)
			m.messages = append(m.messages, uiMsg)
		}

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd)
}

func (m *Model) View() string {
	m.flexBox.ForceRecalculate()
	var cell *flexbox.Cell
	if m.chatRow != nil {
		cell = m.chatRow.GetCell(0)
	}

	if m.chatRow != nil && cell != nil {
		tbl := getTable(getRows(m.messages))
		m.table = tbl
		m.table.SetWidth(cell.GetWidth())
		m.table.SetHeight(cell.GetHeight())
		cell.SetContent(m.table.View())
	}

	if m.inputRow != nil {
		m.inputRow.GetCell(0).SetContent(m.textarea.View())
	}

	return m.flexBox.Render()
}
