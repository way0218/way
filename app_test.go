package way

import (
	"testing"
	"time"

	"github.com/way0218/way/transport/http"
)

func TestApp(t *testing.T) {
	hs := http.NewServer()
	app := New(Name("test"), Version("1.0.0"), Server(hs))
	time.AfterFunc(time.Second, func() {
		_ = app.Stop()
	})
	if err := app.Run(); err != nil {
		t.Fatal(err)
	}
}
