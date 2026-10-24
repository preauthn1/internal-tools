using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Data;
using System.Data.Common;
using System.Data.SqlClient;
using System.Linq;
using System.Reflection;
using System.Text;
using Dapper;
using target.Core.Helper;
using target.Models.Models;
using target.Models.Paging;
using target.Utility;
using target.Utility.Extensions;

namespace target.Core.Repository
{
    public abstract class BaseDapperRepository<T> where T: class
    {
        protected abstract string ConnectionString { get; }

        protected abstract string TableName { get; }
      

        protected virtual string DefaultSortColumn
        {
            get
            {
                return this.PrimaryKeyColumnName;
            }
        }

        protected virtual string PrimaryKeyColumnName
        { 
            get
            {
                return "Id";
            }
        }


        protected string GetConnectionStringByDbName(string dbName)
        {
            string settingName = $"ConnectionStrings:{dbName}";
            var conn = BaseSettingHelper.GetSettingValue(settingName);

            if(string.IsNullOrWhiteSpace(conn))
            {
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_SettingDataNotCompleted, new ErrorDetailModel
                {
                    ColumnName = settingName,
                    IsAlwaysNotify = true
                });
            }               

            return conn;
        }

        protected SqlConnection InstanceSQLConn(string connectString)
        {
            return new SqlConnection(connectString);
        }
        
        protected SqlConnection InstanceSQLConn()
        {
            return new SqlConnection(ConnectionString);
        }


        public DapperSqlQueue GenerateSQL_Create(T model)
        {
            string insertSQL = SQLHelper.GetInsertCommand<T>(TableName);

            return new DapperSqlQueue
            {
                Sql = insertSQL,
                Parameters = model
            };
        }

        public void ExecuteTransactionSQLs(params DapperSqlQueue[] sqlCommands)
        {
            ExecuteTransactionSQLs(sqlCommands?.ToList());
        }


        public void ExecuteTransactionSQLs(List<DapperSqlQueue> sqlCommands)
        {
            using (var conn = InstanceSQLConn())
            {
                conn.Open();
                using (var trans = conn.BeginTransaction())
                {
                    foreach(var item in sqlCommands)
                    {
                        var resultId = conn.Execute(item.Sql, item.Parameters, trans);
                    }

                    trans.Commit();
                }
            }
        }



        public void CreateMany(IEnumerable<T> models)
        {
            if (models == null || models.Count() == 0)
                return;

            string insertSQL = SQLHelper.GetInsertCommand<T>(TableName);

            using (var conn = InstanceSQLConn())
            {
                conn.Open();
                using (var trans = conn.BeginTransaction())
                {
                    var resultId = conn.Execute(insertSQL, models , trans);
                    trans.Commit();
                    return;
                }
            }
        }


        public bool Update(T model)
        {
            var sql = SQLHelper.GetUpdateCommand<T>(TableName);
            using (var conn = InstanceSQLConn())
            {
                return conn.Execute(sql, model) == 1;
            }
        }


        public bool Update(long id, object dynamicObject)
        {
            var allPropertiesOfClass = typeof(T).GetProperties().Select(x => x.Name);

            List<string> updateColumns = new List<string>();

            //先檢查傳入的參數 是否為class定義的db raw entity欄位
            foreach (var item in dynamicObject.GetType().GetProperties())
            {
                string columnName = item.Name;
                bool isExist = allPropertiesOfClass.Any(z => z == columnName);

                if (!isExist)
                {
                    throw new TerproException(ApiErrorCode.Forbidden_NotValidDefineValue);
                }

                updateColumns.Add(columnName);
            }

            string updateColumnsSQL = string.Join(" , ", updateColumns.Select(x => $"[{x}] = @{x}"));
            string updateSQL = $" UPDATE {TableName} SET {updateColumnsSQL} WHERE {PrimaryKeyColumnName} = {id}";

            int result = ExecuteSqlByDapper(updateSQL, dynamicObject);
            return result == 1;
        }

        public bool UpdateMany(IEnumerable<T> models)
        {
            if (models == null || models.Count() == 0)
                return true;

            var sql = SQLHelper.GetUpdateCommand<T>(TableName);
            using (var conn = InstanceSQLConn())
            {
                return conn.Execute(sql, models) == 1;
            }
        }

        public void DeleteByModel(object model)
        {
            string deleteSql = $"DELETE {TableName} WHERE ";

            List<string> columns = new List<string>();

            foreach (var prop in model.GetType().GetProperties())
            {
                string columnName = prop.Name;
                columns.Add($" {columnName} = @{columnName} ");
            }

            deleteSql += string.Join(" AND ", columns);

            using (var conn = InstanceSQLConn())
            {
                conn.Execute(deleteSql, model);
            }
        }

        public void DeleteByKeyValues(Dictionary<string , object> keyValues)
        {
            string deleteSql = $"DELETE {TableName} WHERE ";

            List<string> columns = new List<string>();

            foreach (var eachValue in keyValues)
            {
                string columnName = eachValue.Key;
                columns.Add($" {columnName} = @{columnName} ");
            }

            deleteSql += string.Join(" AND ", columns);

            using (var conn = InstanceSQLConn())
            {
                conn.Execute(deleteSql, keyValues);
            }
        }


        public List<T> GetAll()
        {
            string sql = $@"SELECT * FROM {TableName}";
            return QueryByDapper<T>(sql, new {  }).ToList();
        }


        public List<T> GetByLargeAmountIDs(List<long> ids)
        {
            if (ids == null || ids.Count() == 0)
                return new List<T>();

            string PKColumnName = GetKeyAttributeName();
            string sql = $@"SELECT * FROM {TableName} 
                            WHERE {PKColumnName} in ( SELECT Id FROM #tempIDs ) ";

            return ExecuteQueryWithLargeIDs<T>(sql, ids);
        }

        
        /// <summary>
        /// in query限制最多2100筆，因此使用tempTable的機制查詢大量IDs
        /// </summary>
        private List<T> ExecuteQueryWithLargeIDs<T>(string sql, List<long> ids) where T : class
        {
            using (var db = InstanceSQLConn())
            {
                var idCoverts = ids.Distinct().Select(x => new { Id = x }).ToList();
                db.Open();
                db.Execute("CREATE TABLE #tempIDs(Id int not null primary key);");
                db.Execute("INSERT INTO #tempIDs VALUES(@Id);", idCoverts);
                return db.Query<T>(sql).ToList();
            }
        }
        

        protected List<T> QueryByDapper<T>(string sql , object parameters) 
        {
            using (var conn = InstanceSQLConn())
            {
                var query = conn.Query<T>(sql, parameters);
                return query.ToList();
            }
        }

        protected int ExecuteSqlByDapper(string sql, object parameters)
        {
            using (var conn = InstanceSQLConn())
            {
                var result = conn.Execute(sql, parameters);
                return result;
            }
        }


        protected int ExecuteStoreProcedureByDapper(string spName, object parameters)
        {
            using (var conn = InstanceSQLConn())
            {
                var result = conn.Execute(spName, parameters, commandType: CommandType.StoredProcedure);
                return result;
            }
        }


        protected List<S> ExecuteStoreProcedureResultByDapper<S>(string spName, object parameters)
        {
            using (var conn = InstanceSQLConn())
            {
                var result = conn.Query<S>(spName, parameters, commandType: CommandType.StoredProcedure);
                return result.ToList();
            }
        }


        protected string GenerateColumnEqualWhereQuery(string columnsName) 
        {
            string sql = $@" SELECT * FROM {TableName} 
                             WHERE {columnsName} = @{columnsName} ";

            return sql;
        }


        protected string GenerateColumnInWhereQuery(string columnsName) 
        {
            string sql = $@" SELECT * FROM {TableName} 
                             WHERE {columnsName} in @{columnsName} ";

            return sql;
        }


        protected string GetKeyAttributeName()
        {
            List<string> keyAttr = new List<string>();

            foreach(PropertyInfo item in typeof(T).GetProperties())
            {
                var matchAttr = item.GetCustomAttributes<KeyAttribute>(true);

                if (matchAttr == null || matchAttr.Count() == 0)
                    continue;

                keyAttr.Add(item.Name);
            }

            if(keyAttr.IsEmptyOrNull())
                throw new Exception("not define key columns (PK) in class. ");

            //多key的table禁止使用此function
            if(keyAttr.Count() >= 2)
                throw new Exception("multiple  key columns (PK) define , forbidden. ");

            return keyAttr.FirstOrDefault();
        }



        /// <summary>
        /// 分頁
        /// </summary>
        protected PagedList<S> GetPagingEntitiesBySQLScript<S>(string qrySQL, PaginationParameter pager, object param = null) where S : class
        {
            string sortDirection = pager.IsASC ? "asc" : "desc";
            string[] strArr = qrySQL.Split(new string[] { "FROM", "from", "From" }, StringSplitOptions.None);
            
            if (string.IsNullOrWhiteSpace(pager.OrderByColumn))
            {
                pager.OrderByColumn = DefaultSortColumn;
            }
            else
            {
                //姑且防範下可能的 SQL injection
                if(pager.OrderByColumn.Length >= 25)
                {
                    throw new TerproException(ApiErrorCode.BusinessLogicDeny);
                }
            }

            string str1 = strArr[0] + $",(ROW_NUMBER() OVER (ORDER BY {pager.OrderByColumn} {sortDirection} )) as RowNum ";

            StringBuilder sb = new StringBuilder();
            //跳過第一個
            for (int i = 1; i < strArr.Length; i++)
            {
                sb.Append(@" FROM" + strArr[i]);
            }

            string str2 = sb.ToString();
            qrySQL = str1 + str2;

            string countSQL = $@"
                              SELECT COUNT(0) as cnt 
                              FROM ({qrySQL}) as DataCount";

            return GetPagingEntitiesBySQLScript<S>(qrySQL, countSQL, pager.PageIndex, pager.PageSize, param);
        }


        /// <summary>
        /// 依條件查詢多筆資料
        /// </summary>
        private PagedList<S> GetPagingEntitiesBySQLScript<S>(string sqlScript, string countSQL, int currentPage, int pageSize, object param = null) where S : class
        {
            currentPage = currentPage <= 0 ? 1 : currentPage;
            pageSize = pageSize <= 0 ? 5 : pageSize;
            PagedList<S> pageResult = new PagedList<S>();
            int total = GetScalarBySQLScript<int>(countSQL, param);
            pageResult = new PagedList<S>
            {
                PageIndex = currentPage,
                PageSize = pageSize,
                TotalCount = GetScalarBySQLScript<int>(countSQL, param),
                TotalPage = CalculateTotalPage(total, pageSize)
            };

            //如果要追SQL很慢的問題 
            pageResult.Items = GetPageingEntities<S>(sqlScript, currentPage, pageSize, param);
            return pageResult;
        }


        private int CalculateTotalPage(int totalCount , int pageSize)
        {
            int totalPage = ((int)totalCount / pageSize);

            if (totalCount % pageSize == 0)
                return totalPage;

            return totalPage + 1;
        }


        /// <summary>
        /// 依條件查詢第一筆第一個欄位資料
        /// </summary>
        private S GetScalarBySQLScript<S>(string _SQLScript, object param = null)
        {
            using (SqlConnection conn = InstanceSQLConn())
            {
                return conn.ExecuteScalar<S>(_SQLScript, param);
            }
        }




        /// <summary>
        /// 依條件查詢多筆資料(分頁用)
        /// </summary>
        private List<S> GetPageingEntities<S>(string sqlLScript, int currentPage, int pageSize, object param = null)
        {
            int startRowNum = (currentPage <= 1) ? 1 : 1 + (currentPage - 1) * pageSize;
            int endRowNum = (startRowNum - 1) + pageSize;
            string paggingSQL = @"
                        SELECT *
                        FROM ( " + sqlLScript + @" ) as pageData
                        WHERE RowNum >= " + startRowNum + @" AND RowNum <= " + endRowNum;
            
            return GetEntitiesBySQLScript<S>(paggingSQL, param);
        }



        /// <summary>
        /// 依條件查詢多筆資料
        /// </summary>
        private List<S> GetEntitiesBySQLScript<S>(string _SQLScript, object param = null)
        {
            using (SqlConnection conn = InstanceSQLConn())
            {
                return conn.Query<S>(_SQLScript, param).ToList();
            }
        }


        protected string GenerateLikeQuery(string columnName , string parameterName)
        {
            return $" ({columnName} LIKE CONCAT('%',{parameterName} ,'%')) ";
        }

        protected string GetNoLockSQL(bool isNoLock)
        {
            return isNoLock ? " WITH (NOLOCK) " : "";
        }

        protected DynamicParameters InstanceDynamicParameterWithReturnValue()
        {
            var para = new DynamicParameters();

            //似乎這邊add的順序 有可能導致output參數轉型出問題(?) 感覺是dapper的bug
            para.Add(name: "@RetVal", dbType: DbType.Int32, direction: ParameterDirection.ReturnValue);
            return para;
        }

    }
    

}
