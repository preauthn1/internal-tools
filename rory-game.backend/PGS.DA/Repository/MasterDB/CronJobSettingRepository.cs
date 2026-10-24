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
    public class CronJobSettingRepository : FormatIntegerPrimaryKeyRepository<CronJobSettingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[CronJobSetting]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "JobID"; }
        }

        public List<CronJobSettingEntity> GetByFilter(QueryCronJobSettingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<CronJobSettingEntity>(sql, filter);
            return query;
        }


        public PagedList<CronJobSettingEntity> GetPagingByFilter(QueryCronJobSettingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<CronJobSettingEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryCronJobSettingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND JobID in @IDs ";

            if (filter.IsEnable.HasValue)
                sql += " AND IsEnable = @IsEnable ";
            
            return sql;
        }

    }


}
