package styles

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
)

var Key = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#00ADD8")).
	Width(12)
var Value = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#73FA91"))

var Command = lipgloss.NewStyle().
	Bold(true).
	Background(lipgloss.Color("#00ADD8")).
	Foreground(lipgloss.Color("#F7F7F7"))

//var Head = lipgloss.NewStyle().
//	Bold(true).
//	Background(lipgloss.Color("#00ADD8")).
//	Foreground(lipgloss.Color("#F7F7F7"))

var TableBorder = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#bc94f7")).
	Bold(true)
var TableHeader = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#00ADD8")).
	Bold(true).
	Align(lipgloss.Center)
var TableRow = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#F8F8F8")).
	PaddingLeft(1).PaddingRight(1)

func PrintKV(key, value string) {
	fmt.Println(Key.Render(key) + " " + Value.Render(value))
}
