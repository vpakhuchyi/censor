package censor

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_ToString(t *testing.T) {
	t.Run("default config", func(t *testing.T) {
		// GIVEN a config instance.
		cfg := DefaultConfig()

		// WHEN the ToString method is called on the config instance.
		got := cfg.ToString()

		// THEN the returned string contains the configuration details.
		want := "---------------------------------------------------------------------\n          Censor is configured with the following settings:\n---------------------------------------------------------------------\nprint-config-on-init: true\nmask-value: '[CENSORED]'\nexclude-patterns: []\n---------------------------------------------------------------------\n"
		require.Equal(t, want, got)
	})

	t.Run("custom config", func(t *testing.T) {
		// GIVEN a config instance.
		cfg := Config{
			PrintConfigOnInit: false,
			MaskValue:         "test",
			ExcludePatterns:   []string{"[0-9]"},
		}

		// WHEN the ToString method is called on the config instance.
		got := cfg.ToString()

		// THEN the returned string contains the configuration details.
		want := "---------------------------------------------------------------------\n          Censor is configured with the following settings:\n---------------------------------------------------------------------\nprint-config-on-init: false\nmask-value: test\nexclude-patterns:\n    - '[0-9]'\n---------------------------------------------------------------------\n"
		require.Equal(t, want, got)
	})
}

func TestConfig_FromFile(t *testing.T) {
	t.Run("read config from file", func(t *testing.T) {
		// GIVEN a path to a configuration file.
		path := "testdata/cfg.yml"

		// WHEN the ConfigFromFile function is called with the path.
		got, err := ConfigFromFile(path)

		// THEN the config is read from the file and returned.
		require.NoError(t, err)
		require.Equal(t, Config{
			PrintConfigOnInit: true,
			MaskValue:         "[CENSORED]",
			ExcludePatterns:   []string{"\\d", "^\\w$"},
		}, got)
	})

	t.Run("fail to read config from file", func(t *testing.T) {
		// GIVEN a path to a non-existing configuration file.
		path := "testdata/non-existing.yml"

		// WHEN the ConfigFromFile function is called with the path.
		_, err := ConfigFromFile(path)

		// THEN an error is returned.
		require.Error(t, err)
	})

	t.Run("fail to unmarshal config from file", func(t *testing.T) {
		// GIVEN a path to a configuration file with invalid content.
		path := "testdata/invalid_cfg.yml"

		// WHEN the ConfigFromFile function is called with the path.
		_, err := ConfigFromFile(path)

		// THEN an error is returned.
		require.Error(t, err)
	})
}
