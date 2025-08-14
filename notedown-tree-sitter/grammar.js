// Notedown Tree-sitter Grammar
// Enhanced markdown grammar with wikilink support
// Focus on core functionality with good syntax highlighting

module.exports = grammar({
  name: 'notedown',

  rules: {
    document: $ => repeat($._content),

    _content: $ => choice(
      $.atx_heading,
      $.wikilink,
      $.emphasis,
      $.strong,
      $.strikethrough,
      $.code_span,
      $.link,
      $.image,
      $.fenced_code_block,
      $.text,
      '\n'
    ),

    // ATX Headings - following standard treesitter markdown conventions
    atx_heading: $ => prec.right(seq(
      choice(
        $.atx_h1_marker,
        $.atx_h2_marker,
        $.atx_h3_marker,
        $.atx_h4_marker,
        $.atx_h5_marker,
        $.atx_h6_marker
      ),
      /\s+/,
      repeat(choice($.text, $.wikilink, $.emphasis, $.strong, $.code_span, $.link))
    )),

    atx_h1_marker: $ => '#',
    atx_h2_marker: $ => '##',
    atx_h3_marker: $ => '###',
    atx_h4_marker: $ => '####',
    atx_h5_marker: $ => '#####',
    atx_h6_marker: $ => '######',

    // Wikilinks - our key feature
    wikilink: $ => choice(
      // [[target]]
      seq('[[', field('target', $.wikilink_target), ']]'),
      // [[target|display]]
      seq('[[', field('target', $.wikilink_target), '|', field('display', $.wikilink_display), ']]')
    ),

    wikilink_target: $ => /[^\|\]]+/,
    wikilink_display: $ => /[^\]]+/,

    // Emphasis and strong
    emphasis: $ => choice(
      seq('*', repeat1(choice($.text, $.strong, $.code_span, $.wikilink)), '*'),
      seq('_', repeat1(choice($.text, $.strong, $.code_span, $.wikilink)), '_')
    ),

    strong: $ => choice(
      seq('**', repeat1(choice($.text, $.emphasis, $.code_span, $.wikilink)), '**'),
      seq('__', repeat1(choice($.text, $.emphasis, $.code_span, $.wikilink)), '__')
    ),

    strikethrough: $ => seq(
      '~~',
      repeat1(choice($.text, $.emphasis, $.strong, $.code_span, $.wikilink)),
      '~~'
    ),

    // Code spans
    code_span: $ => seq('`', field('code', /[^`]+/), '`'),

    // Fenced code blocks
    fenced_code_block: $ => seq(
      choice('```', '~~~'),
      optional(field('language', /[^\n]+/)),
      '\n',
      field('code', /[^`~]*/),
      choice('```', '~~~')
    ),

    // Links
    link: $ => seq(
      '[',
      field('text', /[^\]]+/),
      ']',
      '(',
      field('destination', /[^\)]+/),
      ')'
    ),

    // Images
    image: $ => seq(
      '!',
      '[',
      field('alt', /[^\]]*/),
      ']',
      '(',
      field('destination', /[^\)]+/),
      ')'
    ),

    // Basic text - avoid conflicts by being more specific
    text: $ => /[^\[\]#\*_`~!\n]+/
  }
});