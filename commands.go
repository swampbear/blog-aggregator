package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/swampbear/blog-aggregator/internal/config"
	"github.com/swampbear/blog-aggregator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func handlerFeeds(s *state, _ command) error {
	feeds, err := s.db.GetAllFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error: failed to get feeds: %w", err)
	}
	for _, feed := range feeds {
		userId := feed.UserID.UUID
		creator, err := s.db.GetUserById(context.Background(), userId)
		if err != nil {
			return fmt.Errorf("Error: could not find user with id %v: %w", userId, err)
		}
		fmt.Printf("- name: %s, \n  URL: %s, \n  Created by: %s\n", feed.Name.String, feed.Url.String, creator.Name.String)
	}
	return nil

}

func handlerAddFeed(s *state, cmd command) error {

	if len(cmd.arguments) < 2 {
		return fmt.Errorf("Error: Not enough arguments")
	}
	name := sql.NullString{String: cmd.arguments[0], Valid: true}
	url := sql.NullString{String: cmd.arguments[1], Valid: true}
	user, err := s.db.GetUser(context.Background(), sql.NullString{String: s.cfg.Username, Valid: true})
	if err != nil {
		return fmt.Errorf("Error, could not find user")
	}

	userUUid := uuid.NullUUID{UUID: user.ID, Valid: true}
	//create user parameters
	currentTime := sql.NullTime{Time: time.Now()}

	feedParams := database.CreateFeedParams{ID: uuid.New(), CreatedAt: currentTime, UpdatedAt: currentTime, Name: name, Url: url, UserID: userUUid}

	feed, err := s.db.CreateFeed(context.Background(), feedParams)
	if err != nil {
		return fmt.Errorf("Error creating user %w", err)
	}

	fmt.Printf("Created Feed:\nname: %s\nurl: %s", feed.Name.String, feed.Url.String)
	return nil

}

func handlerAgg(_ *state, _ command) error {
	urlFeed := "https://www.wagslane.dev/index.xml"
	rssFeed, err := fetchFeed(context.Background(), urlFeed)
	if err != nil {
		return err
	}
	fmt.Println(rssFeed)
	return nil
}

func handlerUsers(s *state, _ command) error {
	users, err := s.db.GetAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error: Failed to get users %w", err)
	}
	for _, user := range users {
		if user.Name.String == s.cfg.Username {

			fmt.Printf(" * %s (current)\n", user.Name.String)
		} else {
			fmt.Printf(" * %s\n", user.Name.String)
		}
	}

	return nil
}

func handlerReset(s *state, _ command) error {
	// implement deleteAllUsers function generated from sqlc
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error: failed to delete users %w", err)
	}

	fmt.Println("Reset database")
	return nil

}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return fmt.Errorf("Error: Username is required")
	}
	name := sql.NullString{String: cmd.arguments[0], Valid: true}
	// check if user exits and return error if not
	user, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("Error: no user with the name: %s", name.String)
	}

	s.cfg.SetUser(user.Name.String)

	fmt.Println("User has been set")
	return nil
}

func handlerRegister(s *state, cmd command) error {

	if len(cmd.arguments) == 0 {
		return fmt.Errorf("Error: Name is required to register new user")
	}

	name := sql.NullString{String: cmd.arguments[0], Valid: true}

	// check if user already exists
	_, err := s.db.GetUser(context.Background(), name)
	if err == nil {
		return fmt.Errorf("Error: user %s already exists", name.String)
	}

	//create user parameters
	currentTime := sql.NullTime{Time: time.Now()}
	userParams := database.CreateUserParams{ID: uuid.New(), CreatedAt: currentTime, UpdatedAt: currentTime, Name: name}

	user, err := s.db.CreateUser(context.Background(), userParams)
	if err != nil {
		return fmt.Errorf("Error creating user %w", err)
	}
	s.cfg.SetUser(user.Name.String)
	// prints for debugging
	fmt.Println("User was created")
	fmt.Printf("Created user info\n Id: %s \n name: %s\n createdAt: %s\n updatedAt: %s\n", user.ID, user.Name.String, user.CreatedAt.Time, user.UpdatedAt.Time)

	return nil
}
func (c *commands) register(name string, f func(s *state, cmd command) error) {
	c.handlers[name] = f
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("Error, %s command does not exist", cmd.name)
	}
	err := f(s, cmd)
	if err != nil {
		return err
	}

	return nil
}
