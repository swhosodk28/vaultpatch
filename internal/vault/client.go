package vault

import (
	"context"
	"fmt"

	vaultapi "github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client with helper methods.
type Client struct {
	underlying *vaultapi.Client
}

// Config holds configuration for creating a Vault client.
type Config struct {
	Address string
	Token   string
}

// NewClient creates a new Vault client from the given config.
func NewClient(cfg Config) (*Client, error) {
	vcfg := vaultapi.DefaultConfig()
	vcfg.Address = cfg.Address

	c, err := vaultapi.NewClient(vcfg)
	if err != nil {
		return nil, fmt.Errorf("creating vault client: %w", err)
	}

	c.SetToken(cfg.Token)

	return &Client{underlying: c}, nil
}

// ReadSecrets reads all key-value pairs at the given KV v2 path.
func (c *Client) ReadSecrets(ctx context.Context, mount, path string) (map[string]string, error) {
	secret, err := c.underlying.KVv2(mount).Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("reading secrets at %s/%s: %w", mount, path, err)
	}
	if secret == nil || secret.Data == nil {
		return map[string]string{}, nil
	}

	result := make(map[string]string, len(secret.Data))
	for k, v := range secret.Data {
		str, ok := v.(string)
		if !ok {
			str = fmt.Sprintf("%v", v)
		}
		result[k] = str
	}
	return result, nil
}

// WriteSecrets writes the provided key-value pairs to the given KV v2 path.
func (c *Client) WriteSecrets(ctx context.Context, mount, path string, data map[string]string) error {
	payload := make(map[string]interface{}, len(data))
	for k, v := range data {
		payload[k] = v
	}

	_, err := c.underlying.KVv2(mount).Put(ctx, path, payload)
	if err != nil {
		return fmt.Errorf("writing secrets to %s/%s: %w", mount, path, err)
	}
	return nil
}
