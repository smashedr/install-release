package cmd

import (
	"fmt"
	"regexp"
	"strings"
)

func parseRepository(repository string) (owner, repo string, err error) {
	var repoPattern = regexp.MustCompile(`^[a-zA-Z0-9_.-]+/[a-zA-Z0-9_.-]+$`)

	if !repoPattern.MatchString(repository) {
		return "", "", fmt.Errorf("repository must be in format: owner/repo")
	}
	split := strings.Split(repository, "/")
	return split[0], split[1], nil
}

//func headString(text string, length int) string {
//	lines := strings.SplitN(strings.TrimSpace(text), "\n", length+1)
//	if len(lines) > length {
//		lines = lines[:length]
//	}
//	return strings.Join(lines, "\n")
//}
