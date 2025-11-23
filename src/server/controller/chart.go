// Package token
package controller

import (
	"chart-provider/src/interfaces/global"
	"chart-provider/src/interfaces/logger"
	"chart-provider/src/interfaces/server/dto"
	"chart-provider/src/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	clientId      = "charts-rn-desktop"
	clientSecret  = "igljsnfBunGqI706JnQRIkQuJB65iscC"
	deviceAuthUrl = "https://identity.api.navigraph.com/connect/deviceauthorization"
	tokenUrl      = "https://identity.api.navigraph.com/connect/token"
	clientScope   = "userinfo openid offline_access amdb charts email navdata userdata fmsdata tiles simbrief"
)

type ChartController struct {
	logger      logger.Interface
	token       *TokenResponse
	expiresIn   time.Time
	initialized bool
	client      *http.Client
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
}

type DeviceAuthResponse struct {
	DeviceCode              string `json:"device_code"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
	UserCode                string `json:"user_code"`
	VerificationUrl         string `json:"verification_uri"`
	VerificationUrlComplete string `json:"verification_uri_complete"`
}

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func addHeader(req *http.Request) {
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) NavigraphCharts/8.38.3 Chrome/106.0.5249.199 Electron/21.4.1 Safari/537.36")
	req.Header.Add("Accept", "application/json, text/plain, */*")
	req.Header.Add("Host", "identity.api.navigraph.com")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "zh-CN")
}

func NewChartController(
	lg logger.Interface,
) *ChartController {
	manager := &ChartController{
		logger:      logger.NewLoggerAdapter(lg, "TokenManager"),
		initialized: false,
	}
	manager.getCachedFlushToken()
	manager.client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	go func(m *ChartController) {
		if m.token != nil && m.refreshAccessToken() {
			m.logger.Info("Use cached flush token")
			m.initialized = true
			return
		}
		response, verifier := m.requestDeviceAuthorization()
		if response == nil {
			m.logger.Error("Request device authorization fail")
			return
		}
		go m.pollForAccessToken(response.DeviceCode, verifier, response.Interval, response.ExpiresIn)
	}(manager)
	return manager
}

func (m *ChartController) getCachedFlushToken() {
	cachedToken, err := m.readFlushTokenFromFile()
	if err != nil {
		m.logger.Errorf("Read cached flush token error, %s", err.Error())
		m.token = nil
		return
	}
	if cachedToken == "" {
		m.logger.Warn("No cached flush token found")
		m.token = nil
		return
	}
	m.token = &TokenResponse{RefreshToken: cachedToken}
}

func (m *ChartController) readFlushTokenFromFile() (string, error) {
	file, err := os.OpenFile(*global.TokenCacheFile, os.O_RDONLY|os.O_CREATE, global.DefaultFilePermissions)
	if err != nil {
		m.logger.Errorf("Open token cache file fail: %s", err.Error())
		return "", fmt.Errorf("open token cache file fail: %s", err.Error())
	}

	defer func(file *os.File) { _ = file.Close() }(file)

	cachedToken, err := io.ReadAll(file)
	if err != nil {
		m.logger.Errorf("Read token cache file fail: %s", err.Error())
		return "", fmt.Errorf("read token cache file fail: %s", err.Error())
	}

	return string(cachedToken), nil
}

func (m *ChartController) saveFlushTokenToFile(token string) {
	file, err := os.OpenFile(*global.TokenCacheFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, global.DefaultFilePermissions)
	if err != nil {
		m.logger.Errorf("Save token cache file fail: %s", err.Error())
		return
	}
	defer func(file *os.File) { _ = file.Close() }(file)
	_, err = file.WriteString(token)
	if err != nil {
		m.logger.Errorf("Save token cache file fail: %s", err.Error())
		return
	}
}

func (m *ChartController) requestDeviceAuthorization() (*DeviceAuthResponse, string) {
	payload := &url.Values{}
	payload.Add("client_id", clientId)
	payload.Add("client_secret", clientSecret)
	pkceGenerator := utils.NewPKCEGenerator()
	verifier, _ := pkceGenerator.GenerateCodeVerifier()
	challenge := pkceGenerator.GenerateCodeChallenge(verifier)
	payload.Add("code_challenge", challenge)
	payload.Add("code_challenge_method", pkceGenerator.GetCodeChallengeMethod())
	req, _ := http.NewRequest("POST", deviceAuthUrl, strings.NewReader(payload.Encode()))
	addHeader(req)
	res, err := m.client.Do(req)
	if err != nil {
		m.logger.Errorf("device authorization network fail: %s", err.Error())
		return nil, ""
	}
	if res.StatusCode != http.StatusOK {
		m.logger.Errorf("device authorization fail with http status %d", res.StatusCode)
		return nil, ""
	}
	response := &DeviceAuthResponse{}
	data, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if err := json.Unmarshal(data, response); err != nil {
		m.logger.Errorf("device authorization fail: %s", err.Error())
		return nil, ""
	}
	m.logger.Infof("Device authorization, please visit %s to manual authorization", response.VerificationUrlComplete)
	return response, verifier
}

func (m *ChartController) pollForAccessToken(deviceCode string, verifier string, interval int, expiresIn int) {
	intervalDuration := time.Duration(interval) * time.Second
	ticker := time.NewTicker(intervalDuration)
	defer ticker.Stop()

	timeout := time.After(time.Duration(expiresIn) * time.Second)

	for {
		select {
		case <-ticker.C:
			payload := &url.Values{}
			payload.Add("client_id", clientId)
			payload.Add("client_secret", clientSecret)
			payload.Add("code_verifier", verifier)
			payload.Add("device_code", deviceCode)
			payload.Add("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
			payload.Add("scope", clientScope)
			req, _ := http.NewRequest("POST", tokenUrl, strings.NewReader(payload.Encode()))
			addHeader(req)
			res, err := m.client.Do(req)
			if err != nil {
				m.logger.Errorf("pollForAccessToken Error: %s", err.Error())
				_ = res.Body.Close()
				return
			}
			data, _ := io.ReadAll(res.Body)
			_ = res.Body.Close()
			if res.StatusCode != http.StatusOK {
				errorResponse := &ErrorResponse{}
				if err := json.Unmarshal(data, errorResponse); err != nil {
					m.logger.Errorf("pollForAccessToken Unmarshal Error: %s", err.Error())
					return
				}
				if errorResponse.Error == "authorization_pending" {
					continue
				} else if errorResponse.Error == "access_denied" {
					m.logger.Error("Device authorization fail: user denied")
					return
				} else if errorResponse.Error == "slow_down" {
					m.logger.Error("Device authorization fail: slow down")
					intervalDuration += 5 * time.Second
					ticker.Reset(intervalDuration)
					continue
				} else {
					m.logger.Error("Device authorization timeout")
					return
				}
			}
			token := &TokenResponse{}
			if err := json.Unmarshal(data, token); err != nil {
				m.logger.Errorf("pollForAccessToken Unmarshal Error: %s", err.Error())
				return
			}
			m.token = token
			m.initialized = true
			m.saveFlushTokenToFile(token.RefreshToken)
			m.expiresIn = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
			m.logger.Info("Device authorization passed")
			return
		case <-timeout:
			m.logger.Error("Device authorization timeout")
			return
		}
	}
}

func (m *ChartController) refreshAccessToken() bool {
	payload := &url.Values{}
	payload.Add("client_id", clientId)
	payload.Add("client_secret", clientSecret)
	payload.Add("refresh_token", m.token.RefreshToken)
	payload.Add("grant_type", "refresh_token")
	payload.Add("scope", clientScope)
	req, _ := http.NewRequest("POST", tokenUrl, strings.NewReader(payload.Encode()))
	addHeader(req)
	res, err := m.client.Do(req)
	if err != nil {
		m.logger.Errorf("refreshAccessToken Error: %s", err.Error())
		return false
	}
	if res.StatusCode != http.StatusOK {
		m.logger.Errorf("refreshAccessToken StatusCode: %d", res.StatusCode)
		return false
	}
	token := &TokenResponse{}
	data, _ := io.ReadAll(res.Body)
	_ = res.Body.Close()
	if err := json.Unmarshal(data, token); err != nil {
		m.logger.Errorf("refreshAccessToken Unmarshal Error: %s", err.Error())
		return false
	}
	m.token = token
	m.saveFlushTokenToFile(token.RefreshToken)
	m.expiresIn = time.Now().Add(time.Duration(token.ExpiresIn) * time.Second)
	m.logger.Info("Refresh access token success")
	return true
}

var (
	ErrCreateRequest = dto.NewApiStatus("ERR_CREATE_REQUEST", "创建请求失败", dto.HttpCodeInternalError)
	ErrSendRequest   = dto.NewApiStatus("ERR_SEND_REQUEST", "请求目标失败", dto.HttpCodeInternalError)
	ErrCopyRequest   = dto.NewApiStatus("ERR_COPY_REQUEST", "复制目标请求", dto.HttpCodeInternalError)
	ErrNotAvailable  = dto.NewApiStatus("ERR_NOT_AVAILABLE", "航图服务不可用", dto.HttpCodeInternalError)
	ErrTokenExpired  = dto.NewApiStatus("TOKEN_EXPIRED", "令牌已过期，请联系管理员", dto.HttpCodeUnauthorized)
)

func (m *ChartController) HandleProxy(c echo.Context) error {
	if !m.initialized {
		return dto.NewApiResponse[any](ErrNotAvailable, nil).Response(c)
	}

	originalRequest := c.Request()
	targetUrl := c.Param("*")

	req, err := http.NewRequest(originalRequest.Method, targetUrl, originalRequest.Body)
	if err != nil {
		m.logger.Errorf("HandleProxy Error: %s", err.Error())
		return dto.NewApiResponse[any](ErrCreateRequest, nil).Response(c)
	}

	for key, values := range originalRequest.Header {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	req.Header.Set("Authorization", m.getAccessToken())

	resp, err := m.client.Do(req)
	if err != nil {
		m.logger.Errorf("HandleProxy SendRequest Error: %s", err.Error())
		return dto.NewApiResponse[any](ErrSendRequest, nil).Response(c)
	}

	for key, values := range resp.Header {
		for _, value := range values {
			c.Response().Header().Add(key, value)
		}
	}

	if resp.StatusCode == dto.HttpCodeUnauthorized.Code() {
		m.logger.Error("HandleProxy Error: Token expired")
		return dto.NewApiResponse[any](ErrTokenExpired, nil).Response(c)
	}

	c.Response().WriteHeader(resp.StatusCode)

	_, err = io.Copy(c.Response().Writer, resp.Body)
	_ = resp.Body.Close()

	if err != nil {
		m.logger.Errorf("HandleProxy Error: %s", err.Error())
		return dto.NewApiResponse[any](ErrCopyRequest, nil).Response(c)
	}

	return nil
}

func (m *ChartController) getAccessToken() string {
	if !m.initialized {
		return ""
	}
	if time.Now().After(m.expiresIn) {
		m.refreshAccessToken()
	}
	return "Bearer " + m.token.AccessToken
}
