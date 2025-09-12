package compute

import "testing"

func TestParseSET(t *testing.T) {
	cmd, err := Parse("SET key value")
	if err != nil {
		t.Fatal("unexpected error", err)
	}
	if cmd.Type != Set {
		t.Errorf("unexpected SET, got %gv", cmd.Type)
	}
	if len(cmd.Args) != 2 || cmd.Args[0] != "key" || cmd.Args[1] != "value" {
		t.Errorf("unexpected args: %v", cmd.Args)
	}
}

func TestParseGET(t *testing.T) {
	cmd, err := Parse("GET mykey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Type != Get {
		t.Errorf("expected GET, got %v", cmd.Type)
	}
	if cmd.Args[0] != "mykey" {
		t.Errorf("expected arg 'mykey', got %v", cmd.Args[0])
	}
}

func TestParseDEL(t *testing.T) {
	cmd, err := Parse("DEL mykey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cmd.Type != Del {
		t.Errorf("expected DEL, got %v", cmd.Type)
	}
	if cmd.Args[0] != "mykey" {
		t.Errorf("expected arg 'mykey', got %v", cmd.Args[0])
	}
}

func TestParseErrors(t *testing.T) {
	tests := []string{
		"",
		"UNKNOWN x",
		"SET onlyone",
		"GET",
		"DEL",
	}

	for _, input := range tests {
		if _, err := Parse(input); err == nil {
			t.Errorf("expected error for input: %q", input)
		}
	}
}
