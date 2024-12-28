package context

import "context"

// key is the key used to store value in context.
type key string

// Followings are the known context keys.
const (
	keySource         key = "source"
	keySchoolID       key = "school_id"
	keyUserID         key = "user_id"
	keyHTTPStatusCode key = "http_status_code"
)

// SetSource returns a new Context that carries value v as
// source.
func SetSource(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, keySource, v)
}

// GetSource returns the source value stored in the given
// context, if any.
func GetSource(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(keySource).(string)
	return v, ok
}

// SetSchoolID returns a new Context that carries value v as
// school ID.
func SetSchoolID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, keySchoolID, v)
}

// GetSchoolID returns the school ID value stored in the given
// context, if any.
func GetSchoolID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(keySchoolID).(string)
	return v, ok
}

// SetUserID returns a new Context that carries value v as
// user ID.
func SetUserID(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, keyUserID, v)
}

// GetUserID returns the user ID value stored in the given
// context, if any.
func GetUserID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(keyUserID).(string)
	return v, ok
}

// SetHTTPStatusCode returns a new Context that carries
// value v as HTTP status code.
func SetHTTPStatusCode(ctx context.Context, v int) context.Context {
	return context.WithValue(ctx, keyHTTPStatusCode, v)
}

// GetHTTPStatusCode returns the HTTP status code value
// stored in the given context, if any.
func GetHTTPStatusCode(ctx context.Context) (int, bool) {
	v, ok := ctx.Value(keyHTTPStatusCode).(int)
	return v, ok
}
