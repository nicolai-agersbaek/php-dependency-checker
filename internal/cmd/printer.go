package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"strings"
)

type Printer interface {
	Title(title string)
	Line(line string)
	Lines(lines []string)
	LinesMax(lines []string, max int)
	LinesWithTitle(title string, lines []string)
	LinesWithTitleMax(title string, lines []string, max int)
}

type printer struct {
	c *cobra.Command
}

func NewPrinter(c *cobra.Command) Printer {
	return &printer{c: c}
}

func (p *printer) LinesWithTitle(title string, lines []string) {
	if len(lines) > 0 {
		title = fmt.Sprintf("%s (%d)", title, len(lines))
		p.Title(title)
		p.Lines(lines)
	}
}

func (p *printer) LinesWithTitleMax(title string, lines []string, max int) {
	if len(lines) > 0 {
		title = fmt.Sprintf("%s (%d)", title, len(lines))
		p.Title(title)
		p.LinesMax(lines, max)
	}
}

func (p *printer) Title(title string) {
	titleBreak := strings.Repeat("-", len(title))

	p.c.Println(titleBreak)
	p.c.Println(title)
	p.c.Println(titleBreak)
}

func (p *printer) Line(line string) {
	p.c.Println(line)
}

func (p *printer) Lines(lines []string) {
	for _, line := range lines {
		p.Line(line)
	}
}

func (p *printer) LinesMax(lines []string, max int) {
	if len(lines) > max {
		lines = lines[:max]

		defer func() {
			p.c.Println("...")
		}()
	}

	p.Lines(lines)
}

type Verbosity uint

const (
	VerbosityNone     Verbosity = 1 << iota
	VerbosityNormal   Verbosity = 1 << iota
	VerbosityDetailed Verbosity = 1 << iota
	VerbosityDebug    Verbosity = 1 << iota
)

type VerbosePrinter interface {
	Printer

	VTitle(title string, verbosity Verbosity)
	VLine(line string, verbosity Verbosity)
	VLines(lines []string, verbosity Verbosity)
	VLinesMax(lines []string, max int, verbosity Verbosity)
	VLinesWithTitle(title string, lines []string, verbosity Verbosity)
	VLinesWithTitleMax(title string, lines []string, max int, verbosity Verbosity)

	GetVerbosity() Verbosity
	SetVerbosity(verbosity Verbosity)
}

type verbosePrinter struct {
	Printer
	verbosity Verbosity
}

func NewVerbosePrinter(printer Printer, verbosity Verbosity) *verbosePrinter {
	return &verbosePrinter{Printer: printer, verbosity: verbosity}
}

func (p *verbosePrinter) VTitle(title string, verbosity Verbosity) {
	if verbosity <= p.verbosity {
		p.Title(title)
	}
}

func (p *verbosePrinter) VLine(line string, verbosity Verbosity) {
	if verbosity <= p.verbosity {
		p.Line(line)
	}
}

func (p *verbosePrinter) VLines(lines []string, verbosity Verbosity) {
	if verbosity <= p.verbosity {
		p.Lines(lines)
	}
}

func (p *verbosePrinter) VLinesMax(lines []string, max int, verbosity Verbosity) {
	if verbosity <= p.verbosity {
		p.LinesMax(lines, max)
	}
}

func (p *verbosePrinter) VLinesWithTitle(title string, lines []string, verbosity Verbosity) {
	if verbosity <= p.verbosity {
		p.LinesWithTitle(title, lines)
	}
}

func (p *verbosePrinter) VLinesWithTitleMax(title string, lines []string, max int, verbosity Verbosity) {
	if verbosity <= p.verbosity {
		p.LinesWithTitleMax(title, lines, max)
	}
}

func (p *verbosePrinter) GetVerbosity() Verbosity {
	return p.verbosity
}

func (p *verbosePrinter) SetVerbosity(verbosity Verbosity) {
	p.verbosity = verbosity
}
