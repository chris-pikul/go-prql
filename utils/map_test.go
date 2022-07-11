package utils_test

import (
	"testing"

	"github.com/chris-pikul/go-prql/utils"
)

func TestKeyForValue(t *testing.T) {
	tstMap := map[string]int{
		"first":  1,
		"second": 2,
		"third":  3,
	}

	key, ok := utils.KeyForValue(tstMap, 2)
	if !ok {
		t.Error("expected ok to be true on valid value")
	}
	if *key != "second" {
		t.Errorf("returned incorrect key %s", *key)
	}

	key, ok = utils.KeyForValue(tstMap, -1)
	if ok {
		t.Error("expected ok to be false on invalid value")
	}
	if key != nil {
		t.Error("expected nil value on invalid search value")
	}
}

func TestInvertMap(t *testing.T) {
	tstMap := map[string]int{
		"first":  1,
		"second": 2,
		"third":  3,
	}

	newMap := utils.InvertMap(tstMap)
	for key, val := range newMap {
		if tst, ok := tstMap[val]; !ok || tst != key {
			if !ok {
				t.Error("expected valid key after inversion")
			} else {
				t.Errorf("invalid key-value for %s, expected %d received %d", val, key, tst)
			}
		}
	}
}
