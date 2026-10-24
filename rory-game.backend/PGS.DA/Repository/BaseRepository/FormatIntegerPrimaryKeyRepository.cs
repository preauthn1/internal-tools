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
    public abstract class FormatIntegerPrimaryKeyRepository<T> : BaseDapperRepository<T> where T: class
    {     

        public int Create(T model)
        {
            string insertSQL = SQLHelper.GetInsertCommand<T>(TableName);
            insertSQL += " SELECT CAST(SCOPE_IDENTITY() as int) ";

            using (var conn = InstanceSQLConn())
            {
                var resultId = conn.Query<int?>(insertSQL, model);

                if (resultId.IsEmptyOrNull())
                    return 0;

                var myValue = resultId.FirstOrDefault();
                return myValue.HasValue ? myValue.Value : 0;
            }
        }



        public bool IsExist(int id)
        {
            string sql = $"SELECT {PrimaryKeyColumnName} FROM {TableName} WITH (NOLOCK) WHERE {PrimaryKeyColumnName} = @id";

            using (var conn = InstanceSQLConn())
            {
                var resultId = conn.Query<int>(sql, new { id  = id }).FirstOrDefault();

                if (resultId != 0)
                    return true;

                return false;
            }
        }


        protected void DeleteAndInsertNewData(string columnName , int keyValue , List<T> models)
        {
            string deleteSql = $"DELETE {TableName} WHERE {columnName} = @keyValue "; 
            string insertSQL = SQLHelper.GetInsertCommand<T>(TableName);            

            using (var conn = InstanceSQLConn())
            {
                conn.Execute(deleteSql, new { keyValue  = keyValue});

                conn.Execute(insertSQL, models);
            }
        }

        public void Delete(int id)
        {
            string deleteSql = $"DELETE {TableName} WHERE {PrimaryKeyColumnName} = @id ";

            using (var conn = InstanceSQLConn())
            {
                conn.Execute(deleteSql, new { id = id });
            }
        }


        public void Delete(List<int> ids)
        {
            string deleteSql = $"DELETE {TableName} WHERE {PrimaryKeyColumnName} in @ids ";

            using (var conn = InstanceSQLConn())
            {
                conn.Execute(deleteSql, new { ids = ids });
            }
        }


        public T GetById(int id)
        {
            string PKColumnName = GetKeyAttributeName();
            string sql = $@"SELECT * FROM {TableName} 
                            WHERE {PKColumnName} = @id ";

            return QueryByDapper<T>(sql , new  { id }).FirstOrDefault();
        }


        public List<T> GetByIds(List<int> ids)
        {
            ids = ids.Distinct().ToList();

            if (ids == null || ids.Count() == 0)
                return new List<T>();

            if (ids.Count() >= 2100)
                return GetByLargeAmountIDs(ids.Select(x=> (long)x).ToList());

            string PKColumnName = GetKeyAttributeName();
            string sql = $@"SELECT * FROM {TableName} 
                            WHERE {PKColumnName} in @ids ";

            return QueryByDapper<T>(sql , new { ids });
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

    }

}
