package toolv8_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/pushthat/bud/internal/cli/testcli"
	"github.com/pushthat/bud/internal/is"
	"github.com/pushthat/bud/internal/testdir"
)

func TestToolV8(t *testing.T) {
	is := is.New(t)
	ctx := context.Background()
	dir := t.TempDir()
	td := testdir.New(dir)
	cli := testcli.New(dir)
	cli.Stdin = bytes.NewBufferString("2+2")
	result, err := cli.Run(ctx, "tool", "v8")
	is.NoErr(err)
	is.Equal(result.Stderr(), "")
	is.Equal(strings.TrimSpace(result.Stdout()), "4")
	is.NoErr(td.NotExists(
		"bud/internal/app",
		"bud/app",
	))
}
