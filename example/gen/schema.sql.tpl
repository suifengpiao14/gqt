{{define "GetByServiceID"}}
select * from `schema` where `service_id`=:ServiceID and `deleted_at` is null;
{{end}}

{{define "GetBySchemaID"}}
select * from `schema` where `schema_id`=:SchemaID and `deleted_at` is null;
{{end}}

{{define "GetAllBySchemaIDList"}}
select * from `schema` where `schema_id` in ({{in . .SchemaIDList}}) and `deleted_at` is null;
{{end}}

{{define "PaginateWhere"}}
{{end}}

{{define "PaginateTotal"}}
select count(*) as `count` from `schema` where 1=1 {{template "PaginateWhere" .}} and `deleted_at` is null;
{{end}}

{{define "Paginate"}}
select * from `schema` where 1=1 {{template "PaginateWhere" .}} and `deleted_at` is null limit :Offset,:Limit ;
{{end}}

{{define "Insert"}}
insert into `schema` (`schema_id`,`service_id`,`description`,`remark`,`type`,`example`,`deprecated`,`required`,`enum`,`enum_names`,`enum_titles`,`format`,`default`,`nullable`,`multiple_of`,`maxnum`,`exclusive_maximum`,`minimum`,`exclusive_minimum`,`max_length`,`min_length`,`pattern`,`max_items`,`min_items`,`unique_items`,`max_properties`,`min_properties`,`all_of`,`one_of`,`any_of`,`allow_empty_value`,`allow_reserved`,`not`,`additional_properties`,`discriminator`,`read_only`,`write_only`,`xml`,`external_docs`,`external_pros`,`extensions`,`summary`)values
(:SchemaID,:ServiceID,:Description,:Remark,:Type,:Example,:Deprecated,:Required,:Enum,:EnumNames,:EnumTitles,:Format,:Default,:Nullable,:MultipleOf,:Maxnum,:ExclusiveMaximum,:Minimum,:ExclusiveMinimum,:MaxLength,:MinLength,:Pattern,:MaxItems,:MinItems,:UniqueItems,:MaxProperties,:MinProperties,:AllOf,:OneOf,:AnyOf,:AllowEmptyValue,:AllowReserved,:Not,:AdditionalProperties,:Discriminator,:ReadOnly,:WriteOnly,:XML,:ExternalDocs,:ExternalPros,:Extensions,:Summary);
{{end}}

{{define "Update"}}
{{$preComma:=newPreComma}}
update `schema` set {{if .SchemaID}} {{$preComma.PreComma}} `schema_id`=:SchemaID {{end}}
{{if .ServiceID}} {{$preComma.PreComma}} `service_id`=:ServiceID {{end}}
{{if .Description}} {{$preComma.PreComma}} `description`=:Description {{end}}
{{if .Remark}} {{$preComma.PreComma}} `remark`=:Remark {{end}}
{{if .Type}} {{$preComma.PreComma}} `type`=:Type {{end}}
{{if .Example}} {{$preComma.PreComma}} `example`=:Example {{end}}
{{if .Deprecated}} {{$preComma.PreComma}} `deprecated`=:Deprecated {{end}}
{{if .Required}} {{$preComma.PreComma}} `required`=:Required {{end}}
{{if .Enum}} {{$preComma.PreComma}} `enum`=:Enum {{end}}
{{if .EnumNames}} {{$preComma.PreComma}} `enum_names`=:EnumNames {{end}}
{{if .EnumTitles}} {{$preComma.PreComma}} `enum_titles`=:EnumTitles {{end}}
{{if .Format}} {{$preComma.PreComma}} `format`=:Format {{end}}
{{if .Default}} {{$preComma.PreComma}} `default`=:Default {{end}}
{{if .Nullable}} {{$preComma.PreComma}} `nullable`=:Nullable {{end}}
{{if .MultipleOf}} {{$preComma.PreComma}} `multiple_of`=:MultipleOf {{end}}
{{if .Maxnum}} {{$preComma.PreComma}} `maxnum`=:Maxnum {{end}}
{{if .ExclusiveMaximum}} {{$preComma.PreComma}} `exclusive_maximum`=:ExclusiveMaximum {{end}}
{{if .Minimum}} {{$preComma.PreComma}} `minimum`=:Minimum {{end}}
{{if .ExclusiveMinimum}} {{$preComma.PreComma}} `exclusive_minimum`=:ExclusiveMinimum {{end}}
{{if .MaxLength}} {{$preComma.PreComma}} `max_length`=:MaxLength {{end}}
{{if .MinLength}} {{$preComma.PreComma}} `min_length`=:MinLength {{end}}
{{if .Pattern}} {{$preComma.PreComma}} `pattern`=:Pattern {{end}}
{{if .MaxItems}} {{$preComma.PreComma}} `max_items`=:MaxItems {{end}}
{{if .MinItems}} {{$preComma.PreComma}} `min_items`=:MinItems {{end}}
{{if .UniqueItems}} {{$preComma.PreComma}} `unique_items`=:UniqueItems {{end}}
{{if .MaxProperties}} {{$preComma.PreComma}} `max_properties`=:MaxProperties {{end}}
{{if .MinProperties}} {{$preComma.PreComma}} `min_properties`=:MinProperties {{end}}
{{if .AllOf}} {{$preComma.PreComma}} `all_of`=:AllOf {{end}}
{{if .OneOf}} {{$preComma.PreComma}} `one_of`=:OneOf {{end}}
{{if .AnyOf}} {{$preComma.PreComma}} `any_of`=:AnyOf {{end}}
{{if .AllowEmptyValue}} {{$preComma.PreComma}} `allow_empty_value`=:AllowEmptyValue {{end}}
{{if .AllowReserved}} {{$preComma.PreComma}} `allow_reserved`=:AllowReserved {{end}}
{{if .Not}} {{$preComma.PreComma}} `not`=:Not {{end}}
{{if .AdditionalProperties}} {{$preComma.PreComma}} `additional_properties`=:AdditionalProperties {{end}}
{{if .Discriminator}} {{$preComma.PreComma}} `discriminator`=:Discriminator {{end}}
{{if .ReadOnly}} {{$preComma.PreComma}} `read_only`=:ReadOnly {{end}}
{{if .WriteOnly}} {{$preComma.PreComma}} `write_only`=:WriteOnly {{end}}
{{if .XML}} {{$preComma.PreComma}} `xml`=:XML {{end}}
{{if .ExternalDocs}} {{$preComma.PreComma}} `external_docs`=:ExternalDocs {{end}}
{{if .ExternalPros}} {{$preComma.PreComma}} `external_pros`=:ExternalPros {{end}}
{{if .Extensions}} {{$preComma.PreComma}} `extensions`=:Extensions {{end}}
{{if .Summary}} {{$preComma.PreComma}} `summary`=:Summary {{end}} where `schema_id`=:SchemaID;
{{end}}

{{define "Del"}}
update `schema` set `deleted_at`={{currentTime .}} where `schema_id`=:SchemaID;
{{end}}
