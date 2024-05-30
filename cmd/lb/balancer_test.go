package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	suite := new(TestSuite)
	suite.SetupSuite()
	t.Run("TestBalancer", suite.TestBalancer)
	t.Run("TestHealth", suite.TestHealth)
}

type TestSuite struct {
	serversPool []string
}

func (s *TestSuite) SetupSuite() {
	s.serversPool = []string{
		"server1:8080",
		"server2:80",
		"server3:80",
	}
}

func (s *TestSuite) TestBalancer(t *testing.T) {
	// TODO: Implement the unit test for the balancer.
	address1 := getIndex("127.0.0.1:8080")
	address2 := getIndex("192.168.0.0:80")
	address3 := getIndex("26.143.218.9:80")

	assert.Equal(t, 2, address1)
	assert.Equal(t, 0, address2)
	assert.Equal(t, 1, address3)
}

func (s *TestSuite) TestHealth(t *testing.T) {
	result := make([]string, len(s.serversPool))

	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server2.Close()

	parsedURL1, _ := url.Parse(server1.URL)
	hostURL1 := parsedURL1.Host

	parsedURL2, _ := url.Parse(server2.URL)
	hostURL2 := parsedURL2.Host

	servers := []string{
		hostURL1,
		hostURL2,
		"server3:8080",
	}

	healthCheck(servers, result)
	time.Sleep(12 * time.Second)

	assert.Equal(t, hostURL1, result[0])
	assert.Equal(t, hostURL2, result[1])
	assert.Equal(t, "", result[2])
}
