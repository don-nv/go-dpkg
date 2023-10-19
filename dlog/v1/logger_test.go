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

var ctx = dctx.New(dctx.WithNewXRequestID())

func TestLogger_ObjectMarshallerJSON(t *testing.T) {
	var consumer = NewConsumerLogger()

	_ = consumer.LogObjectMarshallerJSON(ctx)
}

func BenchmarkLogger_Consumer(b *testing.B) {
	var consumer = NewConsumerLogger()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = consumer.LogBench(ctx)
		}
	})
}

func BenchmarkZerolog_Consumer(b *testing.B) {
	var consumer = NewConsumerZerolog()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = consumer.LogBench(ctx)
		}
	})
}

func BenchmarkZap_Consumer(b *testing.B) {
	var consumer = NewConsumerZap()

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = consumer.LogBench(ctx)
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

func (c ConsumerLogger) LogBench(ctx context.Context) (err error) {
	var log = c.log.With().Name("logging_bench").Scope(ctx).Build()
	defer log.WriteI("running...").Write("...done")
	defer log.CatchE(&err)

	log = log.With().Bytes("body", TestDataBytes).Build()

	return TestDataError
}

func (c ConsumerLogger) LogObjectMarshallerJSON(ctx context.Context) (err error) {
	var log = c.log.With().Name("logging_object_marshaller_json").Scope(ctx).Build()
	defer log.WriteI("running...").Write("...done")
	defer log.CatchE(&err)

	var object = NewTestDataObjectJSON()
	log = log.With().ObjectMarshallerJSON("object_json", object).Build()

	return TestDataError
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

	log = log.With().Bytes("body", TestDataBytes).Logger()

	return TestDataError
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

	log = log.With(zap.ByteString("body", TestDataBytes))

	return TestDataError
}

var (
	TestDataBytes = []byte("" +
		"{\"glossary\": {\"title\": \"example glossary\",\"GlossDiv\": {\"title\": \"S\",\"GlossList\": " +
		"{\"GlossEntry\": {\"ID\": \"SGML\",\"SortAs\": \"SGML\",\"GlossTerm\": \"Standard Generalized " +
		"Markup Language\",\"Acronym\": \"SGML\",\"Abbrev\": \"ISO 8879:1986\",\"GlossDef\": {\"para\": " +
		"\"A meta-markup language, used to create markup languages such as DocBook.\",\"GlossSeeAlso\": " +
		"[\"GML\",\"XML\"]},\"GlossSee\": \"markup\"}}}}}",
	)
	TestDataError = errors.New("" +
		"error error error error error error error error error " +
		"error error error error error error error error error " +
		"error error error error error error error error error " +
		"error error error error error error error error error",
	)
)

type TestDataObjectJSON struct {
	Field1 string `json:"field_1"`
	Field2 string `json:"field_2"`
	Field3 string `json:"field_3"`
}

func NewTestDataObjectJSON() TestDataObjectJSON {
	return TestDataObjectJSON{
		Field1: "value_1",
		Field2: "value_2",
		Field3: "value_3",
	}
}

func (t TestDataObjectJSON) MarshalJSON() ([]byte, error) {
	type Buff TestDataObjectJSON

	return json.Marshal(Buff(t))
}
