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
    public abstract class FormatLongPrimaryKeyRepository<T> : BaseDapperRepository<T> where T : class
    {
        public long Create(T model)
        {
            string insertSQL = SQLHelper.GetInsertCommand<T>(TableName);
            insertSQL += " SELECT CAST(SCOPE_IDENTITY() as bigint) ";

            using (var conn = InstanceSQLConn())
            {
                var resultId = conn.Query<long?>(insertSQL, model);

                if (resultId.IsEmptyOrNull())
                    return 0;

                var myValue = resultId.FirstOrDefault();
                return myValue.HasValue ? myValue.Value : 0;
            }
        }


        public bool IsExist(long id)
        {
            string sql = $"SELECT {PrimaryKeyColumnName} FROM {TableName} WITH (NOLOCK) WHERE {PrimaryKeyColumnName} = @id";

            using (var conn = InstanceSQLConn())
            {
                var resultId = conn.Query<long>(sql, new { id = id }).FirstOrDefault();

                if (resultId != 0)
                    return true;

                return false;
            }
        }


        protected void DeleteAndInsertNewData(string columnName, long keyValue, List<T> models)
        {
            string deleteSql = $"DELETE {TableName} WHERE {columnName} = @keyValue ";
            string insertSQL = SQLHelper.GetInsertCommand<T>(TableName);

            using (var conn = InstanceSQLConn())
            {
                conn.Execute(deleteSql, new { keyValue = keyValue });

                conn.Execute(insertSQL, models);
            }
        }

        public void Delete(long id)
        {
            string deleteSql = $"DELETE {TableName} WHERE {PrimaryKeyColumnName} = @id ";

            using (var conn = InstanceSQLConn())
            {
                conn.Execute(deleteSql, new { id = id });
            }
        }


        public void Delete(List<long> ids)
        {
            string deleteSql = $"DELETE {TableName} WHERE {PrimaryKeyColumnName} in @ids ";

            using (var conn = InstanceSQLConn())
            {
                conn.Execute(deleteSql, new { ids = ids });
            }
        }


        public T GetById(long id)
        {
            string PKColumnName = GetKeyAttributeName();
            string sql = $@"SELECT * FROM {TableName} 
                            WHERE {PKColumnName} = @id ";

            return QueryByDapper<T>(sql, new { id }).FirstOrDefault();
        }


        public List<T> GetByIds(List<long> ids)
        {
            ids = ids.Distinct().ToList();

            if (ids == null || ids.Count() == 0)
                return new List<T>();

            if (ids.Count() >= 2100)
                return GetByLargeAmountIDs(ids.Select(x => (long)x).ToList());

            string PKColumnName = GetKeyAttributeName();
            string sql = $@"SELECT * FROM {TableName} 
                            WHERE {PKColumnName} in @ids ";

            return QueryByDapper<T>(sql, new { ids });
        }


    }

}
