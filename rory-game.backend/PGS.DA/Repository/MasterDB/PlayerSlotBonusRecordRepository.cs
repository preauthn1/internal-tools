using Dapper;
using game.Entity;
using target.Core.Repository;
using target.Models.Paging;
using target.Models.Table;
using target.Models.CoreLibs;
using target.Models.ViewModels;
using System;
using System.Collections.Generic;
using target.Utility.Extensions;
using System.Linq;

namespace target.DA.Repository
{

    public class PlayerSlotBonusRecordRepository : FormatIntegerPrimaryKeyRepository<PlayerSlotBonusRecordEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[PlayerSlotBonusRecord]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "ID"; }
        }

        public List<PlayerSlotBonusRecordEntity> GetByFilter(QueryPlayerSlotBonusRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<PlayerSlotBonusRecordEntity>(sql, filter);
            return query;
        }

        private string GenerateWhereSQL(QueryPlayerSlotBonusRecordRequestViewModel filter)
        {
            string noLock = filter.IsLockQuery ? "" : " WITH (NOLOCK) ";

            string sql = $@"SELECT * FROM {TableName} {noLock} " + 
                           "WHERE 1=1 ";

            if (filter.PlayerID.HasValue)
                sql += " AND PlayerID = @PlayerID ";

            if (filter.GameID.HasValue)
                sql += " AND GameID = @GameID ";

            if (filter.FlagID.HasValue)
                sql += " AND FlagID = @FlagID ";

            if (filter.OriginID.HasValue)
                sql += " AND OriginID = @OriginID ";

            if (filter.TurnID.HasValue)
                sql += " AND TurnID = @TurnID ";

            if (filter.Location.HasValue)
                sql += " AND Location = @Location AND Location <> 0 ";

            if (filter.BetAmount.HasValue)
                sql += " AND BetAmount = @BetAmount ";

            if (filter.BetRate.HasValue)
                sql += " AND BetRate = @BetRate ";

            return sql;
        }

    }

}
