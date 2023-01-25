package nico

import "testing"

func Test_GetNicoArtwork(t *testing.T) {
	s := New()
	artwork, err := s.Get("im11115799")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%+v\n", artwork)
}
