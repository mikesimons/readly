package readly

//go:generate curl google.com > test.html
import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

// Reader holds state for Readly
type Reader struct {
	Client *http.Client
}

// Credit to http://play.golang.org/p/ps7PXFRI2B
type nopReader struct{}

func (nopReader) Read(p []byte) (n int, err error) { return len(p), nil }

func isHTTP(file string) bool {
	http := regexp.MustCompile("^http(s)?:")
	return http.MatchString(file)
}

// New will return an instance of Readly
func New() *Reader {
	return &Reader{
		Client: &http.Client{},
	}
}

// Reader will return an io.ReadCloser for either a file or http(s) URL
func (readly *Reader) Reader(file string) (io.ReadCloser, error) {
	// TODO: Error wrapping
	// TODO: timeout?
	nilReaderCloser := ioutil.NopCloser(&nopReader{})
	if isHTTP(file) {
		resp, err := readly.Client.Get(file)
		if err != nil {
			return nilReaderCloser, err
		}
		return resp.Body, nil
	}

	fileHandle, err := os.Open(file)
	if err != nil {
		return nilReaderCloser, err
	}
	return fileHandle, nil
}

// Read will take a file or http(s) URL, read all contents and return them
func (readly *Reader) Read(file string) (string, error) {
	r, err := readly.Reader(file)
	defer r.Close()

	if err != nil {
		return "", err
	}

	ret, err := ioutil.ReadAll(r)

	if err != nil {
		return "", err
	}

	return string(ret), nil
}
