package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"os"
	"strings"
	"time"
)

const gpu7000DisplayName = "GPU 7000"
const botnetDisplayName = "Botnet +100"

type application struct {
	gpuCount      int
	botnetCount   int
	canMine       bool
	miningTimeout time.Duration
	miningTimer   timer.Model
	botnetTimeout time.Duration
	botnetTimer   timer.Model
	keymap        keymap
	help          help.Model
	quitting      bool
	cursor        int
	shop          list.Model
	count         int
}

type keymap struct {
	buy  key.Binding
	quit key.Binding
	mine key.Binding
}

type item struct {
	title string
	desc  string
	cost  int
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

var docStyle = lipgloss.NewStyle().Padding(1, 2, 1, 2)

func (m application) Init() tea.Cmd {
	return m.miningTimer.Init()
}

func (m application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		if msg.Timeout {
			m.canMine = true
		}
		m.miningTimer, cmd = m.miningTimer.Update(msg)
		return m, cmd

	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		m.shop.SetSize(200, 20)
		return m, cmd

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keymap.buy):
			i := m.shop.SelectedItem().(item)
			if m.count >= i.cost {
				if i.title == gpu7000DisplayName {
					m.gpuCount++
				} else if i.title == botnetDisplayName {
					m.botnetCount += 100
				}
				m.count -= i.cost
			}
			var cmd tea.Cmd
			return m, cmd

		case key.Matches(msg, m.keymap.mine):
			// if we can mine, start the miningTimer
			if m.canMine {
				m.count += 1 + m.gpuCount
				m.canMine = false
				m.miningTimer = timer.NewWithInterval(m.miningTimeout, time.Millisecond)
				return m, m.miningTimer.Init()
			}
		}
	}

	var cmd tea.Cmd
	m.shop, cmd = m.shop.Update(msg)
	return m, cmd
}

func (m application) helpView() string {
	return "\n" + m.help.ShortHelpView([]key.Binding{
		m.keymap.buy,
		m.keymap.mine,
		m.keymap.quit,
	})
}

func (m application) View() string {

	fullWidth, _, _ := term.GetSize(os.Stdout.Fd())
	//trim our width a little
	fullWidth -= 10
	halfWidth := fullWidth / 2

	var titleStyle = lipgloss.NewStyle().
		Width(fullWidth).
		Bold(true).
		Align(lipgloss.Center, lipgloss.Center)

	var cryptoStyle = lipgloss.NewStyle().
		Width(halfWidth).
		Border(lipgloss.NormalBorder()).
		Align(lipgloss.Left).
		Padding(1)

	var inventoryStyle = lipgloss.NewStyle().
		Width(halfWidth).
		Border(lipgloss.NormalBorder()).
		Align(lipgloss.Left).
		Padding(1)

	var storeStyle = lipgloss.NewStyle().
		Width(fullWidth).
		Align(lipgloss.Left).
		Border(lipgloss.NormalBorder())

	doc := strings.Builder{}
	crypto := strings.Builder{}

	crypto.WriteString(fmt.Sprintf("Crypto: %d\n", m.count))

	if m.miningTimer.Running() {
		crypto.WriteString(fmt.Sprintf("Mining: %s", m.miningTimer.View()))
	} else {
		crypto.WriteString("Collect your Crypto!")
	}

	doc.WriteString(lipgloss.JoinVertical(
		lipgloss.Top,
		titleStyle.Render("Mine them cryptos!!"),
		lipgloss.JoinHorizontal(lipgloss.Top,
			cryptoStyle.Render(crypto.String()),
			inventoryStyle.Render(fmt.Sprintf("GPU: %d\nBots: %d", m.gpuCount, m.botnetCount)),
		),
		storeStyle.Align(lipgloss.Left).Render(m.shop.View()),
	))

	return doc.String()
}

func main() {
	const timeout = time.Second * 5
	items := []list.Item{
		item{gpu7000DisplayName, "The latest in GPU technology", 2},
		item{botnetDisplayName, "Adds a fleet of Thermostats and Printers to your army", 5},
	}

	m := application{
		shop:          list.New(items, list.NewDefaultDelegate(), 0, 0),
		miningTimeout: time.Second * 5,
		miningTimer:   timer.NewWithInterval(timeout, time.Millisecond),
		keymap: keymap{
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
			mine: key.NewBinding(key.WithKeys(" "),
				key.WithHelp("space", "mine"),
			),
			buy: key.NewBinding(key.WithKeys("b", "enter"), key.WithHelp("b", "buy")),
		},
		help: help.New(),
	}

	m.shop.Title = "CryptoBarn Catalog"
	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Uh oh, we encountered an error:", err)
		os.Exit(1)
	}
}
