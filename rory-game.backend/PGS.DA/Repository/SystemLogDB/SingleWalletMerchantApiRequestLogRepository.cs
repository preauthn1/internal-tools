using Dapper;
using game.Entity;
using Lck.Utility.Extensions.Helper;
using target.Core.Repository;
using target.Models.Enums;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.DA.Repository
{
    public class SingleWalletMerchantApiRequestLogRepository : FormatIntegerPrimaryKeyRepository<SingleWalletMerchantApiRequestLogEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.SystemLogDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[SingleWalletMerchantApiRequestLog]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<SingleWalletMerchantApiRequestLogEntity> GetByFilter(QuerySingleWalletMerchantApiRequestLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<SingleWalletMerchantApiRequestLogEntity>(sql, filter);
            return query;
        }


        public PagedList<SingleWalletMerchantApiRequestLogEntity> GetPagingByFilter(QuerySingleWalletMerchantApiRequestLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<SingleWalletMerchantApiRequestLogEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QuerySingleWalletMerchantApiRequestLogRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
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
