using Dapper;
using MSC.Core.Repository;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.DA.Repository
{
    public class SystemWarningNotifyRepository : IDapperRepository<SystemWarningNotifyEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[SystemWarningNotify]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<SystemWarningNotifyEntity> GetByFilter(QuerySystemWarningNotifyRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<SystemWarningNotifyEntity>(sql, filter);
            return query;
        }


        public PagedList<SystemWarningNotifyEntity> GetPagingByFilter(QuerySystemWarningNotifyRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<SystemWarningNotifyEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QuerySystemWarningNotifyRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} WITH (NOLOCK) 
                           WHERE 1=1 ";

            if(filter.IsNotSend)
                sql += " AND isSend = 0 ";

            return sql;
        }

    }


}
