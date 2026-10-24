using Dapper;
using MSC.Core.Repository;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using System;
using System.Collections.Generic;
using MSC.Utility.Extensions;
using System.Linq;
using MSC.Models.Models;
using MSC.Models.Enums;

namespace MSC.DA.Repository
{

    public class MerchantCustomizedSettingRepository : IDapperRepository<MerchantCustomizedSettingEntity>
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

        public MerchantCustomizedSettingConvertModel GetByMerchantSettingKey(int merchantID, MerchantCustomizedSettingKeyEnum key)
        {
            var filter = new QueryMerchantCustomizedSettingRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection(),
                CustomizedSettingKeyIDs = key.ToInteger().ToListCollection()
            };

            var query = GetByFilter(filter).FirstOrDefault();

            if (query == null)
                return null;

            return query.ReturnInjecteObject<MerchantCustomizedSettingConvertModel>();
        }


        public Dictionary<int , List<MerchantCustomizedSettingConvertModel>> GetMultipleMerchants(List<int> merchantIDs, List<int> keys)
        {
            var filter = new QueryMerchantCustomizedSettingRequestViewModel
            {
                MerchantIDs = merchantIDs,
                CustomizedSettingKeyIDs = keys
            };

            var query = GetByFilter(filter);

            if (query.IsEmptyOrNull())
                return new Dictionary<int, List<MerchantCustomizedSettingConvertModel>>();
                       
            //key: merchantID value: 設置List
            return merchantIDs.ToDictionary(x => x, z =>
            {
                var matchThisMerchant = query.Where(y => y.MerchantID == z)
                    .Select(x => x.ReturnInjecteObject<MerchantCustomizedSettingConvertModel>())
                    .ToList();

                return matchThisMerchant;
            });
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

            var query = GetPageingEntitiesBySQLScript<MerchantCustomizedSettingEntity>(sql, filter, filter);
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

            return sql;
        }

    }


}
