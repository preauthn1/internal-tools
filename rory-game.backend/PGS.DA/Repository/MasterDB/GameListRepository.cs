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
using target.Models.Models.StoreProcedures;
using System.Data.SqlClient;
using System.Data;
using Lck.Utility.Extensions.Helper;
using target.Utility;
using target.Models.Models;

namespace target.DA.Repository
{

    public class GameListRepository : FormatIntegerPrimaryKeyRepository<GameListEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[GameList]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "GameID"; }
        }

        public int GetMaxIdByVendor(int vendorID)
        {
            string sql = "SELECT max(GameID) FROM [dbo].[GameList] where GameVendorType = @vendorID";

            using (var conn = InstanceSQLConn())
            {
                var query = conn.Query<int>(sql, new { vendorID });
                return query.FirstOrDefault();
            }
        }

        public GameListEntity GetSingleByCache(int gameID)
        {
            if (gameID <= 0)
                return null;

            string cacheKey = CacheKey.DbEntity_GameList(gameID);

            var queryCachae = CacheHelper.GetCacheOrQuery(cacheKey, () =>
            {
                var querySingle = GetById(gameID);
                return querySingle;

            }, expiredMinute: 30);

            return queryCachae;
        }

        public GameListEntity GetSingleByCacheAndCheckValid(int gameID)
        {
            var queryGame = GetSingleByCache(gameID);

            if (queryGame == null)
                throw new TerproException(ApiErrorCode.NotFoundGame);

            if(queryGame.Status != NormalActiveStatus.Enable.ToInteger())
                throw new TerproException(ApiErrorCode.InvalidStatus_GameStatus);

            return queryGame;
        }


        public GameListEntity GetGameCacheByVendorGameCode(int vendorID, string thirdPartyGameCode)
        {
            string cacheKey = CacheKey.DbEntity_GameList_GameVendorCode(vendorID, thirdPartyGameCode);

            var queryCachae = CacheHelper.GetCacheOrQuery(cacheKey, () =>
            {
                var querySingle = GetGamebyVendorGameCode(vendorID , thirdPartyGameCode);
                return querySingle;

            }, expiredMinute: 10);

            return queryCachae;
        }


        public GameListEntity GetGamebyVendorGameCode(int vendorID , string thirdPartyGameCode)
        {
            var query = GetByFilter(new QueryGameListRequestViewModel
            {
                ThirdPartyGameCode = thirdPartyGameCode,
                GameVendorTypes = vendorID.ToListCollection()
            });

            return query.FirstOrDefault();
        }


        public List<Stat_GameHealthSummarySPResultModel> GetGameHealthResult(DateTime begin , DateTime end , List<int> merchantIDs)
        {
            var para = new
            {
                dateStart = begin,
                dateEnd = end,
                merchantList = string.Join("," , merchantIDs) ?? ""
            };

            string myConn = GetConnectionStringByDbName(DBConnType.StatReportDB);
            using (var conn = new SqlConnection(myConn))
            {
                var user = conn.Query<Stat_GameHealthSummarySPResultModel>("PrWeb_Stat_GameServiceStatus", para, commandType: CommandType.StoredProcedure);
                return user.ToList();
            }
        }



        public List<GameListEntity> GetAllStatusEnableGames()
        {
            var filter = new QueryGameListRequestViewModel
            {
                Status = NormalActiveStatus.Enable.ToInteger()
            };

            var query = GetByFilter(filter);
            return query;
        }


        public bool IsGameExistByCache(int gameID)
        {
            var queryCachae = GetSingleByCache(gameID);
            return queryCachae != null;
        }



        public List<GameListEntity> GetByFilter(QueryGameListRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<GameListEntity>(sql, filter);
            return query;
        }


        public PagedList<GameListEntity> GetPagingByFilter(QueryGameListRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<GameListEntity>(sql, filter, filter);
            return query;
        }

        private string GenerateWhereSQL(QueryGameListRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.GameIDs.IsNotEmptyNull())
                sql += " AND GameID in @GameIDs ";

            if (filter.GameVendorTypes.IsNotEmptyNull())
                sql += " AND GameVendorType in @GameVendorTypes ";

            if (filter.GameCode.IsNotNullEmptyText())
                sql += " AND GameCode = @GameCode ";

            if (filter.ThirdPartyGameCode.IsNotNullEmptyText())
                sql += " AND ThirdPartyGameCode = @ThirdPartyGameCode ";

            if (filter.Status.HasValue)
                sql += " AND Status = @Status ";

            if (filter.IsNotDelete.HasValue)
                sql += $" AND Status != {NormalActiveStatus.PermanentDelete.ToInteger()} ";

            return sql;
        }

    }


}
