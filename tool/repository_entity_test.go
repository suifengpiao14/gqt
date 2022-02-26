package tool

import (
	"fmt"
	"testing"
)

func TestParsSqlTplVariable(t *testing.T) {
	sqlTpl := `
{{$preComma:=newPreComma}}
    update service set
     {{if .Title}} {{$preComma.PreComma}} title=:Title{{end}}
     {{if .Description}} {{$preComma.PreComma}} description=:Description{{end}}
     {{if .Version}}  {{$preComma.PreComma}}version=:Version{{end}}
     {{if .ContactIds}} {{$preComma.PreComma}} contact_ids=:ContactIds{{end}}
     {{if .License}}  {{$preComma.PreComma}}license=:License{{end}}
     {{if .Security}} {{$preComma.PreComma}} security=:Security{{end}}
     {{if .Proxy}} {{$preComma.PreComma}} proxy=:Proxy{{end}}
     {{if .Variables}}  {{$preComma.PreComma}}variables =:Variables {{end}}
	 {{if .Variables}} {{in . .List}} {{end}}
       where {{func .Func}} service_id=:ServiceID limit :Offset,:limit;
  `
	veriableList := ParsSqlTplVariable(sqlTpl)
	fmt.Println(veriableList)
}

type DDLSqlTpl struct {
	DDL    string
	SQLTpl string
}

func ddlSqlTplData() (ddlSqlTpl *DDLSqlTpl) {
	ddlService := `CREATE TABLE if not exists service (
		service_id varchar(64) NOT NULL DEFAULT '' COMMENT '服务标识',
		title varchar(255) NOT NULL DEFAULT '' COMMENT '标题',
		description varchar(255) NOT NULL DEFAULT '' COMMENT '介绍',
		version varchar(20) NOT NULL DEFAULT '' COMMENT '版本',
		contact_ids varchar(100) NOT NULL DEFAULT '' COMMENT '联系人',
		license char(255) NOT NULL DEFAULT '' COMMENT '协议',
		security char(255) NOT NULL DEFAULT '' COMMENT '鉴权',
		proxy char(60) NOT NULL DEFAULT '' COMMENT '代理地址',
		variables text NOT NULL COMMENT 'json字符串',
		created_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
		updated_at datetime DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '修改时间',
		deleted_at datetime  DEFAULT null COMMENT '删除时间',
		PRIMARY KEY (service_id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT="项目/服务组件";
	  `
	sqlTplService := `

{{define "getByServiceId"}}
select * from service where service_id=:ServiceID and  deleted_at is null;
{{end}}

{{define "total"}}
   select count(*) as count from service where 1=1 
   {{if .Title}} and title like "%:Title%" {{end}};
{{end}}

{{define "list"}}
   select * from service where 1=1 
   {{if .Title}} and title like "%:Title%" {{end}}
   limit :Offset,:Limit;
{{end}}

{{define "insert"}}
  insert into service (service_id,title,description,version,contact_ids,license,security,proxy,variables)
  values(:ServiceID,:Title,:Description,:Version,:ContactIds,:License,:Security,:Proxy,:Variables);
{{end}}

{{define "update"}}
{{$preComma:=newPreComma}}
    update service set
     {{if .Title}} {{$preComma.PreComma}} title=:Title{{end}}
     {{if .Description}} {{$preComma.PreComma}} description=:Description{{end}}
     {{if .Version}}  {{$preComma.PreComma}}version=:Version{{end}}
     {{if .ContactIds}} {{$preComma.PreComma}} contact_ids=:ContactIds{{end}}
     {{if .License}}  {{$preComma.PreComma}}license=:License{{end}}
     {{if .Security}} {{$preComma.PreComma}} security=:Security{{end}}
     {{if .Proxy}} {{$preComma.PreComma}} proxy=:Proxy{{end}}
     {{if .Variables}}  {{$preComma.PreComma}}variables =:Variables {{end}}
       where service_id=:ServiceID;
 {{end}}

  {{define "del"}}
 update service set deleted_at={{currentTime . }} where service_id=:ServiceID;
 {{end}}
	`
	ddlSqlTpl = &DDLSqlTpl{
		DDL:    ddlService,
		SQLTpl: sqlTplService,
	}
	return
}
