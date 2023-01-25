package tinami

import "testing"

func Test_GetTinamiArtwork(t *testing.T) {
	s := New()
	artwork, err := s.Get("1110215")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%+v\n", artwork)

}
