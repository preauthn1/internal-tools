using Dapper;
using game.Entity;
using System.Linq;
using target.Core.Repository;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using System;
using System.Collections.Generic;
using target.Utility.Extensions;
using target.Models.Enums;
using Lck.Utility.Extensions.Helper;

namespace target.DA.Repository
{

    public class MerchantRepository : FormatIntegerPrimaryKeyRepository<MerchantEntity>
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

        public MerchantBriefEntity GetByTokenAndSecret(int merchantID , string secretKey)
        {
            var filter = new QueryMerchantRequestViewModel
            {
                ApiSecretKey = secretKey,
                MerchantIDs = merchantID.ToListCollection()
            };

            var query = GetBriefByFilter(filter);
            return query.FirstOrDefault();
        }

        public MerchantEntity GetSingleByCache(int merchantID)
        {
            string cacheKey = CacheKey.AdminSys_MerchantInfo(merchantID);

            var queryCachae = CacheHelper.GetCacheOrQuery(cacheKey, () =>
               {
                   var querySingle = GetById(merchantID);
                   return querySingle;
               });

            return queryCachae;
        }
        
        public List<MerchantBriefEntity> GetBriefByFilter(QueryMerchantRequestViewModel filter)
        {
            string selectColumns = " MerchantID , Currency , Status , DefaultLanguage , ApiSecretKey ";
            string sql = GenerateWhereSQL(filter , selectColumns);
            var query = QueryByDapper<MerchantBriefEntity>(sql, filter);
            return query;
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

            var query = GetPagingEntitiesBySQLScript<MerchantEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryMerchantRequestViewModel filter , string columns = "*")
        {
            string sql = $@"SELECT {columns} FROM {TableName} 
                           WHERE 1=1 ";

            if(!string.IsNullOrWhiteSpace(filter.ApiSecretKey))
                sql += " AND ApiSecretKey = @ApiSecretKey ";

            if (!string.IsNullOrWhiteSpace(filter.Currency))
                sql += " AND (Currency = @Currency OR ChargeCurrency = @Currency) ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if(!filter.IsIncludeDeleteData)
                sql += $" AND Status != {NormalActiveStatus.PermanentDelete.ToInteger()} ";

            if (filter.Status.HasValue)
                sql += $" AND Status = @Status ";

            if (filter.WalletMode.HasValue)
                sql += $" AND WalletMode = @WalletMode ";

            if (filter.Type.HasValue)
                sql += $" AND Type = @Type ";

            if (filter.OrganizationIDs.IsNotEmptyNull())
                sql += $" AND MerchantID in (SELECT MerchantID FROM [PGS_Operation].[dbo].[OrganizationMerchantMapping] WHERE organizationID in @OrganizationIDs) ";

            return sql;
        }

    }


}
