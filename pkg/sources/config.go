package sources

type SourceConfig map[string]interface{}

func (s SourceConfig) GetString(key string) string {
	if value, ok := s[key].(string); ok {
		return value
	}

	return ""
}
