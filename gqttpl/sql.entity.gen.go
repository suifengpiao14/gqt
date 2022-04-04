package gqttpl

type GenAPISQLDelEntity struct {
	APIID string
	TplEmptyEntity
}

func (t *GenAPISQLDelEntity) TplName() string {
	return "gen.api.sql.Del"
}
func (t *GenAPISQLDelEntity) TplType() string {
	return "sql_update"
}

type GenAPISQLGetAllByAPIIDListEntity struct {
	APIIDList []string
	TplEmptyEntity
}

func (t *GenAPISQLGetAllByAPIIDListEntity) TplName() string {
	return "gen.api.sql.GetAllByAPIIDList"
}
func (t *GenAPISQLGetAllByAPIIDListEntity) TplType() string {
	return "sql_select"
}

type GenAPISQLGetByAPIIDEntity struct {
	APIID string
	TplEmptyEntity
}

func (t *GenAPISQLGetByAPIIDEntity) TplName() string {
	return "gen.api.sql.GetByAPIID"
}
func (t *GenAPISQLGetByAPIIDEntity) TplType() string {
	return "sql_select"
}

type GenAPISQLGetByServiceIDEntity struct {
	ServiceID string
	TplEmptyEntity
}

func (t *GenAPISQLGetByServiceIDEntity) TplName() string {
	return "gen.api.sql.GetByServiceID"
}
func (t *GenAPISQLGetByServiceIDEntity) TplType() string {
	return "sql_select"
}

type GenAPISQLInsertEntity struct {
	APIID       string
	Description string
	Name        string
	ServiceID   string
	Summary     string
	Tags        string
	Title       string
	URI         string
	TplEmptyEntity
}

func (t *GenAPISQLInsertEntity) TplName() string {
	return "gen.api.sql.Insert"
}
func (t *GenAPISQLInsertEntity) TplType() string {
	return "sql_insert"
}

type GenAPISQLPaginateEntity struct {
	Limit  int
	Offset int
	GenAPISQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenAPISQLPaginateEntity) TplName() string {
	return "gen.api.sql.Paginate"
}
func (t *GenAPISQLPaginateEntity) TplType() string {
	return "sql_select"
}

type GenAPISQLPaginateTotalEntity struct {
	GenAPISQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenAPISQLPaginateTotalEntity) TplName() string {
	return "gen.api.sql.PaginateTotal"
}
func (t *GenAPISQLPaginateTotalEntity) TplType() string {
	return "sql_select"
}

type GenAPISQLPaginateWhereEntity struct {
	TplEmptyEntity
}

func (t *GenAPISQLPaginateWhereEntity) TplName() string {
	return "gen.api.sql.PaginateWhere"
}
func (t *GenAPISQLPaginateWhereEntity) TplType() string {
	return "text"
}

type GenAPISQLUpdateEntity struct {
	APIID       string
	Description string
	Name        string
	ServiceID   string
	Summary     string
	Tags        string
	Title       string
	URI         string
	TplEmptyEntity
}

func (t *GenAPISQLUpdateEntity) TplName() string {
	return "gen.api.sql.Update"
}
func (t *GenAPISQLUpdateEntity) TplType() string {
	return "sql_update"
}

type GenExampleSQLDelEntity struct {
	ExampleID string
	TplEmptyEntity
}

func (t *GenExampleSQLDelEntity) TplName() string {
	return "gen.example.sql.Del"
}
func (t *GenExampleSQLDelEntity) TplType() string {
	return "sql_update"
}

type GenExampleSQLGetAllByExampleIDListEntity struct {
	ExampleIDList []string
	TplEmptyEntity
}

func (t *GenExampleSQLGetAllByExampleIDListEntity) TplName() string {
	return "gen.example.sql.GetAllByExampleIDList"
}
func (t *GenExampleSQLGetAllByExampleIDListEntity) TplType() string {
	return "sql_select"
}

type GenExampleSQLGetByExampleIDEntity struct {
	ExampleID string
	TplEmptyEntity
}

func (t *GenExampleSQLGetByExampleIDEntity) TplName() string {
	return "gen.example.sql.GetByExampleID"
}
func (t *GenExampleSQLGetByExampleIDEntity) TplType() string {
	return "sql_select"
}

type GenExampleSQLGetByServiceIDEntity struct {
	ServiceID string
	TplEmptyEntity
}

func (t *GenExampleSQLGetByServiceIDEntity) TplName() string {
	return "gen.example.sql.GetByServiceID"
}
func (t *GenExampleSQLGetByServiceIDEntity) TplType() string {
	return "sql_select"
}

type GenExampleSQLInsertEntity struct {
	APIID            string
	Auth             string
	Body             string
	ContentType      string
	ExampleID        string
	Headers          string
	Method           string
	Parameters       string
	PreRequestScript string
	Response         string
	ServiceID        string
	Summary          string
	Tag              string
	TestScript       string
	Title            string
	URL              string
	TplEmptyEntity
}

func (t *GenExampleSQLInsertEntity) TplName() string {
	return "gen.example.sql.Insert"
}
func (t *GenExampleSQLInsertEntity) TplType() string {
	return "sql_insert"
}

type GenExampleSQLPaginateEntity struct {
	Limit  int
	Offset int
	GenExampleSQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenExampleSQLPaginateEntity) TplName() string {
	return "gen.example.sql.Paginate"
}
func (t *GenExampleSQLPaginateEntity) TplType() string {
	return "sql_select"
}

type GenExampleSQLPaginateTotalEntity struct {
	GenExampleSQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenExampleSQLPaginateTotalEntity) TplName() string {
	return "gen.example.sql.PaginateTotal"
}
func (t *GenExampleSQLPaginateTotalEntity) TplType() string {
	return "sql_select"
}

type GenExampleSQLPaginateWhereEntity struct {
	TplEmptyEntity
}

func (t *GenExampleSQLPaginateWhereEntity) TplName() string {
	return "gen.example.sql.PaginateWhere"
}
func (t *GenExampleSQLPaginateWhereEntity) TplType() string {
	return "text"
}

type GenExampleSQLUpdateEntity struct {
	APIID            string
	Auth             string
	Body             string
	ContentType      string
	ExampleID        string
	Headers          string
	Method           string
	Parameters       string
	PreRequestScript string
	Response         string
	ServiceID        string
	Summary          string
	Tag              string
	TestScript       string
	Title            string
	URL              string
	TplEmptyEntity
}

func (t *GenExampleSQLUpdateEntity) TplName() string {
	return "gen.example.sql.Update"
}
func (t *GenExampleSQLUpdateEntity) TplType() string {
	return "sql_update"
}

type GenMarkdownSQLDelEntity struct {
	MarkdownID string
	TplEmptyEntity
}

func (t *GenMarkdownSQLDelEntity) TplName() string {
	return "gen.markdown.sql.Del"
}
func (t *GenMarkdownSQLDelEntity) TplType() string {
	return "sql_update"
}

type GenMarkdownSQLGetAllByMarkdownIDListEntity struct {
	MarkdownIDList []string
	TplEmptyEntity
}

func (t *GenMarkdownSQLGetAllByMarkdownIDListEntity) TplName() string {
	return "gen.markdown.sql.GetAllByMarkdownIDList"
}
func (t *GenMarkdownSQLGetAllByMarkdownIDListEntity) TplType() string {
	return "sql_select"
}

type GenMarkdownSQLGetByMarkdownIDEntity struct {
	MarkdownID string
	TplEmptyEntity
}

func (t *GenMarkdownSQLGetByMarkdownIDEntity) TplName() string {
	return "gen.markdown.sql.GetByMarkdownID"
}
func (t *GenMarkdownSQLGetByMarkdownIDEntity) TplType() string {
	return "sql_select"
}

type GenMarkdownSQLGetByServiceIDEntity struct {
	ServiceID string
	TplEmptyEntity
}

func (t *GenMarkdownSQLGetByServiceIDEntity) TplName() string {
	return "gen.markdown.sql.GetByServiceID"
}
func (t *GenMarkdownSQLGetByServiceIDEntity) TplType() string {
	return "sql_select"
}

type GenMarkdownSQLInsertEntity struct {
	APIID      string
	Content    string
	Markdown   string
	MarkdownID string
	Name       string
	OwnerID    int
	OwnerName  string
	ServiceID  string
	Title      string
	TplEmptyEntity
}

func (t *GenMarkdownSQLInsertEntity) TplName() string {
	return "gen.markdown.sql.Insert"
}
func (t *GenMarkdownSQLInsertEntity) TplType() string {
	return "sql_insert"
}

type GenMarkdownSQLPaginateEntity struct {
	Limit  int
	Offset int
	GenMarkdownSQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenMarkdownSQLPaginateEntity) TplName() string {
	return "gen.markdown.sql.Paginate"
}
func (t *GenMarkdownSQLPaginateEntity) TplType() string {
	return "sql_select"
}

type GenMarkdownSQLPaginateTotalEntity struct {
	GenMarkdownSQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenMarkdownSQLPaginateTotalEntity) TplName() string {
	return "gen.markdown.sql.PaginateTotal"
}
func (t *GenMarkdownSQLPaginateTotalEntity) TplType() string {
	return "sql_select"
}

type GenMarkdownSQLPaginateWhereEntity struct {
	TplEmptyEntity
}

func (t *GenMarkdownSQLPaginateWhereEntity) TplName() string {
	return "gen.markdown.sql.PaginateWhere"
}
func (t *GenMarkdownSQLPaginateWhereEntity) TplType() string {
	return "text"
}

type GenMarkdownSQLUpdateEntity struct {
	APIID      string
	Content    string
	Markdown   string
	MarkdownID string
	Name       string
	OwnerID    int
	OwnerName  string
	ServiceID  string
	Title      string
	TplEmptyEntity
}

func (t *GenMarkdownSQLUpdateEntity) TplName() string {
	return "gen.markdown.sql.Update"
}
func (t *GenMarkdownSQLUpdateEntity) TplType() string {
	return "sql_update"
}

type GenParameterSQLDelEntity struct {
	ParameterID string
	TplEmptyEntity
}

func (t *GenParameterSQLDelEntity) TplName() string {
	return "gen.parameter.sql.Del"
}
func (t *GenParameterSQLDelEntity) TplType() string {
	return "sql_update"
}

type GenParameterSQLGetAllByParameterIDListEntity struct {
	ParameterIDList []string
	TplEmptyEntity
}

func (t *GenParameterSQLGetAllByParameterIDListEntity) TplName() string {
	return "gen.parameter.sql.GetAllByParameterIDList"
}
func (t *GenParameterSQLGetAllByParameterIDListEntity) TplType() string {
	return "sql_select"
}

type GenParameterSQLGetByParameterIDEntity struct {
	ParameterID string
	TplEmptyEntity
}

func (t *GenParameterSQLGetByParameterIDEntity) TplName() string {
	return "gen.parameter.sql.GetByParameterID"
}
func (t *GenParameterSQLGetByParameterIDEntity) TplType() string {
	return "sql_select"
}

type GenParameterSQLGetByServiceIDEntity struct {
	ServiceID string
	TplEmptyEntity
}

func (t *GenParameterSQLGetByServiceIDEntity) TplName() string {
	return "gen.parameter.sql.GetByServiceID"
}
func (t *GenParameterSQLGetByServiceIDEntity) TplType() string {
	return "sql_select"
}

type GenParameterSQLInsertEntity struct {
	APIID           string
	AllowEmptyValue string
	AllowReserved   string
	Deprecated      string
	Description     string
	Example         string
	Explode         string
	FullName        string
	HTTPStatus      string
	Method          string
	Name            string
	ParameterID     string
	Position        string
	Required        string
	SchemaID        string
	Serialize       string
	ServiceID       string
	Tag             string
	Title           string
	Type            string
	TplEmptyEntity
}

func (t *GenParameterSQLInsertEntity) TplName() string {
	return "gen.parameter.sql.Insert"
}
func (t *GenParameterSQLInsertEntity) TplType() string {
	return "sql_insert"
}

type GenParameterSQLPaginateEntity struct {
	Limit  int
	Offset int
	GenParameterSQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenParameterSQLPaginateEntity) TplName() string {
	return "gen.parameter.sql.Paginate"
}
func (t *GenParameterSQLPaginateEntity) TplType() string {
	return "sql_select"
}

type GenParameterSQLPaginateTotalEntity struct {
	GenParameterSQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenParameterSQLPaginateTotalEntity) TplName() string {
	return "gen.parameter.sql.PaginateTotal"
}
func (t *GenParameterSQLPaginateTotalEntity) TplType() string {
	return "sql_select"
}

type GenParameterSQLPaginateWhereEntity struct {
	TplEmptyEntity
}

func (t *GenParameterSQLPaginateWhereEntity) TplName() string {
	return "gen.parameter.sql.PaginateWhere"
}
func (t *GenParameterSQLPaginateWhereEntity) TplType() string {
	return "text"
}

type GenParameterSQLUpdateEntity struct {
	APIID           string
	AllowEmptyValue string
	AllowReserved   string
	Deprecated      string
	Description     string
	Example         string
	Explode         string
	FullName        string
	HTTPStatus      string
	Method          string
	Name            string
	ParameterID     string
	Position        string
	Required        string
	SchemaID        string
	Serialize       string
	ServiceID       string
	Tag             string
	Title           string
	Type            string
	TplEmptyEntity
}

func (t *GenParameterSQLUpdateEntity) TplName() string {
	return "gen.parameter.sql.Update"
}
func (t *GenParameterSQLUpdateEntity) TplType() string {
	return "sql_update"
}

type GenSchemaSQLDelEntity struct {
	SchemaID string
	TplEmptyEntity
}

func (t *GenSchemaSQLDelEntity) TplName() string {
	return "gen.schema.sql.Del"
}
func (t *GenSchemaSQLDelEntity) TplType() string {
	return "sql_update"
}

type GenSchemaSQLGetAllBySchemaIDListEntity struct {
	SchemaIDList []string
	TplEmptyEntity
}

func (t *GenSchemaSQLGetAllBySchemaIDListEntity) TplName() string {
	return "gen.schema.sql.GetAllBySchemaIDList"
}
func (t *GenSchemaSQLGetAllBySchemaIDListEntity) TplType() string {
	return "sql_select"
}

type GenSchemaSQLGetBySchemaIDEntity struct {
	SchemaID string
	TplEmptyEntity
}

func (t *GenSchemaSQLGetBySchemaIDEntity) TplName() string {
	return "gen.schema.sql.GetBySchemaID"
}
func (t *GenSchemaSQLGetBySchemaIDEntity) TplType() string {
	return "sql_select"
}

type GenSchemaSQLGetByServiceIDEntity struct {
	ServiceID string
	TplEmptyEntity
}

func (t *GenSchemaSQLGetByServiceIDEntity) TplName() string {
	return "gen.schema.sql.GetByServiceID"
}
func (t *GenSchemaSQLGetByServiceIDEntity) TplType() string {
	return "sql_select"
}

type GenSchemaSQLInsertEntity struct {
	AdditionalProperties string
	AllOf                string
	AllowEmptyValue      string
	AllowReserved        string
	AnyOf                string
	Default              string
	Deprecated           string
	Description          string
	Discriminator        string
	Enum                 string
	EnumNames            string
	EnumTitles           string
	Example              string
	ExclusiveMaximum     string
	ExclusiveMinimum     string
	Extensions           string
	ExternalDocs         string
	ExternalPros         string
	Format               string
	MaxItems             int
	MaxLength            int
	MaxProperties        int
	Maxnum               int
	MinItems             int
	MinLength            int
	MinProperties        int
	Minimum              int
	MultipleOf           int
	Not                  string
	Nullable             string
	OneOf                string
	Pattern              string
	ReadOnly             string
	Remark               string
	Required             string
	SchemaID             string
	ServiceID            string
	Summary              string
	Type                 string
	UniqueItems          string
	WriteOnly            string
	XML                  string
	TplEmptyEntity
}

func (t *GenSchemaSQLInsertEntity) TplName() string {
	return "gen.schema.sql.Insert"
}
func (t *GenSchemaSQLInsertEntity) TplType() string {
	return "sql_insert"
}

type GenSchemaSQLPaginateEntity struct {
	Limit  int
	Offset int
	GenSchemaSQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenSchemaSQLPaginateEntity) TplName() string {
	return "gen.schema.sql.Paginate"
}
func (t *GenSchemaSQLPaginateEntity) TplType() string {
	return "sql_select"
}

type GenSchemaSQLPaginateTotalEntity struct {
	GenSchemaSQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenSchemaSQLPaginateTotalEntity) TplName() string {
	return "gen.schema.sql.PaginateTotal"
}
func (t *GenSchemaSQLPaginateTotalEntity) TplType() string {
	return "sql_select"
}

type GenSchemaSQLPaginateWhereEntity struct {
	TplEmptyEntity
}

func (t *GenSchemaSQLPaginateWhereEntity) TplName() string {
	return "gen.schema.sql.PaginateWhere"
}
func (t *GenSchemaSQLPaginateWhereEntity) TplType() string {
	return "text"
}

type GenSchemaSQLUpdateEntity struct {
	AdditionalProperties string
	AllOf                string
	AllowEmptyValue      string
	AllowReserved        string
	AnyOf                string
	Default              string
	Deprecated           string
	Description          string
	Discriminator        string
	Enum                 string
	EnumNames            string
	EnumTitles           string
	Example              string
	ExclusiveMaximum     string
	ExclusiveMinimum     string
	Extensions           string
	ExternalDocs         string
	ExternalPros         string
	Format               string
	MaxItems             int
	MaxLength            int
	MaxProperties        int
	Maxnum               int
	MinItems             int
	MinLength            int
	MinProperties        int
	Minimum              int
	MultipleOf           int
	Not                  string
	Nullable             string
	OneOf                string
	Pattern              string
	ReadOnly             string
	Remark               string
	Required             string
	SchemaID             string
	ServiceID            string
	Summary              string
	Type                 string
	UniqueItems          string
	WriteOnly            string
	XML                  string
	TplEmptyEntity
}

func (t *GenSchemaSQLUpdateEntity) TplName() string {
	return "gen.schema.sql.Update"
}
func (t *GenSchemaSQLUpdateEntity) TplType() string {
	return "sql_update"
}

type GenServerSQLDelEntity struct {
	ServerID string
	TplEmptyEntity
}

func (t *GenServerSQLDelEntity) TplName() string {
	return "gen.server.sql.Del"
}
func (t *GenServerSQLDelEntity) TplType() string {
	return "sql_update"
}

type GenServerSQLGetAllByServerIDListEntity struct {
	ServerIDList []string
	TplEmptyEntity
}

func (t *GenServerSQLGetAllByServerIDListEntity) TplName() string {
	return "gen.server.sql.GetAllByServerIDList"
}
func (t *GenServerSQLGetAllByServerIDListEntity) TplType() string {
	return "sql_select"
}

type GenServerSQLGetByServerIDEntity struct {
	ServerID string
	TplEmptyEntity
}

func (t *GenServerSQLGetByServerIDEntity) TplName() string {
	return "gen.server.sql.GetByServerID"
}
func (t *GenServerSQLGetByServerIDEntity) TplType() string {
	return "sql_select"
}

type GenServerSQLGetByServiceIDEntity struct {
	ServiceID string
	TplEmptyEntity
}

func (t *GenServerSQLGetByServiceIDEntity) TplName() string {
	return "gen.server.sql.GetByServiceID"
}
func (t *GenServerSQLGetByServiceIDEntity) TplType() string {
	return "sql_select"
}

type GenServerSQLInsertEntity struct {
	Description  string
	ExtensionIds string
	Proxy        string
	ServerID     string
	ServiceID    string
	URL          string
	TplEmptyEntity
}

func (t *GenServerSQLInsertEntity) TplName() string {
	return "gen.server.sql.Insert"
}
func (t *GenServerSQLInsertEntity) TplType() string {
	return "sql_insert"
}

type GenServerSQLPaginateEntity struct {
	Limit  int
	Offset int
	GenServerSQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenServerSQLPaginateEntity) TplName() string {
	return "gen.server.sql.Paginate"
}
func (t *GenServerSQLPaginateEntity) TplType() string {
	return "sql_select"
}

type GenServerSQLPaginateTotalEntity struct {
	GenServerSQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenServerSQLPaginateTotalEntity) TplName() string {
	return "gen.server.sql.PaginateTotal"
}
func (t *GenServerSQLPaginateTotalEntity) TplType() string {
	return "sql_select"
}

type GenServerSQLPaginateWhereEntity struct {
	TplEmptyEntity
}

func (t *GenServerSQLPaginateWhereEntity) TplName() string {
	return "gen.server.sql.PaginateWhere"
}
func (t *GenServerSQLPaginateWhereEntity) TplType() string {
	return "text"
}

type GenServerSQLUpdateEntity struct {
	Description  string
	ExtensionIds string
	Proxy        string
	ServerID     string
	ServiceID    string
	URL          string
	TplEmptyEntity
}

func (t *GenServerSQLUpdateEntity) TplName() string {
	return "gen.server.sql.Update"
}
func (t *GenServerSQLUpdateEntity) TplType() string {
	return "sql_update"
}

type GenServiceSQLDelEntity struct {
	ServiceID string
	TplEmptyEntity
}

func (t *GenServiceSQLDelEntity) TplName() string {
	return "gen.service.sql.Del"
}
func (t *GenServiceSQLDelEntity) TplType() string {
	return "sql_update"
}

type GenServiceSQLGetAllByServiceIDListEntity struct {
	ServiceIDList []string
	TplEmptyEntity
}

func (t *GenServiceSQLGetAllByServiceIDListEntity) TplName() string {
	return "gen.service.sql.GetAllByServiceIDList"
}
func (t *GenServiceSQLGetAllByServiceIDListEntity) TplType() string {
	return "sql_select"
}

type GenServiceSQLGetByServiceIDEntity struct {
	ServiceID string
	TplEmptyEntity
}

func (t *GenServiceSQLGetByServiceIDEntity) TplName() string {
	return "gen.service.sql.GetByServiceID"
}
func (t *GenServiceSQLGetByServiceIDEntity) TplType() string {
	return "sql_select"
}

type GenServiceSQLInsertEntity struct {
	ContactIds  string
	Description string
	License     string
	Proxy       string
	Security    string
	ServiceID   string
	Title       string
	Variables   string
	Version     string
	TplEmptyEntity
}

func (t *GenServiceSQLInsertEntity) TplName() string {
	return "gen.service.sql.Insert"
}
func (t *GenServiceSQLInsertEntity) TplType() string {
	return "sql_insert"
}

type GenServiceSQLPaginateEntity struct {
	Limit  int
	Offset int
	GenServiceSQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenServiceSQLPaginateEntity) TplName() string {
	return "gen.service.sql.Paginate"
}
func (t *GenServiceSQLPaginateEntity) TplType() string {
	return "sql_select"
}

type GenServiceSQLPaginateTotalEntity struct {
	GenServiceSQLPaginateWhereEntity
	TplEmptyEntity
}

func (t *GenServiceSQLPaginateTotalEntity) TplName() string {
	return "gen.service.sql.PaginateTotal"
}
func (t *GenServiceSQLPaginateTotalEntity) TplType() string {
	return "sql_select"
}

type GenServiceSQLPaginateWhereEntity struct {
	TplEmptyEntity
}

func (t *GenServiceSQLPaginateWhereEntity) TplName() string {
	return "gen.service.sql.PaginateWhere"
}
func (t *GenServiceSQLPaginateWhereEntity) TplType() string {
	return "text"
}

type GenServiceSQLUpdateEntity struct {
	ContactIds  string
	Description string
	License     string
	Proxy       string
	Security    string
	ServiceID   string
	Title       string
	Variables   string
	Version     string
	TplEmptyEntity
}

func (t *GenServiceSQLUpdateEntity) TplName() string {
	return "gen.service.sql.Update"
}
func (t *GenServiceSQLUpdateEntity) TplType() string {
	return "sql_update"
}
