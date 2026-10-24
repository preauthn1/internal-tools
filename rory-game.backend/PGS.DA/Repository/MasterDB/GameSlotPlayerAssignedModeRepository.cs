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

namespace target.DA.Repository
{

    public class GameSlotPlayerAssignedModeRepository : FormatIntegerPrimaryKeyRepository<GameSlotPlayerAssignedModeEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[GameSlotPlayerAssignedMode]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "ID"; }
        }

        public List<GameSlotPlayerAssignedModeEntity> GetByFilter(QueryPlayerSlotAssignedModeRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<GameSlotPlayerAssignedModeEntity>(sql, filter);
            return query;
        }

        private string GenerateWhereSQL(QueryPlayerSlotAssignedModeRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} " + 
                           "WHERE 1 = 1 ";

            if (filter.PlayerID.HasValue)
                sql += " AND PlayerID = @PlayerID ";

            if (filter.GameID.HasValue)
                sql += " AND GameID = @GameID ";

            if (filter.FlagID.HasValue)
                sql += " AND FlagID = @FlagID ";

            return sql;
        }

    }

}
