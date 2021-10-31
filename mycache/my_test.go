package mycache

import (
	"fmt"
	"testing"
	"time"
)

// Tweet is a structure used for serializing/deserializing data in Elasticsearch.
type Tweet struct {
	User     string    `json:"user"`
	Message  string    `json:"message"`
	Retweets int       `json:"retweets"`
	Image    string    `json:"image,omitempty"`
	Created  time.Time `json:"created,omitempty"`
	Tags     []string  `json:"tags,omitempty"`
	Location string    `json:"location,omitempty"`
}

func Test_cache(t *testing.T) {
	//获取
	cache := NewCache("_cache")
	tp := Tweet{
		User:     "aaaa",
		Retweets: 12,
	}
	cache.Add("key", tp, 24*time.Hour)

	var tmp Tweet
	err := cache.Value("key", &tmp)
	if err != nil {
		fmt.Println(tmp)
	} else {
		fmt.Println(tmp)
	}

	return
}
