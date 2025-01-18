package main

import (
	"os"
	"strings"

	//	pb_common "base-station/protobuf/generated/go"
	//	pb_column "base-station/protobuf/generated/go/column"
	//	pb_node "base-station/protobuf/generated/go/node"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	//"github.com/golang/protobuf/proto"
	"golang.org/x/term"
)

type ActiveComponent int

const (
	ComponentMessages ActiveComponent = iota
	ComponentList
	ComponentEditor
)

// For list items
type Option struct {
	title       string
	description string
}

func (o Option) Title() string       { return o.title }
func (o Option) Description() string { return o.description }
func (o Option) FilterValue() string { return o.title }

type MessageCmd string

type Model struct {
	msgView    viewport.Model
	optionList list.Model
	jsonEditor textarea.Model

	messages        []string
	jsonTemplates   map[string]string
	activeComponent ActiveComponent
	chosen_command  string

	width  int
	height int
}

func NewModel() Model {
	msgView := viewport.New(0, 0)
	//msgView.MouseWheelEnabled = true
	//msgView.YOffset = 15
	msgView.SetContent("This is where messages will come from")

	optItems := []list.Item{
		Option{title: "SetPumpStateCommand", description: "sets the pump state"},
		Option{title: "PumpUpdateScheduleCommand", description: "sets the pump schedule"},
	}
	optionList := list.New(optItems, list.NewDefaultDelegate(), 0, 0)

	jsonEditor := textarea.New()
	jsonEditor.Placeholder = "This is where JSON will go"

	return Model{
		msgView:         msgView,
		jsonEditor:      jsonEditor,
		optionList:      optionList,
		activeComponent: ComponentList,
		messages:        make([]string, 0),
		chosen_command:  "",
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.optionList.StartSpinner(),
		textarea.Blink,
		m.checkWindowSize(),
	)
}

func (m Model) checkWindowSize() tea.Cmd {
	return func() tea.Msg {
		w, h, _ := term.GetSize(int(os.Stdin.Fd()))
		return tea.WindowSizeMsg{Width: w, Height: h}
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height

		m.msgView.Width = m.width/2 - 2
		m.msgView.Height = m.height - 2

		m.optionList.SetWidth(m.width/2 - 2)
		m.optionList.SetHeight(m.height/2 - 2)

		m.jsonEditor.SetWidth(m.width/2 - 2)
		m.jsonEditor.SetHeight(m.height/2 - 2)
		return m, nil
	case MessageCmd:
		m.messages = append(m.messages, string(msg))
		var sb strings.Builder
		for _, message := range m.messages {
			sb.WriteString(message)
			sb.WriteString("\n")
		}

		// Update the viewport content
		m.msgView.SetContent(sb.String())
		m.msgView.GotoBottom()
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+n":
			// cycle the actively selected view
			m.cycleFocus()
			return m, nil
		}

		switch m.activeComponent {
		case ComponentList:
			newListModel, cmd := m.optionList.Update(msg)
			m.optionList = newListModel
			cmds = append(cmds, cmd)
			// If an item was selected (usually enter key)
			if msg.String() == "enter" {
				if selected, ok := m.optionList.SelectedItem().(Option); ok {
					// do something here on selection of item
					m.chosen_command = selected.Title()
					marshalled, err := ProtoToJSON(Commands[m.chosen_command])
					if err != nil {
						m.jsonEditor.SetValue(err.Error())
						// probably send an error to the msg view
						return m, nil
					}

					m.jsonEditor.SetValue(string(marshalled))
					m.activeComponent = ComponentEditor
					m.jsonEditor.Focus()
				}
			}
		case ComponentEditor:
			newEditor, cmd := m.jsonEditor.Update(msg)
			m.jsonEditor = newEditor
			cmds = append(cmds, cmd)
			if msg.String() == "ctrl+s" {
				// do something on submit
				json_string := m.jsonEditor.Value()
				tx := &TXMessage{
					proto_type:   m.chosen_command,
					json_content: json_string,
				}

				TxMessages <- tx

				m.activeComponent = ComponentMessages
				m.jsonEditor.Blur()
			}
		case ComponentMessages:
			newView, cmd := m.msgView.Update(msg)
			m.msgView = newView
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m Model) View() string {

	componentWidth := m.width/2 - 2
	listHeight := m.height/2 - 2
	editorHeight := m.height/2 - 2

	msgBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Width(componentWidth).
		Height(m.height - 2)
	if m.activeComponent == ComponentMessages {
		msgBox = msgBox.BorderForeground(lipgloss.Color("86"))
	} else {
		msgBox = msgBox.BorderForeground(lipgloss.Color("238"))
	}

	listBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Width(componentWidth).
		Height(listHeight)
	if m.activeComponent == ComponentList {
		listBox = listBox.BorderForeground(lipgloss.Color("86"))
	} else {
		listBox = listBox.BorderForeground(lipgloss.Color("238"))
	}

	editorBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1).
		Width(componentWidth).
		Height(editorHeight)
	if m.activeComponent == ComponentEditor {
		editorBox = editorBox.BorderForeground(lipgloss.Color("86"))
	} else {
		editorBox = editorBox.BorderForeground(lipgloss.Color("238"))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top,
		msgBox.Render(m.msgView.View()),
		lipgloss.JoinVertical(lipgloss.Left,
			listBox.Render(m.optionList.View()),
			editorBox.Render(m.jsonEditor.View()),
		),
	)
}

func (m *Model) cycleFocus() {
	switch m.activeComponent {
	case ComponentMessages:
		m.activeComponent = ComponentList
	case ComponentList:
		m.activeComponent = ComponentEditor
		m.jsonEditor.Focus()
	case ComponentEditor:
		m.activeComponent = ComponentMessages
		m.jsonEditor.Blur()
	}

}
