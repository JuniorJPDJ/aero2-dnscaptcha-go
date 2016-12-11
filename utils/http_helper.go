package utils

import (
	"net/http"
	"net/url"
	"time"
	"bytes"
)

var (
	httpcl http.Client = http.Client{
		Timeout: time.Duration(3 * time.Second),
	}
)

func HttpGet(url string)(*bytes.Buffer, error){
	resp, err := httpcl.Get(url)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(resp.Body)
	if err != nil{
		return nil, err
	}
	return &buf, nil
}

func HttpPost(url string, form url.Values)(*bytes.Buffer, error){
	resp, err := httpcl.PostForm(url, form)
	if err != nil{
		return nil, err
	}
	defer resp.Body.Close()
	buf := bytes.Buffer{}
	_, err = buf.ReadFrom(resp.Body)
	if err != nil{
		return nil, err
	}
	return &buf, nil
}
