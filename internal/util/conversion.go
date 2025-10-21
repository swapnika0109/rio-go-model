package util

// LanguageMapper maps language and country to language code
func LanguageMapper(theme string) string {
	if theme == "2" {
		return "en-IN"
	}
	return "en-US"
}

// SafeStringSlice converts interface{} to []string safely.
// It supports both []string and []interface{} inputs and returns nil for unsupported types.
func SafeStringSlice(val interface{}) []string {
	if val == nil {
		return nil
	}

	switch v := val.(type) {
	case []string:
		return v
	case []interface{}:
		result := make([]string, 0, len(v))
		for _, item := range v {
			if str, ok := item.(string); ok {
				result = append(result, str)
			}
		}
		return result
	default:
		return nil
	}
}
