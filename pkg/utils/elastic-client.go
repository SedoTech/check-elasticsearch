package utils

import (
	"fmt"
	"github.com/olivere/elastic"
)

const (
	Elk01	= "http://elk01.i.sedorz.net:9200"
	Elk02	= "http://elk02.i.sedorz.net:9200"
	Elk03	= "http://elk03.i.sedorz.net:9200"
)

// NewResult creates a new instance of Result
func NewElasticClient(url string) (*elastic.Client, error) {

	if len(url) == 0 {
		url = Elk01
	}
	client, err := elastic.NewClient(elastic.SetURL(url))

	if err != nil {
		return nil, fmt.Errorf("could not get ElasticSearch client: %s", err)
	}

	return client, err
}
