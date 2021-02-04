package xmlhelper

import (
	"bufio"
	"bytes"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"os"
)

func DetermineEncodingFromReader(r io.Reader) (e encoding.Encoding, name string, certain bool, err error) {
	b, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		return
	}

	e, name, certain = charset.DetermineEncoding(b, "")
	return
}

func ReadFileWithBadUTF8(path string) ([]byte, error) {
	file1, err := os.Open(path)
	defer file1.Close()
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return nil, err
	}

	e, _, _, err := DetermineEncodingFromReader(file1)
	var data []byte
	if e != nil {
		data, err = ioutil.ReadAll(transform.NewReader(file, e.NewDecoder()))
	} else {
		data, err = ioutil.ReadAll(file)
	}

	illegalUtf8SequencesToPurge := []string{
		"&#x03;", "&#x05;", "&#x00;", "&#x10;", "&#x1E;", "&#x0F;",
		"\u0000", "\u000F", "\u0003", "\u0001"}

	for _, s := range illegalUtf8SequencesToPurge {
		data = bytes.ReplaceAll(data, []byte(s), []byte(""))
	}

	return data, nil
}
