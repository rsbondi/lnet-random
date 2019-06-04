package main

import (
	"bytes"
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"os/exec"
	"strings"
)

type node struct {
	Buff     string
	Cmds     []string
	CmdIndex *int
}

func main() {
	nodes := make(map[string]*node)
	s := -1
	sleepy := &node{"", []string{}, &s}
	nodes["sleepy"] = sleepy
	g := -1
	grumpy := &node{"", []string{}, &g}
	nodes["grumpy"] = grumpy
	currentnode := "sleepy"

	app := tview.NewApplication()
	middle := tview.NewTextView().SetDynamicColors(true)

	cli := tview.NewInputField()
	middle.SetBorder(true).SetTitle("CLI Result")
	list := tview.NewList().
		AddItem("Sleepy", "2 in 1 out 92000 sats", 'a', func() {
			cli.SetText("")
			currentnode = "sleepy"
			middle.SetText(nodes["sleepy"].Buff)
			app.SetFocus(cli)
		}).
		AddItem("Grumpy", "Some explanatory text", 'b', func() {
			cli.SetText("")
			currentnode = "grumpy"
			middle.SetText(nodes["grumpy"].Buff)
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
			text := cli.GetText()
			cmdfmt := fmt.Sprintf("[#ff0000]# %s[white]\n", text)
			fmt.Fprintf(middle, cmdfmt)
			if text == "" {
				fmt.Fprintf(middle, "Please provide a command to execute\n")
				return key
			}
			args, err := parseCommandLine(text)
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
			nodes[currentnode].Buff += cmdfmt
			nodes[currentnode].Buff += out.String()
			nodes[currentnode].Cmds = append(nodes[currentnode].Cmds, cli.GetText())
			*nodes[currentnode].CmdIndex = len(nodes[currentnode].Cmds)

			cli.SetText("")
		} else if key.Key() == tcell.KeyUp {
			index := nodes[currentnode].CmdIndex
			if *index > 0 {
				*index = *index - 1
			}
			if *index >= 0 && *index < len(nodes[currentnode].Cmds) {
				cli.SetText(nodes[currentnode].Cmds[*index])
			}
		} else if key.Key() == tcell.KeyDown {
			index := nodes[currentnode].CmdIndex
			if *index == len(nodes[currentnode].Cmds)-1 {
				cli.SetText("")
				*index = *index + 1
				return key
			}
			if *index < len(nodes[currentnode].Cmds)-1 {
				*index = *index + 1
			}
			if *index >= 0 && *index < len(nodes[currentnode].Cmds) {
				cli.SetText(nodes[currentnode].Cmds[*index])
			}
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
