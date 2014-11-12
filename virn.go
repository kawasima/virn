package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"unicode/utf8"
)

import goncurses "github.com/gbin/goncurses"

func findFile(filename string) string {
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if pair[0] == "HOME" {
			return filepath.Join(pair[1], "/.virn/", filename)
		}
	}
	return filename
}

func renderInitialWindow(w *goncurses.Window) {
	y, _ := w.MaxYX()
	for i := 1; i < y-1; i++ {
		w.MoveAddChar(i, 0, '~')
	}
	w.Move(0, 0)
	w.Refresh()
}

func terminateWindow() {
	log.Println("TERMINATE")
	goncurses.CBreak(false)
	goncurses.NewLines(false)
	goncurses.End()
}

func main() {
	filename := os.Args[len(os.Args)-1]
	filepath := findFile(filename)
	bs, _ := ioutil.ReadFile(filepath)
	s := string(bs[:])

	scr, err := goncurses.Init()
	if err != nil {

	}
	defer terminateWindow()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		terminateWindow()
		os.Exit(1)
	}()

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
	renderInitialWindow(scr)

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
				scr.Print("\n")
				scr.DelChar()
				for true {
					i++
					if i >= len(chars) {
						break
					}
					r, _ := utf8.DecodeRuneInString(chars[i])
					if r == '\t' || r == ' ' {
						scr.Printf("%c", r)
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

	ioutil.WriteFile(filename, []byte(s), 0644)
}
