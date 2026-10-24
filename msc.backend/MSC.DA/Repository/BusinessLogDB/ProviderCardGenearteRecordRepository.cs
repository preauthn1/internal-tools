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
    public class ProviderCardGenearteRecordRepository : IDapperRepository<ProviderCardGenearteRecordEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.BusinessLogDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[ProviderCardGenearteRecord]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public  bool IsProviderOrderExist(int providerID , string transactionNo)
        {
            var query = GetByFilter(new QueryProviderCardGenearteRecordRequestViewModel
            {
                ProviderIDs = providerID.ToListCollection(),
                TransactionNo = transactionNo
            });

            return query.FirstOrDefault() != null;
        }


        public List<ProviderCardGenearteRecordEntity> GetByFilter(QueryProviderCardGenearteRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<ProviderCardGenearteRecordEntity>(sql, filter);
            return query;
        }


        public PagedList<ProviderCardGenearteRecordEntity> GetPagingByFilter(QueryProviderCardGenearteRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<ProviderCardGenearteRecordEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryProviderCardGenearteRecordRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.ProviderIDs.IsNotEmptyNull())
                sql += " AND ProviderID in @ProviderIDs ";

            if (filter.TransactionNo.IsNotNullNotWhiteSpace())
                sql += " AND TransactionNo = @TransactionNo ";

            return sql;
        }

    }


}
