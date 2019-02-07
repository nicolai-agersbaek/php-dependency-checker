package names

import (
	"fmt"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var cOut = newCmdOut()

type cmdOut struct {
	out   io.Writer
	lines []string
}

func newCmdOut() *cmdOut {
	return &cmdOut{nil, make([]string, 0)}
}

func (o *cmdOut) reset() {
	o.lines = make([]string, 0)
}

func (o *cmdOut) prune() {
	o.lines = slices.UniqueString(o.lines)
}

func (o *cmdOut) Write(p []byte) (n int, err error) {
	if o.out != nil {
		n, err = o.out.Write(p)
	} else {
		n = len(p)
	}

	if n > 0 {
		s := string(p)
		lines := strings.Split(s, "\n")
		o.lines = append(o.lines, lines...)

		return len(p), nil
	}

	return n, err
}

var namePattern = regexp.MustCompile("([A-Za-z0-9_\\\\]+)")

var trim slices.StringMapping = func(s string) string {
	return strings.Trim(s, " \n")
}

var nameFilter slices.StringFilter = func(s string) bool {
	return namePattern.MatchString(s)
}

func (o *cmdOut) append(lines []string) {
	trimmed := slices.MapString(lines, trim)
	filtered := slices.FilterAllString(trimmed, IsEmpty, nameFilter)

	o.lines = append(o.lines, filtered...)
}

var constantsListFunc = `
function getDefinedConstants()
{
    $byCategory = get_defined_constants(true);
    
    // Remove user-defined constants and merge the built-in constants into a single array
    unset($byCategory['user']);
    
    $allConstants = [];
    
    foreach ($byCategory as $C) {
        $allConstants[] = array_keys($C);
    }
    
    return array_merge(...$allConstants);
}
`

var printLinesFunc = `
function printLines(array $lines) : void
{
    \printf(\implode("\n", $lines));
}
`

var constantsListCode = `
printLines(getDefinedConstants())
`

func GetBuiltInNames() *Names {
	nameTypes := map[string]string{
		"functions":  "get_defined_functions(true)[\"internal\"]",
		"classes":    "get_declared_classes()",
		"traits":     "get_declared_traits()",
		"interfaces": "get_declared_interfaces()",
	}

	names := NewNames()

	functions := getNames(nameTypes["functions"])
	classes := getNames(nameTypes["classes"])
	interfaces := getNames(nameTypes["interfaces"])
	traits := getNames(nameTypes["traits"])
	constants := getConstants()

	names.Functions = functions
	names.Classes = classes
	names.Interfaces = interfaces
	names.Traits = traits
	names.Constants = constants

	names.Classes = append(names.Classes, interfaces...)
	names.Classes = append(names.Classes, traits...)

	names.Clean()

	return names
}

func getConstants() []string {
	cOut.reset()

	phpCode := strings.Join([]string{printLinesFunc, constantsListFunc, constantsListCode}, ";")
	phpCode += ";"
	listFunctions := command("php", "-r", phpCode)

	cmd.CheckError(listFunctions.Run())

	cOut.prune()
	return cOut.lines
}

func getNames(listCode string) []string {
	cOut.reset()

	printCode := phpPrintLines(listCode)
	listFunctions := command("php", "-r", printCode)

	cmd.CheckError(listFunctions.Run())

	cOut.prune()
	return cOut.lines
}

func phpPrintLines(linesCode string) string {
	printLines := "foreach(%s as $l){printf(\"$l\n\");}"

	return fmt.Sprintf(printLines, linesCode)
}

func command(name string, args ...string) *exec.Cmd {
	return commandOut(cOut, name, args...)
}

func commandOut(out *cmdOut, name string, args ...string) *exec.Cmd {
	command := exec.Command(name, args...)
	command.Stdout = out
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	return command
}
