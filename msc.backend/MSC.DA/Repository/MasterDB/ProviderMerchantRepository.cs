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
    public class ProviderMerchantRepository : IDapperRepository<ProviderMerchantEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[ProviderMerchant]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "ProviderMerchantID"; }
        }

        public List<ProviderMerchantEntity> GetByFilter(QueryProviderMerchantRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<ProviderMerchantEntity>(sql, filter);
            return query;
        }


        public PagedList<ProviderMerchantEntity> GetPagingByFilter(QueryProviderMerchantRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<ProviderMerchantEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryProviderMerchantRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.ProviderMerchantIDs.IsNotEmptyNull())
                sql += " AND ProviderMerchantID in @ProviderMerchantIDs ";

            return sql;
        }

    }


}
