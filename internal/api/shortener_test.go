package api

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
	"urlShortener/internal/db"
	"urlShortener/internal/db/mock"
)


const (
	host          = "http://127.0.0.1:1323/"
	validURL = "http://hostname/testpath"
	serverFailURL = "http://hostname/failme"
)

func TestShortener_ShortenURL(t *testing.T) {
	testCases := []struct {
		url              string
		expectedStatus   int
		expectedCode     string
		expectedLocation string
		repository       db.Repository
	}{
		{
			url:            validURL,
			expectedStatus: http.StatusOK,
			expectedCode: encode(validURL),
			repository: &mock.RepositoryMock{
				InsertFunc: func(key string, url string) error {
					return nil
				},
				GetFunc: func(key string) (string, error) {
					return validURL, nil
				},
			},
		},
		{
			url:            validURL,
			expectedStatus: http.StatusOK,
			expectedCode: encode(validURL),
			repository: &mock.RepositoryMock{
				InsertFunc: func(key string, url string) error {
					return nil
				},
				GetFunc: func(key string) (string, error) {
					return validURL, nil
				},
			},
		},
		{
			url:            "127.0.0.1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			url:            serverFailURL,
			expectedStatus: http.StatusServiceUnavailable,
			repository: &mock.RepositoryMock{
				InsertFunc: func(key string, url string) error {
					if url == serverFailURL {
						return errors.New("emulate db unavailable")
					}
					return nil
				},
				GetFunc: func(key string) (string, error) {
					return "", nil
				},
			},
		},
	}

	// client ignoring redirects
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	for _, tc := range testCases {
		closer := configureRouter(tc.repository)

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
			require.Equal(t, tc.expectedCode, string(body))

			req, err = http.NewRequest(
				echo.GET,
				fmt.Sprintf("%s%s", host, tc.expectedCode),
				strings.NewReader(``),
			)

			resp, err = client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, http.StatusMovedPermanently, resp.StatusCode)
			require.Equal(t, tc.url, resp.Header.Get("Location"))
		}

		err = closer()
		if err != nil {
			continue
		}
	}
}

func configureRouter(repository db.Repository) func() error {
	router := NewRouter(
		NewShortener(repository),
	)

	router.echo.TRACE("/", func(context echo.Context) error {
		return nil
	})

	go router.Start()
	awaitServerStart(5*time.Second)

	return router.Close
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
