package main

import (
	"bytes"
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	"strconv"
)

type model struct {
	ready     bool
	list      list.Model
	finished  bool
	endResult chan<- string
	jump      string
}

type stringItem string

func (s stringItem) String() string {
	return string(s)
}

func main() {
	itemList := []string{
		"Welcome to the bubbles-list example!",
		"Use 'q' or 'ctrl-c' to quit!",
		"You can move the highlighted index up and down with the (arrow keys 'k' and 'j'.",
		"Move to the beginning with 'g' and to the end with 'G'.",
		"Sort the entrys with 's', but be carefull you can't unsort it again.",
		"The list can handel linebreaks,\nand has wordwrap enabled if the line gets to long.",
		"You can select items with the space key which will select the line and mark it as such.",
		"Ones you hit 'enter', the selected lines will be printed to StdOut and the program exits.",
		"When you print the items there will be a loss of information,",
		"since one can not say what was a line break within an item or what is a new item",
		"Use '+' or '-' to move the item under the curser up and down.",
		"The key 'v' inverts the selected state of each item.",
		"To toggle betwen only absolute itemnumbers and also relativ numbers, the 'r' key is your friend.",
		"","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","","",
	}
	stringerList := list.MakeStringerList(itemList)

	endResult := make(chan string, 1)
	list := list.NewModel()
	list.AddItems(stringerList)
	// uncomment the following lines for fancy check (selected) box :-)
	// l.WrapPrefix = false
	// l.SelectedPrefix = " [x]"
	// l.UnSelectedPrefix = "[ ]"

	// Since in this example we only use UNIQUE string items we can use a String Comparison for the equals methode
	// but be aware that different items in your case can have the same string -> false-positiv
	// Better: Assert back to your struct and test on something unique within it!
	list.SetEquals(func(first, second fmt.Stringer) bool { return first.String() == second.String() })
	m := model{}
	m.list = list


	m.endResult = endResult


	p := tea.NewProgram(m)

	// Use the full size of the terminal in its "alternate screen buffer"
	fullScreen := false // change to true if you want fullscreen

	if fullScreen {
		p.EnterAltScreen()
	}

	if err := p.Start(); err != nil {
		fmt.Println("could not run program:", err)
		os.Exit(1)
	}
	if fullScreen {
		p.ExitAltScreen()
	}

	res := <-endResult
	// allways print a newline even on empty string result
	fmt.Println(res)
}


func (m model) Init() tea.Cmd {
	return nil
}

// View waits till the terminal sizes is knowen to the model and than,
// pipes the model to the list View for rendering the list
func (m model) View() string {
	if !m.ready {
		return "\n  Initalizing...\n\n  Waiting for info about window size.\n"
	}

	listString := m.list.View()
	return listString
}

// update recives messages and the model and changes the model accordingly to the messages
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Ctrl+c exits
		if msg.Type == tea.KeyCtrlC {
			m.endResult <- ""
			return m, tea.Quit
		}
		switch keyString := msg.String(); keyString {
		case "q":
			m.endResult <- ""
			return m, tea.Quit
		case "1", "2", "3", "4", "5", "6", "7", "8", "9", "0":
			m.jump += keyString
			return m, nil
		case "down", "j":
			j := 1
			if m.jump != "" {
				j, _ = strconv.Atoi(m.jump)
				m.jump = ""
			}
			m.list.Move(j)
			return m, nil
		case "up", "k":
			j := 1
			if m.jump != "" {
				j, _ = strconv.Atoi(m.jump)
				m.jump = ""
			}
			m.list.Move(-j)
			return m, nil
		case "r":
			m.list.NumberRelative = !m.list.NumberRelative
			return m, nil
		case "m":
			j := 1
			if m.jump != "" {
				j, _ = strconv.Atoi(m.jump)
				m.jump = ""
			}
			m.list.MarkSelected(j, true)
			return m, nil
		case "M":
			j := 1
			if m.jump != "" {
				j, _ = strconv.Atoi(m.jump)
				m.jump = ""
			}
			m.list.MarkSelected(j, false)
			return m, nil
		case " ":
			j := 1
			if m.jump != "" {
				j, _ = strconv.Atoi(m.jump)
				m.jump = ""
			}
			m.list.ToggleSelect(j)
			m.list.Move(1)
			return m, nil
		case "enter":
			// Enter prints the selected lines to StdOut
			var result bytes.Buffer
			for _, item := range m.list.GetSelected() {
				result.WriteString(item.String())
				result.WriteString("\n")
			}
			m.endResult <- result.String()
			return m, tea.Quit
		default:
			// resets jump buffer to prevent confusion
			m.jump = ""

			// pipe all other commands to the update from the list
			list, newMsg := list.Update(msg, m.list)
			m.list = list
			return m, newMsg
		}

	case tea.WindowSizeMsg:

		m.list.Width = msg.Width
		m.list.Height = msg.Height

		if !m.ready {
			// Since this program can use the full size of the viewport we need
			// to wait until we've received the window dimensions before we
			// can initialize the viewport. The initial dimensions come in
			// quickly, though asynchronously, which is why we wait for them
			// here.
			m.ready = true
		}

		return m, nil

	}

	return m, nil
}
