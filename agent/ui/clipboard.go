package ui

import (
	"github.com/atotto/clipboard"
)

// CopyToClipboard copies the provided content string to the system clipboard
func CopyToClipboard(content string) error {
	// Copy content to clipboard
	return clipboard.WriteAll(content)

}
