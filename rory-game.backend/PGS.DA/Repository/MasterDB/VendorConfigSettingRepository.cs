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
using target.Models.Models;
using target.Models.Enums;

namespace target.DA.Repository
{

    public class VendorConfigSettingRepository : FormatIntegerPrimaryKeyRepository<VendorConfigSettingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[ThirdParty].[VendorConfigSetting]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }
        
            
        /// <summary>
        /// 取得指定的 merchantID + vendorID的 所有設置
        /// </summary>
        public List<MerchantVendorSettingModel> GetTargetMerchantVendorSettings(int merchantID , GameVendorTypeEnum vendorID)
        {
            string sql = @" SELECT [SettingKey] ,[Value] as SettingValue
                            FROM [ThirdParty].[VendorConfigMerchantMapping] as mapping
                            INNER JOIN [ThirdParty].[VendorConfigSetting] as setting on mapping.VendorConfigID = setting.VendorConfigID
                            WHERE MerchantID = @merchantID AND [VendorID] = @vendorID";

            var para = new { merchantID = merchantID, vendorID = vendorID.ToInteger() };
            var query = QueryByDapper<MerchantVendorSettingModel>(sql, para);
            return query;
        }
        
        
        public List<VendorConfigSettingEntity> GetByFilter(QueryVendorConfigSettingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<VendorConfigSettingEntity>(sql, filter);
            return query;
        }


        public PagedList<VendorConfigSettingEntity> GetPagingByFilter(QueryVendorConfigSettingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<VendorConfigSettingEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryVendorConfigSettingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.VendorConfigIDs.IsNotEmptyNull())
                sql += " AND VendorConfigID in @VendorConfigIDS ";

            if (filter.SettingKeys.IsNotEmptyNull())
                sql += " AND SettingKey in @SettingKeys ";

            return sql;
        }

    }


}
