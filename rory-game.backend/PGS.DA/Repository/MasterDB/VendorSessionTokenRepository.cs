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
using target.Utility.Helper;
using Lck.Utility.Extensions.Helper;

namespace target.DA.Repository
{

    public class VendorSessionTokenRepository : FormatIntegerPrimaryKeyRepository<VendorSessionTokenEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[ThirdParty].[VendorSessionToken]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public void ResetAllPlayerSessionToken(int playerID)
        {
            DateTime emptyDate = DateTimeHelper.GetDefaultTime();

            //查出還存活的cache
            var queryExistToken = [REDACTED] QueryVendorSessionTokenRequestViewModel
            {
                PlayerIDs = playerID.ToListCollection(),
                IsNotExpired = true
            });

            if (queryExistToken.IsEmptyOrNull())
                return;

            string sql = @$"UPDATE {TableName} 
                            SET ExpiredTime = @emptyDate
                            WHERE playerID = @playerID AND ExpiredTime > GETUTCDATE()";

            var para = new
            {
                emptyDate = emptyDate,
                playerID = playerID
            };

            foreach(var item in queryExistToken)
            {
                //事實上採用memory cache 這邊因為機器不同 清cache毫無意義
                //待未來優化成redisCache吧
                string cacheKey = CacheKey.DbEntity_VendorSessionToken(item.SessionToken);
                CacheHelper.RemoveCacheAsync(cacheKey);
            }

            ExecuteSqlByDapper(sql, para);
        }

     
        public VendorSessionTokenEntity GetTokenByCache(string sessionToken)
        {
            string cacheKey = CacheKey.DbEntity_VendorSessionToken(sessionToken);

            int expiredMinute = 1; //這個cache 不要設太久，為了避免token該失效的情況沒有反應

            var queryCacche = CacheHelper.GetCacheOrQuery(cacheKey, () =>
              {
                  var query = GetByFilter(new QueryVendorSessionTokenRequestViewModel
                  {
                      SessionToken = [REDACTED]
                  });
                  return query.FirstOrDefault();
              }, expiredMinute);

            return queryCacche;
        }

        public VendorSessionTokenEntity GetByToken(string sessionToken)
        {
            var query = GetByFilter(new QueryVendorSessionTokenRequestViewModel
            {
                SessionToken = [REDACTED]
            });

            return query.FirstOrDefault();
        }



        public List<VendorSessionTokenEntity> GetByFilter(QueryVendorSessionTokenRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<VendorSessionTokenEntity>(sql, filter);
            return query;
        }


        public PagedList<VendorSessionTokenEntity> GetPagingByFilter(QueryVendorSessionTokenRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<VendorSessionTokenEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryVendorSessionTokenRequestViewModel filter , string selectColumns = "*")
        {
            string sql = $@"SELECT {selectColumns} FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.PlayerIDs.IsNotEmptyNull())
                sql += " AND PlayerID in @PlayerIDs ";

            if (filter.Status.IsNotEmptyNull())
                sql += " AND Status in @Status ";

            if (filter.SessionToken.IsNotNullEmptyText())
                sql += " AND SessionToken = [REDACTED] ";

            if (filter.IsNotExpired)
                sql += " AND  ExpiredTime > GETUTCDATE() ";

            return sql;
        }

    }


}
