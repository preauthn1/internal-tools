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

    public class SuperToolSystemPoolWinLossRepository : FormatIntegerPrimaryKeyRepository<SuperToolSystemPoolWinLossEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[SuperToolSystemPoolWinLoss]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "GameID"; }
        }


        public void UpdateOrCreateSystemPoolWinLoss(int gameID, int flagID , decimal poolAmount , string gameData)
        {
            string spName = "PrWeb_SuperTool_UpdateOrCreateSystemPoolWinLoss";

            var para = new
            {
                gameID = gameID,
                flagID = flagID,
                poolAmount = poolAmount,
                gameData = gameData,
            };

            ExecuteStoreProcedureByDapper(spName, para);
        }


        public SuperToolSystemPoolWinLossEntity GetByGameID_And_FlagID(int gameID , int flagID)
        {
            var filter = new QuerySuperToolSystemPoolWinLossRequestViewModel
            {
                GameIDs = gameID.ToListCollection(),
                FlagID = flagID
            };

            var query = GetByFilter(filter);
            return query.FirstOrDefault();
        }


        public List<SuperToolSystemPoolWinLossEntity> GetByFilter(QuerySuperToolSystemPoolWinLossRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<SuperToolSystemPoolWinLossEntity>(sql, filter);
            return query;
        }


        public PagedList<SuperToolSystemPoolWinLossEntity> GetPagingByFilter(QuerySuperToolSystemPoolWinLossRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<SuperToolSystemPoolWinLossEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QuerySuperToolSystemPoolWinLossRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.GameIDs.IsNotEmptyNull())
                sql += " AND GameID in @GameIDs ";

            if (filter.FlagID.HasValue)
                sql += " AND FlagID = @FlagID ";

            return sql;
        }

    }


}
