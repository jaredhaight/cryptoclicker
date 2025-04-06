package internal

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"os"
	"strings"
)

func (app Application) Init() tea.Cmd {
	return app.MiningTimer.Init()
}

func (app Application) helpView() string {
	return "\n" + app.Help.ShortHelpView([]key.Binding{
		app.Keymap.buy,
		app.Keymap.collect,
		app.Keymap.quit,
	})
}

func (app Application) MiningProgressPercentage() float64 {
	elapsed := app.MiningTimeout - app.MiningTimer.Timeout
	e := float64(elapsed)
	t := float64(app.MiningTimeout)
	percentage := e / t
	return percentage
}

func (app Application) View() string {

	fullWidth, _, _ := term.GetSize(os.Stdout.Fd())
	//trim our width a little
	fullWidth -= 10
	halfWidth := fullWidth / 2

	app.MiningProgress.Width = fullWidth

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

	var progressStyle = lipgloss.NewStyle().
		Width(fullWidth).
		Padding(1)

	var storeStyle = lipgloss.NewStyle().
		Width(fullWidth).
		Align(lipgloss.Left).
		Border(lipgloss.NormalBorder())

	doc := strings.Builder{}
	crypto := strings.Builder{}

	crypto.WriteString(fmt.Sprintf("Available Crypto: %d\n", app.AvailableCrypto))
	crypto.WriteString(fmt.Sprintf("Crypto Per Click: %d\n", app.CryptoPerClick))
	crypto.WriteString(fmt.Sprintf("Mining Time: %s\n", app.MiningTimeout))

	if app.MiningTimer.Running() {
		crypto.WriteString(fmt.Sprintf("Status: Mining"))
	} else {
		crypto.WriteString("Status: Waiting for Collection")
	}

	doc.WriteString(lipgloss.JoinVertical(
		lipgloss.Top,
		titleStyle.Render("Mine them cryptos!!"),
		lipgloss.JoinHorizontal(lipgloss.Top,
			cryptoStyle.Render(crypto.String()),
			inventoryStyle.Render(fmt.Sprintf("GPU: %d\nBots: %d", app.GpuCount, app.BotnetCount)),
		),
		progressStyle.Render(app.MiningProgress.ViewAs(app.MiningProgressPercentage())),
		storeStyle.Align(lipgloss.Left).Render(app.Shop.View()),
		app.helpView(),
	))

	return doc.String()
}
