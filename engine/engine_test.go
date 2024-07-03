package engine

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
	"time"
	"ts-game/room"
)

func TestServerRuns(t *testing.T) {
	server := New()
	go func() {
		err := server.Start()
		if err != nil {
			t.Errorf(err.Error())
		}
	}()
	input := "TestPlayer\n"
	r := strings.NewReader(input)
	var output bytes.Buffer
	w := bufio.NewWriter(&output)
	for len(room.Rooms) <= 0 {
		time.Sleep(time.Millisecond * 10)
	}
	server.Connect(r, w, func() {})
	time.Sleep(time.Millisecond * 10)

	w.Flush()

	expected := "Welcome to my very professional game, TestPlayer!"
	if !strings.Contains(output.String(), expected) {
		t.Errorf("expected %q to be in output %q", expected, output.String())
	}
	expected = "Town Center"
	if !strings.Contains(output.String(), expected) {
		t.Errorf("expected %q to be in output %q", expected, output.String())
	}
	expected = "You are standing in the bustling town center."
	if !strings.Contains(output.String(), expected) {
		t.Errorf("expected %q to be in output %q", expected, output.String())
	}
}
