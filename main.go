package main

import (
	"blog_aggregator/internal/config"
	"fmt"
	"log"
	"os"
)

type state struct {
	cfg *config.Config
}

func main() {
	cfg, err := config.ReadJson()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}

	fmt.Printf("Read config: %+v\n", cfg)

	programState := &state{
		cfg: &cfg,
	}

	commands := commands{
		registeredCommands: make(map[string]func(*state, command) error),
	}

	commands.register("login", handlerLogin)

	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
	}

	cmd := command{
		Name: args[1],
		Args: args[2:],
	}

	if err := commands.run(programState, cmd); err != nil {
		log.Fatal(err)
	}

	/*err = cfg.SetUser("Steffen")
	if err != nil {
		log.Fatalf("couldn't set current user: %v", err)
	}
	*/

	cfg, err = config.ReadJson()

	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	fmt.Printf("Read config again: %+v\n", cfg)

}
