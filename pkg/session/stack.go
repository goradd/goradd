package session

import "context"

// PushStack pushes the given value onto the given named stack in the session.
// The name becomes a variable in the current session. The value must be JsonSerializable.
func PushStack(ctx context.Context, stack string, value string) {
	var a []string
	if Has(ctx, stack) {
		a, _ = Get(ctx, stack).([]string)
	}
	a = append(a, value)
	Set(ctx, stack, a)
}

// PopStack pops the given value off of the named stack and returns it.
func PopStack(ctx context.Context, stack string) (value string) {
	var a []string
	if Has(ctx, stack) {
		a, _ = Get(ctx, stack).([]string)
	}
	if len(a) > 0 {
		value = a[len(a)-1]
		a = a[:len(a)-1]
	}

	if len(a) == 0 {
		Remove(ctx, stack)
	} else {
		Set(ctx, stack, a)
	}
	return
}

// ClearStack clears the named stack
func ClearStack(ctx context.Context, stack string) {
	Remove(ctx, stack)
}
