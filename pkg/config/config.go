package config

type config struct {
	Listen string `yaml:"listen"`

	Source SourceConfiguration `yaml:"source"`
}

type SourceType string

type SourceConfiguration struct {
	Type   SourceType             `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}

func (c *config) Validate() error {
	return nil
}

func (c *config) loadDefaults() {
	c.Listen = "169.254.169.254:80"
	c.Source = SourceConfiguration{
		Type: "gpcloud",
		Config: map[string]interface{}{
			"client_id":     "",
			"client_secret": "",
		},
	}
}
