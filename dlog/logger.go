package dlog

import (
	"context"
)

/*
	TODO?:
		Global logger implementation:
			* Logger do self-nesting for each unique function or method once;
			* During nesting Logger invokes WithName() according to function or struct{}.method() name;
			* Each subsequent call returns preconfigured Logger that is suited for that place;
			* There are may be two getters: dlog.Log() and dlog.LogCtx(context.WithContext);
			* Has order to omit callers check (to decide which one Logger has to be returned) for each Logger getting,
			may be added code-generation;
*/

/*
Logger - writes logs. This is a generalized interface. Overall behaviour:
  - Each method returning Logger instance returns original instance, not a copy, except for the Copy();
  - It isn't safe to be used concurrently;
  - Logs writing is asynchronous;
  - Logging context format may vary and depends on an implementation;
*/
type Logger interface {
	// E - writes log at LevelError. Nil error is acceptable.
	E(err error) Logger

	// W - writes log at LevelWarn.
	W(msg string) Logger

	// I - writes log at LevelInfo.
	I(msg string) Logger

	// D - writes log at LevelDebug.
	D(msg string) Logger

	/*
		WithName - adds `names` to the logging context. This should be used by method "owner" on its invocation or at
		functions when there is an interest to add additional logging context.

		Example:
			func New(log Logger) Server {
				return Server{log: log.WithName("server")}
			}

			func (s Server) HandleRegisterUserRequest(r *http.Request) {
				var log = s.log.WithName("handling_register_user_request")
				log.I("served")
				// [I] server.handling_register_user_request	served
			}
	*/
	WithName(names ...string) Logger

	/*
		WithContext - looks for context keys, reads their respective values and populates logging context with them.
		Behaviour and a result is similar to Logger.WithKV(). This should be used when there is a context available
		and its data should be logged.

		Example:
			log.WithContext(ctx).I("served")
			// [I] served {"x_req_id":"req_123","go_id":"go_123"}
	*/
	WithContext(ctx context.Context) Logger

	/*
		WithKV - adds key/value pair to the logging context. This should be used to populate logging context with
		interesting data.

		Example:
			log.WithKV("id", "id_1").WithKV("amount", 100).I("done")
			// [I] done {"id":"id_1","amount":100}
	*/
	WithKV(key string, value interface{}) Logger

	// Copy - returns Logger copy.
	Copy() Logger

	/*
		CatchE - catches an `err` and if `*err != nil`, logs error value at LevelError. If `knownErrs` is not nil, then
		`*err != nil` gets compared against them and if found, nothing will be logged. This should be used with
		defer not to invoke write log.E() in many places. Also, it's convenient to be invoked alongside
		Logger.WithKV() at potentially buggy or critical places in order to reveal operation intermediaries on error.
		This method might be tricky to use in cases when it comes to error shadowing or redeclaration.

		Example:
			func F(log Logger) {
				var err error
				defer log.CatchE(&err)
				err = errors.New("an error occurred")
				// [E] an error occurred
			}
	*/
	CatchE(err *error, knownErrs ...error) Logger

	/*
		CatchED - is the same as Logger.CatchE(), but logs predefined message at LevelDebug if `*err == nil` or there is
		no match with `knownErrs`.
	*/
	CatchED(err *error, knownErrs ...error) Logger

	// Flushing - is a long-running operation awaiting pending logs to be written. Context controls max duration.
	Flushing(ctx context.Context) Logger
}
