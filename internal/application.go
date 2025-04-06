package internal

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/timer"
	"time"
)

const gpu7000DisplayName = "GPU 7000"
const botnetDisplayName = "Botnet +100"
const timeout = time.Second * 5

type Application struct {
	CryptoPerClick  int
	GpuCount        int
	BotnetCount     int
	CanCollect      bool
	BaseTimeout     time.Duration
	MiningProgress  progress.Model
	MiningTimeout   time.Duration
	MiningTimer     timer.Model
	BotnetTimeout   time.Duration
	BotnetTimer     timer.Model
	Keymap          keymap
	Help            help.Model
	Quitting        bool
	Cursor          int
	Shop            list.Model
	AvailableCrypto int
}

func NewApplication() *Application {
	items := []list.Item{
		Item{gpu7000DisplayName, "The latest in GPU technology", 2},
		Item{botnetDisplayName, "Adds a fleet of Thermostats and Printers to your army", 5},
	}

	a := Application{
		CryptoPerClick: 1,
		Shop:           list.New(items, list.NewDefaultDelegate(), 0, 0),
		BaseTimeout:    timeout,
		MiningProgress: progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C")),
		MiningTimeout:  time.Second * 5,
		MiningTimer:    timer.NewWithInterval(timeout, time.Millisecond),
		Keymap: keymap{
			quit: key.NewBinding(
				key.WithKeys("q", "ctrl+c"),
				key.WithHelp("q", "quit"),
			),
			collect: key.NewBinding(key.WithKeys(" "),
				key.WithHelp("space", "collect"),
			),
			buy: key.NewBinding(key.WithKeys("b", "enter"), key.WithHelp("enter", "buy")),
		},
		Help: help.New(),
	}

	a.Shop.Title = "CryptoBarn Catalog"
	return &a
}

type keymap struct {
	buy     key.Binding
	quit    key.Binding
	collect key.Binding
}

type Item struct {
	title string
	desc  string
	cost  int
}

func (i Item) Title() string       { return i.title }
func (i Item) Description() string { return i.desc }
func (i Item) FilterValue() string { return i.title }
