package aigc

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	transportLock sync.Mutex
	transportLink *transportNode
)

type transportNode struct {
	next      *transportNode
	transport *http.Transport
	proxy     string
}

type HttpClient struct {
	Proxy   string
	Retries int

	client *http.Client
	inited bool
}

type roundtripRetryer struct {
	Retries   int
	transport *http.Transport
}

func dialContext(ctx context.Context, network string, addr string) (net.Conn, error) {
	var dialer = net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	return dialer.DialContext(ctx, network, addr)
}

func GetHttpTransport(proxy string) (*http.Transport, error) {
	var err error
	transportLock.Lock()
	defer transportLock.Unlock()

	if transportLink == nil {
		transportLink = &transportNode{
			transport: &http.Transport{
				// Proxy:                 http.ProxyFromEnvironment,
				Proxy:                 func(*http.Request) (*url.URL, error) { return nil, nil },
				DialContext:           dialContext,
				ForceAttemptHTTP2:     true,
				MaxIdleConns:          200,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   30 * time.Second,
				ExpectContinueTimeout: 5 * time.Second,
			},
			proxy: "",
		}
	}

	var node *transportNode = transportLink
	for {
		if node.proxy == proxy {
			return node.transport, nil
		}
		if node.next == nil {
			break
		}
		node = node.next
	}

	var proxyURL *url.URL

	proxyURL, err = url.Parse(proxy)
	if err != nil {
		return nil, fmt.Errorf("[GetHttpTransport] invalid proxy url '%s' %w", proxy, err)
	}

	node.next = &transportNode{
		transport: &http.Transport{
			Proxy:                 http.ProxyURL(proxyURL),
			DialContext:           dialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          200,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   30 * time.Second,
			ExpectContinueTimeout: 5 * time.Second,
		},
		proxy: proxy,
	}

	return node.next.transport, nil
}

func (r roundtripRetryer) RoundTrip(request *http.Request) (*http.Response, error) {
	var err error
	var response *http.Response

	if r.Retries > 0 {
		for i := 0; i < r.Retries; i++ {
			response, err = r.transport.RoundTrip(request)
			if err == nil && response.StatusCode < 500 {
				return response, nil
			}
			time.Sleep(time.Second * 2)
		}
		return response, err
	}

	return r.transport.RoundTrip(request)
}

func (c *HttpClient) Client() (*http.Client, error) {
	if c.inited {
		return c.client, nil
	}

	var err error
	var transport *http.Transport

	transport, err = GetHttpTransport(c.Proxy)
	if err != nil {
		return nil, fmt.Errorf("[HttpClient.Client] %w", err)
	}

	c.client = &http.Client{Transport: roundtripRetryer{transport: transport}}

	c.inited = true

	return c.client, nil
}

func (c *HttpClient) Do(request *http.Request) (*http.Response, error) {
	var err error
	var client *http.Client
	var response *http.Response

	client, err = c.Client()
	if err != nil {
		return nil, fmt.Errorf("[HttpClient.Do] %w", err)
	}

	response, err = client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("[HttpClient.Do] %w", err)
	}

	return response, nil
}

func HttpResponseText(response *http.Response) string {
	var errorText *bytes.Buffer

	errorText = AllocBuffer()
	defer FreeBuffer(errorText)

	errorText.WriteString(response.Status)
	errorText.WriteString(", ")
	io.Copy(errorText, response.Body)

	return errorText.String()
}
