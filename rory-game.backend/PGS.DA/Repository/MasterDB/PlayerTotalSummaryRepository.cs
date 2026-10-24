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

    public class PlayerTotalSummaryRepository : FormatIntegerPrimaryKeyRepository<PlayerTotalSummaryEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[PlayerTotalSummary]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "PlayerID"; }
        }

        public void UpdateLoginTimes(int playerID)
        {
            string sql = $@" UPDATE {TableName} SET SumLoginTimes = SumLoginTimes + 1 
                             WHERE PlayerID = @playerID ";

            ExecuteSqlByDapper(sql, new { playerID  = playerID });
        }


        public List<PlayerTotalSummaryEntity> GetByFilter(QueryPlayerTotalSummaryRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<PlayerTotalSummaryEntity>(sql, filter);
            return query;
        }


        public PagedList<PlayerTotalSummaryEntity> GetPagingByFilter(QueryPlayerTotalSummaryRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<PlayerTotalSummaryEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryPlayerTotalSummaryRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.PlayerIDs.IsNotEmptyNull())
                sql += " AND PlayerID in @PlayerIDs ";

            return sql;
        }

    }


}
