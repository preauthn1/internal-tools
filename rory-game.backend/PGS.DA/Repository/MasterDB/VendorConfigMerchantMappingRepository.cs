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
    public class VendorConfigMerchantMappingRepository : FormatIntegerPrimaryKeyRepository<VendorConfigMerchantMappingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[ThirdParty].[VendorConfigMerchantMapping]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }


        public List<int> GetNotSettingMerchants(int vendorID)
        {
            int deleteStatus = NormalActiveStatus.PermanentDelete.ToInteger();
            string sql = $@"SELECT MerchantID FROM[dbo].[Merchant]
                            WHERE Status != {deleteStatus} AND MerchantID not in
                            (
                                SELECT MerchantID 
                                FROM {TableName}
                                WHERE VendorID = @vendorID
                            )";

            var query = QueryByDapper<int>(sql, new { vendorID  });
            return query;
        }


        /// <summary>
        /// 幫新建的商戶 預設塞入所有 遊戲商配置預設值
        /// </summary>
        public void CreateMerchantAllDefaultValue(int merchantID)
        {
            int deleteStatus = NormalActiveStatus.PermanentDelete.ToInteger();
            string sql = $@"INSERT INTO {TableName} ([MerchantID] , [VendorID] , [VendorConfigID])
                            SELECT @merchantID , VendorID , VendorConfigID 
                            FROM [ThirdParty].[VendorConfig] 
                            WHERE IsDefault = 1";

            ExecuteSqlByDapper(sql , new { merchantID });
        }


        public List<VendorConfigMerchantMappingEntity> GetByFilter(QueryVendorConfigMerchantMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<VendorConfigMerchantMappingEntity>(sql, filter);
            return query;
        }


        public PagedList<VendorConfigMerchantMappingEntity> GetPagingByFilter(QueryVendorConfigMerchantMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<VendorConfigMerchantMappingEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryVendorConfigMerchantMappingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.VendorIDs.IsNotEmptyNull())
                sql += " AND VendorID in @VendorIDs ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.VendorConfigIDs.IsNotEmptyNull())
                sql += " AND VendorConfigID in @VendorConfigIDs ";

            if(!filter.IsIncludeDeleteMerchant)
                sql += $" AND MerchantID not in ( SELECT MerchantID FROM [dbo].[Merchant] WHERE Status = {NormalActiveStatus.PermanentDelete.ToInteger()}) ";

            return sql;
        }

    }


}
