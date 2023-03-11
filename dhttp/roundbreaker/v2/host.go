package roundbreaker

import (
	"context"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/derr/v1"
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
			scheme:      config.Scheme,
			hostPort:    config.HostPort,
			attempts:    dstruct.NewAttemptsSync(config.Attempts),
			disabledFor: config.DisabledFor,
		})
	}

	return hosts
}

// Host - represents remote host to send http requests to. It is safe to be used concurrently.
type Host struct {
	mu          dsync.RWMutex
	attempts    dstruct.AttemptsSync
	enableNow   context.CancelFunc
	disabled    bool
	scheme      string
	hostPort    string
	disabledFor time.Duration
}

/*
attempt - attempts 'f' at least once if host enabled.

Errors:
  - derr.ErrDisabled
  - derr.ErrExceeded
*/
func (h *Host) attempt(f func() error) error {
	if !h.enabled() {
		return derr.ErrDisabled
	}

	err := f()
	if err != nil {
		// Inc attempt.
		if !h.attempts.Next() {
			h.disableFor()

			return derr.Join(err, derr.ErrExceeded)
		}

		return err
	}

	// Was disabled concurrently, but robustness has been just proven (err == nil).
	if !h.enabled() {
		h.revokeDisableFor()

		return nil
	}

	// Prefer read over write lock for infrequent cases.
	if h.attempts.AtLeastOnceAttempted() {
		h.attempts.Reset()
	}

	return nil
}

/*
disableFor - disables enabled host for 'disabledFor' duration or until 'enableNow' call. Enabling is done in
a separate goroutine.
*/
func (h *Host) disableFor() {
	h.mu.LockF(func() {
		// Prevent multiple enabling.
		if h.enabledUnsafe() {
			var ctx, cancel = dctx.NewTTL()

			h.enableNow = cancel

			go func() {
				select {
				case <-ctx.Done():
				case <-time.After(h.disabledFor):
				}

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

// revokeDisableFor - calls for host enabling and resets attempts.
func (h *Host) revokeDisableFor() {
	h.mu.LockF(func() {
		h.attempts.Reset()

		if h.enableNow != nil {
			h.enableNow()
		}
	})
}

func (h *Host) enabled() bool {
	defer h.mu.RLock().RUnlock()

	return h.enabledUnsafe()
}

func (h *Host) enabledUnsafe() bool {
	return !h.disabled
}

type HostConfig struct {
	Scheme string `json:"scheme" yaml:"scheme"`

	// HostPort - port is optional.
	HostPort string `json:"host_port" yaml:"host_port"`

	// Attempts - once, all attempts were failed, host will not be attempted for DisabledFor.
	Attempts dstruct.Attempts `json:"attempts" yaml:"attempts"`

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
