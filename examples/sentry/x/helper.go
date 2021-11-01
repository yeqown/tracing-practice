package x

import (
	"context"
	"fmt"

	"github.com/getsentry/sentry-go"
)

func LogWithContext(ctx context.Context, msg string, args ...interface{}) {
	sp := sentry.TransactionFromContext(ctx)
	fmt.Printf("called traceId=%s %s\n", sp.ToSentryTrace(), fmt.Sprintf(msg, args...))
}
