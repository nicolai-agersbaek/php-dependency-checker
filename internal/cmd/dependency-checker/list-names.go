package dependency_checker

import (
	"fmt"
	"github.com/spf13/cobra"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/cmd"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/dependency-checker"
	"gitlab.zitcom.dk/smartweb/proj/php-dependency-checker/internal/util/slices"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func init() {
	rootCmd.AddCommand(listNamesCmd)
}

var listNamesCmd = &cobra.Command{
	Use:   "list-names",
	Short: "List built-in PHP functions, classes and constants.",
	Args:  cobra.NoArgs,
	Run:   listNames,
}

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
	filtered := slices.FilterAllString(trimmed, dependency_checker.IsEmpty, nameFilter)

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

func listNames(c *cobra.Command, args []string) {
	nameTypes := map[string]string{
		//"functions" : "get_defined_functions(true)[\"internal\"]",
		//"classes" : "get_declared_classes()",
		//"traits" : "get_declared_traits()",
		"interfaces": "get_declared_interfaces()",
	}

	printConstants(c)

	for nameType, listCode := range nameTypes {
		printNames(c, nameType, listCode)
	}
	//
	for _, interfaceName := range cOut.lines {
		c.Printf("#%v\n", interfaceName)
	}
}

func printConstants(c *cobra.Command) {
	cOut.reset()

	phpCode := strings.Join([]string{printLinesFunc, constantsListFunc, constantsListCode}, ";")
	phpCode += ";"
	listFunctions := command("php", "-r", phpCode)

	cmd.CheckError(listFunctions.Run())

	c.Printf("Defined constants: %d\n", len(cOut.lines))

	//for k, v := range cOut.lines {
	//	c.Printf("[%d]: %s\n", k, strings.Replace(v, "\n", "%NEWLINE", -1))
	//}

	//p := newPrinter(c)
	//p.lines(cOut.lines)
}

func phpEval(phpCode string, out *cmdOut) *cmdOut {
	if out != nil {
		out = newCmdOut()
	}

	p := commandOut(out, "php", "-r", phpCode)

	cmd.CheckError(p.Run())

	return out
}

func printNames(c *cobra.Command, nameType, listCode string) {
	cOut.reset()

	printFuncs := phpPrintLines(listCode)
	listFunctions := command("php", "-r", printFuncs)

	cmd.CheckError(listFunctions.Run())

	c.Printf("Defined %s: %d\n", nameType, len(cOut.lines))
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
