package memcached

import (
	"testing"
	"time"
)

func Test_Memcached(t *testing.T) {
	m := New(Localhost)

	var key = "abc"
	var d = time.Second * 10
	if err := m.Set(key, []byte(key), d); err != nil {
		t.Fatalf(err.Error())
	}
	o, err := m.Get(key)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if string(o) != key {
		t.Fatalf("not equal")
	}
	if err := m.Delete(key); err != nil {
		t.Fatalf(err.Error())
	}
}
