package pixiv

import (
	"github.com/Netflix/go-env"
	"github.com/faryne/api-server/config"
	"github.com/joho/godotenv"
	"testing"
)

func Test_GetPixivArtwork(t *testing.T) {
	_ = godotenv.Load("../../../.env")
	env.UnmarshalFromEnviron(&config.Config)

	s := New()
	// 92817663 - 一般向
	// 94937757 - R-18 多圖
	// 104001276 - R-18 多圖
	artwork, err := s.Get("104001276")
	if err != nil {
		t.Fatalf(err.Error())
	}
	t.Logf("%+v\n", artwork)
}
