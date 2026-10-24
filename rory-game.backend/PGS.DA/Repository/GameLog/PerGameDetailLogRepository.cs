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

namespace target.DA.Repository
{

    public class PerGameDetailLogRepository : FormatLongPrimaryKeyRepository<PerGameDetailLogEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.GameLogDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[PerGameDetailLog]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }


        public void UpdateCurrentAmount(long perGameDetailID , decimal currentAmount)
        {
            Update(perGameDetailID, new
            {
                CurrentAmount = currentAmount
            });
        }


        public List<PerGameDetailLogEntity> GetByFilter(QueryDetailLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<PerGameDetailLogEntity>(sql, filter);
            return query;
        }


        public PagedList<PerGameDetailLogEntity> GetPagingByFilter(QueryDetailLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<PerGameDetailLogEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryDetailLogRequestViewModel filter)
        {
            string noLock = filter.IsLockQuery ? "" : " WITH (NOLOCK) ";

            string sql = $@"SELECT * FROM {TableName} {noLock}
                           WHERE 1=1 ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.PlayerIDs.IsNotEmptyNull())
                sql += " AND PlayerID in @PlayerIDs ";

            if (!string.IsNullOrWhiteSpace(filter.PlayerName))
                sql += " AND PlayerName = @PlayerName ";

            if (filter.GameIDs.IsNotEmptyNull())
                sql += " AND GameID in @GameIDs ";

            if (filter.ServerIDs.IsNotEmptyNull())
                sql += " AND ServerID in @ServerIDs ";

            if (filter.PerGameLogIDs.IsNotEmptyNull())
                sql += " AND PerGameLogID in @PerGameLogIDs ";

            if (filter.GameStartTimeBegin.HasValue)
                sql += " AND GameTime >= @GameStartTimeBegin ";

            if (filter.GameStartTimeEnd.HasValue)
                sql += " AND GameTime <= @GameStartTimeEnd ";

            if (filter.TransactionNo.IsNotNullEmptyText())
                sql += " AND TransactionNo = @TransactionNo ";

            if (filter.IsJackpotOnly)
            {
                int tagID = SystemTagEnum.PerGameDetailLog_Jackpot.ToInteger();
                sql += $" AND Id in ( SELECT PerGameDetailLogID FROM [dbo].[PerGameDetailLogTagMapping] WHERE TagID = {tagID} ) ";
            }

            return sql;
        }

    }


}
