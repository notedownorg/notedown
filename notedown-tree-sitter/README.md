# Tree-sitter Grammar for Notedown Flavored Markdown

This is a Tree-sitter grammar for Notedown Flavored Markdown (NFM), focusing on the core syntax elements including wikilinks.

## Features

- **Headings**: Standard markdown headings (`#`, `##`, etc.)
- **Wikilinks**: Notedown-specific wikilink syntax
  - Simple: `[[target]]`
  - With display text: `[[target|display text]]`
- **Text**: Basic text content

## Building

```bash
npm run generate
```

## Testing

```bash
npm test
```

## Installation for Neovim

1. Build the grammar:
   ```bash
   npm run generate
   ```

2. Copy to Neovim's parser directory:
   ```bash
   cp -r . ~/.local/share/nvim/site/pack/packer/start/nvim-treesitter/parser/notedown/
   ```

3. Add to your Neovim Tree-sitter configuration:
   ```lua
   require'nvim-treesitter.configs'.setup {
     ensure_installed = { "notedown" },
     highlight = {
       enable = true,
       additional_vim_regex_highlighting = false,
     },
   }
   ```

## Grammar Structure

The grammar currently supports:

- `document`: Root node containing all content
- `heading`: Markdown headings with level and content
- `wikilink`: Notedown wikilinks with target and optional display
- `text`: Regular text content

## Syntax Highlighting

The grammar includes highlight queries for Neovim that provide:

- Heading markers (`#`, `##`, etc.)
- Heading content
- Wikilink brackets (`[[`, `]]`)
- Wikilink separators (`|`)
- Wikilink targets and display text

## Development

To extend the grammar:

1. Edit `grammar.js`
2. Run `npm run generate`
3. Add tests to `test/corpus/`
4. Run `npm test`
5. Update highlight queries in `queries/highlights.scm`