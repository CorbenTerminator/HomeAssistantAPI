package homeassistant

import (
	"github.com/valyala/fasthttp"
	"net"
	"time"
	"fmt"
	"errors"
)

//HomeAssistant - main struct that the user will use
type HomeAssistant struct {
	URL 	string
	bearerToken	string
	httpClient	*fasthttp.Client
}

//this http Client will communicate with Home Assistant in LAN, so https is not required
var urlTemplate = "http://%s:%d/api/"
var bearerTemplate = "Bearer %s"

//Creates Home Asisstnat object, initialize httpClient
func NewHomeAssistant(host string, port int, bearer string) (ha *HomeAssistant) {
	fastClient := &fasthttp.Client{
		Dial: func(addr string) (net.Conn, error) {
            return fasthttp.DialTimeout(addr, time.Second*10)
        },
	}
	url := fmt.Sprintf(urlTemplate, host, port)
	ha = &HomeAssistant{
		URL: url,
		bearerToken: bearer,
		httpClient: fastClient,
	}
	return
}

//Connect - Checks connection with HomeAssistant and its availability
func (ha *HomeAssistant) Connect() error {
	//'/api/' - Returns a message if the HA API is up and running
	req, resp := ha.createReqResp(ha.URL)
	err := ha.httpClient.Do(req, resp)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		//Unathorized
		if resp.StatusCode() == 401 {
			return errors.New("Bearer token is invalid")
		} else {
			return errors.New(string(resp.Body()))
		}
	}
	return nil
}

func (ha *HomeAssistant) createReqResp(url string) (*fasthttp.Request, *fasthttp.Response) {
	req := fasthttp.AcquireRequest()
	//create request with header 'Authorization: Bearer ABCDEFG...'
	req.SetRequestURI(url)
	req.Header.Set("Authorization", fmt.Sprintf(bearerTemplate, ha.bearerToken))
	req.Header.SetContentType("application/json")
	//create response
	resp := fasthttp.AcquireResponse()
	return req, resp
}

func (ha *HomeAssistant) httpGet(path string) ([]byte, error) {
	url := ha.URL + path
	req, resp := ha.createReqResp(url)
	req.Header.SetMethod("GET")
	err := ha.httpClient.Do(req, resp)
	if err != nil {
		return nil, err
	}
	return resp.Body(), nil
}

func (ha *HomeAssistant) httpPost(path string, json []byte) ([]byte, error) {
	url := ha.URL + path
	req, resp := ha.createReqResp(url)
	req.Header.SetMethod("POST")
	req.SetBody(json)
	err := ha.httpClient.Do(req, resp)
	if err != nil {
		return nil, err
	}
	if (resp.StatusCode() != 200) && (resp.StatusCode() != 201) {
		return nil, errors.New(string(resp.Body()))
	}
	return resp.Body(), nil
}