package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/seancrowe1/gator/internal/database"
)

func handlerFollow(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})

	fmt.Printf("Feed followed succesfully!")
	printFeedFollow(feedFollow)

	return nil
}

func handlerFollowing(s *state, cmd command) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return err
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds followed.")
		return nil
	}

	fmt.Printf("Currently following: \n\n")

	for _, feed := range feeds {
		fmt.Printf(" * %v\n", feed.FeedName)
	}

	return nil
}

func printFeedFollow(ff database.CreateFeedFollowRow) {
	fmt.Printf("\n * User Name: %v\n", ff.UserName)
	fmt.Printf(" * Feed Name: %v\n\n", ff.FeedName)
	fmt.Println("=====================================")
}
