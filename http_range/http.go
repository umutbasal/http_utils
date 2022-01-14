package http_range

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// Client proxies http.Client
type Client struct {
	*http.Client
}

// getInMultiParts downloads a file with a range
func getInMultiParts(url string, contentLength int, chunks int) (*http.Response, error) {
	var responses []*http.Response

	chunkSize := contentLength / chunks

	for i := 0; i < chunks; i++ {
		start := i * chunkSize
		end := start + chunkSize - 1
		if i == chunks-1 {
			end = contentLength - 1
		}
		resp, err := rangeRequest(url, start, end)
		if err != nil {
			return nil, errors.Wrap(err, "Requesting part of file failed")
		}

		if resp.StatusCode != http.StatusPartialContent {
			return nil, errors.New("Range not supported on server")
		}

		responses = append(responses, resp)
	}

	if responses == nil {
		return nil, errors.New("No bytes received")
	}

	point := func(p1, p2 *http.Response) bool {
		point1, err := parseRange(p1.Request.Header.Get("Range"), int64(contentLength))
		if err != nil {
			return false
		}
		point2, err := parseRange(p2.Request.Header.Get("Range"), int64(contentLength))
		if err != nil {
			return false
		}

		return point1[0].start < point2[0].start
	}

	By(point).Sort(responses)

	var body []byte
	for _, r := range responses {
		bytes, err := io.ReadAll(r.Body)

		if err != nil {
			return nil, errors.Wrap(err, "Read response body failed")
		}

		defer r.Body.Close()

		body = append(body, bytes...)
	}

	// Fake a response
	ret := responses[0]
	ret.StatusCode = http.StatusOK
	ret.Status = http.StatusText(http.StatusOK)
	ret.Body = ioutil.NopCloser(bytes.NewReader(body))
	ret.ContentLength = int64(len(body))
	ret.Request.Header.Set("Range", "bytes=0-")

	return ret, nil
}

// rangeRequest requests a range of bytes from a url
func rangeRequest(url string, start int, end int) (*http.Response, error) {
	log.Printf("Requesting range %d-%d", start, end)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest failed")
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Do request failed")
	}

	return resp, nil
}

// GetWithByteRange takes default get arguments and threads
func (c *Client) GetWithByteRange(url string, chunks int) (*http.Response, error) {
	resp, err := http.Head(url)
	if err != nil {
		return nil, errors.Wrap(err, "Head request failed")
	}
	contentLength := int(resp.ContentLength)
	acceptbyteRanges := resp.Header.Get("Accept-Ranges") == "bytes"

	if acceptbyteRanges {
		return getInMultiParts(url, contentLength, chunks)
	}
	return c.Get(url)
}
