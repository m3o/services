package microtime

import (
	"context"
	"strconv"
	"time"

	"github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/metadata"
)

// TimeFromContext provides the current time optionally parsed in the "Time" request header
func TimeFromContext(ctx context.Context) (time.Time, error) {
	// There isn't any context
	md, ok := metadata.FromContext(ctx)
	if !ok {
		return time.Now(), errors.InternalServerError("MISSING_CONTEXT", "No metadata exists on the current context")
	}

	// No time was provided, use the default
	timeStr, ok := md["Time"]
	if !ok {
		return time.Now(), nil
	}

	// Parse the Unix timestamp
	unix, err := strconv.Atoi(timeStr)
	if err != nil {
		return time.Now(), errors.BadRequest("INVALID_TIME", "The Time request header was invalid")
	}
	return time.Unix(int64(unix), 0), nil
}

// ContextWithTime returns a context with the time set
func ContextWithTime(ctx context.Context, t time.Time) context.Context {
	md := metadata.Metadata{
		"Time": strconv.Itoa(int(t.Unix())),
	}

	return metadata.MergeContext(ctx, md, true)
}
