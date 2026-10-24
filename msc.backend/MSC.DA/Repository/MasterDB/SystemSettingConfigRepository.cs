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
using MSC.Models.Enums;
using MSC.Models.Models;

namespace MSC.DA.Repository
{

    public class SystemSettingConfigRepository : IDapperRepository<SystemSettingConfigEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[SystemSettingConfig]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "SettingKey"; }
        }

        public CustomizedSettingParseModel GetSingleByCache(SystemSettingConfigEnum key , int cacheTimeMinutes = 5)
        {
            int keyInt = key.ToInteger();

            string cacheKey = CacheKey.DB_Data_SystemSettingConfig(keyInt);

            var queryCachae = CacheHelper.GetCacheOrQuery(cacheKey, () =>
            {
                var querySingle = GetById(keyInt);
                return querySingle;

            }, expiredMinute : cacheTimeMinutes);

            if (queryCachae == null)
                throw new MSCException(ErrorCode.BusinessLogicDeny_SettingDataNotCompleted);

            return new CustomizedSettingParseModel
            {
                Key = queryCachae.SettingKey,
                SettingValue = queryCachae.SettingValue
            };
        }


        public List<SystemSettingConfigEntity> GetByFilter(QuerySystemSettingConfigRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<SystemSettingConfigEntity>(sql, filter);
            return query;
        }


        public PagedList<SystemSettingConfigEntity> GetPagingByFilter(QuerySystemSettingConfigRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<SystemSettingConfigEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QuerySystemSettingConfigRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.SettingKeyIDs.IsNotEmptyNull())
                sql += " AND SettingKey in @SettingKeyIDs ";

            return sql;
        }

    }


}
