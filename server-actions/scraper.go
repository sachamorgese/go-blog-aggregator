package server_actions

import (
	"BlogAggregator/internal/database"
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"github.com/google/uuid"
	"io"
	"net/http"
	"sync"
	"time"
)

type RSS struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		LastBuild   string `xml:"lastBuildDate"`
		Items       []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

var possibleLayouts = []string{
	time.RFC822,                       // Tue, 20 Dec 2023 12:05:03 +0100
	time.RFC3339,                      // 2023-12-20T12:05:03+01:00
	"Mon, 02 Jan 2006 15:04:05 -0700", // Add more if needed
}

func Fetcher(dbQueries *database.Queries, feed_id uuid.UUID, url string) ([]byte, error) {
	// Send an HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the HTTP response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rss := RSS{}
	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	for _, item := range rss.Channel.Items {
		var timeObj time.Time
		var err error

		for _, layout := range possibleLayouts {
			timeObj, err = time.Parse(layout, item.PubDate)
			if err == nil {
				break // We found a matching layout!
			}
		}

		if err != nil {
			fmt.Println("Error parsing time:", err)
		} else {
			fmt.Println(timeObj)
		}

		_, err = dbQueries.CreatePost(ctx, database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title[:min(len(item.Title), 255)],
			Url:         item.Link[:min(len(item.Link), 255)],
			Description: item.Description[:min(len(item.Description), 255)],
			PublishedAt: timeObj,
			FeedID:      feed_id,
		})

		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return body, nil
}

func FetchFeeds(db *sql.DB, n int32) {
	dbQueries := database.New(db)
	ctx := context.Background()

	feedsToFetch, err := dbQueries.GetNextFeedsToFetch(ctx, n)

	if err != nil {
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup

	for _, feed := range feedsToFetch {
		wg.Add(1)

		go func(feed database.Feed) {
			defer wg.Done()
			_, err := Fetcher(dbQueries, feed.ID, feed.Url)

			if err != nil {
				fmt.Println(err)
				return
			}

			currentTime := time.Now()

			_ = dbQueries.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
				ID:            feed.ID,
				LastFetchedAt: sql.NullTime{Time: currentTime, Valid: true},
				UpdatedAt:     currentTime,
			})

			if err != nil {
				fmt.Println(err)
				return
			}
		}(feed)
	}
}

func StartFetchWorker(db *sql.DB, interval time.Duration, n int32) {
	ticker := time.NewTicker(interval * time.Second)

	FetchFeeds(db, n)

	for {
		select {
		case <-ticker.C:
			FetchFeeds(db, n)
		}
	}
}
