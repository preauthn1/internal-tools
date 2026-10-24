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

    public class GameWarningLogRepository : FormatIntegerPrimaryKeyRepository<GameWarningLogEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.TransactionDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[GameWarningLog]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<GameWarningLogEntity> GetByFilter(QueryGameWarningLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<GameWarningLogEntity>(sql, filter);
            return query;
        }


        public PagedList<GameWarningLogEntity> GetPagingByFilter(QueryGameWarningLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<GameWarningLogEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryGameWarningLogRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.WarningTypes.IsNotEmptyNull())
                sql += " AND WarningType in @WarningTypes ";

            if (filter.PlayerIDs.IsNotEmptyNull())
                sql += " AND PlayerID in @PlayerIDs ";

            return sql;
        }

    }


}
