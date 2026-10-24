using Dapper;
using MSC.Core.Repository;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using System;
using System.Collections.Generic;
using MSC.Utility.Extensions;
using System.Linq;
using Lck.Utility.Extensions.Helper;

namespace MSC.DA.Repository
{

    public class MerchantRepository : IDapperRepository<MerchantEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[Merchant]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "MerchantID"; }
        }

        public MerchantEntity GetSingleByCache(int merchantID)
        {
            string cacheKey = CacheKey.DB_Data_Merchant(merchantID);

            var queryCachae = CacheHelper.GetCacheOrQuery(cacheKey, () =>
            {
                var querySingle = GetById(merchantID);
                return querySingle;
            });

            return queryCachae;
        }



        public List<MerchantEntity> GetByFilter(QueryMerchantRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<MerchantEntity>(sql, filter);
            return query;
        }


        public PagedList<MerchantEntity> GetPagingByFilter(QueryMerchantRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<MerchantEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryMerchantRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.Status.HasValue)
                sql += " AND Status = @Status ";

            return sql;
        }

    }


}
