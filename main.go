package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type byteMap struct {
	bytes     []byte
	byteRange byteRange
}

type byteRange struct {
	start int
	end   int
}

type File struct {
	Name             string
	URL              string
	ContentLength    int
	AcceptByteRanges bool
	Bytes            []byteMap
}

func (f *File) headRequest() error {
	headers, err := http.Head(f.URL)
	if err != nil {
		return errors.Wrap(err, "Head request failed")
	}
	f.ContentLength = int(headers.ContentLength)
	f.AcceptByteRanges = headers.Header.Get("Accept-Ranges") == "bytes"

	return nil
}

// DownloadWithRange downloads a file with a range
func (f *File) DownloadWithRange() error {
	length := f.ContentLength
	chunks := 10
	chunkSize := length / chunks

	for i := 0; i < chunks; i++ {
		start := i * chunkSize
		end := start + chunkSize - 1
		if i == chunks-1 {
			end = length - 1
		}
		bytes, err := f.RequestRange(start, end)
		if err != nil {
			return errors.Wrap(err, "Requesting part of file failed")
		}

		f.Bytes = append(f.Bytes, byteMap{bytes, byteRange{start, end}})
	}
	return nil
}

// RequestRange requests a range of bytes from a url
func (f *File) RequestRange(start int, end int) ([]byte, error) {
	req, err := http.NewRequest("GET", f.URL, nil)
	if err != nil {
		return nil, errors.Wrap(err, "NewRequest failed")
	}
	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "Do request failed")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return nil, errors.New("Request failed")
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "Read response body failed")
	}
	f.Bytes = append(f.Bytes, byteMap{bytes, byteRange{start, end}})

	return bytes, nil
}

// Download downloads a file
func (f *File) Download() error {
	resp, err := http.Get(f.URL)
	if err != nil {
		return errors.Wrap(err, "Get request failed")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "Read response body failed")
	}
	f.Bytes = append(f.Bytes, byteMap{bytes, byteRange{0, len(bytes) - 1}})

	return nil
}

func (file *File) StartDownload() error {
	if file.AcceptByteRanges {
		return file.DownloadWithRange()
	}
	return file.Download()
}

// String converts bytes of a file to a string
func (f File) String() string {
	bytes := make([]byte, f.ContentLength)
	for _, b := range f.Bytes {
		bytes = append(bytes, b.bytes...)
	}
	return string(bytes)
}

func NewFile(url string) (file File, err error) {
	file.URL = url

	err = file.headRequest()
	if err != nil {
		return file, err
	}
	return file, nil
}

func main() {

	url := "https://sample-videos.com/video123/mp4/720/big_buck_bunny_720p_30mb.mp4"

	// time calc
	start2 := time.Now()
	ok2, err := down(url, false)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ok2)
	fmt.Println(time.Since(start2))

	// time calc
	start := time.Now()
	ok, err := down(url, true)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ok)
	fmt.Println(time.Since(start))

}

func down(url string, acceptByteRanges bool) (bool, error) {
	file, err := NewFile(url)
	fmt.Println(file)
	if err != nil {
		return false, err
	}

	if !acceptByteRanges {
		file.AcceptByteRanges = false
	}

	err = file.StartDownload()
	if err != nil {
		return false, err
	}
	if file.Bytes[len(file.Bytes)-1].byteRange.end+1 == file.ContentLength {
		return true, nil
	}
	return false, fmt.Errorf("size mismatch")
}
