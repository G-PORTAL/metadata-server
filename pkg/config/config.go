package config

type Sources []SourceConfiguration

func (s Sources) GetConfig(t SourceType) map[string]interface{} {
	for _, source := range s {
		if source.Type == t {
			return source.Config
		}
	}

	return make(map[string]interface{})
}

type config struct {
	Listen string `yaml:"listen"`

	Sources Sources `yaml:"sources"`
}

type SourceType string

type SourceConfiguration struct {
	Type     SourceType             `yaml:"type"`
	Priority uint8                  `yaml:"priority"`
	Config   map[string]interface{} `yaml:"config"`
}

func (c *config) Validate() error {
	return nil
}

func (c *config) loadDefaults() {
	c.Listen = "169.254.169.254:80"
}
