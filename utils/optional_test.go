package utils_test

import (
	"testing"

	"github.com/chris-pikul/go-prql/utils"
)

func TestOptionalDefault(t *testing.T) {
	opt := utils.Optional[int]{}

	// Check Value
	if opt.Value() != nil {
		t.Error("expected default optional Value() to be nil")
	}

	// Check OK
	if opt.Ok() {
		t.Error("expected default optional Ok() to be false")
	}

	// Check Get()
	if val, ok := opt.Get(); val != nil || ok {
		t.Error("expected default optional Get() to be nil, false")
	}
}

func TestOptionalNewExplicit(t *testing.T) {
	opt := utils.NewOptional[int](nil, true)

	if val, ok := opt.Get(); val != nil || !ok {
		if val != nil {
			t.Error("NewOptional should have nil value")
		}

		if !ok {
			t.Error("NewOptional was explicit, should be ok")
		}
	}
}

func TestOptionalNewImplied(t *testing.T) {
	val := 123
	opt := utils.NewOptional(&val, false)

	if val, ok := opt.Get(); val == nil || !ok {
		if val != nil {
			t.Error("NewOptional should not have nil value with implicit")
		}

		if !ok {
			t.Error("NewOptional was implicit, should be ok")
		}
	}
}

func TestOptionalSet(t *testing.T) {
	val := 123
	opt := utils.NewOptional(&val, true)

	// Change the root value
	val += 123
	if val, ok := opt.Get(); *val != 246 || !ok {
		if !ok {
			t.Error("Optional should still be ok")
		}
		if *val != 246 {
			t.Error("Optional had it's underlying object change, but not reflected")
		}
	}

	newVal := 321
	opt.Set(&newVal)
	if opt.Value() != &newVal {
		t.Error("Optional had it's value Set but did not reflect the change")
	}

	nilOpt := utils.NewOptional[int](nil, false)
	nilOpt.Set(nil)
	if !nilOpt.Ok() {
		t.Error("Optional.Set() should enforce the Ok flag")
	}
}

func TestOptionalClear(t *testing.T) {
	val := 123
	opt := utils.NewOptional(&val, true)

	opt.Clear()
	if opt.Value() != nil {
		t.Error("Optional.Clear() should set value to nil")
	}
	if opt.Ok() {
		t.Error("Optional.Clear() should clear Ok flag")
	}
}
