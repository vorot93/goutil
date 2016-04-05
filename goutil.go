package goutil

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

var ErrMismatch = errors.New("Mismatch")
var ErrPanic = errors.New("Procedure panic")

func ErrorOut(err error, expectation, result interface{}) string {
	return fmt.Sprintf("Error: %s\nExpected %#v\nReceived %#v\n", err.Error(), expectation, result)
}

func ErrorOutJSON(err error, expectation, result interface{}) string {
	expJSON, _ := json.Marshal(expectation)
	resJSON, _ := json.Marshal(result)
	return fmt.Sprintf("Error: %s\nExpected %s\nReceived %s\n", err.Error(), expJSON, resJSON)
}

func SprintfCompare(expectation, result interface{}) bool {
	return fmt.Sprintf("%#v", expectation) == fmt.Sprintf("%#v", result)
}

func JSONcompare(expectation, result interface{}) bool {
	var e, _ = json.Marshal(expectation)
	var r, _ = json.Marshal(result)

	if string(e) != string(r) {
		return false
	}

	return true
}

func DownloadURL(urlData url.URL) (data string, err error) { return download(urlData.String(), true) }

func DownloadURLNoRedirect(urlData url.URL) (data string, err error) {
	return download(urlData.String(), false)
}

// Download retrieves data from the specified HTTP address.
func Download(url string) (data string, err error) { return download(url, true) }

func download(url string, allowRedirect bool) (data string, err error) {
	var redirectHandler func(req *http.Request, via []*http.Request) error
	if !allowRedirect {
		redirectHandler = func(req *http.Request, via []*http.Request) error {
			return errors.New("Redirects are not allowed.")
		}
	}
	var dlClient = http.Client{CheckRedirect: redirectHandler}

	var resp *http.Response
	resp, err = dlClient.Get(url)
	if err == nil {
		defer resp.Body.Close()
		var status = resp.StatusCode
		if status == 200 {
			var readallContents, _ = ioutil.ReadAll(resp.Body)
			data = string(readallContents)
		} else {
			err = errors.New(fmt.Sprintf("Received status code: %d", status))
		}

	}
	return data, err
}

func IntP(i int) *int {
	return &i
}

func GetIntP(data string) (p *int) {
	var t, vErr = strconv.ParseInt(data, 10, 64)
	if vErr == nil {
		var v = int(t)
		p = &v
	}

	return p
}

func GetInt64P(data string) (p *int64) {
	var v, vErr = strconv.ParseInt(data, 10, 64)
	if vErr == nil {
		p = &v
	}

	return p
}

func GetFloatP(data string) (p *float64) {
	var v, vErr = strconv.ParseFloat(data, 64)
	if vErr == nil {
		p = &v
	}

	return p
}

func FieldsToMap(header []string, records [][]string) []map[string]string {
	var outMap = make([]map[string]string, 0, len(records))

	var hLen = len(header)

	for _, record := range records {
		var fieldNum int

		var recLen = len(record)

		if hLen > recLen {
			fieldNum = recLen
		} else {
			fieldNum = hLen
		}
		var newMap = make(map[string]string, fieldNum)
		for i := 0; i < fieldNum; i++ {
			var k = header[i]
			if k != "" {
				newMap[k] = record[i]
			}
		}
		outMap = append(outMap, newMap)
	}

	return outMap
}

type Document map[string]interface{}

func (self *Document) ToValue(factory func() interface{}) (interface{}, error) {
	var s, merr = json.Marshal(self)
	if merr != nil {
		return nil, merr
	}

	var v = factory()
	var umerr = json.Unmarshal(s, v)
	if umerr != nil {
		return nil, umerr
	}

	return v, nil
}
