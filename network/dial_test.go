package network

import (
	"testing"
)

func TestDial(t *testing.T) {
	err := Dial("rtmps", "a.rtmps.youtube.com:1935")
	if err != nil {
		t.Error(err)
	}

	err = Dial("rtmp", "a.rtmp.youtube.com:1935")
	if err != nil {
		t.Error(err)
	}
}
