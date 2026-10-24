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

    public class PerGameLogRepository : FormatLongPrimaryKeyRepository<PerGameLogEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.GameLogDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[PerGameLog]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "PerGameLogID"; }
        }

        public List<PerGameLogEntity> GetByFilter(QueryPerGameLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<PerGameLogEntity>(sql, filter);
            return query;
        }


        public PagedList<PerGameLogEntity> GetPagingByFilter(QueryPerGameLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<PerGameLogEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryPerGameLogRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND PerGameLogID in @IDs ";

            if (filter.ServerIDs.IsNotEmptyNull())
                sql += " AND ServerID in @ServerIDs ";
        
            if (filter.GameIDs.IsNotEmptyNull())
                sql += " AND GameID in @GameIDs ";

            return sql;
        }

    }


}
