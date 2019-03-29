package spanner

import (
	"context"

	spannerV1 "cloud.google.com/go/spanner/apiv1"
	"google.golang.org/api/option"
)

type Client struct {
	config             *Config
	spannerNomalClient *spannerV1.Client
}

func NewClient(ctx context.Context, config *Config) (*Client, error) {
	opts := make([]option.ClientOption, 0)
	if config.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(config.CredentialsFile))
	}

	spannerClient, err := spannerV1.NewClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		config:             config,
		spannerNomalClient: spannerClient,
	}, nil
}

type Config struct {
	Project         string
	Instance        string
	Database        string
	CredentialsFile string
}
