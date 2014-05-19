// redditnews package implements a basic client for the Reddit API.
package redditnews

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

// Item describes a RedditNews item.
type Item struct {
	Author string `json:"author"`
	Score  int    `json:"score"`
	URL    string `json:"url"`
	Title  string `json:"title"`
}

type response struct {
        Data1 struct {
                Children []struct {
                        Data2 Item `json:"data"`
                } `json:"children"`
        } `json:"data"`
}

// Get fetches the most recent Items posted to the specified subreddit.
func Get(reddit string) ([]Item, error) {
        url := fmt.Sprintf("http://reddit.com/r/%s.json", reddit)
	resp, err := http.Get(url)
	if err != nil {
	        return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
	        return nil, errors.New(resp.Status)
	}

	r := new(response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
	        return nil, err
	}

	items := make([]Item, len(r.Data1.Children))
	for i, child := range r.Data1.Children {
	        items[i] = child.Data2
	}
	return items, nil
}

func (i Item) String() string {
	return fmt.Sprintf(
		"Author: %s\nScore: %d\nURL: %s\nTitle: %s\n\n",
		i.Author,
		i.Score,
		i.URL,
		i.Title)
}

// Email prepares the body of an email
func Email() string {
	var buffer bytes.Buffer

	items, err := Get("golang")
	if err != nil {
	        log.Fatal(err)
	}

	// Need to build strings from items
	for _, item := range items {
	        buffer.WriteString(item.String())
	}

	return buffer.String()
}

