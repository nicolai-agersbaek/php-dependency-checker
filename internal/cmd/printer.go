package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

type Printer struct {
	c *cobra.Command
}

func NewPrinter(c *cobra.Command) *Printer {
	return &Printer{c: c}
}

func (p *Printer) LinesWithTitle(title string, lines []string) {
	if len(lines) > 0 {
		title = fmt.Sprintf("%s (%d)", title, len(lines))
		p.Title(title)
		p.Lines(lines)
	}
}

func (p *Printer) LinesWithTitleMax(title string, lines []string, max int) {
	if len(lines) > 0 {
		title = fmt.Sprintf("%s (%d)", title, len(lines))
		p.Title(title)
		p.LinesMax(lines, max)
	}
}

func (p *Printer) Title(title string) {
	titleBreak := strings.Repeat("-", len(title))

	p.c.Println(titleBreak)
	p.c.Println(title)
	p.c.Println(titleBreak)
}

func (p *Printer) Lines(lines []string) {
	for _, line := range lines {
		p.c.Println(line)
	}
}

func (p *Printer) LinesMax(lines []string, max int) {
	if len(lines) > max {
		lines = lines[:max]

		defer func() {
			p.c.Println("...")
		}()
	}

	p.Lines(lines)
}
