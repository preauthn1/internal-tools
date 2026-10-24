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

    public class PlayerGameRecordRepository : FormatLongPrimaryKeyRepository<PlayerGameRecordEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.TransactionDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[PlayerGameRecord]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public PlayerGameRecordEntity GetBy_PlayerID_TradeID(int playerID , long tradeID)
        {
            var filter = new QueryPlayerGameRecordRequestViewModel
            {
                TradeID = tradeID,
                PlayerIDs = playerID.ToListCollection()
            };

            var query = GetByFilter(filter);
            return query.FirstOrDefault();
        }
        

        public List<PlayerGameRecordEntity> GetByFilter(QueryPlayerGameRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<PlayerGameRecordEntity>(sql, filter);
            return query;
        }


        public PagedList<PlayerGameRecordEntity> GetPagingByFilter(QueryPlayerGameRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<PlayerGameRecordEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryPlayerGameRecordRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} WITH (NOLOCK) 
                           WHERE 1=1 ";

            if (filter.PlayerIDs.IsNotEmptyNull())
                sql += " AND PlayerID in @PlayerIDs ";

            if (filter.ServerIDs.IsNotEmptyNull())
                sql += " AND ServerID in @ServerIDs ";

            if (filter.GameIDs.IsNotEmptyNull())
                sql += " AND GameID in @GameIDs ";

            if (filter.TradeID.HasValue)
                sql += " AND TradeID = @TradeID ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (!string.IsNullOrWhiteSpace(filter.PlayerName))
                sql += " AND PlayerName = @PlayerName ";

            if (!string.IsNullOrWhiteSpace(filter.IPAddress))
                sql += " AND IPAddress = @IPAddress ";

            if (filter.LoginTimeBegin.HasValue)
                sql += " AND LoginTime >= @LoginTimeBegin ";

            if (filter.LoginTimeEnd.HasValue)
                sql += " AND LoginTime <= @LoginTimeEnd ";

            if (filter.LogoutTimeBegin.HasValue)
                sql += " AND LogoutTime >= @LogoutTimeBegin ";

            if (filter.LogoutTimeEnd.HasValue)
                sql += " AND LogoutTime <= @LogoutTimeEnd ";

            return sql;
        }

    }


}
