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
using target.Models.Enums;
using target.Models.Models;
using Lck.Utility.Extensions.Helper;

namespace target.DA.Repository
{

    public class MerchantCustomizedSettingRepository : FormatIntegerPrimaryKeyRepository<MerchantCustomizedSettingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[MerchantCustomizedSetting]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<MerchantCustomizedSettingModel> GetMultipleByMerchantSettingKey(List<int> merchantIDs, MerchantCustomizedSettingKeyEnum key)
        {
            var filter = new QueryMerchantCustomizedSettingRequestViewModel
            {
                MerchantIDs = merchantIDs,
                CustomizedSettingKeyIDs = key.ToInteger().ToListCollection()
            };

            var query = GetByFilter(filter);
            return query.Select(x=> x.ReturnInjecteObject<MerchantCustomizedSettingModel>()).ToList();
        }

        public bool IsInternalTestMerchant(int merchantID)
        {
            var query = GetMerchantSettingByCache(merchantID, MerchantCustomizedSettingKeyEnum.IsInternalTestMerchant);

            if (query == null)
                return false;

            return query.SettingValueOfBoolean;
        }


        public decimal GetMerchantGamePointConvertRateByCache(int merchantID, MerchantCustomizedSettingKeyEnum key)
        {
            var querySetting = GetMerchantSettingByCache(merchantID , key);

            if (querySetting == null || string.IsNullOrWhiteSpace(querySetting.SettingValue))
                return 1;

            return querySetting.SettingValueOfDecimal;
        }

        public MerchantCustomizedSettingModel GetMerchantSettingByCache(int merchantID, MerchantCustomizedSettingKeyEnum key)
        {
            string cacheKey = CacheKey.DbEntity_ParseModel_MerchantCustomizedSetting(merchantID , key.ToInteger());

            var queryCachae = CacheHelper.GetCacheOrQuery(cacheKey, () =>
            {
                var querySingle = GetByMerchantSettingKey(merchantID , key);
                return querySingle;

            }, expiredMinute: 5);

            return queryCachae;
        }

        public MerchantCustomizedSettingModel GetByMerchantSettingKey(int merchantID , MerchantCustomizedSettingKeyEnum key)
        {
            var filter = new QueryMerchantCustomizedSettingRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection(),
                CustomizedSettingKeyIDs = key.ToInteger().ToListCollection()
            };

            var query = GetByFilter(filter).FirstOrDefault();

            if (query == null)
                return null;

            return query.ReturnInjecteObject<MerchantCustomizedSettingModel>();
        }


        public List<MerchantCustomizedSettingModel> GetMultipleMerchantSettings(int merchantID, List<MerchantCustomizedSettingKeyEnum> keys)
        {
            var filter = new QueryMerchantCustomizedSettingRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection(),
                CustomizedSettingKeyIDs = keys.Select(x=> x.ToInteger()).ToList()
            };

            var query = GetByFilter(filter);

            if (query.IsEmptyOrNull())
                return new List<MerchantCustomizedSettingModel>();

            return query.Select(x => x.ReturnInjecteObject<MerchantCustomizedSettingModel>()).ToList(); ;
        }


        public List<MerchantCustomizedSettingEntity> GetByFilter(QueryMerchantCustomizedSettingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<MerchantCustomizedSettingEntity>(sql, filter);
            return query;
        }


        public PagedList<MerchantCustomizedSettingEntity> GetPagingByFilter(QueryMerchantCustomizedSettingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<MerchantCustomizedSettingEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryMerchantCustomizedSettingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.CustomizedSettingKeyIDs.IsNotEmptyNull())
                sql += " AND CustomizedSettingKey in @CustomizedSettingKeyIDs ";

            if (filter.SettingValueMatch.IsNotNullEmptyText())
                sql += " AND SettingValue = @SettingValueMatch ";

            return sql;
        }

    }


}
