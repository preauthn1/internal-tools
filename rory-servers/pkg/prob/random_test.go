package prob

import (
	"math"
	"math/rand"
	"testing"
	"time"
)

func TestRandom(t *testing.T) {

	// rand.Seed(int64(time.Now().Nanosecond()))

	r := new(Random)

	for i := 0; i < 100; i++ {
		r.Add(i, 1)
	}

	for i := 0; i < 100; i++ {
		t.Log(r.Choose().(int))
	}

}

func TestRandom_Case(t *testing.T) {

	rand.Seed(int64(time.Now().Nanosecond()))

	r := new(Random)

	r.Add(true, 4998000)
	r.Add(false, 9995002000)

	for i := 0; i < 10; i++ {
		t.Logf("%v, %+v\n", r.Choose(), r)
	}

	t.Log(math.MaxInt64)
}
