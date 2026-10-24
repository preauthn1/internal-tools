using Dapper;
using MSC.Core.Repository;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using System;
using System.Collections.Generic;
using MSC.Utility.Extensions;
using System.Linq;

namespace MSC.DA.Repository
{

    public class ExternalApiRequestLogRepository : IDapperRepository<ExternalApiRequestLogEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.SystemLogDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[ExternalApiRequestLog]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<ExternalApiRequestLogEntity> GetByFilter(QueryExternalApiRequestLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<ExternalApiRequestLogEntity>(sql, filter);
            return query;
        }


        public PagedList<ExternalApiRequestLogEntity> GetPagingByFilter(QueryExternalApiRequestLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<ExternalApiRequestLogEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryExternalApiRequestLogRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND ID in @IDs ";

            if(filter.Group.IsNotNullNotWhiteSpace())
                sql += " AND [Group] = @Group ";

            if (filter.ApiCategory.IsNotNullNotWhiteSpace())
                sql += " AND ApiCategory = @ApiCategory ";

            return sql;

        }

    }


}
