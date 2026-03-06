package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/swampbear/blog-aggregator/internal/rss"
)

func fetchFeed(ctx context.Context, feedURL string) (*rss.RSSFeed, error) {
	body := strings.NewReader("")
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, body)
	if err != nil {
		return &rss.RSSFeed{}, fmt.Errorf("Error, failed to create request with context, %w", err)
	}
	req.Header.Set("User-Agent", "gator")
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return &rss.RSSFeed{}, fmt.Errorf("Error, failed to do request: %w", err)
	}
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &rss.RSSFeed{}, fmt.Errorf("Error: Failed to read response body: %w", err)
	}
	rssFeed := rss.RSSFeed{}
	if err = xml.Unmarshal(data, &rssFeed); err != nil {
		return &rss.RSSFeed{}, fmt.Errorf("Error: failed to unmarshal rss feed: %w", err)
	}
	rssFeed.CleanRSSFeed()

	return &rssFeed, nil
}
