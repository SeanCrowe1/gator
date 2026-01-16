package main

import (
	"context"
	"fmt"
)

func handlerAgg(s *state, cmd command) error {
	feed, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return fmt.Errorf("couldn't fetch feed: %w", err)
	}

	printFeed(feed)
	return nil
}

func printFeed(f *RSSFeed) {
	fmt.Printf(" * Title: 		%v\n", f.Channel.Title)
	fmt.Printf(" * Link: 		%v\n", f.Channel.Link)
	fmt.Printf(" * Description: %v\n\n", f.Channel.Description)
	for _, item := range f.Channel.Item {
		fmt.Printf(" ** Title: 			%v\n", item.Title)
		fmt.Printf(" ** Link: 			%v\n", item.Link)
		fmt.Printf(" ** Description: 	\n%v\n", item.Description)
		fmt.Printf(" ** PubDate: 		%v\n\n", item.PubDate)
	}
}
