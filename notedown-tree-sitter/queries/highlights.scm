; Headings - based on actual AST structure
(heading) @markup.heading
(heading (text) @markup.heading.content)

; Wikilinks - based on actual AST structure  
(wikilink "[[" @punctuation.bracket)
(wikilink "]]" @punctuation.bracket)
(wikilink "|" @punctuation.separator)
(wikilink_target) @markup.link.url
(wikilink_display) @markup.link.label

; Regular text
(text) @text