
{{define "where"}}
 and `name`="test"
{{end}}

{{define "list"}}
    select * from `t` where 1=1 {{template "where" .}}
{{end}}