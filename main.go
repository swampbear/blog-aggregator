package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/swampbear/blog-aggregator/internal/config"
	"github.com/swampbear/blog-aggregator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println(err)
	}
	s := state{cfg: &cfg}

	// create database and add it to state
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer db.Close()
	dbQueries := database.New(db)
	s.db = dbQueries

	//prapare for incoming commands and arguments
	cmds := commands{map[string]func(*state, command) error{}}
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Not enough arguments were provided")
		os.Exit(1)
	}

	//register handler fumctions from commands.go
	cmds.register("login", handlerLogin)
	cmds.register("register", handlerRegister)
	cmds.register("reset", handlerReset)
	cmds.register("users", handlerUsers)
	cmds.register("agg", handlerAgg)
	commandName := args[1]
	arguments := args[2:]

	cmd := command{name: commandName, arguments: arguments}

	err = cmds.run(&s, cmd)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

}
