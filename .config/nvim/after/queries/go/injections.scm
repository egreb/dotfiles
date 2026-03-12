; extends

; SQL injection in Go string literals

; regex-based: match common SQL patterns
([
  (interpreted_string_literal_content)
  (raw_string_literal_content)
] @injection.content
  (#match? @injection.content "(SELECT|select|INSERT|insert|UPDATE|update|DELETE|delete).+(FROM|from|INTO|into|VALUES|values|SET|set).*(WHERE|where|GROUP BY|group by)?")
  (#set! injection.language "sql"))

; keyword-based: match SQL keywords and comments
([
  (interpreted_string_literal_content)
  (raw_string_literal_content)
] @injection.content
  (#contains? @injection.content "-- sql" "--sql" "ADD CONSTRAINT" "ALTER TABLE" "ALTER COLUMN"
    "DATABASE" "FOREIGN KEY" "GROUP BY" "HAVING" "CREATE INDEX" "INSERT INTO"
    "NOT NULL" "PRIMARY KEY" "UPDATE SET" "TRUNCATE TABLE" "LEFT JOIN" "WITH"
    "add constraint" "alter table" "alter column" "database" "foreign key" "group by"
    "having" "create index" "insert into" "not null" "primary key" "update set"
    "truncate table" "left join")
  (#set! injection.language "sql"))
