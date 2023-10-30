package dlog_test

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/don-nv/go-dpkg/dlog/v1"
	"github.com/rs/zerolog"
	"go.uber.org/zap"
	"os"
	"testing"
)

func BenchmarkLogger_Consumer(b *testing.B) {
	var consumer = NewConsumerLogger()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = consumer.LogBench(ctx) //nolint:errcheck
		}
	})
}

func BenchmarkZerolog_Consumer(b *testing.B) {
	var consumer = NewConsumerZerolog()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = consumer.LogBench(ctx) //nolint:errcheck
		}
	})
}

func BenchmarkZap_Consumer(b *testing.B) {
	var consumer = NewConsumerZap()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = consumer.LogBench(ctx) //nolint:errcheck
		}
	})
}

type ConsumerLogger struct {
	log dlog.Logger
}

func NewConsumerLogger() ConsumerLogger {
	return ConsumerLogger{
		log: dlog.New().With().Name("consumer").Build(),
	}
}

func (c *ConsumerLogger) LogBench(ctx context.Context) (err error) {
	var log = c.log.With().Name("logging_bench").Scope(ctx).Build()
	defer log.CatchE(&err)

	var info = log.I()
	defer info.Write("running...").Write("...done")

	log = log.With().Bytes("body", BytesJSON).Build()

	return Err
}

func (c ConsumerLogger) LogObjectMarshallerJSON(ctx context.Context) (err error) {
	var log = c.log.With().Name("logging_object_marshaller_json").Scope(ctx).Build()
	defer log.CatchE(&err)

	var info = log.I()
	defer info.Write("running...").Write("...done")

	var object = NewDataObjectJSON()
	log = log.With().ObjectMarshallerJSON("object_json", object).Build()

	return Err
}

type ConsumerZerolog struct {
	log zerolog.Logger
}

func NewConsumerZerolog() ConsumerZerolog {
	return ConsumerZerolog{
		log: zerolog.New(os.Stdout).With().Str("name", "consumer").Logger(),
	}
}

func (c ConsumerZerolog) LogBench(ctx context.Context) (err error) {
	var log = c.log.With().Str("name", "logging_bench").Logger()

	id := dctx.GoID(ctx)
	if id != "" {
		log = log.With().Str("go_id", id).Logger()
	}

	id = dctx.XRequestID(ctx)
	if id != "" {
		log = log.With().Str("x_req_id", id).Logger()
	}

	log.Info().Msg("running...")
	defer log.Info().Msg("...done")

	defer func() {
		if err != nil {
			log.Error().Msg(err.Error())
		}
	}()

	log = log.With().Bytes("body", BytesJSON).Logger()

	return Err
}

type ConsumerZap struct {
	log *zap.Logger
}

func NewConsumerZap() ConsumerZap {
	config := zap.NewProductionConfig()
	config.Sampling = nil
	config.DisableCaller = true

	z, err := config.Build()
	if err != nil {
		panic(err)
	}

	return ConsumerZap{
		log: z,
	}
}

func (c ConsumerZap) LogBench(ctx context.Context) (err error) {
	var log = c.log.Named("name").Named("logging_bench")

	id := dctx.GoID(ctx)
	if id != "" {
		log = log.With(zap.String("go_id", id))
	}

	id = dctx.XRequestID(ctx)
	if id != "" {
		log = log.With(zap.String("go_id", id))
	}

	log.Info("running...")
	defer log.Info("...done")

	defer func() {
		if err != nil {
			log.Error(err.Error())
		}
	}()

	log = log.With(zap.ByteString("body", BytesJSON))

	return Err
}

var (
	BytesJSON = []byte("" +
		"{\"glossary\": {\"title\": \"example glossary\",\"GlossDiv\": {\"title\": \"S\",\"GlossList\": " +
		"{\"GlossEntry\": {\"ID\": \"SGML\",\"SortAs\": \"SGML\",\"GlossTerm\": \"Standard Generalized " +
		"Markup Language\",\"Acronym\": \"SGML\",\"Abbrev\": \"ISO 8879:1986\",\"GlossDef\": {\"para\": " +
		"\"A meta-markup language, used to create markup languages such as DocBook.\",\"GlossSeeAlso\": " +
		"[\"GML\",\"XML\"]},\"GlossSee\": \"markup\"}}}}}",
	)
	Err = errors.New("" +
		"purus viverra accumsan in nisl nisi scelerisque eu ultrices vitae auctor eu augue ut lectus arcu bibendum at " +
		"varius vel pharetra vel turpis nunc eget lorem dolor sed viverra ipsum nunc aliquet bibendum enim facilisis " +
		"gravida neque convallis a cras",
	)
)

type DataObjectJSON struct {
	Field1 string `json:"field_1"`
	Field2 string `json:"field_2"`
	Field3 string `json:"field_3"`
}

func NewDataObjectJSON() DataObjectJSON {
	return DataObjectJSON{
		Field1: "value_1",
		Field2: "value_2",
		Field3: "value_3",
	}
}

func (d DataObjectJSON) MarshalJSON() ([]byte, error) {
	type Buff DataObjectJSON

	return json.Marshal(Buff(d))
}
