module.exports = grammar({
  name: 'notedown',

  rules: {
    document: $ => repeat($._content),

    _content: $ => choice(
      $.heading,
      $.wikilink,
      $.text,
      '\n'
    ),

    // Headings
    heading: $ => prec.right(seq(
      field('level', repeat1('#')),
      /\s+/,
      field('content', repeat(choice($.text, $.wikilink)))
    )),

    // Wikilinks - Notedown-specific syntax
    wikilink: $ => choice(
      // [[target]]
      seq('[[', field('target', $.wikilink_target), ']]'),
      // [[target|display]]
      seq('[[', field('target', $.wikilink_target), '|', field('display', $.wikilink_display), ']]')
    ),

    wikilink_target: $ => /[^\|\]]+/,
    wikilink_display: $ => /[^\]]+/,

    // Basic text
    text: $ => /[^\[\]#\n]+/
  }
});