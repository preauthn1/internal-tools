using Dapper;
using game.Entity;
using target.Core.Repository;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Models.CoreLibs;
using System;
using System.Collections.Generic;
using target.Utility.Extensions;
using System.Linq;


namespace target.DA.Repository
{

    /// <summary>
    /// 管理某些老虎機遊戲的機率與項目設置 (大概?)
    /// </summary>
    public class GameSlotProbabilitySettingRepository : FormatIntegerPrimaryKeyRepository<GameSlotProbabilitySettingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[GameSlotProbabilitySetting]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "ID"; }
        }

        public List<GameSlotProbabilitySettingEntity> GetByFilter(QueryGameSlotProbabilitySettingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<GameSlotProbabilitySettingEntity>(sql, filter);
            return query;
        }

        private string GenerateWhereSQL(QueryGameSlotProbabilitySettingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} " + 
                           "WHERE 1=1 ";

            if (filter.PlayerID.HasValue)
                sql += " AND PlayerID = @PlayerID ";

            if (filter.GameID.HasValue)
                sql += " AND GameID = @GameID ";

            if (filter.FlagID.HasValue)
                sql += " AND FlagID = @FlagID ";

            if (filter.OriginID.HasValue)
                sql += " AND OriginID = @OriginID ";

            return sql;
        }

    }

}
