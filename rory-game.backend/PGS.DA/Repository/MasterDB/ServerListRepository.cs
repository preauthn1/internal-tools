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
using Lck.Utility.Extensions.Helper;
using target.Utility.Helper;

namespace target.DA.Repository
{

    public class ServerListRepository : FormatIntegerPrimaryKeyRepository<ServerListEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[ServerList]"; }
        }


        protected override string PrimaryKeyColumnName
        {
            get { return "ServerID"; }
        }
        

        public void DeleteRelationDataWhenServerDelete(int serverID)
        {
            ExecuteStoreProcedureByDapper("PrWeb_DeleteRelationDataWhenServerDelete", new { serverID = serverID });
        }


        public void UpdateLastTime(int serverID)
        {
            string sql = $@" UPDATE {TableName} 
                             SET UpdateTime = GETUTCDATE()   
                             WHERE ServerID = @serverID ";

            ExecuteSqlByDapper(sql , new { serverID  = serverID });
        }


        public ServerBriefCacheModel GetBriefByCache(int serverID)
        {
            string cacheKey = CacheKey.AdminSys_ServerBriefInfo(serverID);

            var queryCachae = CacheHelper.GetCacheOrQuery(cacheKey, () =>
            {
                string sql = @" SELECT s.ServerID , s.ServerName , s.FlagLevel ,  s.Status  , g.GameID , g.GameType , g.GameCode
                                FROM ServerList as s
                                LEFT JOIN GameList as g on s.GameID = g.GameID
                                WHERE ServerID = @serverID";

                var query = QueryByDapper<ServerBriefCacheModel>(sql, new { serverID = serverID });

                if (query == null || query.IsEmptyOrNull())
                    return null;

                return query.FirstOrDefault();

            } , expiredMinute: 15);

            if (queryCachae == null)
                return queryCachae;

            return queryCachae;
        }


        public List<ServerListEntity> GetByFilter(QueryServerListRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<ServerListEntity>(sql, filter);
            return query;
        }


        public PagedList<ServerListEntity> GetPagingByFilter(QueryServerListRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<ServerListEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryServerListRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.ServerIDs.IsNotEmptyNull())
                sql += " AND ServerID in @ServerIDs ";

            if (filter.GameIDs.IsNotEmptyNull())
                sql += " AND GameID in @GameIDs ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND ServerID in ( SELECT ServerID from [dbo].[ServerMerchantMapping] WHERE MerchantID in  @MerchantIDs )";

            if (filter.Status.HasValue)
                sql += " AND Status = @Status ";

            if (!string.IsNullOrWhiteSpace(filter.Currency))
                sql += $" AND Currency = @Currency ";

            if (!string.IsNullOrWhiteSpace(filter.ServerNameFullMatch))
                sql += $" AND ServerName = @ServerNameFullMatch ";

            if (!string.IsNullOrWhiteSpace(filter.ConnIPKeyword))
                sql += $" AND ({GenerateLikeQuery("ClientAddr" , "@ConnIPKeyword")} OR {GenerateLikeQuery("NotifyAddr", "@ConnIPKeyword")} OR {GenerateLikeQuery("HttpAddr", "@ConnIPKeyword")} ) ";


            return sql;
        }

    }


}
