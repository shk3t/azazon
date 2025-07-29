package authtest

import (
	"base/pkg/helper"
	"testing"
)

func assertEqual[V comparable](
	t *testing.T,
	currentValue V,
	expectedValue V,
	info string,
) {
	if currentValue != expectedValue {
		t.Fatalf(
			"Unexpected %s: %v\nExpected: %v",
			info, currentValue, expectedValue,
		)
	}
}

func requireNoError(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func requireOk(t *testing.T, ok bool, err error) {
	if !ok {
		t.Fatal(err)
	}
}

func requireNotNil(t *testing.T, value any, info string) {
	if value == nil {
		t.Fatalf(
			"%s is <nil>",
			helper.Capitalize(info),
		)
	}
}