package main

import (
	"bytes"
	"fmt"
	"github.com/atotto/clipboard"
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

type alias struct {
	Name *string
	Path *string
}

var app = tview.NewApplication()

func main() {

	cliresult := tview.NewTextView().SetDynamicColors(true) //.SetWrap(false)
	cli := tview.NewInputField()
	list := tview.NewList()

	cmd := exec.Command("lnet-cli", "alias")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Fprintf(cliresult, "%s\n", err.Error())
	}

	lines := strings.Split(out.String(), "\n")
	currentnode := ""
	aliases := make(map[string]*alias)
	for i, line := range lines {
		args, _ := parseCommandLine(line)
		if len(args) > 0 {
			name := strings.Split(args[1], "-")[1]
			cmd := args[2][1:]
			aliases[name] = &alias{&name, &cmd}
			if i == 0 {
				currentnode = name
			}
		}
	}

	nodes := make(map[string]*node)

	cliresult.SetBorder(true).SetTitle("CLI Result (ctrl+r)")

	i := 0
	for _, a := range aliases {
		s := -1
		anode := &node{"", []string{}, &s}
		nodes[*a.Name] = anode

		name := a.Name
		list.AddItem(*a.Name, "", rune('a'+byte(i)), func() {
			cli.SetText("")
			currentnode = *name
			cliresult.SetText(nodes[*name].Buff)
			app.SetFocus(cli)
		})
		i++
	}

	firstpath := aliases[currentnode].Path
	runpathsegs := strings.Split(*firstpath, "/")
	conf := append(runpathsegs[:len(runpathsegs)-1], "bitcoin.conf")
	confpath := strings.Split(strings.Join(conf, "/"), "=")[1]
	confcmd := fmt.Sprintf("bitcoin-cli --conf=%s", confpath)
	name := "Regtest"
	aliases[name] = &alias{&name, &confcmd}
	s := -1
	anode := &node{"", []string{}, &s}
	nodes[name] = anode

	list.AddItem(name, "", 'r', func() {
		cli.SetText("")
		currentnode = name
		cliresult.SetText(nodes[name].Buff)
		app.SetFocus(cli)
	})

	list.AddItem("Quit", "Press to exit", 'q', func() {
		app.Stop()
	})

	list.SetBorder(true).SetTitle("Nodes (Ctrl+n)")
	cli.
		SetPlaceholder("Enter cli command - use Ctrl+v to paste (no shift)").
		SetFieldBackgroundColor(tcell.ColorBlack).
		SetFieldWidth(0).SetBorder(true).SetTitle("CLI (Clrl+l)")

	cli.SetInputCapture(func(key *tcell.EventKey) *tcell.EventKey {
		if key.Key() == tcell.KeyEnter {
			text := cli.GetText()
			cmdfmt := fmt.Sprintf("[#ff0000]# %s[white]\n", text)
			fmt.Fprintf(cliresult, cmdfmt)
			if text == "" {
				fmt.Fprintf(cliresult, "Please provide a command to execute\n")
				return key
			}
			args, err := parseCommandLine(text)
			if err != nil {
				fmt.Fprintf(cliresult, "%s\n", err.Error())
			}
			clicmd := strings.Split(*aliases[currentnode].Path, " ")
			cliarg := []string{clicmd[1]}
			cliargs := append(cliarg, args...)
			cmd := exec.Command(clicmd[0], cliargs...)
			cmd.Stdin = strings.NewReader("some input")
			var out bytes.Buffer
			cmd.Stdout = &out
			err = cmd.Run()
			if err != nil {
				fmt.Fprintf(cliresult, "%s\n", err.Error())
			}

			fmt.Fprintf(cliresult, "%s\n", out.String())
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
		} else if key.Key() == tcell.KeyCtrlV {
			clip, err := clipboard.ReadAll()
			if err != nil {
				fmt.Fprintf(cliresult, "%s\n", err.Error())
			} else {
				full := strings.Replace(clip, "\n", "", -1)
				hack := strings.Split(full, "││") // TODO: fix if selected with result pane active "│║" "║│"
				if len(hack) == 3 {
					clip = fmt.Sprintf("%s%s", hack[0], hack[2])
				}
				cli.SetText(fmt.Sprintf("%s%s", cli.GetText(), clip)) // TODO: this only paste to end, fix for insert
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
			app.SetFocus(cliresult)
		}
		return key
	})

	flex := tview.NewFlex().
		AddItem(list, 40, 1, true).
		AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
			AddItem(cli, 3, 1, false).
			AddItem(cliresult, 0, 3, false), 0, 2, false)

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
