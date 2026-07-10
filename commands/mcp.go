package commands

import (
	"github.com/njayp/ophis"
	"github.com/spf13/cobra"
)

// excludedFromMCP are command-name substrings kept out of the MCP tool surface: setup/meta
// commands an agent should not drive, and the raw `api` escape hatch (which would bypass
// the per-command read-only/write/destructive annotations). The `mcp` and `agent` subtrees
// are excluded too so an agent can neither re-enter the server nor disable its own
// guardrails.
var excludedFromMCP = []string{
	"agent", "auth", "config", "alias", "init", "doctor", "completion", "version", "api",
}

// secretFlags must never reach the MCP tool schema: an agent must not read the token,
// switch profiles, retarget the instance, or hop accounts (DECISIONS.md #14). The server
// uses whatever profile/account is active at startup.
var secretFlags = []string{"show-token", "profile", "base-url", "account-id"}

func init() {
	metaRegistrars = append(metaRegistrars, func(_ *deps) *cobra.Command {
		// ophis walks the command tree and exposes each runnable leaf as an MCP tool,
		// replaying the cobra command on invocation so tools reuse the same client, keyring,
		// and profile.
		return ophis.Command(&ophis.Config{
			ToolNamePrefix: "cw",
			Selectors: []ophis.Selector{{
				CmdSelector:           ophis.ExcludeCmdsContaining(excludedFromMCP...),
				InheritedFlagSelector: ophis.ExcludeFlags(secretFlags...),
			}},
		})
	})
}
