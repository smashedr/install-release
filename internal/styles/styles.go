package styles

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var Success = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#00ff00"))
var Failure = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#ff0000"))

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

func PrintS(key, format string, args ...interface{}) {
	fmt.Println(Success.Render(key) + " " + fmt.Sprintf(format, args...))
}
func PrintF(key, format string, args ...interface{}) {
	fmt.Println(Failure.Render(key) + " " + fmt.Sprintf(format, args...))
}
func PrintKV(key, value string) {
	fmt.Println(Key.Render(key) + " " + Value.Render(value))
}

func RenderTable(rows [][]string, headers ...string) {
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(TableBorder).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return TableHeader
			}
			return TableRow
		}).
		Headers(headers...).
		Rows(rows...)
	fmt.Println(t)
}
