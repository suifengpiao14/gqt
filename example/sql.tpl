
{{define "where"}}
 and `name`in ({{in . .IDS}})
{{end}}

{{define "list"}}
    select * from `t` where 1=1 {{template "where" .}}
{{end}}