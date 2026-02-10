package cmd

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/smashedr/install-release/internal/styles"
	"regexp"
	"strings"
)

func renderTable(rows [][]string, headers ...string) {
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(styles.TableBorder).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return styles.TableHeader
			}
			return styles.TableRow
		}).
		Headers(headers...).
		Rows(rows...)
	fmt.Println(t)
}

func formatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func parseRepository(repository string) (owner, repo string, err error) {
	var repoPattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+$`)

	if !repoPattern.MatchString(repository) {
		return "", "", fmt.Errorf("repository must be in format: owner/repo")
	}
	split := strings.Split(repository, "/")
	return split[0], split[1], nil
}
