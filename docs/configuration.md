# Configuration

Notedown can be configured at the program level or per-workspace. To reduce complexity settings are either available at the program level or the workspace level never both.

## Program Level Configuration

Program-level configuration is located at `$HOME/.config/notedown/config.yaml`.

### Workspaces

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

## Workspace Level Configuration

Workspace-level configuration is located at `${WORKSPACE_ROOT}/.config/notedown.yaml`

### Sources

Sources are external content (articles, videos, podcasts, etc.) added to your workspace. The only thing that can be configured is the directory they are added to.

```yaml
sources:
    default_directory: sources
```


