package liblog

import "context"

func RequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(cRequestId).(string); ok {
		return requestID
	}

	return ""
}
