{{define "GetByServiceID"}}
select * from `example` where `service_id`=:ServiceID and `deleted_at` is null;
{{end}}

{{define "GetByExampleID"}}
select * from `example` where `example_id`=:ExampleID and `deleted_at` is null;
{{end}}

{{define "GetAllByExampleIDList"}}
select * from `example` where `example_id` in ({{in . .ExampleIDList}}) and `deleted_at` is null;
{{end}}

{{define "PaginateWhere"}}
{{end}}

{{define "PaginateTotal"}}
select count(*) as `count` from `example` where 1=1 {{template "PaginateWhere" .}} and `deleted_at` is null;
{{end}}

{{define "Paginate"}}
select * from `example` where 1=1 {{template "PaginateWhere" .}} and `deleted_at` is null limit :Offset,:Limit ;
{{end}}

{{define "Insert"}}
insert into `example` (`example_id`,`service_id`,`api_id`,`tag`,`title`,`summary`,`url`,`method`,`pre_request_script`,`auth`,`headers`,`parameters`,`content_type`,`body`,`test_script`,`response`)values
(:ExampleID,:ServiceID,:APIID,:Tag,:Title,:Summary,:URL,:Method,:PreRequestScript,:Auth,:Headers,:Parameters,:ContentType,:Body,:TestScript,:Response);
{{end}}

{{define "Update"}}
{{$preComma:=newPreComma}}
update `example` set {{if .ExampleID}} {{$preComma.PreComma}} `example_id`=:ExampleID {{end}}
{{if .ServiceID}} {{$preComma.PreComma}} `service_id`=:ServiceID {{end}}
{{if .APIID}} {{$preComma.PreComma}} `api_id`=:APIID {{end}}
{{if .Tag}} {{$preComma.PreComma}} `tag`=:Tag {{end}}
{{if .Title}} {{$preComma.PreComma}} `title`=:Title {{end}}
{{if .Summary}} {{$preComma.PreComma}} `summary`=:Summary {{end}}
{{if .URL}} {{$preComma.PreComma}} `url`=:URL {{end}}
{{if .Method}} {{$preComma.PreComma}} `method`=:Method {{end}}
{{if .PreRequestScript}} {{$preComma.PreComma}} `pre_request_script`=:PreRequestScript {{end}}
{{if .Auth}} {{$preComma.PreComma}} `auth`=:Auth {{end}}
{{if .Headers}} {{$preComma.PreComma}} `headers`=:Headers {{end}}
{{if .Parameters}} {{$preComma.PreComma}} `parameters`=:Parameters {{end}}
{{if .ContentType}} {{$preComma.PreComma}} `content_type`=:ContentType {{end}}
{{if .Body}} {{$preComma.PreComma}} `body`=:Body {{end}}
{{if .TestScript}} {{$preComma.PreComma}} `test_script`=:TestScript {{end}}
{{if .Response}} {{$preComma.PreComma}} `response`=:Response {{end}} where `example_id`=:ExampleID;
{{end}}

{{define "Del"}}
update `example` set `deleted_at`={{currentTime .}} where `example_id`=:ExampleID;
{{end}}
