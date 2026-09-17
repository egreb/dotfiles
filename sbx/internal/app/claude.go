package app

import (
	"context"
	"fmt"

	"sbx/internal/process"
)

// Runs only during creation, before Claude can read or write its config.
// Theme is a user preference in .claude.json, not .claude/settings.json.
const initializeClaudeThemeScript = `import json, os, pathlib, sys, tempfile
directory = pathlib.Path(os.environ.get("CLAUDE_CONFIG_DIR") or pathlib.Path.home())
path = directory / ".claude.json"
directory.mkdir(parents=True, exist_ok=True)
data = json.loads(path.read_text()) if path.exists() else {}
data["theme"] = sys.argv[1]
fd, temporary = tempfile.mkstemp(prefix=".claude-theme-", dir=directory)
try:
    with os.fdopen(fd, "w") as output:
        json.dump(data, output, indent=2)
        output.write("\n")
    os.replace(temporary, path)
finally:
    if os.path.exists(temporary):
        os.unlink(temporary)
`

func (app *App) initializeClaudeTheme(ctx context.Context, name string) error {
	args := []string{"exec", name, "--", "python3", "-c", initializeClaudeThemeScript, app.config.ClaudeTheme}
	if err := app.runner.Run(ctx, app.native, args, process.Options{}); err != nil {
		return fmt.Errorf("cannot initialize Claude theme in sandbox %s: %w", name, err)
	}
	return nil
}
