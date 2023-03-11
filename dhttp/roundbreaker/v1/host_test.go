package roundbreaker_test

import (
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/dhttp/roundbreaker/v1"
	"github.com/don-nv/go-dpkg/djson/v1"
	"github.com/don-nv/go-dpkg/dstruct/v1"
	"github.com/sethvargo/go-envconfig"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
	"testing"
	"time"
)

func TestHostEnvConfigJSON_EnvDecode(t *testing.T) {
	env, err := djson.Marshal([]roundbreaker.HostConfig{
		{
			Scheme:   "https",
			HostPort: "host_1:port",
			Attempts: dstruct.AttemptsV1{
				Delays: []time.Duration{
					time.Millisecond,
					time.Microsecond,
					time.Second,
				},
			},
			DisabledFor: time.Second,
		},
		{
			Scheme:   "http",
			HostPort: "host_2:port",
			Attempts: dstruct.AttemptsV1{
				Delays: []time.Duration{
					time.Millisecond,
					time.Microsecond,
				},
			},
			DisabledFor: time.Millisecond,
		},
	})
	require.NoError(t, err)

	t.Setenv("CLIENT_HOSTS", string(env))
	t.Logf("\n%s", env)

	var config struct {
		Hosts roundbreaker.HostsEnvConfigsJSON `env:"CLIENT_HOSTS"`
	}

	err = envconfig.Process(dctx.New(), &config)
	require.NoError(t, err)

	env, err = djson.MarshalPretty(&config)
	require.NoError(t, err)

	t.Logf("\n%s", env)
}

func TestHostEnvConfigYAML_EnvDecode(t *testing.T) {
	env, err := yaml.Marshal([]roundbreaker.HostConfig{
		{
			Scheme:   "https",
			HostPort: "host_1:port",
			Attempts: dstruct.AttemptsV1{
				Delays: []time.Duration{
					time.Millisecond,
					time.Microsecond,
					time.Second,
				},
			},
			DisabledFor: time.Second,
		},
		{
			Scheme:   "http",
			HostPort: "host_2:port",
			Attempts: dstruct.AttemptsV1{
				Delays: []time.Duration{
					time.Millisecond,
					time.Microsecond,
				},
			},
			DisabledFor: time.Millisecond,
		},
	})
	require.NoError(t, err)

	t.Setenv("CLIENT_HOSTS", string(env))
	t.Logf("\n%s", env)

	var config struct {
		Hosts roundbreaker.HostsEnvConfigsYAML `env:"CLIENT_HOSTS"`
	}

	err = envconfig.Process(dctx.New(), &config)
	require.NoError(t, err)

	env, err = djson.MarshalPretty(&config)
	require.NoError(t, err)

	t.Logf("\n%s", env)
}
