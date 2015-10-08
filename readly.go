package readly

import(
  "io"
  "io/ioutil"
  "net/http"
  "os"
  "regexp"
)

var httpClient *http.Client;

// Credit to http://play.golang.org/p/ps7PXFRI2B
type nopReader struct{}
func (nopReader) Read(p []byte) (n int, err error) { return len(p), nil }

func isHttp(file string) bool {
  http := regexp.MustCompile("^http(s)?:")
  return http.MatchString(file)
}

func SetClient(c *http.Client) {
  httpClient = c
}

func client() *http.Client {
  var nilClient *http.Client
  if httpClient == nilClient {
    httpClient = &http.Client{}
  }
  return httpClient
}

// TODO: Error wrapping
// TODO: timeout?
func Reader(file string) (io.ReadCloser, error) {
  nilReaderCloser := ioutil.NopCloser(&nopReader{})
  if isHttp(file) {
    resp, err := client().Get(file)
    if err != nil {
      return nilReaderCloser, err
    }
    return resp.Body, nil
  } else {
    file, err := os.Open(file)
    if err != nil {
      return nilReaderCloser, err
    }
    return file, nil
  }
}

func Read(file string) (string, error) {
  r, err := Reader(file)
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
