# Configuration

Notedown can be configured at the program level or per-workspace with the more specific (per-workspace) taking precedence when settings can be configured in multiple places.

## Program Level Configuration

Program-level configuration can be found at `$HOME/.notedown/config.yaml`.

### Workspaces (Program-Level Only)

Workspaces tells Notedown where the roots of each of your sources of notes can be found. 

```yaml
# Map of names -> location on disk
workspaces:
    personal:
        location: ~/notes
    work:
        location: ~/worknotes

# Which workspace tools should default to if no override is set
default_workspace: personal
```
