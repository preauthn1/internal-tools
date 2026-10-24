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

    public class PlayerGameSummaryRepository : FormatIntegerPrimaryKeyRepository<PlayerGameSummaryEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[PlayerGameSummary]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "PlayerID"; }
        }

        public PlayerGameSummaryEntity GetByPlayerAndGameID(int playerID , int gameID)
        {
            var filter = new QueryPlayerGameSummaryRequestViewModel
            {
                GameIDs = gameID.ToListCollection(),
                PlayerIDs = playerID.ToListCollection()
            };

            var query = GetByFilter(filter);
            return query.FirstOrDefault();
        }


        public List<PlayerGameSummaryEntity> GetByFilter(QueryPlayerGameSummaryRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<PlayerGameSummaryEntity>(sql, filter);
            return query;
        }


        public PagedList<PlayerGameSummaryEntity> GetPagingByFilter(QueryPlayerGameSummaryRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = GetPagingEntitiesBySQLScript<PlayerGameSummaryEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryPlayerGameSummaryRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.PlayerIDs.IsNotEmptyNull())
                sql += " AND PlayerID in @PlayerIDs ";

            if (filter.GameIDs.IsNotEmptyNull())
                sql += " AND GameID in @GameIDs ";

            return sql;
        }

    }


}
