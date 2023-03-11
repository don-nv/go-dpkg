package dprom

import (
	"context"
	"fmt"
	"github.com/sethvargo/go-envconfig"
	"time"
)

var configEnv = mustNewConfigEnv()

type config struct {
	Namespace           string    `env:"DPROM_NAMESPACE,default=dprom_default_namespace"`
	Subsystem           string    `env:"DPROM_SUBSYSTEM,default=dprom_default_subsystem"`
	BucketsPartitionsMs []float64 `env:"DPROM_BUCKETS_PARTITIONS_MS,default=0,25,50,100,200,400,800,1600,3200"`
	GORuntime           struct {
		Enabled bool          `env:"DPROM_GO_RUNTIME_ENABLED,default=true"`
		Delay   time.Duration `env:"DPROM_GO_RUNTIME_DELAY,default=30s"`
	}
}

func mustNewConfigEnv() config {
	var config config

	err := envconfig.Process(context.Background(), &config)
	if err != nil {
		panic(fmt.Errorf("prcessing environment configuration: %w", err))
	}

	const minGoRuntimeDelay = time.Second
	if config.GORuntime.Enabled && config.GORuntime.Delay < minGoRuntimeDelay {
		panic(fmt.Errorf("too short DPROM_GO_RUNTIME_DELAY, < %q", minGoRuntimeDelay))
	}

	return config
}
