package menu

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/devonboyer/airbot"
	"github.com/devonboyer/airbot/bot"
)

const prompt = "airbot [? for menu]: "

type menu []command

type command struct {
	name    string
	help    string
	aliases []string
	runFunc func(args []string) error
}

var (
	mainMenu menu
	secrets  *airbot.Secrets
)

func init() {
	mainMenu = menu{
		{"?", "  Print menu", []string{"h", "help"}, func(_ []string) error { fmt.Print(mainMenu); return nil }},
		{"bot", "Start a bot", nil, botCommand},
		{"q", "  Quit and shut down", []string{"quit"}, func(_ []string) error { os.Exit(0); return nil }},
	}
}

func Run(args []string) {
	// This is a hack to avoid trying to parse secrets in production.
	secrets = airbot.MustReadSecrets("config")
	mainMenu.run(args)
}

func (m menu) run(args []string) {
	for {
		fmt.Print(prompt)
		buf := bufio.NewReader(os.Stdin)
		line, _ := buf.ReadString('\n')
		line = strings.TrimSpace(line)

		for _, cmd := range strings.Split(line, ";") {
			if err := m.dispatch(cmd); err != nil {
				fmt.Printf("Menu command failed: %s\n", err)
			}
		}
	}
}

func (m menu) dispatch(cmd string) error {
	// Parse args.
	args := strings.Split(cmd, " ")
	if len(args) == 0 {
		return fmt.Errorf("No command")
	}
	// Run command.
	for _, cmd := range m {
		if args[0] == cmd.name {
			return cmd.runFunc(args)
		}
		for _, alias := range cmd.aliases {
			if args[0] == alias {
				return cmd.runFunc(args)
			}
		}
	}
	return fmt.Errorf("%s: Command not found", args[0])
}

func (m menu) String() string {
	buf := &bytes.Buffer{}
	buf.WriteString("\n")
	for _, cmd := range m {
		fmt.Fprintf(buf, "[%s] %s\n", cmd.name, cmd.help)
	}
	buf.WriteString("\n")
	return buf.String()
}

func botCommand(args []string) error {
	fmt.Println("Starting bot...")

	source := newCLISource()
	b := airbot.NewBot(secrets, source)
	b.Run()
	defer b.Stop()

	for {
		fmt.Print("> ")
		buf := bufio.NewReader(os.Stdin)
		line, _ := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if line == "q" {
			fmt.Println("Stopping bot...")
			return nil
		}
		source.msgs <- bot.Message{Text: line}
	}
}
