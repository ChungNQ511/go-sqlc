package clients

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"example.com/m/pkg/utils/app"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// ClientConfig ClientConfig
type ClientConfig struct {
	KeycloakURL       string // keycloak url
	ClientID          string // client id
	ClientSecret      string // client secret
	Audience          string // audience
	Scope             string // scope
	RedisTokenKey     string // redis token key
	TokenCacheSeconds int    // token cache seconds
	URL               string // url
	ServiceName       string // service name
}

// Client Client
type Client struct {
	config     ClientConfig
	httpClient *http.Client
	logger     *log.Logger
}

// RequestOptions RequestOptions
type RequestOptions struct {
	Headers map[string]string
	Accept  string
}

// New New
// @param cfg ClientConfig
// @return *Client
func New(cfg ClientConfig) *Client {
	return &Client{
		config:     cfg,
		httpClient: &http.Client{Timeout: 15 * time.Second},
		logger:     log.New(os.Stdout, fmt.Sprintf("[%s] ", strings.ToUpper(cfg.ServiceName)), log.LstdFlags),
	}
}

// Get Get
// @param url string
// @param headers *RequestOptions
// @return *http.Response, error
func (cl *Client) Get(url string, headers *RequestOptions) (*http.Response, error) {
	return cl.doRequest("GET", url, nil, headers)
}

// Post Post
// @param url string
// @param body interface{}
// @param headers *RequestOptions
// @return *http.Response, error
func (cl *Client) Post(url string, body interface{}, headers *RequestOptions) (*http.Response, error) {
	return cl.doRequest("POST", url, body, headers)
}

// Put Put
// @param url string
// @param body interface{}
// @param headers *RequestOptions
// @return *http.Response, error
func (cl *Client) Put(url string, body interface{}, headers *RequestOptions) (*http.Response, error) {
	return cl.doRequest("PUT", url, body, headers)
}

// Patch Patch
// @param url string
// @param body interface{}
// @param headers *RequestOptions
// @return *http.Response, error
func (cl *Client) Patch(url string, body interface{}, headers *RequestOptions) (*http.Response, error) {
	return cl.doRequest("PATCH", url, body, headers)
}

// Delete Delete
// @param url string
// @param headers *RequestOptions
// @return *http.Response, error
func (cl *Client) Delete(url string, headers *RequestOptions) (*http.Response, error) {
	return cl.doRequest("DELETE", url, nil, headers)
}

// doRequest doRequest
// @param method string
// @param url string
// @param body interface{}
// @param opts *RequestOptions
// @return *http.Response, error
func (cl *Client) doRequest(method, url string, body interface{}, opts *RequestOptions) (*http.Response, error) {
	token, err := cl.getToken()
	if err != nil {
		return nil, fmt.Errorf("get token error: %w", err)
	}
	fmt.Println("🔁 Start doRequest with URL: ", cl.config.URL+url)
	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal body error: %w", err)
		}
		fmt.Println("🔁 Body: ", string(jsonBody))

		bodyReader = bytes.NewReader(jsonBody)
	}
	req, err := http.NewRequest(method, cl.config.URL+url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	// 🔀 Allow override Accept
	if opts != nil && opts.Accept != "" {
		req.Header.Set("Accept", opts.Accept)
	} else {
		req.Header.Set("Accept", "application/json")
	}

	// 👌 Add custom headers
	if opts != nil && opts.Headers != nil {
		for k, v := range opts.Headers {
			req.Header.Set(k, v)
		}
	}

	cl.logger.Printf("🚀 Start Request %s %s", method, cl.config.URL+url)

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		cl.logger.Printf("❌ Request error: %v", err)
		return nil, err
	}

	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			cl.logger.Printf("❌ Error reading response body: %v", err)
			return nil, err
		}
		cl.logger.Printf("❌ HTTP %d: %s", resp.StatusCode, string(body))
		return nil, errors.New("❌ HTTP Status " + strconv.Itoa(resp.StatusCode))
	}
	return resp, nil
}

// getToken getToken - get token from redis (optional when using client_credentials provider is keycloak)
// @return string, error
func (cl *Client) getToken() (string, error) {
	// use cached token if production mode
	token, err := app.RGetValue(cl.config.RedisTokenKey, "")
	if err == nil && tokenValid(token) {
		cl.logger.Println("✅ Using cached token from Redis")
		return token, nil
	}

	cl.logger.Println("🔁 Cached token missing or expired, requesting new token...")

	accessToken, err := cl.getClientCredentialsToken()
	if err != nil {
		return "", err
	}

	finalToken, err := cl.exchangeUMAToken(accessToken)
	if err != nil {
		return "", err
	}

	cl.logger.Println("Final token: ", finalToken)
	cl.logger.Println("✅ Token acquired, caching to Redis")
	_ = app.RSet(cl.config.RedisTokenKey, finalToken, cl.config.TokenCacheSeconds)
	return finalToken, nil
}

// getClientCredentialsToken getClientCredentialsToken
func (cl *Client) getClientCredentialsToken() (string, error) {
	resp, err := http.PostForm(cl.config.KeycloakURL, map[string][]string{
		"grant_type":    {"client_credentials"},
		"client_id":     {cl.config.ClientID},
		"client_secret": {cl.config.ClientSecret},
	})
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var parsed map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&parsed)
	if err != nil {
		return "", err
	}
	token := parsed["access_token"].(string)
	cl.logger.Println("🔐 Got token via client_credentials")
	return token, nil
}

// exchangeUMAToken exchangeUMAToken - optional when using client_credentials (provider is keycloak)
// @param accessToken string
// @return string, error
func (cl *Client) exchangeUMAToken(accessToken string) (string, error) {
	req, _ := http.NewRequest("POST", cl.config.KeycloakURL, strings.NewReader(
		"grant_type=urn:ietf:params:oauth:grant-type:uma-ticket&audience="+cl.config.Audience+"&scope="+strings.ReplaceAll(cl.config.Scope, " ", "%20"),
	))
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := cl.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var parsed map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&parsed)
	if err != nil {
		return "", err
	}
	token := parsed["access_token"].(string)
	cl.logger.Println("🎟  Got token via UMA ticket exchange")
	return token, nil
}

// tokenValid tokenValid
// @param tokenStr string
// @return bool
func tokenValid(tokenStr string) bool {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenStr, jwt.MapClaims{})
	if err != nil {
		return false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok { // use jwt.MapClaims to avoid type assertion
		if exp, ok := claims["exp"].(float64); ok {
			return int64(exp) > time.Now().Unix()+60 // add 60 seconds to avoid token expired
		}
	}
	return false
}

// HandleResponse HandleResponse
// @param c *gin.Context
// @param resp *http.Response
func (cl *Client) HandleResponse(c *gin.Context, resp *http.Response) {
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Error reading response: %v", err)
		return
	}

	// ✅ If JSON → return object
	if strings.Contains(contentType, "application/json") {
		var result interface{}
		if err := json.Unmarshal(bodyBytes, &result); err != nil {
			c.String(http.StatusInternalServerError, "Error parsing JSON: %v", err)
			return
		}
		c.JSON(resp.StatusCode, result)
		return
	}

	// ✅ If HTML or text → return raw
	if strings.Contains(contentType, "text/html") || strings.Contains(contentType, "text/plain") {
		c.Data(resp.StatusCode, contentType, bodyBytes)
		return
	}

	// ❓ If content-type is not clear, return raw
	c.Data(resp.StatusCode, "application/octet-stream", bodyBytes)
}

// ParseJSONResponse ParseJSONResponse
// @param resp *http.Response
// @return *T, error
func ParseJSONResponse[T any](resp *http.Response) (*T, error) {
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf("unexpected content type: %s", contentType)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var result T
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON: %w", err)
	}

	return &result, nil
}
