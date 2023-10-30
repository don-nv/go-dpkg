package dlog

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/don-nv/go-dpkg/dctx/v1"
	"github.com/rs/zerolog"
	"time"
)

type Data struct {
	zctx   zerolog.Context
	logger Logger
}

func newData(logger Logger) Data {
	return Data{
		logger: logger,
		/*
			TODO?
				This may be replaced with custom initialization. Under the hood, zerolog makes []byte escaping to heap
				(spill). Custom initialization may use sync pool of []byte.
		*/
		zctx: logger.zero.With(),
	}
}

func (d Data) Build() Logger {
	d.logger.zero = d.zctx.Logger()

	return d.logger
}

// Scope - reads context key-value pairs and populates data with it.
func (d Data) Scope(ctx context.Context) Data {
	d.logger = d.logger.readScope(ctx, d.logger)

	//d.logger.D().Write(dctx.GoID(ctx) + " foo")
	//d.Build().D().Write(dctx.GoID(ctx) + " foo")

	id := dctx.GoID(ctx)
	if id != "" {
		d = d.String("go_id", id)
	}

	id = dctx.XRequestID(ctx)
	if id != "" {
		d = d.String("x_req_id", id)
	}

	return d
}

// Name - adds `names` - a separate field of each log.
func (d Data) Name(names ...string) Data {
	if len(names) < 1 {
		return d
	}

	var (
		oldNames = d.logger.names
		allNames = make([][]byte, 0, len(oldNames)+len(names))
	)
	allNames = append(allNames, oldNames...)
	for _, name := range names {
		allNames = append(allNames, []byte(name))
	}
	d.logger.names = allNames

	return d
}

/*
Any - this should be used if a `value` can be of several types depending on a circumstances this method was called.
Otherwise, it's better to use type specific methods, because they don't have a performance drawback this method has.
*/
//nolint:funlen,cyclop
func (d Data) Any(key string, value interface{}) Data {
	switch v := value.(type) {
	case string:
		return d.String(key, v)
	case []string:
		return d.Strings(key, v)

	case error:
		return d.Error(key, v)
	case []error:
		return d.Errors(key, v)

	case bool:
		return d.Bool(key, v)
	case []bool:
		return d.Bools(key, v)

	case int:
		return d.Int(key, v)
	case []int:
		return d.Ints(key, v)
	case int8:
		return d.Int8(key, v)
	case []int8:
		return d.Ints8(key, v)
	case int16:
		return d.Int16(key, v)
	case []int16:
		return d.Ints16(key, v)
	case int32:
		return d.Int32(key, v)
	case []int32:
		return d.Ints32(key, v)
	case int64:
		return d.Int64(key, v)
	case []int64:
		return d.Ints64(key, v)

	case uint:
		return d.Uint(key, v)
	case []uint:
		return d.Uints(key, v)
	case uint8:
		return d.Uint8(key, v)
	case []byte:
		return d.Bytes(key, v)
	case uint16:
		return d.Uint16(key, v)
	case []uint16:
		return d.Uints16(key, v)
	case uint32:
		return d.Uint32(key, v)
	case []uint32:
		return d.Uints32(key, v)
	case uint64:
		return d.Uint64(key, v)
	case []uint64:
		return d.Uints64(key, v)

	case float32:
		return d.Float32(key, v)
	case []float32:
		return d.Floats32(key, v)
	case float64:
		return d.Float64(key, v)
	case []float64:
		return d.Floats64(key, v)

	case time.Time:
		return d.Time(key, v)
	case []time.Time:
		return d.Times(key, v)
	case time.Duration:
		return d.Duration(key, v)
	case []time.Duration:
		return d.Durations(key, v)

	/*
		TODO
			- Add known types registration with respective json marshal functions. So, each time a known object
			passed, get respective marshal function, marshal object to bytes and add them as bytes;
			- Ensure correct order: first registered, then json marshallers type assertion;
	*/
	case json.Marshaler:
		return d.ObjectMarshallerJSON(key, v)

	default:
		d.zctx = d.zctx.Interface(key, v)

		return d
	}
}

func (d Data) String(key string, value string) Data {
	d.zctx = d.zctx.Str(key, value)

	return d
}

func (d Data) Strings(key string, value []string) Data {
	d.zctx = d.zctx.Strs(key, value)

	return d
}

func (d Data) Error(key string, value error) Data {
	d.zctx = d.zctx.AnErr(key, value)

	return d
}

func (d Data) Errors(key string, value []error) Data {
	d.zctx = d.zctx.Errs(key, value)

	return d
}

func (d Data) Bool(key string, value bool) Data {
	d.zctx = d.zctx.Bool(key, value)

	return d
}

func (d Data) Bools(key string, value []bool) Data {
	d.zctx = d.zctx.Bools(key, value)

	return d
}

func (d Data) Int(key string, value int) Data {
	d.zctx = d.zctx.Int(key, value)

	return d
}

func (d Data) Ints(key string, value []int) Data {
	d.zctx = d.zctx.Ints(key, value)

	return d
}

func (d Data) Int8(key string, value int8) Data {
	d.zctx = d.zctx.Int8(key, value)

	return d
}

func (d Data) Ints8(key string, value []int8) Data {
	d.zctx = d.zctx.Ints8(key, value)

	return d
}

func (d Data) Int16(key string, value int16) Data {
	d.zctx = d.zctx.Int16(key, value)

	return d
}

func (d Data) Ints16(key string, value []int16) Data {
	d.zctx = d.zctx.Ints16(key, value)

	return d
}

func (d Data) Int32(key string, value int32) Data {
	d.zctx = d.zctx.Int32(key, value)

	return d
}

func (d Data) Ints32(key string, value []int32) Data {
	d.zctx = d.zctx.Ints32(key, value)

	return d
}

func (d Data) Int64(key string, value int64) Data {
	d.zctx = d.zctx.Int64(key, value)

	return d
}

func (d Data) Ints64(key string, value []int64) Data {
	d.zctx = d.zctx.Ints64(key, value)

	return d
}

func (d Data) Uint(key string, value uint) Data {
	d.zctx = d.zctx.Uint(key, value)

	return d
}

func (d Data) Uints(key string, value []uint) Data {
	d.zctx = d.zctx.Uints(key, value)

	return d
}

func (d Data) Uint8(key string, value uint8) Data {
	d.zctx = d.zctx.Uint8(key, value)

	return d
}

func (d Data) Bytes(key string, value []byte) Data {
	d.zctx = d.zctx.Bytes(key, value)

	return d
}

func (d Data) Uint16(key string, value uint16) Data {
	d.zctx = d.zctx.Uint16(key, value)

	return d
}

func (d Data) Uints16(key string, value []uint16) Data {
	d.zctx = d.zctx.Uints16(key, value)

	return d
}

func (d Data) Uint32(key string, value uint32) Data {
	d.zctx = d.zctx.Uint32(key, value)

	return d
}

func (d Data) Uints32(key string, value []uint32) Data {
	d.zctx = d.zctx.Uints32(key, value)

	return d
}

func (d Data) Uint64(key string, value uint64) Data {
	d.zctx = d.zctx.Uint64(key, value)

	return d
}

func (d Data) Uints64(key string, value []uint64) Data {
	d.zctx = d.zctx.Uints64(key, value)

	return d
}

func (d Data) Float32(key string, value float32) Data {
	d.zctx = d.zctx.Float32(key, value)

	return d
}

func (d Data) Floats32(key string, value []float32) Data {
	d.zctx = d.zctx.Floats32(key, value)

	return d
}

func (d Data) Float64(key string, value float64) Data {
	d.zctx = d.zctx.Float64(key, value)

	return d
}

func (d Data) Floats64(key string, value []float64) Data {
	d.zctx = d.zctx.Floats64(key, value)

	return d
}

func (d Data) Time(key string, value time.Time) Data {
	d.zctx = d.zctx.Time(key, value)

	return d
}

func (d Data) Times(key string, value []time.Time) Data {
	d.zctx = d.zctx.Times(key, value)

	return d
}

func (d Data) Duration(key string, value time.Duration) Data {
	d.zctx = d.zctx.Dur(key, value)

	return d
}

func (d Data) Durations(key string, value []time.Duration) Data {
	d.zctx = d.zctx.Durs(key, value)

	return d
}

func (d Data) ObjectMarshallerJSON(key string, value json.Marshaler) Data {
	data, err := value.MarshalJSON()
	if err != nil {
		d.zctx = d.zctx.Err(fmt.Errorf(key+" marshalling value as json: %w", err))

		return d
	}

	d.zctx = d.zctx.Bytes(key, data)

	return d
}
