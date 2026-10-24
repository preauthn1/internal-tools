using target.Core.Helper;
using target.Core.Repository;
using target.Models.Table;
using target.Utility.Extensions;
using System.Collections.Generic;
using System.Linq;
using Dapper;

namespace target.DA.Repository
{
    public class ApiRequestLogRepository : FormatIntegerPrimaryKeyRepository<ApiRequestLogEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.SystemLogDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[ApiRequestLog]"; }
        }

        protected override string PrimaryKeyColumnName 
        {
            get { return "Id"; }
        }

        public void CreateToTargetTable(string tableName , ApiRequestLogEntity model)
        {
            string insertSQL = SQLHelper.GetInsertCommand<ApiRequestLogEntity>(tableName);
            ExecuteSqlByDapper(insertSQL , model);
        }


        public void CreateToTargetTable(string tableName, IEnumerable<ApiRequestLogEntity> models)
        {
            if (models.IsEmptyOrNull())
                return;

            string insertSQL = SQLHelper.GetInsertCommand<ApiRequestLogEntity>(tableName);

            using (var conn = InstanceSQLConn())
            {
                conn.Open();
                using (var trans = conn.BeginTransaction())
                {
                    var resultId = conn.Execute(insertSQL, models, trans);
                    trans.Commit();
                    return;
                }
            }
        }     


        public List<ApiRequestLogEntity> DebugTest_GetList()
        {
            string sql = "SELECT * FROM  [dbo].[ApiRequestLogGameAPI] WHERE Id <= 1500";
            var query = QueryByDapper<ApiRequestLogEntity>(sql , null);
            return query;
        }


    }


}
