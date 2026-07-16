; extends

; Highlight the non-action text of Go html/template files as HTML.
; combined: parse all text chunks as one document, so tags that span
; template actions (e.g. <title>{{ block ... }}</title>) still work.
((text) @injection.content
  (#set! injection.language "html")
  (#set! injection.combined))
