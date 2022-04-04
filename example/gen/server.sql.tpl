{{define "GetByServiceID"}}
select * from `server` where `service_id`=:ServiceID and `deleted_at` is null;
{{end}}

{{define "GetByServerID"}}
select * from `server` where `server_id`=:ServerID and `deleted_at` is null;
{{end}}

{{define "GetAllByServerIDList"}}
select * from `server` where `server_id` in ({{in . .ServerIDList}}) and `deleted_at` is null;
{{end}}

{{define "PaginateWhere"}}
{{end}}

{{define "PaginateTotal"}}
select count(*) as `count` from `server` where 1=1 {{template "PaginateWhere" .}} and `deleted_at` is null;
{{end}}

{{define "Paginate"}}
select * from `server` where 1=1 {{template "PaginateWhere" .}} and `deleted_at` is null limit :Offset,:Limit ;
{{end}}

{{define "Insert"}}
insert into `server` (`server_id`,`service_id`,`url`,`description`,`proxy`,`extension_ids`)values
(:ServerID,:ServiceID,:URL,:Description,:Proxy,:ExtensionIds);
{{end}}

{{define "Update"}}
{{$preComma:=newPreComma}}
update `server` set {{if .ServerID}} {{$preComma.PreComma}} `server_id`=:ServerID {{end}}
{{if .ServiceID}} {{$preComma.PreComma}} `service_id`=:ServiceID {{end}}
{{if .URL}} {{$preComma.PreComma}} `url`=:URL {{end}}
{{if .Description}} {{$preComma.PreComma}} `description`=:Description {{end}}
{{if .Proxy}} {{$preComma.PreComma}} `proxy`=:Proxy {{end}}
{{if .ExtensionIds}} {{$preComma.PreComma}} `extension_ids`=:ExtensionIds {{end}} where `server_id`=:ServerID;
{{end}}

{{define "Del"}}
update `server` set `deleted_at`={{currentTime .}} where `server_id`=:ServerID;
{{end}}
