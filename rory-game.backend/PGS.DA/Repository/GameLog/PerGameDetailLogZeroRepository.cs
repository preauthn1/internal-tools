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

    public class PerGameDetailLogZeroRepository : FormatLongPrimaryKeyRepository<PerGameDetailLogZeroEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.GameLogDB); }
        }

        protected override string TableName
        {
            get { return "PerGameDetailLogZero"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<PerGameDetailLogZeroEntity> GetByFilter(QueryPerGameDetailLogZeroRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<PerGameDetailLogZeroEntity>(sql, filter);
            return query;
        }


        public PagedList<PerGameDetailLogZeroEntity> GetPagingByFilter(QueryPerGameDetailLogZeroRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<PerGameDetailLogZeroEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryPerGameDetailLogZeroRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND ID in @IDs ";

            if (filter.ServerIDs.IsNotEmptyNull())
                sql += " AND ServerID in @ServerIDs ";

            if (filter.GameIDs.IsNotEmptyNull())
                sql += " AND GameID in @GameIDs ";

            return sql;
        }

    }


}
