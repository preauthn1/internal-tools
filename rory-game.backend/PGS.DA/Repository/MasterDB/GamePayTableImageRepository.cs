using Dapper;
using game.Entity;
using Lck.Utility.Extensions.Helper;
using target.Core.Repository;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.DA.Repository
{
    public class GamePayTableImageRepository : FormatIntegerPrimaryKeyRepository<GamePayTableImageEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[GamePayTableImage]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<GamePayTableImageEntity> GetByFilter(QueryGamePayTableImageRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<GamePayTableImageEntity>(sql, filter);
            return query;
        }


        public PagedList<GamePayTableImageEntity> GetPagingByFilter(QueryGamePayTableImageRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<GamePayTableImageEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryGamePayTableImageRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND ID in @IDs ";

            if (filter.GameIDs.IsNotEmptyNull())
                sql += " AND GameID in @GameIDs ";

            if (filter.Code.IsNotNullEmptyText())
                sql += " AND Code = @Code ";

            return sql;
        }

    }


}
