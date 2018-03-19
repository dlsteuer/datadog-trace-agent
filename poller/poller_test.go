package poller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/DataDog/datadog-trace-agent/config"
	"github.com/stretchr/testify/assert"
)

func newTestConfigServer() *httptest.Server {
	server := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			payload := config.ServerConfig{
				ModifyIndex: 1000,
				AnalyzedServices: map[string]float64{
					"web": 1.0,
				},
			}
			bytes, _ := json.Marshal(payload)
			w.WriteHeader(http.StatusOK)
			w.Write(bytes)
		}))

	return server
}

func TestPoller(t *testing.T) {
	assert := assert.New(t)
	server := newTestConfigServer()
	defer server.Close()

	url, err := url.Parse(server.URL)
	done := make(chan struct{})

	assert.NotNil(url)
	assert.Nil(err)

	p := &Poller{
		defaultInterval, "", http.DefaultClient,
		url.String(), make(chan *config.ServerConfig), "apikey_2", 0,
	}

	go func() {
		for {
			select {
			case conf := <-p.updates:
				assert.Equal(conf.ModifyIndex, int64(1000))
				assert.Equal(len(conf.AnalyzedServices), 1)
				assert.Equal(conf.AnalyzedServices["web"], 1.0)

				done <- struct{}{}
			default:
			}
		}
	}()

	err = p.update()
	assert.Nil(err)
	<-done
}