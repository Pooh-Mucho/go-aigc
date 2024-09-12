package chat

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

import (
	"github.com/Pooh-Mucho/go-aigc"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
)

var (
	// Do not use constants because some APIs want to take address
	httpAcceptAll       = "*/*"
	httpContentTypeJson = "application/json"
)

// bedrockCredentialsProvider is an implementation of aws.CredentialsProvider
type bedrockCredentialsProvider struct {
	AccessKey string
	SecretKey string
}

type bedrockClient struct {
	Region    string
	AccessKey string
	SecretKey string
	Proxy     string
	Retries   int

	client *bedrockruntime.Client
}

func (p *bedrockCredentialsProvider) Retrieve(ctx context.Context) (aws.Credentials, error) {
	_ = ctx
	return aws.Credentials{
		AccessKeyID:     p.AccessKey,
		SecretAccessKey: p.SecretKey,
		SessionToken:    "",
		Source:          "",
		CanExpire:       false,
		Expires:         time.Time{},
	}, nil
}

func (c *bedrockClient) Client() (*bedrockruntime.Client, error) {
	if c.client != nil {
		return c.client, nil
	}

	var err error
	var transport *http.Transport
	var bedrockOpts bedrockruntime.Options

	transport, err = aigc.GetHttpTransport(c.Proxy)
	if err != nil {
		return nil, fmt.Errorf("[bedrockClient.Client] %w", err)
	}

	bedrockOpts = bedrockruntime.Options{
		Region:      c.Region,
		Credentials: &bedrockCredentialsProvider{AccessKey: c.AccessKey, SecretKey: c.SecretKey},
		HTTPClient:  &http.Client{Transport: transport},
	}

	if c.Retries > 0 {
		bedrockOpts.RetryMaxAttempts = c.Retries
	}

	c.client = bedrockruntime.New(bedrockOpts)

	return c.client, nil
}

func (c *bedrockClient) InvokeModel(
	ctx context.Context,
	params *bedrockruntime.InvokeModelInput,
) (*bedrockruntime.InvokeModelOutput, error) {
	var err error
	var client *bedrockruntime.Client

	client, err = c.Client()
	if err != nil {
		return nil, fmt.Errorf("[bedrockClient.InvokeModel] %w", err)
	}

	return client.InvokeModel(ctx, params)
}
