package log

import "testing"

func TestFormat(t *testing.T) {
	if r, e := format("upload", "completed"), "upload\tcompleted"; r != e {
		t.Errorf("Expecting %v, got %v", e, r)
	}
	if r, e := format("upload", "completed", "filename", "a.png"), "upload\tcompleted\tfilename=a.png"; r != e {
		t.Errorf("Expecting %v, got %v", e, r)
	}
}
