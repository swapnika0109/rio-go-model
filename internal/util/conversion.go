package util

// LanguageMapper maps language and country to language code
func LanguageMapper(language string) string {
	if language == "" {
		return "en-US"
	}

	switch language {
	case "English":
		return "en-US"
	case "Telugu":
		return "te-IN"
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
