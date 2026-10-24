using Dapper;
using MSC.Core.Repository;
using MSC.Models.Table;
using System;
using System.Collections.Generic;
using System.Data;
using System.Data.Common;
using System.Linq;
using System.Reflection;
using MSC.Utility.Extensions;
using MSC.Core.Helper;

namespace MSC.DA.Repository
{
    public class ApiRequestLogRepository : IDapperRepository<ApiRequestLogEntity>
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

        /// <summary>
        /// 取得時間段內，每一隻API的請求量
        /// </summary>
        public List<QueryApiRequestLogEachApiTotalCountModel> GetEachApiRequestCountDuration(string tableName, DateTime begin , DateTime end)
        {
            string sql = $@" SELECT RelativeUrl , count(*) as totalCount 
                             FROM {tableName} with (NOLOCK)
                             WHERE CreateTime >= @begin AND CreateTime < @end
                             GROUP BY RelativeUrl";

            var para = new 
            {
                begin = begin,
                end = end
            };

            var query = QueryByDapper<QueryApiRequestLogEachApiTotalCountModel>(sql, para);
            return query.ToList();
        }


    }


}
