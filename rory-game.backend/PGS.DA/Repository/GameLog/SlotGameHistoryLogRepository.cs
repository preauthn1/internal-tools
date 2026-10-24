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

    public class SlotGameHistoryLogRepository : FormatLongPrimaryKeyRepository<SlotGameHistoryLogEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.GameLogDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[SlotGameHistoryLog]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "ID"; }
        }

        public List<SlotGameHistoryLogEntity> GetByFilter(QuerySlotGameHistoryLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<SlotGameHistoryLogEntity>(sql, filter);
            return query;
        }

        public PagedList<SlotGameHistoryLogEntity> GetPagingByFilter(QuerySlotGameHistoryLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = GetPagingEntitiesBySQLScript<SlotGameHistoryLogEntity>(sql, filter, filter);
            return query;
        }

        private string GenerateWhereSQL(QuerySlotGameHistoryLogRequestViewModel filter)
        {
            string sql = $@" SELECT * FROM {TableName} WITH (NOLOCK) 
                             WHERE 1 = 1 ";

            if (filter.PlayerID.HasValue)
                sql += "AND PlayerID = @PlayerID ";

            if (filter.GameID.HasValue)
                sql += "AND GameID = @GameID ";

            if (filter.GameTimeBegin.HasValue)
                sql += "AND GameTime >= @GameTimeBegin ";

            if (filter.GameTimeEnd.HasValue)
                sql += "AND GameTime <= @GameTimeEnd ";

            return sql;
        }

    }


}
