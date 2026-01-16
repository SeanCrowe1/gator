package main

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/seancrowe1/gator/internal/database"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]

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

func handlerFollowing(s *state, cmd command, user database.User) error {
	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
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

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <url>", cmd.Name)
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return err
	}

	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})

	if err != nil {
		return err
	}

	fmt.Printf("Feed unfollowed successfully!")
	return nil
}

func printFeedFollow(ff database.CreateFeedFollowRow) {
	fmt.Printf("\n * User Name: %v\n", ff.UserName)
	fmt.Printf(" * Feed Name: %v\n\n", ff.FeedName)
	fmt.Println("=====================================")
}
