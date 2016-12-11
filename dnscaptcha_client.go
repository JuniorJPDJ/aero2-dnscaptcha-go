package main

import (
	_ "image/jpeg"
	"github.com/juniorjpdj/aero2-dnscaptcha-go/dnscaptcha"
	"github.com/juniorjpdj/aero2-dnscaptcha-go/utils"
	"time"
	"fmt"
)

const(
	//NS string = "212.2.96.51:53"	// aero2 ns1
	NS string = "212.2.96.52:53"  // aero2 ns2
	api_base_domain string = "io.juniorjpdj.pl"
	api_key string = "ZTI2F"
	captcha_server string = "http://10.2.37.78:8080/"

	COMPRESSION_ENABLED = true
)

var (
	sender = dnscaptcha.CaptchaSender{
		APIKey: api_key,
		BaseDomain: api_base_domain,
		NS: NS,
	}
)

func Log(str string){
	t := time.Now()
	fmt.Printf("[%s][INFO] DNSCaptcha.ClientGO: %s\n", t.Format("2006.01.02 15:04:05"), str)
}

func LogErr(err error){
	Log(fmt.Sprintf("Error occured: %v", err))
}

func main() {
	Log("Started")

	for {
		time.Sleep(2 * time.Second)
		img, sessid := utils.BlockForCaptcha(captcha_server)
		Log("Captcha detected and loaded, uploading")

		cid, err := sender.CaptchaUpload(*img, COMPRESSION_ENABLED)
		if err != nil{
			LogErr(err)
			continue
		}
		Log("Captcha uploaded, waiting for answer")

		status := 0
		for{
			status, err := sender.CaptchaStatusQuery(cid)
			if err != nil || status <= dnscaptcha.StatusN["FAILED"] || status >= dnscaptcha.StatusN["RESOLVED"]{
				break
			}
			time.Sleep(time.Second)
		}
		if err != nil{
			LogErr(err)
			continue
		} else if status <= dnscaptcha.StatusN["FAILED"] {
			Log("Captcha resolve failed")
			continue
		}

		Log("Captcha resolved, loading answer")

		ctxt, err := sender.CaptchaTextQuery(cid)
		if err != nil{
			LogErr(err)
			continue
		}
		Log(fmt.Sprintf("Captcha answer is %s", ctxt))

		cvalid, err := utils.SendCaptchaResponse(captcha_server, ctxt, sessid)
		if err != nil{
			LogErr(err)
			continue
		}

		tmp := ""
		if !cvalid{
			tmp = "not "
		}
		Log(fmt.Sprintf("Answer was %svalid, informing server", tmp))
		err = sender.CaptchaMarkValidQuery(cid, cvalid)
		if err != nil{
			LogErr(err)
			continue
		}
		if cvalid{
			time.Sleep(time.Minute)
		}
	}
}
