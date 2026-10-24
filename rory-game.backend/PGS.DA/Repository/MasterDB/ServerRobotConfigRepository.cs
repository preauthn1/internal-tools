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

    public class ServerRobotConfigRepository : FormatIntegerPrimaryKeyRepository<ServerRobotConfigEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[ServerRobotConfig]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "ServerID"; }
        }

        public List<ServerRobotConfigEntity> GetByFilter(QueryServerRobotConfigRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<ServerRobotConfigEntity>(sql, filter);
            return query;
        }


        public PagedList<ServerRobotConfigEntity> GetPagingByFilter(QueryServerRobotConfigRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = GetPagingEntitiesBySQLScript<ServerRobotConfigEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryServerRobotConfigRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.ServerIDs.IsNotEmptyNull())
                sql += " AND ServerID in @ServerIDs ";

            if (filter.GameIDs.IsNotEmptyNull())
                sql += " AND ServerID in (SELECT ServerID FROM [dbo].[ServerList] WHERE GameID in @GameIDs ) ";
            
            return sql;
        }

    }


}
