package loader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

type MemStatsLoader struct {
	url    string
	client *http.Client
}

func NewMemStatsLoader(url string) *MemStatsLoader {
	return &MemStatsLoader{
		url: url,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}
func (m *MemStatsLoader) Load() (*runtime.MemStats, error) {
	res, err := m.client.Get(m.url)
	if err != nil {
		return nil, fmt.Errorf("connect err %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status %w", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read err %w", err)
	}

	result := &struct {
		Stats *runtime.MemStats `json:"memstats"`
	}{}
	if err := json.Unmarshal(body, result); err != nil {
		return nil, fmt.Errorf("fetch memstat, json  err %w", err)
	}

	return result.Stats, nil

}
