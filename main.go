package main

import (
	"bytes"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"os/exec"
	"strings"
)

func main() {
	buffs := make(map[string]string)
	buffs["sleepy"] = ""
	buffs["grumpy"] = ""
	currentbuff := "sleepy"

	app := tview.NewApplication()
	middle := tview.NewTextView()
	cli := tview.NewInputField()
	middle.SetBorder(true).SetTitle("CLI Result")
	list := tview.NewList().
		AddItem("Sleepy", "2 in 1 out 92000 sats", 'a', func() {
			cli.SetText("")
			currentbuff = "sleepy"
			middle.SetText(buffs["sleepy"])
			app.SetFocus(cli)
		}).
		AddItem("Grumpy", "Some explanatory text", 'b', func() {
			cli.SetText("")
			currentbuff = "grumpy"
			middle.SetText(buffs["grumpy"])
			app.SetFocus(cli)
		}).
		AddItem("Dopy", "Some explanatory text", 'c', nil).
		AddItem("Sneezy", "Some explanatory text", 'd', nil).
		AddItem("Quit", "Press to exit", 'q', func() {
			app.Stop()
		})

	list.SetBorder(true).SetTitle("Nodes (Ctrl+n)")
	cli.
		SetPlaceholder("Enter cli command").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldWidth(0).SetBorder(true).SetTitle("CLI (Clrl+l)")

	cli.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		if key.Key() == tcell.KeyEnter {

			args, err := parseCommandLine(cli.GetText())
			if err != nil {
				fmt.Fprintf(middle, "%s\n", err.Error())
			}
			cmd := exec.Command(args[0], args[1:]...)
			cmd.Stdin = strings.NewReader("some input")
			var out bytes.Buffer
			cmd.Stdout = &out
			err = cmd.Run()
			if err != nil {
				fmt.Fprintf(middle, "%s\n", err.Error())
			}

			fmt.Fprintf(middle, "%s\n", out.String())
			buffs[currentbuff] += out.String()

			cli.SetText("")
		}
		return key
	})

	app.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		if key.Key() == tcell.KeyCtrlN {
			app.SetFocus(list)
		} else if key.Key() == tcell.KeyCtrlL {
			cli.SetText("")
			app.SetFocus(cli)
		} else if key.Key() == tcell.KeyCtrlR {
			app.SetFocus(middle)
		}
		return key
	})

	flex := tview.NewFlex().
		AddItem(list, 40, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(cli, 3, 1, false).
			AddItem(middle, 0, 3, false), 0, 2, false).
		AddItem(tview.NewBox().SetBorder(true).SetTitle("Node info"), 40, 1, false)
	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
