package config

import "githup.com/htl/tcclienttest/internal/appserver/options"

//appserver config
type Config struct {
	*options.Options
}

// CreateConfigFromOptions creates a running configuration instance based
// on a given IAM pump command line or configuration file option.
func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
