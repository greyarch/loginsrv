package login

import (
	"flag"
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfig_ReadConfigDefaults(t *testing.T) {
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	assert.Equal(t, &DefaultConfig, ReadConfig())
}

func TestConfig_ReadConfig(t *testing.T) {
	input := []string{
		"--host=host",
		"--port=port",
		"--log-level=loglevel",
		"--text-logging=true",
		"--jwt-secret=jwtsecret",
		"--success-url=successurl",
		"--cookie-name=cookiename",
		"--cookie-http-only=false",
		"--backend=provider=simple",
		"--backend=provider=foo",
	}

	expected := &Config{
		Host:           "host",
		Port:           "port",
		LogLevel:       "loglevel",
		TextLogging:    true,
		JwtSecret:      "jwtsecret",
		SuccessUrl:     "successurl",
		CookieName:     "cookiename",
		CookieHttpOnly: false,
		Backends: BackendOptions{
			map[string]string{
				"provider": "simple",
			},
			map[string]string{
				"provider": "foo",
			},
		},
	}

	cfg, err := readConfig(flag.NewFlagSet("", flag.ContinueOnError), input)
	assert.NoError(t, err)
	assert.Equal(t, expected, cfg)
}

func TestConfig_ReadConfigFromEnv(t *testing.T) {
	assert.NoError(t, os.Setenv("LOGINSRV_HOST", "host"))
	assert.NoError(t, os.Setenv("LOGINSRV_PORT", "port"))
	assert.NoError(t, os.Setenv("LOGINSRV_LOG_LEVEL", "loglevel"))
	assert.NoError(t, os.Setenv("LOGINSRV_TEXT_LOGGING", "true"))
	assert.NoError(t, os.Setenv("LOGINSRV_JWT_SECRET", "jwtsecret"))
	assert.NoError(t, os.Setenv("LOGINSRV_SUCCESS_URL", "successurl"))
	assert.NoError(t, os.Setenv("LOGINSRV_COOKIE_NAME", "cookiename"))
	assert.NoError(t, os.Setenv("LOGINSRV_COOKIE_HTTP_ONLY", "false"))
	assert.NoError(t, os.Setenv("LOGINSRV_BACKEND", "provider=simple,foo=bar"))
	assert.NoError(t, os.Setenv("LOGINSRV_BACKEND_FOO", "provider=foo"))
	assert.NoError(t, os.Setenv("LOGINSRV_BACKEND_BAR", "provider=bar"))

	expected := &Config{
		Host:           "host",
		Port:           "port",
		LogLevel:       "loglevel",
		TextLogging:    true,
		JwtSecret:      "jwtsecret",
		SuccessUrl:     "successurl",
		CookieName:     "cookiename",
		CookieHttpOnly: false,
		Backends: BackendOptions{
			map[string]string{
				"provider": "simple",
				"foo":      "bar",
			},
			map[string]string{
				"provider": "foo",
			},
			map[string]string{
				"provider": "bar",
			},
		},
	}

	cfg, err := readConfig(flag.NewFlagSet("", flag.ContinueOnError), []string{})
	assert.NoError(t, err)
	assert.Equal(t, expected, cfg)
}

func TestConfig_ParseBackendOptions(t *testing.T) {
	testCases := []struct {
		input       []string
		expected    BackendOptions
		expectError bool
	}{
		{
			[]string{},
			BackendOptions{},
			false,
		},
		{
			[]string{"name=p1,key1=value1,key2=value2"},
			BackendOptions{},
			true, // no provider name specified
		},
		{
			[]string{
				"provider=simple,name=p1,key1=value1,key2=value2",
				"provider=simple,name=p2,key3=value3,key4=value4",
			},
			BackendOptions{
				map[string]string{
					"provider": "simple",
					"name":     "p1",
					"key1":     "value1",
					"key2":     "value2",
				},
				map[string]string{
					"provider": "simple",
					"name":     "p2",
					"key3":     "value3",
					"key4":     "value4",
				},
			},
			false,
		},
		{
			[]string{"foo"},
			BackendOptions{},
			true,
		},
	}
	for i, test := range testCases {
		t.Run(fmt.Sprintf("test %v", i), func(t *testing.T) {
			options := &BackendOptions{}
			for _, input := range test.input {
				err := options.Set(input)
				if test.expectError {
					assert.Error(t, err)
				} else {
					if err != nil {
						assert.NoError(t, err)
						continue
					}
				}
			}
			assert.Equal(t, test.expected, *options)
		})
	}
}
