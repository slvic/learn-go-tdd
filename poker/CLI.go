package poker

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const PlayerPrompt = "Please enter the number of players: "
const BadPlayerInputErrMsg = "Bad value received for number of players, please try again with a number"
const BadWinnerInputMsg = "invalid winner input, expect format of 'PlayerName wins'"

type CLI struct {
	in          *bufio.Scanner
	out         io.Writer
	game        Game
}

func NewCLI(in io.Reader, out io.Writer, game Game) *CLI {
	return &CLI{
		in:  bufio.NewScanner(in),
		out: out,
		game: game,
	}
}

func (cli *CLI) PlayPoker() {
	_, cliOutputErr := fmt.Fprint(cli.out, PlayerPrompt)

	numberOfPlayers, err := strconv.Atoi(cli.readLine())

	if err != nil {
		_, cliOutputErr = fmt.Fprint(cli.out, BadPlayerInputErrMsg)
		return
	}

	cli.game.Start(numberOfPlayers, cli.out)

	winnerInput := cli.readLine()
	winner, err := extractWinner(winnerInput)

	if err != nil {
		_, cliOutputErr = fmt.Fprint(cli.out, BadWinnerInputMsg)
		return
	}

	if cliOutputErr != nil {
		fmt.Printf("could not print to cli.out: %v\n", cliOutputErr)
	}

	cli.game.Finish(winner)
}

func extractWinner(userInput string) (string, error) {
	if !strings.Contains(userInput, " wins") {
		return "", errors.New(BadWinnerInputMsg)
	}
	return strings.Replace(userInput, " wins", "", 1), nil
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}

