package config

import "net"

type Sources []SourceConfiguration

type NetworkList []string

func (n NetworkList) Contains(ip net.IP) bool {
	for _, whitelist := range n {
		if _, net, err := net.ParseCIDR(whitelist); err == nil && net.Contains(ip) {
			return true
		}
	}

	return false
}

func (s Sources) GetConfig(t SourceType) map[string]interface{} {
	for _, source := range s {
		if source.Type == t {
			return source.Config
		}
	}

	return make(map[string]interface{})
}

func (s Sources) ShouldLoad(t SourceType) bool {
	for _, source := range s {
		if source.Type == t {
			return true
		}
	}

	return false
}

type Config struct {
	Debug bool `yaml:"debug"`

	Listen string `yaml:"listen"`

	ForwardedForWhitelist NetworkList `yaml:"forwardedForWhitelist"`

	MetricsWhitelist NetworkList `yaml:"metricsWhitelist"`

	Sources Sources `yaml:"sources"`
}

type SourceType string

type SourceConfiguration struct {
	Type     SourceType             `yaml:"type"`
	Priority uint8                  `yaml:"priority"`
	Config   map[string]interface{} `yaml:"config"`
}

func (c *Config) Validate() error {
	return nil
}

func (c *Config) loadDefaults() {
	c.Listen = "169.254.169.254:80"
}
