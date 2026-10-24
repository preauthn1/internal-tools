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

    public class SystemConfigSettingRepository : FormatIntegerPrimaryKeyRepository<SystemConfigSettingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[SystemConfigSetting]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public bool IsSystemMaintain()
        {
            var query = GetParseModelByEnum(SystemConfigSettingEnum.SystemMaintain_IsSwtichToMaintain);

            if (query == null)
                return false;

            return query.SettingValueOfInt != (int)SystemMaintainStatusEnum.NormalStatus;
        }

        public void CheckSystemMaintainAndThrowException()
        {
            bool isSystemMaintain = IsSystemMaintain();

            if (isSystemMaintain)
                throw new TerproException(ApiErrorCode.SystemMaintain);

            return;
        }


        public SystemConfigSettingParseModel GetParseModelByCache(SystemConfigSettingEnum value)
        {
            string cacheKey2 = CacheKey.DbEntity_ParseModel_SystemConfigSetting(value.ToInteger());

            var cacheValue = CacheHelper.GetCacheOrQuery(cacheKey2, () =>
               {
                   var query = GetParseModelByEnum(value);
                   return query;
               });

            return cacheValue;
        }


        public SystemConfigSettingParseModel GetParseModelByEnum(SystemConfigSettingEnum value)
        {
            var query = GetById(value.ToInteger());

            if (query == null)
                return null;

            return new SystemConfigSettingParseModel
            {
                SettingKey = query.ID,
                SettingValue = query.SettingValue
            };
        }


        public SystemConfigSettingEntity GetEntityByEnum(SystemConfigSettingEnum value)
        {
            var query = GetById(value.ToInteger());
            return query;
        }

        
        public List<SystemConfigSettingEntity> GetByFilter(QuerySystemConfigSettingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<SystemConfigSettingEntity>(sql, filter);
            return query;
        }


        public PagedList<SystemConfigSettingEntity> GetPagingByFilter(QuerySystemConfigSettingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<SystemConfigSettingEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QuerySystemConfigSettingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND ID in @IDs ";

            return sql;
        }

    }


}
