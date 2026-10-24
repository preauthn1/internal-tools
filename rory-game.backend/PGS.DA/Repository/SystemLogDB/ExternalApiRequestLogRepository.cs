using Dapper;
using game.Entity;
using target.Core.Repository;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using System;
using System.Collections.Generic;
using target.Utility.Extensions;
using System.Linq;

namespace target.DA.Repository
{

    public class ExternalApiRequestLogRepository : FormatLongPrimaryKeyRepository<ExternalApiRequestLogEntity>
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

            var query = GetPagingEntitiesBySQLScript<ExternalApiRequestLogEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryExternalApiRequestLogRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} WITH (NOLOCK)
                           WHERE 1=1 ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND ID in @IDs ";

            if (filter.RequestID.IsNotNullEmptyText())
                sql += " AND RequestID = @RequestID ";

            if (filter.Category.IsNotNullEmptyText())
                sql += " AND Category = @Category ";

            if (filter.SubCategory.IsNotNullEmptyText())
                sql += " AND SubCategory = @SubCategory ";

            if (filter.CreateTimeBegin.HasValue)
                sql += " AND CreateTime >= @CreateTimeBegin ";

            if (filter.CreateTimeEnd.HasValue)
                sql += " AND CreateTime < @CreateTimeEnd ";

            return sql;
        }

    }


}
