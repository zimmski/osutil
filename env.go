package osutil

import (
	"fmt"
	"os"
	"strings"
)

// EnvironMap returns a map of the current environment variables.
func EnvironMap() (environMap map[string]string) {
	environ := os.Environ()
	environMap = make(map[string]string, len(environ))
	for _, e := range environ {
		kv := strings.SplitN(e, "=", 2)
		environMap[kv[0]] = kv[1]
	}

	return environMap
}

// EnvOrDefault returns the environment variable with the given key, or the default value if the key is not defined.
func EnvOrDefault(key string, defaultValue string) (value string) {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return defaultValue
}

// RequireEnv returns the environment variable with the given key, or an error if the key is not defined.
func RequireEnv(key string) (value string, err error) {
	if v, ok := os.LookupEnv(key); ok {
		return v, nil
	}

	return "", fmt.Errorf("environment variable %q needs to be set", key)
}

// IsEnvEnabled checks if the environment variable is enabled.
// By default an environment variable is considered enabled if it is set to "1", "true", "on" or "yes". Further such values can be provided as well. Capitalization is ignored.
func IsEnvEnabled(key string, additionalEnabledValues ...string) bool {
	value, ok := os.LookupEnv(key)
	if !ok {
		return false
	}
	value = strings.ToLower(value)

	if value == "1" || value == "true" || value == "on" || value == "yes" {
		return true
	}

	for _, match := range additionalEnabledValues {
		if value == match {
			return true
		}
	}

	return false
}
