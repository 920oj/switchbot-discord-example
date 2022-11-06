package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

const (
	endPoint string = "https://api.switch-bot.com/v1.1"
)

type requestInfo struct {
	URL    string
	Method string
	Body   interface{}
}

type deviceListResponse struct {
	StatusCode int `json:"statusCode"`
	Body       struct {
		DeviceList []struct {
			DeviceID           string `json:"deviceId"`
			DeviceName         string `json:"deviceName"`
			DeviceType         string `json:"deviceType,omitempty"`
			HubDeviceID        string `json:"hubDeviceId"`
			EnableCloudService bool   `json:"enableCloudService,omitempty"`
		} `json:"deviceList"`
		InfraredRemoteList []struct {
			DeviceID    string `json:"deviceId"`
			DeviceName  string `json:"deviceName"`
			RemoteType  string `json:"remoteType"`
			HubDeviceID string `json:"hubDeviceId"`
		} `json:"infraredRemoteList"`
	} `json:"body"`
	Message string `json:"message"`
}

type botDeviceStatusResponse struct {
	StatusCode int `json:"statusCode"`
	Body       struct {
		DeviceID    string `json:"deviceId"`
		DeviceType  string `json:"deviceType"`
		HubDeviceID string `json:"hubDeviceId"`
		Power       string `json:"power"`
	} `json:"body"`
	Message string `json:"message"`
}

type botCommandRequest struct {
	Command     string `json:"command"`
	Parameter   string `json:"parameter"`
	CommandType string `json:"commandType"`
}

type botCommandResponse struct {
	StatusCode int `json:"statusCode"`
	Body       struct {
	} `json:"body"`
	Message string `json:"message"`
}

func requestApi(info requestInfo) (resp *http.Response, err error) {
	var req *http.Request

	// BodyがあればRequestに含める
	if info.Body != nil {
		reqJson, err := json.Marshal(info.Body)
		if err != nil {
			return nil, err
		}
		req, _ = http.NewRequest(info.Method, info.URL, bytes.NewBuffer(reqJson))
	} else {
		req, _ = http.NewRequest(info.Method, info.URL, nil)
	}

	// nonceヘッダ用にUUIDを作成
	u, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	requestId := u.String()

	// tヘッダ用に現在時刻をUNIXタイムスタンプで出力
	miliTime := time.Now().UTC().UnixNano() / int64(time.Millisecond)
	t := strconv.FormatInt(miliTime, 10)

	// 署名作成
	sign, err := makeRequestSign(t, requestId, switchbotToken, switchbotSecret)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", switchbotToken)
	req.Header.Set("sign", sign)
	req.Header.Set("nonce", requestId)
	req.Header.Set("t", t)
	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}

func makeRequestSign(t, nonce, token, secret string) (string, error) {
	data := token + t + nonce

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func requestGetDeviceList() (res deviceListResponse, err error) {
	url := endPoint + "/devices"
	info := requestInfo{
		URL:    url,
		Method: http.MethodGet,
	}

	resp, err := requestApi(info)
	if err != nil {
		return deviceListResponse{}, err
	}

	b, _ := io.ReadAll(resp.Body)

	var r deviceListResponse
	json.Unmarshal(b, &r)

	return r, nil
}

func requestGetBotDeviceStatus(deviceId string) (res botDeviceStatusResponse, err error) {
	url := endPoint + "/devices/" + deviceId + "/status"
	info := requestInfo{
		URL:    url,
		Method: http.MethodGet,
	}

	resp, err := requestApi(info)
	if err != nil {
		return botDeviceStatusResponse{}, err
	}

	b, _ := io.ReadAll(resp.Body)

	var r botDeviceStatusResponse
	json.Unmarshal(b, &r)

	return r, nil
}

func requestPostBotCommand(deviceId, command string) (res botCommandResponse, err error) {
	url := endPoint + "/devices/" + deviceId + "/commands"
	reqBody := botCommandRequest{
		Command:     command,
		Parameter:   "default",
		CommandType: "command",
	}
	info := requestInfo{
		URL:    url,
		Method: http.MethodPost,
		Body:   reqBody,
	}

	resp, err := requestApi(info)
	if err != nil {
		return botCommandResponse{}, err
	}

	b, _ := io.ReadAll(resp.Body)

	var r botCommandResponse
	json.Unmarshal(b, &r)

	return r, nil
}
