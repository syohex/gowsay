package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"io"
	"log"
	"text/template"

	"github.com/mattn/go-runewidth"
	"github.com/mitchellh/go-wordwrap"
	flag "github.com/ogier/pflag"
)

type face struct {
	Eyes     string
	Tongue   string
	Thoughts string
	cowfile  string
}

var borg, dead, greedy, paranoid, stoned, tired, wired, young, think *bool
var list *bool
var columns *int32
var cowfile *string

func newFace() *face {
	f := &face{
		Eyes:    "oo",
		Tongue:  "  ",
		cowfile: *cowfile,
	}

	if *borg {
		f.Eyes = "=="
	}
	if *dead {
		f.Eyes = "xx"
		f.Tongue = "U "
	}
	if *greedy {
		f.Eyes = "$$"
	}
	if *paranoid {
		f.Eyes = "@@"
	}
	if *stoned {
		f.Eyes = "**"
		f.Tongue = "U "
	}
	if *tired {
		f.Eyes = "--"
	}
	if *wired {
		f.Eyes = "OO"
	}
	if *young {
		f.Eyes = ".."
	}

	return f
}

func readInput(args []string) []string {
	var tmps []string
	if len(args) == 0 {
		s := bufio.NewScanner(os.Stdin)

		for s.Scan() {
			tmps = append(tmps, s.Text())
		}
	} else {
		tmps = args
	}

	var msgs []string
	for i := 0; i < len(tmps); i++ {
		expand := strings.Replace(tmps[i], "\t", "        ", -1)

		tmp := wordwrap.WrapString(expand, uint(*columns))
		for _, s := range strings.Split(tmp, "\n") {
			msgs = append(msgs, s)
		}
	}

	return msgs
}

func setPadding(msgs []string, width int) []string {
	var ret []string
	for _, m := range msgs {
		s := m + strings.Repeat(" ", width - runewidth.StringWidth(m))
		ret = append(ret, s)
	}

	return ret
}

func constructBallon(f *face, msgs []string, width int) string {
	var borders []string
	line := len(msgs)

	if *think {
		f.Thoughts = "o"
		borders = []string{"(", ")", "(", ")", "(", ")"}
	} else {
		f.Thoughts = "\\"
		if line == 1 {
			borders = []string{"<", ">"}
		} else {
			borders = []string{"/", "\\", "\\", "/", "|", "|"}
		}
	}

	var lines []string

	topBorder := " " + strings.Repeat("_", width+2)
	bottomBoder := " " + strings.Repeat("-", width+2)

	lines = append(lines, topBorder)
	if line == 1 {
		s := fmt.Sprintf("%s %s %s", borders[0], msgs[0], borders[1])
		lines = append(lines, s)
	} else {
		s := fmt.Sprintf(`%s %s %s`, borders[0], msgs[0], borders[1])
		lines = append(lines, s)
		i := 1
		for ; i < line-1; i++ {
			s = fmt.Sprintf(`%s %s %s`, borders[4], msgs[i], borders[5])
			lines = append(lines, s)
		}
		s = fmt.Sprintf(`%s %s %s`, borders[2], msgs[i], borders[3])
		lines = append(lines, s)
	}

	lines = append(lines, bottomBoder)
	return strings.Join(lines, "\n")
}

func maxWidth(msgs []string) int {
	max := -1
	for _, m := range msgs {
		l := runewidth.StringWidth(m)
		if l > max {
			max = l
		}
	}

	return max
}

func renderCow(f *face, w io.Writer) {
	t := template.Must(template.New("cow").Parse(cows[f.cowfile]))

	if err := t.Execute(w, f); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func main() {
	borg = flag.BoolP("borg", "b", false, "borg eyes")
	dead = flag.BoolP("dead", "d", false, "dead eyes and tongue")
	greedy = flag.BoolP("greedy", "g", false, "greedy eyes")
	paranoid = flag.BoolP("paranoid", "p", false, "paranoid eyes")
	stoned = flag.BoolP("stoned", "s", false, "stoned eyes and tongue")
	tired = flag.BoolP("tired", "t", false, "tired eyes")
	wired = flag.BoolP("wired", "w", false, "wired eyes")
	young = flag.BoolP("young", "y", false, "young eyes")
	think = flag.Bool("think", false, "think version")

	list = flag.BoolP("list", "l", false, "list cow files")
	cowfile = flag.StringP("type", "f", "default", "specify cow file")
	columns = flag.Int32P("columns", "W", 40, "columns")

	flag.Parse()

	if *list {
		displayCows()
		os.Exit(0)
	}

	inputs := readInput(flag.Args())
	width := maxWidth(inputs)
	messages := setPadding(inputs, width)

	f := newFace()
	balloon := constructBallon(f, messages, width)

	fmt.Println(balloon)
	renderCow(f, os.Stdout)
}
