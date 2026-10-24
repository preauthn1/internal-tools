package util

import (
	"testing"
)

func TestNetworkFilter_CheckIPContains(t *testing.T) {
	nf := NewNetworkFilter("0.0.0.0,0.0.0.0,0.0.0.0/24")
	notContains := []string{"192", "0.0.0.0/24", "0.0.0.0", "abc", "0.0.0.0"}
	for _, rAddr := range notContains {
		if nf.CheckIPContains(rAddr) {
			t.Errorf("[%s] check wrong type ip format but contains.", rAddr)
		}
	}
	isContains := []string{"0.0.0.0", "0.0.0.0", "0.0.0.0", "0.0.0.0"}
	for _, rAddr := range isContains {
		if !nf.CheckIPContains(rAddr) {
			t.Errorf("[%s] check ip is not contains.", rAddr)
		}
	}
}
