package dnscaptcha

import (
	"github.com/juniorjpdj/aero2-dnscaptcha-go/utils"
	"github.com/miekg/dns"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"image/color"
	"image"
)

const (
	MAX_DOMAIN_NAME_LEN = 253
)

var (
	Status map[int]string = map[int]string{
		-1: "FAILED",
		0: "CREATED",
		1: "UPLOADING",
		2: "UPLOADED",
		3: "RESOLVING",
		4: "RESOLVED",
		5: "VERIFIED",
	}

	StatusN map[string]int = map[string]int{
		"FAILED": -1,
		"CREATED": 0,
		"UPLOADING": 1,
		"UPLOADED": 2,
		"RESOLVING": 3,
		"RESOLVED": 4,
		"VERIFIED": 5,
	}

	dnscl dns.Client = dns.Client{}
)

func QueryNSAPI(query, base_domain, ns string)(string, error){
	m := dns.Msg{}
	m.SetQuestion(query + "." + base_domain + ".", dns.TypeCNAME)
	m.RecursionDesired = true

	resp, _, err := dnscl.Exchange(&m, ns)
	if err != nil{
		return "", err
	}

	if len(resp.Answer) != 1 {
		return "", errors.New("dnscaptcha: server error")
	}

	info := resp.Answer[0].(*dns.CNAME).Target

	if info == "error." {
		return info, errors.New("dnscaptcha: server returned error")
	}

	return info, nil
}

type CaptchaSender struct {
	APIKey     string
	BaseDomain string
	NS         string
}

func (s *CaptchaSender) NewCaptchaQuery(parts int) (string, error){
	query := fmt.Sprintf("%d.%s.new", parts, s.APIKey)
	resp, err := QueryNSAPI(query, s.BaseDomain, s.NS)
	if err != nil{
		return "", err
	}
	return resp[:len(resp) - 1], nil
}

func (s *CaptchaSender) CaptchaUpload(img image.Image, compress bool) (string, error){
	if compress{
		img = utils.ConvertImgPalette(img, color.Palette{color.Black, color.White})
	}
	data, err := utils.Base32Image(img, true)
	if err != nil{
		fmt.Println(err)
	}

	remaining := len(data)
	part_nr := -1
	parts := [][]string{}
	available := 0
	data_len := 0
	for remaining > 0{
		part_nr += 1
		parts = append(parts, []string{})
		// available = MAX_DOMAIN_NAME_LEN - 1(dot) - len(domain) - 1(dot) - 2(opid="ul") - 1(dot) - len(api_key) - 1(dot) - 5(captchaid) - len(strconv.Itoa(part_nr))
		available = MAX_DOMAIN_NAME_LEN - 11 - len(s.BaseDomain) - len(s.APIKey) - len(strconv.Itoa(part_nr))
		for available > 0 && remaining > 0{
			if available > 64{
				data_len = 63
				available -= 64
			} else {
				data_len = available - 1
				available = 0
			}
			if remaining <= data_len{
				data_len = remaining
				remaining = 0
			} else {
				remaining -= data_len
			}
			tmp := data[:data_len]
			parts[part_nr] = append(parts[part_nr], tmp)
			data = data[data_len:]
		}
	}

	captcha_id, err := s.NewCaptchaQuery(part_nr + 1)
	if err != nil{
		return "", err
	}

	for part_nr, part := range parts{
		query := fmt.Sprintf("%s.%d.%s.%s.ul", strings.Join(part, "."), part_nr, captcha_id, s.APIKey)
		resp, err := QueryNSAPI(query, s.BaseDomain, s.NS)
		if err != nil{
			return "", err
		} else if resp != "ok."{
			return "", errors.New(fmt.Sprintf("ul: failed on part #%d", part_nr))
		}
	}

	return captcha_id, nil
}

func (s *CaptchaSender) CaptchaStatusQuery(captcha_id string) (int, error) {
	resp, err := QueryNSAPI(fmt.Sprintf("%s.%s.st", captcha_id, s.APIKey), s.BaseDomain, s.NS)
	if err != nil{
		return -1, err
	}
	return StatusN[resp[:len(resp) - 1]], nil
}

func (s *CaptchaSender) CaptchaTextQuery(captcha_id string) (string, error) {
	resp, err := QueryNSAPI(fmt.Sprintf("%s.%s.get", captcha_id, s.APIKey), s.BaseDomain, s.NS)
	if err != nil{
		return "", err
	}
	if !strings.HasSuffix(resp, ".ok."){
		return "", errors.New("get: " + resp)
	}
	resp = resp[:len(resp) - 4]
	buf, err := utils.UnBase32(resp)
	if err != nil{
		return "", err
	}
	return buf.String(), nil
}

func (s *CaptchaSender) CaptchaMarkValidQuery(captcha_id string, valid bool) error {
	str_valid := "bad"
	if valid{
		str_valid = "ok"
	}
	resp, err := QueryNSAPI(fmt.Sprintf("%s.%s.%s.val", str_valid, captcha_id, s.APIKey), s.BaseDomain, s.NS)
	if err != nil{
		return err
	}
	if resp != "ok."{
		return errors.New("val: " + resp)
	}
	return nil
}
