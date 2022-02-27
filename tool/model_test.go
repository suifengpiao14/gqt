package tool

import (
	"fmt"
	"testing"
)

func TestGenerateTable(t *testing.T) {
	ddlList := GetDDL()
	tableList, err := GenerateTable(ddlList)
	if err != nil {
		panic(err)
	}
	fmt.Println(tableList)
}

func GetDDL() (ddlList []string) {
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
	ddlApi := `
	CREATE TABLE  if not exists api (
		api_id varchar(64) NOT NULL DEFAULT '' COMMENT 'api 文档标识',
		service_id varchar(64) NOT NULL DEFAULT '' COMMENT '服务标识',
		name varchar(255) NOT NULL COMMENT '路由名称(英文)',
		title varchar(255) NOT NULL COMMENT '标题',
		tags char(255) NOT NULL DEFAULT '' COMMENT '标签ID集合',
		uri varchar(255) NOT NULL DEFAULT '' COMMENT '路径',
		summary varchar(255) NOT NULL DEFAULT '' COMMENT '摘要',
		description varchar(255) NOT NULL DEFAULT '' COMMENT '介绍',
		created_at datetime DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
		updated_at datetime DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '修改时间',
		deleted_at datetime  DEFAULT null COMMENT '删除时间',
		PRIMARY KEY (api_id),
		KEY name (name),
		key server_id (service_id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT="接口";
	`
	ddlSchema := `
	CREATE TABLE  if not exists validate_schema (
		validate_schema_id varchar(64) NOT NULL DEFAULT '' COMMENT 'api schema 标识',
		service_id varchar(64) NOT NULL DEFAULT '' COMMENT '所属服务标识',
		remark varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
		example varchar(1024) NOT NULL DEFAULT '' COMMENT '案例',
		deprecated enum("true","false")  NOT NULL DEFAULT 'false' COMMENT '是否弃用',
		enum_names char(100) NOT NULL DEFAULT '' COMMENT '枚举名称',
		enum_titles char(100) NOT NULL DEFAULT '' COMMENT '枚举标题',
		nullable  enum('true','false') NOT NULL DEFAULT 'false' COMMENT '是否可以为空',
		multiple_of int(11) unsigned NOT NULL DEFAULT '0' COMMENT '倍数',
		maxnum int(11) unsigned NOT NULL DEFAULT '0' COMMENT '最大值',
		exclusive_maximum enum('true','false') NOT NULL DEFAULT 'false' NOT NULL DEFAULT 'false' COMMENT '是否不包含最大项',
		minimum int(11) unsigned NOT NULL DEFAULT '0' COMMENT '最小值',
		exclusive_minimum enum('true','false') NOT NULL DEFAULT 'false' NOT NULL DEFAULT 'true' COMMENT '是否不包含最小项',
		max_length int(11) unsigned NOT NULL DEFAULT '0' COMMENT '最大长度',
		min_length int(11) unsigned NOT NULL DEFAULT '0' COMMENT '最小长度',
		pattern varchar(255) NOT NULL DEFAULT '' COMMENT '正则表达式',
		max_items int(11) unsigned NOT NULL DEFAULT '0' COMMENT '最大项数',
		min_items int(11) unsigned NOT NULL DEFAULT '0' COMMENT '最小项数',
		unique_items enum('true','false') NOT NULL DEFAULT 'false' COMMENT '所有项是否需要唯一',
		max_properties int(11) unsigned NOT NULL DEFAULT '0' COMMENT '最多属性项',
		min_properties int(11) unsigned NOT NULL DEFAULT '0' COMMENT '最少属性项',
		all_of varchar(255) NOT NULL DEFAULT '' COMMENT '所有',
		one_of varchar(255) NOT NULL DEFAULT '' COMMENT '只满足一个',
		any_of varchar(255) NOT NULL DEFAULT '' COMMENT '任何一个SchemaID',
		allow_empty_value enum('true','false') NOT NULL DEFAULT 'true' COMMENT '是否容许空值',
		allow_reserved enum('true','false') NOT NULL DEFAULT 'false' COMMENT '特殊字符是否容许出现在uri参数中',
		additional_properties varchar(255) NOT NULL DEFAULT '' COMMENT 'boolean',
		discriminator varchar(255) NOT NULL DEFAULT '' COMMENT 'schema鉴别',
		write_only enum('true','false') NOT NULL DEFAULT 'false' COMMENT '是否只写',
		external_docs varchar(255)  NOT NULL DEFAULT '' COMMENT '附加文档',
		external_pros varchar(255)  NOT NULL DEFAULT '' COMMENT '附加文档',
		extensions varchar(255) NOT NULL DEFAULT '' COMMENT '扩展字段',
		created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
		updated_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP on update CURRENT_TIMESTAMP COMMENT '修改时间',
		deleted_at timestamp  DEFAULT null COMMENT '删除时间',
		summary varchar(255) NOT NULL DEFAULT '' COMMENT '简介',
		PRIMARY KEY (validate_schema_id),
		key server_id (service_id)
	  ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='响应格式';
	`
	ddlList = make([]string, 0)
	ddlList = append(ddlList, ddlService)
	ddlList = append(ddlList, ddlApi)
	ddlList = append(ddlList, ddlSchema)
	return
}
