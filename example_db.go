package gqt

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/pkg/errors"
	"github.com/suifengpiao14/errorformatter"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB
var dbOnce sync.Once

var DB_SOURCE = ""
var gormConfig = gorm.Config{
	Logger: logger.Default.LogMode(logger.Info),
}

// GetDb is a signal DB
func GetDb() *gorm.DB {
	if db == nil {
		InitDB()
	}
	return db
}

// GetDb 获取db实例
func InitDB() *gorm.DB {
	if DB_SOURCE == "" {
		panic("var DB_SOURCE must be set value")
	}
	dbOnce.Do(func() {
		dbCon, err := gorm.Open(mysql.Open(DB_SOURCE), &gormConfig)
		if err != nil {
			panic(err)
		}
		db = dbCon

	})
	return db
}

// DBCreateTable 初始化数据表
func DBCreateTable(getrepositorySQL func() *RepositorySQL) {
	repositorySQL := getrepositorySQL()

	ddlNamespace, err := repositorySQL.GetDDLNamespace()
	if err != nil {
		panic(err)
	}
	ddlSQLRow, err := repositorySQL.GetByNamespace(ddlNamespace, nil)
	if err != nil {
		panic(err)
	}

	ddlMap := make(map[string]string)
	for _, sqlRaw := range ddlSQLRow {
		fullname := fmt.Sprintf("%s.%s", sqlRaw.Namespace, sqlRaw.Name)
		ddlMap[fullname] = sqlRaw.SQL
	}
	if err := DBBatchExec(GetDb, ddlMap); err != nil {
		panic(err)
	}
}

func DBExec(sqlRepository func() *RepositorySQL, db func() *gorm.DB, entity TplEntityInterface) (err error) {
	var sql string
	err = errorformatter.NewErrorChain().
		SetError(sqlRepository().GetSQLRef(entity, &sql)).
		SetError(Flight(sql, nil, func() (interface{}, error) {
			err = db().Exec(sql).Error
			return nil, err
		})).
		Error()
	return
}

func DBRawScan(sqlRepository func() *RepositorySQL, db func() *gorm.DB, entity TplEntityInterface, output interface{}) (err error) {
	var sql string
	err = errorformatter.NewErrorChain().
		SetError(sqlRepository().GetSQLRef(entity, &sql)).
		SetError(Flight(sql, output, func() (interface{}, error) {
			// 闭包函数中的output 只是复用外部output 的类型，最终Flight 中使用 匿名函数返回值 第一个作为最终返回值
			err = db().Raw(sql).Scan(output).Error
			return output, err
		})).
		Error()
	return
}
func DBCount(sqlRepository func() *RepositorySQL, db func() *gorm.DB, entity TplEntityInterface, count *int) (err error) {
	var sql string
	err = errorformatter.NewErrorChain().
		SetError(sqlRepository().GetSQLRef(entity, &sql)).
		SetError(Flight(sql, count, func() (interface{}, error) {
			var count64 int64
			err = db().Raw(sql).Count(&count64).Error
			*count = int(count64)
			return count, err
		})).
		Error()
	return
}

func DBBatchExec(db func() *gorm.DB, sqlMap map[string]string) (err error) {
	err = db().Transaction(func(tx *gorm.DB) (err error) {
		for _, sql := range sqlMap {
			err = tx.Exec(sql).Error
			if err != nil {
				return err
			}
		}
		return
	})
	return
}

func WrapDBScanSQL(db func() *gorm.DB, throwNotFoundErr bool) func(sqlRowList []*SQLRow) (err error) {
	return func(sqlRowList []*SQLRow) (err error) {
		dbInst := db()
		errorChain := errorformatter.NewErrorChain()

		for _, sqlRow := range sqlRowList {
			if sqlRow.Result == nil {
				err = errors.Errorf("not foud tplName:%s result data struct", sqlRow.Name)
				errorChain.SetError(err)
				return err
			}
			if sqlRow.SQL == "" {
				err = errors.Errorf("tplName:%s  sql string required", sqlRow.Name)
				errorChain.SetError(err)
				return err
			}
			err = dbInst.Raw(sqlRow.SQL).Scan(sqlRow.Result).Error
			if err == gorm.ErrRecordNotFound && !throwNotFoundErr {
				err = nil
			}

			if err != nil {
				errorChain.SetError(err)
				return err
			}
		}
		return errorChain.Error()
	}
}

func WrapDBExecSQL(db func() *gorm.DB) func(sqlRowList []*SQLRow) (err error) {
	return func(sqlRowList []*SQLRow) (err error) {
		errorChain := errorformatter.NewErrorChain()
		err = db().Transaction(func(tx *gorm.DB) (err error) {
			for _, sqlRow := range sqlRowList {
				if sqlRow.SQL == "" {
					err = errors.Errorf("tag %s sql required", sqlRow.Name)
					return err
				}
				err = tx.Exec(sqlRow.SQL).Error
				if err != nil {
					errorChain.SetError(err)
					return errorChain.Error()
				}
			}
			return errorChain.Error()
		})
		errorChain.SetError(err)
		return errorChain.Error()
	}
}

func DBTryFind(sqlRepository func() *RepositorySQL, db func() *gorm.DB, entity TplEntityInterface, output interface{}) (err error) {
	var sql string
	err = errorformatter.NewErrorChain().
		SetError(sqlRepository().GetSQLRef(entity, &sql)).
		SetError(Flight(sql, output, func() (interface{}, error) {
			err = db().Raw(sql).Scan(output).Error
			if err == gorm.ErrRecordNotFound {
				err = nil
			}
			return output, err
		})).
		Error()
	return
}

type DBBatchSaveArgs struct {
	ModelList       interface{}
	PrimaryKeyCamel string
	UpdateEntity    TplEntityInterface // 无需填充
	InsertEntity    TplEntityInterface // 无需填充
	DelEntity       TplEntityInterface // 无需填充
}

func DBBatchSave(sqlRepository func() *RepositorySQL, db func() *gorm.DB, getByIDsEntity TplEntityInterface, args *DBBatchSaveArgs) (err error) {

	sqlChain := NewSQLChain(sqlRepository)
	var dbModelList []interface{}
	sqlChain.SetError(DBRawScan(sqlRepository, db, getByIDsEntity, &dbModelList))
	err = sqlChain.Error()
	if err != nil {
		return err
	}
	batchSqlArgs := &BatchInsertUpdateDelSQLArgs{
		ModelList:       args.ModelList,
		DBModelList:     dbModelList,
		PrimaryKeyCamel: args.PrimaryKeyCamel,
		UpdateEntity:    args.UpdateEntity,
		InsertEntity:    args.InsertEntity,
		DelEntity:       nil,
		SqlChain:        sqlChain,
	}
	BatchInsertUpdateDelSQL(batchSqlArgs)
	err = sqlChain.Error()
	if err != nil {
		return
	}
	err = sqlChain.Exec(WrapDBExecSQL(db))
	return
}

type BatchInsertUpdateDelSQLArgs struct {
	ModelList       interface{}
	DBModelList     interface{}
	PrimaryKeyCamel string
	UpdateEntity    TplEntityInterface
	InsertEntity    TplEntityInterface
	DelEntity       TplEntityInterface
	SqlChain        *SQLChain
}

func BatchInsertUpdateDelSQL(args *BatchInsertUpdateDelSQLArgs) {
	if args.SqlChain.Error() != nil {
		return
	}
	dt := reflect.TypeOf(args.DBModelList)
	if dt.Kind() == reflect.Ptr {
		dt = dt.Elem()
	}
	if !(dt.Kind() == reflect.Array || dt.Kind() == reflect.Slice) {
		panic(fmt.Errorf("BatchInsertUpdateDelSQLArgsource.DBModelList want array/slice ,have %s ", dt.Kind().String()))
	}

	mt := reflect.TypeOf(args.ModelList)
	if mt.Kind() == reflect.Ptr {
		mt = mt.Elem()
	}
	if !(mt.Kind() == reflect.Array || mt.Kind() == reflect.Slice) {
		panic(fmt.Errorf("BatchInsertUpdateDelSQLArgs.ModelList want array/slice ,have %s ", mt.Kind().String()))
	}

	var dbModelMap = make(map[string]interface{})
	dv := reflect.Indirect(reflect.ValueOf(args.DBModelList))
	dl := dv.Len()
	for i := 0; i < dl; i++ {
		dbModelV := dv.Index(i)
		id := dbModelV.FieldByName(args.PrimaryKeyCamel).String()
		dbModelMap[id] = dbModelV
	}

	var updateMap = make(map[string]interface{})
	var insertMap = make(map[string]interface{})
	var delMap = make(map[string]interface{})
	mv := reflect.Indirect(reflect.ValueOf(args.ModelList))
	ml := mv.Len()
	for i := 0; i < ml; i++ {
		modelV := mv.Index(i)
		if modelV.Type().Kind() == reflect.Ptr {
			modelV = modelV.Elem()
		}
		id := modelV.FieldByName(args.PrimaryKeyCamel).String()
		_, ok := dbModelMap[id]
		model := modelV.Interface()
		if ok {
			updateMap[id] = model
		} else {
			insertMap[id] = model
		}
	}

	for _, model := range updateMap {
		ConvertStruct(model, args.UpdateEntity)
		args.SqlChain.ParseSQL(args.UpdateEntity, nil)
	}
	for _, model := range insertMap {
		ConvertStruct(model, args.InsertEntity)
		args.SqlChain.ParseSQL(args.InsertEntity, nil)
	}
	if args.DelEntity == nil {
		return
	}
	for id, dmV := range dbModelMap {
		_, ok := updateMap[id]
		if ok {
			continue
		}
		delMap[id] = dmV
	}
	for _, model := range delMap {
		ConvertStruct(model, &args.DelEntity)
		args.SqlChain.ParseSQL(args.DelEntity, nil)
	}
}
