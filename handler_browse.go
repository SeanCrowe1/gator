package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/seancrowe1/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	if len(cmd.Args) > 1 {
		return fmt.Errorf("usage: %s <optional: limit>", cmd.Name)
	}

	var err error

	var limit int
	if len(cmd.Args) == 1 {
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return err
		}
	} else {
		limit = 2
	}

	posts, err := s.db.GetPostsforUser(context.Background(), database.GetPostsforUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		fmt.Printf("Title: %v\nPublished At: %v\n", post.Title, post.PublishedAt)
		fmt.Printf("Description: %v\n\n", post.Description)
	}

	return nil
}
