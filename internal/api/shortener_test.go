package api

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
	"urlShortener/internal/db/mock"
)


const (
	host          = "http://127.0.0.1:1323/"
	validURL = "http://127.0.0.1/testpath"
	serverFailURL = "http://127.0.0.1/failme"
)

func TestMain(m *testing.M) {
	setUp()
	retCode := m.Run()
	os.Exit(retCode)
}

func setUp() {
	router := NewRouter(
		NewShortener(&mock.RepositoryMock{
			InsertFunc: func(key string, url string) error {
				if url == serverFailURL {
					return errors.New("emulate  service unavailable")
				}
				return nil
			},
			GetFunc: func(key string) (string, error) {
				return validURL, nil
			},
		},),
	)

	router.echo.TRACE("/", func(context echo.Context) error {
		return nil
	})

	go router.Start()

	awaitServerStart(5*time.Second)
}

func TestShortener_ShortenURL(t *testing.T) {
	testCases := []struct {
		url            string
		expectedStatus int
		expectedLocation string
	}{
		{
			url:            validURL,
			expectedStatus: http.StatusOK,
		},
		{
			url:            validURL,
			expectedStatus: http.StatusOK,
		},
		{
			url:            "127.0.0.1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			url:            serverFailURL,
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	// client ignoring redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, tc := range testCases {
		req, err := http.NewRequest(
			echo.POST,
			fmt.Sprintf("%sshorten?url=%s", host, tc.url),
			strings.NewReader(``),
		)

		resp, err := client.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)


		require.Equal(t, tc.expectedStatus, resp.StatusCode)
		if resp.StatusCode == http.StatusOK {
			expectedCode := encode(tc.url)
			require.Equal(t, expectedCode, string(body))

			req, err = http.NewRequest(
				echo.GET,
				fmt.Sprintf("%s%s", host, expectedCode),
				strings.NewReader(``),
			)

			resp, err = client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusMovedPermanently, resp.StatusCode)
			require.Equal(t, tc.url, resp.Header.Get("Location"))
		}
	}
}

func awaitServerStart(timeout time.Duration)  {
	fmt.Printf("Awaiting for server to start in %s", timeout)

	req, _ := http.NewRequest(
		echo.TRACE,
		host,
		strings.NewReader(``),
	)

	start := time.Now()
	for {
		resp, _ := http.DefaultClient.Do(req)

		if resp.StatusCode == http.StatusOK {
			break
		}

		if time.Now().Sub(start) > timeout {
			panic(fmt.Sprintf("Failed to start echo server for %s", timeout))
		}

		time.Sleep(1 * time.Second)
	}
}
