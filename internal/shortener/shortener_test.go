package shortener_test

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
	"urlShortener/internal/api"
	"urlShortener/internal/db/mock"
	"urlShortener/internal/router"
	"urlShortener/internal/shortener"
)

const (
	host          = "http://127.0.0.1:1323/"
	validURL = "http://hostname/testpath"
	serverFailURL = "http://hostname/failme"
)

func TestShortener_ShortenURL(t *testing.T) {
	testCases := []struct {
		testName		 string
		url              string
		expectedStatus   int
		expectedCode     string
		expectedLocation string
		repository       api.Repository
	}{
		{
			testName:       "regular valid URL is encoded, return 200",
			url:            validURL,
			expectedStatus: http.StatusOK,
			expectedCode:   shortener.Encode(validURL),
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
			testName: 		"Invalid URL without protocol, return 400",
			url:            "127.0.0.1",
			expectedStatus: http.StatusBadRequest,
		},
		{
			testName:       "DB is unavailable return 503",
			url:            serverFailURL,
			expectedStatus: http.StatusServiceUnavailable,
			repository: &mock.RepositoryMock{
				InsertFunc: func(key string, url string) error {
					if url == serverFailURL {
						return errors.New("emulate api.unavailable")
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
		t.Run(tc.testName, func(t *testing.T) {
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
			require.NoError(t, err)
		},
		)
	}
}

func configureRouter(repository api.Repository) func() error {
	router := router.NewRouter(
		shortener.NewShortener(repository),
	)

	router.Trace("/", func(context echo.Context) error {
		return nil
	})

	routerc := make(chan error, 1)
	go router.Start(routerc)
	awaitServerStart(5 * time.Second)

	return router.Close
}

func awaitServerStart(timeout time.Duration) {
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