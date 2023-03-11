package roundbreaker

import (
	"github.com/don-nv/go-dpkg/djson/v1"
	"github.com/don-nv/go-dpkg/dstruct/v1"
	"github.com/don-nv/go-dpkg/dsync/v1"
	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v3"
	"time"
)

func newHostsFromEnv[T HostEnvConfigJSON | HostEnvConfigYAML](ts []T) []Host {
	if len(ts) < 1 {
		return nil
	}

	var (
		n     = len(ts)
		hosts = make([]Host, 0, n)
	)
	for _, t := range ts {
		var config = HostConfig(t)

		hosts = append(hosts, Host{
			Scheme:      config.Scheme,
			HostPort:    config.HostPort,
			Attempts:    dstruct.NewAttemptsV1Sync(config.Attempts),
			DisabledFor: config.DisabledFor,
		})
	}

	return hosts
}

// Host - represents remote host to send http requests to. It is safe to be used concurrently.
type Host struct {
	mu          dsync.RWMutex
	disabled    bool
	Attempts    dstruct.AttemptsV1Sync
	Scheme      string
	HostPort    string
	DisabledFor time.Duration
}

func (h *Host) disableFor() {
	h.mu.LockF(func() {
		// Prevent multiple enabling.
		if !h.disabled && h.Attempts.Exceeded() {
			go func() {
				<-time.After(h.DisabledFor)

				h.enable()
			}()

			h.disabled = true
		}
	})
}

func (h *Host) enable() {
	h.mu.LockF(func() {
		h.disabled = false
	})
}

func (h *Host) enabled() bool {
	defer h.mu.RLock().RUnlock()

	return !h.disabled
}

type HostConfig struct {
	Scheme string `json:"scheme" yaml:"scheme"`

	// HostPort - port is optional.
	HostPort string `json:"host_port" yaml:"host_port"`

	// Attempts - once, all attempts were failed, host will not be attempted for DisabledFor.
	Attempts dstruct.AttemptsV1 `json:"attempts" yaml:"attempts"`

	// DisabledFor - controls duration for which host gets disabled.
	DisabledFor time.Duration `json:"disabled_for" yaml:"disabled_for"`
}

type (
	iHostsEnvConfigs interface {
		envconfig.Decoder
		Hosts() []Host
	}

	iHostEnvConfig interface {
		envconfig.Decoder
	}
)

// HostsEnvConfigsJSON - is the same as HostEnvConfigJSON, but is a slice.
type HostsEnvConfigsJSON []HostEnvConfigJSON

var _ iHostsEnvConfigs = (*HostsEnvConfigsJSON)(nil)

func (c *HostsEnvConfigsJSON) EnvDecode(val string) error {
	return djson.Unmarshal([]byte(val), c)
}

func (c *HostsEnvConfigsJSON) Hosts() []Host {
	if c == nil {
		return nil
	}

	return newHostsFromEnv(*c)
}

// HostsEnvConfigsYAML - is the same as HostEnvConfigYAML, but is a slice.
type HostsEnvConfigsYAML []HostEnvConfigYAML

var _ iHostsEnvConfigs = (*HostsEnvConfigsYAML)(nil)

func (c *HostsEnvConfigsYAML) EnvDecode(val string) error {
	return yaml.Unmarshal([]byte(val), c)
}

func (c *HostsEnvConfigsYAML) Hosts() []Host {
	if c == nil {
		return nil
	}

	return newHostsFromEnv(*c)
}

// HostEnvConfigJSON - expects "github.com/sethvargo/go-envconfig" usage of envconfig.Process().
type HostEnvConfigJSON HostConfig

var _ iHostEnvConfig = (*HostEnvConfigJSON)(nil)

func (c *HostEnvConfigJSON) EnvDecode(val string) error {
	return djson.Unmarshal([]byte(val), c)
}

// HostEnvConfigYAML - expects "github.com/sethvargo/go-envconfig" usage of envconfig.Process().
type HostEnvConfigYAML HostConfig

var _ iHostEnvConfig = (*HostEnvConfigYAML)(nil)

func (c *HostEnvConfigYAML) EnvDecode(val string) error {
	return yaml.Unmarshal([]byte(val), c)
}
