
# Basic Syntax

Notedown flavored Markdown (NFM) is an optionated Markdown subset focused on readability and symantic meaning rather than rendering to HTML. It features several extensions to the Commonmark spec (tasks, frontmatter, etc.) the documentation of which can be found [here](./extended.md).

NFM supports all Commonmark Markdown features except [Link Reference Definitions](https://spec.commonmark.org/0.31.2/#link-reference-definitions) and lazy continuation (both lists and block quotes).

## Headings

Headings are created by starting a line with a series of `#` followed by a space and then a word or phrase. The number of `#` used corresponds to the heading level. Heading level cannot exceed six and they cannot span multiple lines. Although not required it is considered best practice to put a blank line before and after a heading.

```
# Level 1

## Level 2 (child of level 1)

### Level 3 (child of level 2)

#### Level 4 (child of level 3)

##### Level 5 (child of level 4)

###### Level 6 (child of level 5)

### Second level 3 (child of level 2)
```

Setext headings are also supported. Any number of `=` or `-` characters can be used as the underline.

```
Level 1
=======

Level 2
-------
```

## Paragraphs

Paragraphs are the default way text is interpreted. Use a blank line to separate paragraphs. Indentation is allowed but is discouraged.

```
This is a paragraph.

This is another paragraph.

  So is this but try to avoid it.
```

## Lists



## Thematic Breaks

Thematic breaks are created with three or more dashes (`---`), asterisks (`***`), or underscores (`___`) on a line. There can be up to three spaces before the first dash/asterisk/underscore and any number inbetween. No other characters are allowed on the line.

```
---

 *   *   *

   __  __  __  __
```

## Code Blocks

Code blocks can be created either by indenting every line of the block by four spaces.

```
    func main() {
        fmt.Println("Notedown!")
    }
```

Or by wrapping it in three backticks (```) or tildes (`~~~`). No indentation required!

```
func main() {
    fmt.Println("Notedown!")
}
```

Fenced code blocks can optionally contain an infostring after the opening fence. Typically, this is used to denote the language of the code for syntax highlighting.

```go
func main() {
    fmt.Println("Notedown!")
}
```

## Block Quotes

Block quotes are created by adding a `>` followed by a space at the start of a line.

```
> This is a block quote
```

To add two paragraphs in the same quote add a `>` to the blank line between paragraphs.

```
> This is a paragraph
>
> This is a second paragraph, but part of the same block quote.
```

Block quotes can be nested in block quotes. The level of nesting is determined by the number of `>`.

```
> Block quote
>> Nested block quote
> Back to the top level block quote
```

Other elements can be nested in block quotes.

```
> #### Heading
>
> - List item
> - List item 2
```

Unlike Markdown, Notedown does not support lazy continuation. Lazy continuation is somewhat unintuitive for users and not supporting it allows the parser implementation to be much simpler.

```
> block quote
This is now outside the block quote but markdown would consider this still part of the block quote.
```

