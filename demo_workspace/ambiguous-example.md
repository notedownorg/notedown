# Ambiguous Wikilink Example

This document demonstrates how ambiguous wikilinks are detected and flagged.

## The Problem

When you write [[config]], which file should it link to?

- `/config.md` - The basic configuration overview
- `/docs/config.md` - The advanced configuration guide

This ambiguity should trigger a diagnostic warning in your editor.

## How to Fix

To resolve the ambiguity, be more specific:

- Use [[docs/config]] for the advanced configuration
- Or rename one of the files to avoid the conflict (e.g., `basic-config.md` and `advanced-config.md`)

## Other Examples

These links are NOT ambiguous:

- [[README]] - Only one README file exists
- [[docs/architecture]] - Full path is specific
- [[notes/ideas]] - Full path is specific
- [[future-document]] - Non-existent files are perfectly valid