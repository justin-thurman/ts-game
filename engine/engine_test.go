package engine

import (
	"bufio"
	"bytes"
	"strings"
	"testing"
)

func TestServerRuns(t *testing.T) {
	server := New()
	go server.Start()
	input := "TestPlayer\n"
	r := strings.NewReader(input)
	var output bytes.Buffer
	w := bufio.NewWriter(&output)
	server.Connect(r, w, func() {})

	w.Flush()

	expected := "Welcome to my very professional game, TestPlayer!"
	if !strings.Contains(output.String(), expected) {
		t.Errorf("expected %q to be in output %q", expected, output.String())
	}
}
