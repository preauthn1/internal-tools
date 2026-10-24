using Dapper;
using game.Entity;
using Lck.Utility.Extensions.Helper;
using target.Core.Repository;
using target.Models.Enums;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.Utility.Utility;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.DA.Repository
{
    public class LobbyTokenLogRepository : FormatLongPrimaryKeyRepository<LobbyTokenLogEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.TransactionDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[LobbyTokenLog]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public void ResetAllPlayerToken(int playerID)
        {
            DateTime now = DateTimeHelper.GetCurrentTime();
            DateTime emptyTime = new DateTime(1970, 1, 1);

            string sql = @$" UPDATE {TableName} SET ExpiredTime = @emptyTime 
                             WHERE playerID = @playerID AND ExpiredTime > @now";

            var para = new
            {
                emptyTime = emptyTime,
                now = now,
                playerID = playerID
            };

            ExecuteSqlByDapper(sql , para);
        }


        public LobbyTokenLogEntity GetByHash(string tokenHash)
        {
            var query = GetByFilter(new QueryLobbyTokenLogRequestViewModel
            {
                IdentityHash = tokenHash
            });

            return query.FirstOrDefault();
        }


        public LobbyTokenLogEntity GetByGuid(Guid uuid)
        {
            var query = GetByFilter(new QueryLobbyTokenLogRequestViewModel
            {
                IdentityGuid = uuid
            });

            return query.FirstOrDefault();
        }


        public List<LobbyTokenLogEntity> GetByFilter(QueryLobbyTokenLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<LobbyTokenLogEntity>(sql, filter);
            return query;
        }


        public PagedList<LobbyTokenLogEntity> GetPagingByFilter(QueryLobbyTokenLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<LobbyTokenLogEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryLobbyTokenLogRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.IdentityHash.IsNotNullEmptyText())
                sql += " AND IdentityHash = @IdentityHash ";

            if (filter.IdentityGuid.HasValue)
                sql += " AND IdentityGuid = @IdentityGuid ";

            return sql;
        }

    }


}
