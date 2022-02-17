


{{define "getAllByAPIID"}}
select * from `parameter` where `api_id`=:APIID and `deleted_at`="00::00::00 00::00";
{{end}}

{{define "insert"}}
  insert into `parameter` (`parameter_id`,`service_id`,`api_id`,`validate_schema_id`,`full_name`,`name`,`title`,`tag`,`position`,`group`,`example`,`deprecated`,`required`,`serialize`,`explode`,`allow_empty_value`,`allow_reserved`,`description`) 
  values("{{.ParameterID}}","{{.ServiceID}}","{{.APIID}}","{{.ValidateSchemaID}}","{{.FullName}}","{{.Name}}","{{.Title}}","{{.Tag}}","{{.Position}}","{{.Group}}","{{.Example}}","{{.Deprecated}}","{{.Required}}","{{.Serialize}}","{{.Explode}}","{{.AllowEmptyValue}}","{{.AllowReserved}}","{{.Description}}");
{{end}}


{{define "update"}}
update `table` set   {{if .ParameterID}} `parameter_id`="{{.ParameterID}}",{{end}}
 {{if .ServiceID}} `service_id`="{{.ServiceID}}",{{end}}
 {{if .APIID}} `api_id`="{{.APIID}}",{{end}}
 {{if .ValidateSchemaID}} `validate_schema_id`="{{.ValidateSchemaID}}",{{end}}
 {{if .FullName}} `full_name`="{{.FullName}}",{{end}}
 {{if .Name}} `name`="{{.Name}}",{{end}}
 {{if .Title}} `title`="{{.Title}}",{{end}}
 {{if .Tag}} `tag`="{{.Tag}}",{{end}}
 {{if .Position}} `position`="{{.Position}}",{{end}}
 {{if .Group}} `group`="{{.Group}}",{{end}}
 {{if .Example}} `example`="{{.Example}}",{{end}}
 {{if .Deprecated}} `deprecated`="{{.Deprecated}}",{{end}}
 {{if .Required}} `required`="{{.Required}}",{{end}}
 {{if .Serialize}} `serialize`="{{.Serialize}}",{{end}}
 {{if .Explode}} `explode`="{{.Explode}}",{{end}}
 {{if .AllowEmptyValue}} `allow_empty_value`="{{.AllowEmptyValue}}",{{end}}
 {{if .AllowReserved}} `allow_reserved`="{{.AllowReserved}}",{{end}}
 {{if .Description}} `description`="{{.Description}}",{{end}}
 {{if .CreatedAt}} `created_at`="{{.CreatedAt}}",{{end}}
 {{if .UpdatedAt}} `updated_at`="{{.UpdatedAt}}",{{end}}
 {{if .DeletedAt}} `deleted_at`="{{.DeletedAt}}",{{end}}
where `parameter_id`="{{.ParameterID}}";
 {{end}}


{{define "del"}}
 update `request` set `deleted_at`="{{currentTime}}" where `parameter_id`="{{.ParameterID}}";
{{end}}