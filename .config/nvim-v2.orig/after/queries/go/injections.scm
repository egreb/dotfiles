; extends

; Detect SQL statements after whitespace and optional leading SQL line comments.
; A bare -- on the first line is common in multiline Go queries.
; Set magic mode explicitly: #match? otherwise defaults to very-magic regexes.
([
  (interpreted_string_literal_content)
  (raw_string_literal_content)
] @injection.content
  (#match? @injection.content "\\m\\c^[ \t\r\n]*\\(--[^\n]*\n[ \t\r\n]*\\)*\\(SELECT\\|INSERT\\|UPDATE\\|DELETE\\|WITH\\|CREATE\\|ALTER\\|DROP\\|TRUNCATE\\)[ \t\r\n]")
  (#set! injection.language "sql"))

; An explicit marker also supports SQL fragments and leading comments.
([
  (interpreted_string_literal_content)
  (raw_string_literal_content)
] @injection.content
  (#any-contains? @injection.content "-- sql" "--sql")
  (#set! injection.language "sql"))
