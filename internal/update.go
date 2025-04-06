package internal

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"time"
)

func (app Application) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		if msg.Timeout {
			app.CanCollect = true
		}
		app.MiningTimer, cmd = app.MiningTimer.Update(msg)
		return app, cmd

	case tea.WindowSizeMsg:
		var cmd tea.Cmd
		app.Shop.SetSize(200, 20)
		return app, cmd

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, app.Keymap.quit):
			app.Quitting = true
			return app, tea.Quit
		case key.Matches(msg, app.Keymap.buy):
			i := app.Shop.SelectedItem().(Item)
			if app.AvailableCrypto >= i.cost {
				if i.title == gpu7000DisplayName {
					app.GpuCount++
					app.CryptoPerClick++
				} else if i.title == botnetDisplayName {
					app.BotnetCount += 100
					app.CryptoPerClick += 2
				}
				app.AvailableCrypto -= i.cost
			}
			var cmd tea.Cmd
			return app, cmd

		case key.Matches(msg, app.Keymap.collect):
			// if we can collect, start the MiningTimer
			if app.CanCollect {
				app.CanCollect = false
				app.AvailableCrypto += app.CryptoPerClick

				// timeout adjustments
				// gpu = 100ms per gpu
				gpuAdjustment := app.GpuCount * 500
				adjustment, err := time.ParseDuration(fmt.Sprintf("%dms", gpuAdjustment))

				if err != nil {
					log.Fatal(err)
				}
				seconds := adjustment.Seconds()
				// we don't want GPUs to do too much for us
				if seconds > 1 {
					adjustment = time.Second
				}

				app.MiningTimeout = app.BaseTimeout - adjustment

				app.MiningTimer = timer.NewWithInterval(app.MiningTimeout, time.Millisecond)
				return app, app.MiningTimer.Init()
			}
		}
	}

	var cmd tea.Cmd
	app.Shop, cmd = app.Shop.Update(msg)
	return app, cmd
}
