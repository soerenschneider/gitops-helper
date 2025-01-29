package tui

import (
	"errors"
	"fmt"
	"gitops-helper/internal"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type UserData struct {
	ClusterName string
	Components  []string
	GitOpsTool  string
}

func (u *UserData) Validate() error {
	if u.ClusterName == "" {
		u.ClusterName = internal.DefaultClusterName
	}

	if u.GitOpsTool == "" {
		return errors.New("no gitops tool selected")
	}

	return nil
}

var docStyle = lipgloss.NewStyle().Margin(3, 4)

type item struct {
	title, desc string
	selected    bool
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	str := fmt.Sprintf("%d. %s", index+1, i.title)
	if i.selected {
		str = fmt.Sprintf("[X] %s", str)
	} else {
		str = fmt.Sprintf("[ ] %s", str)
	}

	_, _ = fmt.Fprint(w, fn(str))
}

type model struct {
	clusterName textinput.Model

	componentsList list.Model
	components     map[string]struct{}

	gitopsList   list.Model
	gitopsChoice string

	screen int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.screen {
	case 0:
		m.clusterName, cmd = m.clusterName.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				m.screen = 1
			}
			if msg.String() == "q" {
				return m, tea.Quit
			}
		}
		return m, cmd
	case 1:
		m.componentsList, cmd = m.componentsList.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				m.screen = 2
			}
			if msg.String() == " " {
				i, ok := m.componentsList.SelectedItem().(item)
				if ok {
					// Toggle the checked state
					i.selected = !i.selected

					if i.selected {
						m.components[i.title] = struct{}{}
					} else {
						delete(m.components, i.title)
					}

					// Update the item in the componentsList slice
					m.componentsList.Items()[m.componentsList.Index()] = i
				}
			}
		case tea.WindowSizeMsg:
			h, v := docStyle.GetFrameSize()
			m.componentsList.SetSize(msg.Width-h, msg.Height-v)
		}

		return m, cmd
	case 2:
		m.gitopsList, cmd = m.gitopsList.Update(msg)

		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
			if msg.String() == "enter" {
				i, ok := m.gitopsList.SelectedItem().(item)
				if ok {
					m.gitopsChoice = i.title
				}
				return m, tea.Quit
			}
		case tea.WindowSizeMsg:
			h, v := docStyle.GetFrameSize()
			m.gitopsList.SetSize(msg.Width-h, msg.Height-v)
		}

		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	switch m.screen {
	case 0:
		return fmt.Sprintf("Step 1: Enter a name for the new management cluster\n\n%s\n\n(press Enter to confirm, q to quit)\n", m.clusterName.View())
	case 1:
		return fmt.Sprintf("Step 2: Select Components\n\n%s\n\n(press Enter to confirm)\n", m.componentsList.View())
	case 2:
		return fmt.Sprintf("Step 2: Select GitOps Tool\n\n%s\n\n(press Enter to confirm)\n", m.gitopsList.View())
	}
	return "Done!"
}

func RunWizard(components []string) (UserData, error) {
	componentItems := make([]list.Item, 0, len(components))

	for _, component := range components {
		componentItems = append(componentItems, item{
			title: component,
		})
	}

	gitopsTools := []list.Item{
		item{title: internal.ArgoCD},
		item{title: internal.FluxCD},
	}

	//x, y := docStyle.GetFrameSize()
	m := model{
		clusterName:    textinput.New(),
		componentsList: list.New(componentItems, itemDelegate{}, 70, 20),
		components:     map[string]struct{}{},
		gitopsList:     list.New(gitopsTools, list.NewDefaultDelegate(), 70, 20),
	}
	m.clusterName.Placeholder = internal.DefaultClusterName
	m.clusterName.Focus()
	m.componentsList.Title = "Pick components to install"
	m.gitopsList.Title = "Which GitOps tool to use?"
	m.gitopsList.SetShowStatusBar(false)
	m.gitopsList.SetShowHelp(false)

	p := tea.NewProgram(m, tea.WithAltScreen())

	res, err := p.Run()
	if err != nil {
		return UserData{}, err
	}

	converted := res.(model)

	selectedComponents := make([]string, 0, len(converted.components))
	for component := range converted.components {
		selectedComponents = append(selectedComponents, component)
	}
	return UserData{
		ClusterName: converted.clusterName.Value(),
		Components:  selectedComponents,
		GitOpsTool:  converted.gitopsChoice,
	}, nil
}
