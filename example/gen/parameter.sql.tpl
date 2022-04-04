{{define "GetByServiceID"}}
select * from `parameter` where `service_id`=:ServiceID and `deleted_at` is null;
{{end}}

{{define "GetByParameterID"}}
select * from `parameter` where `parameter_id`=:ParameterID and `deleted_at` is null;
{{end}}

{{define "GetAllByParameterIDList"}}
select * from `parameter` where `parameter_id` in ({{in . .ParameterIDList}}) and `deleted_at` is null;
{{end}}

{{define "PaginateWhere"}}
{{end}}

{{define "PaginateTotal"}}
select count(*) as `count` from `parameter` where 1=1 {{template "PaginateWhere" .}} and `deleted_at` is null;
{{end}}

{{define "Paginate"}}
select * from `parameter` where 1=1 {{template "PaginateWhere" .}} and `deleted_at` is null limit :Offset,:Limit ;
{{end}}

{{define "Insert"}}
insert into `parameter` (`parameter_id`,`service_id`,`api_id`,`schema_id`,`full_name`,`name`,`title`,`type`,`tag`,`method`,`http_status`,`position`,`example`,`deprecated`,`required`,`serialize`,`explode`,`allow_empty_value`,`allow_reserved`,`description`)values
(:ParameterID,:ServiceID,:APIID,:SchemaID,:FullName,:Name,:Title,:Type,:Tag,:Method,:HTTPStatus,:Position,:Example,:Deprecated,:Required,:Serialize,:Explode,:AllowEmptyValue,:AllowReserved,:Description);
{{end}}

{{define "Update"}}
{{$preComma:=newPreComma}}
update `parameter` set {{if .ParameterID}} {{$preComma.PreComma}} `parameter_id`=:ParameterID {{end}}
{{if .ServiceID}} {{$preComma.PreComma}} `service_id`=:ServiceID {{end}}
{{if .APIID}} {{$preComma.PreComma}} `api_id`=:APIID {{end}}
{{if .SchemaID}} {{$preComma.PreComma}} `schema_id`=:SchemaID {{end}}
{{if .FullName}} {{$preComma.PreComma}} `full_name`=:FullName {{end}}
{{if .Name}} {{$preComma.PreComma}} `name`=:Name {{end}}
{{if .Title}} {{$preComma.PreComma}} `title`=:Title {{end}}
{{if .Type}} {{$preComma.PreComma}} `type`=:Type {{end}}
{{if .Tag}} {{$preComma.PreComma}} `tag`=:Tag {{end}}
{{if .Method}} {{$preComma.PreComma}} `method`=:Method {{end}}
{{if .HTTPStatus}} {{$preComma.PreComma}} `http_status`=:HTTPStatus {{end}}
{{if .Position}} {{$preComma.PreComma}} `position`=:Position {{end}}
{{if .Example}} {{$preComma.PreComma}} `example`=:Example {{end}}
{{if .Deprecated}} {{$preComma.PreComma}} `deprecated`=:Deprecated {{end}}
{{if .Required}} {{$preComma.PreComma}} `required`=:Required {{end}}
{{if .Serialize}} {{$preComma.PreComma}} `serialize`=:Serialize {{end}}
{{if .Explode}} {{$preComma.PreComma}} `explode`=:Explode {{end}}
{{if .AllowEmptyValue}} {{$preComma.PreComma}} `allow_empty_value`=:AllowEmptyValue {{end}}
{{if .AllowReserved}} {{$preComma.PreComma}} `allow_reserved`=:AllowReserved {{end}}
{{if .Description}} {{$preComma.PreComma}} `description`=:Description {{end}} where `parameter_id`=:ParameterID;
{{end}}

{{define "Del"}}
update `parameter` set `deleted_at`={{currentTime .}} where `parameter_id`=:ParameterID;
{{end}}
