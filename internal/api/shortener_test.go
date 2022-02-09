package api

import (
	"crypto/sha256"
	"encoding/hex"
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


const address = "http://127.0.0.1:1323/"

func TestMain(m *testing.M) {
	setUp()

	retCode := m.Run()

	os.Exit(retCode)
}

var e *echo.Echo
var sh *Shortener

func setUp() {
	e = echo.New()

	sh = NewShortener(address, &mock.RepositoryMock{
		InsertFunc: func(key []byte, url string) error {
			return nil
		},
	})

	e.TRACE("/", func(context echo.Context) error {
		return nil
	})
	e.POST("/shorten", sh.ShortenUrl)
	go e.Start(":1323")

	awaitServerStart(5*time.Second)
}

func TestShortener_ShortenUrl(t *testing.T) {
	url := "https://127.0.0.1/testpath"
	hash := sha256.Sum256([]byte(url))
	req, err := http.NewRequest(
		echo.POST,
		fmt.Sprintf("%sshorten?url=%s", address, url),
		strings.NewReader(``),
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.Equal(t, address + hex.EncodeToString(hash[:4]), string(b))

}

func awaitServerStart(timeout time.Duration)  {
	fmt.Printf("Awaiting for server to start in %s", timeout)

	req, _ := http.NewRequest(
		echo.TRACE,
		address,
		strings.NewReader(``),
	)

	start := time.Now()
	for {
		resp, _ := http.DefaultClient.Do(req)

		if resp.StatusCode == http.StatusOK {
			break
		} else {
			if time.Now().Sub(start) > timeout {
				panic(fmt.Sprintf("Failed to start echo server for %s", timeout))
			}
		}

		time.Sleep(1 * time.Second)
	}
}
