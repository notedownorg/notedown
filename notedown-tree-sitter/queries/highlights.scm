; ATX Headings - entire heading gets level-specific color
((atx_heading
  (atx_h1_marker)) @markup.heading.1
  (#set! "priority" 105))

((atx_heading
  (atx_h2_marker)) @markup.heading.2
  (#set! "priority" 105))

((atx_heading
  (atx_h3_marker)) @markup.heading.3
  (#set! "priority" 105))

((atx_heading
  (atx_h4_marker)) @markup.heading.4
  (#set! "priority" 105))

((atx_heading
  (atx_h5_marker)) @markup.heading.5
  (#set! "priority" 105))

((atx_heading
  (atx_h6_marker)) @markup.heading.6
  (#set! "priority" 105))

; Wikilinks - our key feature
(wikilink "[[" @punctuation.bracket)
(wikilink "]]" @punctuation.bracket)
(wikilink "|" @punctuation.separator)
(wikilink_target) @markup.link.url
(wikilink_display) @markup.link.label

; Emphasis and strong
(emphasis "*" @markup.emphasis.marker)
(emphasis "_" @markup.emphasis.marker)
(emphasis) @markup.emphasis

(strong "**" @markup.strong.marker)
(strong "__" @markup.strong.marker)  
(strong) @markup.strong

(strikethrough "~~" @markup.strikethrough.marker)
(strikethrough) @markup.strikethrough

; Code
(code_span "`" @punctuation.delimiter)
(fenced_code_block "```" @punctuation.delimiter)
(fenced_code_block "~~~" @punctuation.delimiter)

; Links
(link "[" @punctuation.bracket)
(link "]" @punctuation.bracket)
(link "(" @punctuation.bracket)
(link ")" @punctuation.bracket)

; Images
(image "!" @punctuation.special)
(image "[" @punctuation.bracket)
(image "]" @punctuation.bracket)
(image "(" @punctuation.bracket)
(image ")" @punctuation.bracket)

; Text - lower priority so heading colors take precedence
((text) @text
  (#set! "priority" 90))