# Wikilink Diagnostics and Code Actions Test

This workspace demonstrates Notedown's comprehensive diagnostic system and disambiguation code actions.

## Broken Wikilinks (Should show error diagnostics)

The following wikilinks point to non-existent files:

- [[missing-file]] - File doesn't exist, should show broken wikilink diagnostic
- [[nonexistent-document]] - Another missing file
- [[invalid/path]] - Path that doesn't exist

## Ambiguous Wikilinks (Should show disambiguation options)

These wikilinks match multiple files and can be resolved with code actions:

- [[meeting]] - Matches current and archived meetings (can be disambiguated)
- [[project]] - Matches multiple project files
- [[api]] - Matches different API files

## Valid Wikilinks (Should be clean)

These should show no diagnostic errors:

- [[valid-document]] - Points to existing file
- [[docs/guide]] - Valid path-based wikilink

## Real-time Testing Area

Add new wikilinks below to test real-time diagnostic updates:
