package dlog

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog"
	"time"
)

type Builder struct {
	ctx zerolog.Context
	log Logger
}

func newBuilder(log Logger) Builder {
	return Builder{
		log: log,
		/*
			TODO?
				This may be replaced with custom initialization. Under the hood, zerolog makes []byte escaping to heap
				(spill). Custom initialization may use sync pool of []byte.
		*/
		ctx: log.zero.With(),
	}
}

func (b Builder) Build() Logger {
	b.log.zero = b.ctx.Logger()

	return b.log
}

// Scope - reads context key-value pairs and populates data with it.
func (b Builder) Scope(ctx context.Context) Builder {
	b.log = b.log.readCtx(ctx, b.log)

	return b
}

// Name - adds `names` - a separate field of each log.
func (b Builder) Name(names ...string) Builder {
	for _, name := range names {
		b.log.names = append(b.log.names, name)
	}

	return b
}

/*
Any - this should be used if a `value` can be of several types depending on a circumstances this method was called.
Otherwise, it's better to use type specific methods, because they don't have a performance drawback this method has.
*/
func (b Builder) Any(key string, value interface{}) Builder {
	switch v := value.(type) {
	case string:
		return b.String(key, v)
	case []string:
		return b.Strings(key, v)

	case error:
		return b.Error(key, v)
	case []error:
		return b.Errors(key, v)

	case bool:
		return b.Bool(key, v)
	case []bool:
		return b.Bools(key, v)

	case int:
		return b.Int(key, v)
	case []int:
		return b.Ints(key, v)
	case int8:
		return b.Int8(key, v)
	case []int8:
		return b.Ints8(key, v)
	case int16:
		return b.Int16(key, v)
	case []int16:
		return b.Ints16(key, v)
	case int32:
		return b.Int32(key, v)
	case []int32:
		return b.Ints32(key, v)
	case int64:
		return b.Int64(key, v)
	case []int64:
		return b.Ints64(key, v)

	case uint:
		return b.Uint(key, v)
	case []uint:
		return b.Uints(key, v)
	case uint8:
		return b.Uint8(key, v)
	case []byte:
		return b.Bytes(key, v)
	case uint16:
		return b.Uint16(key, v)
	case []uint16:
		return b.Uints16(key, v)
	case uint32:
		return b.Uint32(key, v)
	case []uint32:
		return b.Uints32(key, v)
	case uint64:
		return b.Uint64(key, v)
	case []uint64:
		return b.Uints64(key, v)

	case float32:
		return b.Float32(key, v)
	case []float32:
		return b.Floats32(key, v)
	case float64:
		return b.Float64(key, v)
	case []float64:
		return b.Floats64(key, v)

	case time.Time:
		return b.Time(key, v)
	case []time.Time:
		return b.Times(key, v)
	case time.Duration:
		return b.Duration(key, v)
	case []time.Duration:
		return b.Durations(key, v)

	/*
		TODO
			- Add known types registration with respective json marshal functions. So, each time a known object
			passed, get respective marshal function, marshal object to bytes and add them as bytes;
			- Ensure correct order: first registered, then json marshallers type assertion;
	*/
	case json.Marshaler:
		return b.ObjectMarshallerJSON(key, v)

	default:
		b.ctx = b.ctx.Interface(key, v)

		return b
	}
}

func (b Builder) String(key string, value string) Builder {
	b.ctx = b.ctx.Str(key, value)

	return b
}

func (b Builder) Strings(key string, value []string) Builder {
	b.ctx = b.ctx.Strs(key, value)

	return b
}

func (b Builder) Error(key string, value error) Builder {
	b.ctx = b.ctx.AnErr(key, value)

	return b
}

func (b Builder) Errors(key string, value []error) Builder {
	b.ctx = b.ctx.Errs(key, value)

	return b
}

func (b Builder) Bool(key string, value bool) Builder {
	b.ctx = b.ctx.Bool(key, value)

	return b
}

func (b Builder) Bools(key string, value []bool) Builder {
	b.ctx = b.ctx.Bools(key, value)

	return b
}

func (b Builder) Int(key string, value int) Builder {
	b.ctx = b.ctx.Int(key, value)

	return b
}

func (b Builder) Ints(key string, value []int) Builder {
	b.ctx = b.ctx.Ints(key, value)

	return b
}

func (b Builder) Int8(key string, value int8) Builder {
	b.ctx = b.ctx.Int8(key, value)

	return b
}

func (b Builder) Ints8(key string, value []int8) Builder {
	b.ctx = b.ctx.Ints8(key, value)

	return b
}

func (b Builder) Int16(key string, value int16) Builder {
	b.ctx = b.ctx.Int16(key, value)

	return b
}

func (b Builder) Ints16(key string, value []int16) Builder {
	b.ctx = b.ctx.Ints16(key, value)

	return b
}

func (b Builder) Int32(key string, value int32) Builder {
	b.ctx = b.ctx.Int32(key, value)

	return b
}

func (b Builder) Ints32(key string, value []int32) Builder {
	b.ctx = b.ctx.Ints32(key, value)

	return b
}

func (b Builder) Int64(key string, value int64) Builder {
	b.ctx = b.ctx.Int64(key, value)

	return b
}

func (b Builder) Ints64(key string, value []int64) Builder {
	b.ctx = b.ctx.Ints64(key, value)

	return b
}

func (b Builder) Uint(key string, value uint) Builder {
	b.ctx = b.ctx.Uint(key, value)

	return b
}

func (b Builder) Uints(key string, value []uint) Builder {
	b.ctx = b.ctx.Uints(key, value)

	return b
}

func (b Builder) Uint8(key string, value uint8) Builder {
	b.ctx = b.ctx.Uint8(key, value)

	return b
}

func (b Builder) Bytes(key string, value []byte) Builder {
	b.ctx = b.ctx.Bytes(key, value)

	return b
}

func (b Builder) Uint16(key string, value uint16) Builder {
	b.ctx = b.ctx.Uint16(key, value)

	return b
}

func (b Builder) Uints16(key string, value []uint16) Builder {
	b.ctx = b.ctx.Uints16(key, value)

	return b
}

func (b Builder) Uint32(key string, value uint32) Builder {
	b.ctx = b.ctx.Uint32(key, value)

	return b
}

func (b Builder) Uints32(key string, value []uint32) Builder {
	b.ctx = b.ctx.Uints32(key, value)

	return b
}

func (b Builder) Uint64(key string, value uint64) Builder {
	b.ctx = b.ctx.Uint64(key, value)

	return b
}

func (b Builder) Uints64(key string, value []uint64) Builder {
	b.ctx = b.ctx.Uints64(key, value)

	return b
}

func (b Builder) Float32(key string, value float32) Builder {
	b.ctx = b.ctx.Float32(key, value)

	return b
}

func (b Builder) Floats32(key string, value []float32) Builder {
	b.ctx = b.ctx.Floats32(key, value)

	return b
}

func (b Builder) Float64(key string, value float64) Builder {
	b.ctx = b.ctx.Float64(key, value)

	return b
}

func (b Builder) Floats64(key string, value []float64) Builder {
	b.ctx = b.ctx.Floats64(key, value)

	return b
}

func (b Builder) Time(key string, value time.Time) Builder {
	b.ctx = b.ctx.Time(key, value)

	return b
}

func (b Builder) Times(key string, value []time.Time) Builder {
	b.ctx = b.ctx.Times(key, value)

	return b
}

func (b Builder) Duration(key string, value time.Duration) Builder {
	b.ctx = b.ctx.Dur(key, value)

	return b
}

func (b Builder) Durations(key string, value []time.Duration) Builder {
	b.ctx = b.ctx.Durs(key, value)

	return b
}

func (b Builder) ObjectMarshallerJSON(key string, value json.Marshaler) Builder {
	data, err := value.MarshalJSON()
	if err != nil {
		b.ctx.Err(fmt.Errorf(key+" marshalling value as json: %w", err))

		return b
	}

	b.ctx = b.ctx.Bytes(key, data)

	return b
}
