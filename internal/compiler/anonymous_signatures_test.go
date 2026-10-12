package compiler

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestAnonymousMethodSignatureTypes(t *testing.T) {
    src := `package main
class Holder {
    Wrap(input struct { Name string }) struct {
        Name string
        Meta struct { OK bool }
    } {
        return struct {
            Name string
            Meta struct { OK bool }
        }{Name: input.Name, Meta: struct{ OK bool }{OK: true}}
    }
    Lookup() interface { String() string } { return nil }
    Pair() (struct { Count int }, error) {
        return struct { Count int }{Count: 2}, nil
    }
}
func main() {
    h := new Holder()
    _ = h.Wrap(struct { Name string }{Name: "example"})
    _ = h.Lookup()
    _, _ = h.Pair()
}
	`
	out, err := Compile("anonymous.goop", []byte(src))
	if err != nil {
		t.Fatalf("transpile anonymous signatures: %v", err)
	}
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.org/anonymous\n\ngo 1.22\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "anonymous.go"), out, 0o644); err != nil {
		t.Fatal(err)
	}
	cmd := exec.Command("go", "run", ".")
	cmd.Dir = dir
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generated code does not build: %v\n%s", err, output)
	}
}
