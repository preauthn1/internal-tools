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

    public class JackpotPoolRepository : FormatIntegerPrimaryKeyRepository<JackpotPoolEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[JackpotPool]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "GameID"; }
        }

        public List<JackpotPoolEntity> GetByFilter(QueryJackpotPoolRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<JackpotPoolEntity>(sql, filter);
            return query;
        }

        public PagedList<JackpotPoolEntity> GetPagingByFilter(QueryJackpotPoolRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<JackpotPoolEntity>(sql, filter, filter);

            return query;
        }

        private string GenerateWhereSQL(QueryJackpotPoolRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} " + 
                           "WHERE 1 = 1 ";

            if (filter.GameID.HasValue)
                sql += "AND GameID = @GameID ";

            if (filter.FlagID.HasValue)
                sql += "AND FlagID = @FlagID ";

            if (filter.LevelID.HasValue)
                sql += "AND LevelID = @LevelID ";

            if (filter.TableNo.HasValue)
                sql += "AND TableNo = @TableNo ";

            if (filter.IsLink.HasValue)
                sql += "AND IsLink = @IsLink ";
            else
                sql += "AND IsLink = 0 ";

            return sql;
        }

    }

}
