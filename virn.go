package main

import (
	"io/ioutil"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

import goncurses "github.com/gbin/goncurses"

func FindFile(filename string) string {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if pair[0] == "HOME" {
			return pair[1] + "/.virn/" + filename
		}
	}
	return filename
}

func main() {
	filename := os.Args[len(os.Args)-1]
	filepath := FindFile(filename)
	bs, _ := ioutil.ReadFile(filepath)
	s := string(bs[:])

	scr, err := goncurses.Init()
	if err != nil {

	}
	defer goncurses.End()
	goncurses.CBreak(true)
	scr.ScrollOk(true)
	scr.Keypad(true)
	goncurses.Echo(false)
	goncurses.NewLines(true)

	y, x := scr.MaxYX()
	scr.Resize(y-1, x)

	status_win, err := goncurses.NewWindow(1, x, y-1, 0)
	if err != nil {
	}
	status_win.ClearOk(true)
	status_win.Printf("\"%s\" [New File]", filename)
	status_win.Refresh()

	scr.GetChar()
	status_win.Erase()
	status_win.Move(0, 0)
	status_win.AttrOn(goncurses.A_BOLD)
	status_win.Print("--INSERT--")
	status_win.Refresh()
	scr.Move(0, 0)
	scr.Refresh()

	chars := strings.Split(s, "")
	i := 0
	for true {
		ch := scr.GetChar()
		rune, _ := utf8.DecodeRuneInString(chars[i])
		if rune == '\n' {
			if ch == '\n' {
				for true {
					r, _ := utf8.DecodeRuneInString(chars[i])
					if unicode.IsSpace(r) {
						scr.Printf("%c", r)
						i++
						if i >= len(chars) {
							break
						}
					} else {
						break
					}
				}
			}
		} else {
			scr.Printf("%c", rune)
			i++
		}

		if i >= len(chars) {
			break
		}
	}
	scr.GetChar()
	status_win.Erase()
	status_win.Refresh()
	scr.Refresh()
	scr.GetChar()
	scr.GetChar()
	goncurses.End()
}
