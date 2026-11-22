package llmrec

import "context"

// providerAdmissionContextKey is scoped to one task/agent execution. The hook
// runs synchronously before EVERY physical model request (including the
// Recorder's automatic retries), not merely before the logical completion.
type providerAdmissionContextKey struct{}

// WithProviderAdmission installs a preflight check. A non-nil error rejects
// the provider request before any network call. A durable budget integration
// must perform its reservation transaction inside this callback and fail
// closed on any database/owner-fence error.
//
// Multiple calls compose in registration order; the first error stops admission.
// Without this hook the Recorder retains its existing passthrough behavior.
// A successful admission is not a usage settlement; actual attempts continue
// to be recorded through the usage sink after the provider returns.
func WithProviderAdmission(ctx context.Context, admit func(context.Context) error) context.Context {
	if admit == nil {
		return ctx
	}
	if previous, ok := ctx.Value(providerAdmissionContextKey{}).(func(context.Context) error); ok && previous != nil {
		combined := func(ctx context.Context) error {
			if err := previous(ctx); err != nil {
				return err
			}
			return admit(ctx)
		}
		return context.WithValue(ctx, providerAdmissionContextKey{}, combined)
	}
	return context.WithValue(ctx, providerAdmissionContextKey{}, admit)
}

func checkProviderAdmission(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if admit, ok := ctx.Value(providerAdmissionContextKey{}).(func(context.Context) error); ok {
		return admit(ctx)
	}
	return nil
}
