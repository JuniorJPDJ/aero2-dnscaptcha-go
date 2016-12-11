package utils

import (
	"net/url"
	"image"
	"time"
	"strings"
)

func LoadCaptcha(server string)(*image.Image, string){
	form := url.Values{}
	form.Set("viewForm", "true")
	buf, err := HttpPost(server, form)
	if err != nil{
		return nil, ""
	}

	sessid := StringBetween((*buf).String(), "name=\"PHPSESSID\" value=\"", "\"")

	form = url.Values{}
	form.Set("PHPSESSID", sessid)
	buf, err = HttpGet(server + "getCaptcha.html?" + form.Encode())
	if err != nil{
		return nil, ""
	}

	img, _, err := image.Decode(buf)
	if err != nil{
		return nil, ""
	}

	return &img, sessid
}

func BlockForCaptcha(server string)(*image.Image, string){
	for {
		img, sessid := LoadCaptcha(server)
		if img != nil{
			return img, sessid
		}
		time.Sleep(time.Second)
	}
}

func SendCaptchaResponse(server, response, sessid string) (bool, error){
	form := url.Values{}
	form.Set("PHPSESSID", sessid)
	form.Set("viewForm", "true")
	form.Set("captcha", response)

	buf, err := HttpPost(server, form)
	if err != nil{
		return false, err
	}

	return !strings.Contains(buf.String(), "getCaptcha"), nil
}
