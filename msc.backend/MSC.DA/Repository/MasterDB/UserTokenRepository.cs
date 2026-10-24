using Dapper;
using Lck.Utility.Extensions.Helper;
using MSC.Core.Repository;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.DA.Repository
{
    public class UserTokenRepository : IDapperRepository<UserTokenEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[AdminSys].[UserToken]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public UserTokenEntity GetCacheDataByHashToken(string tokenHash)
        {
            string cacheKey = CacheKey.DB_Data_UserToken(tokenHash);

            var queryCache = CacheHelper.GetCacheOrQuery(cacheKey, () =>
               {
                   var query = GetByFilter(new QueryUserTokenRequestViewModel
                   {
                       Token = [REDACTED]
                   });

                   return query.FirstOrDefault();
               });

            return queryCache;
        }

        public void ForceTokenExpired(List<int> userIDs)
        {
            string sql = $" UPDATE {TableName} SET ExpiredTime = '1970-01-01' ";

            //如果沒帶userIDs 會把全系統的user都登出
            if(userIDs.IsNotEmptyNull())
            {
                sql += " WHERE UserID in @userIDs ";
            }

            ExecuteSqlByDapper(sql, new { userIDs = userIDs });
        }


        public List<UserTokenEntity> GetByFilter(QueryUserTokenRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<UserTokenEntity>(sql, filter);
            return query;
        }


        public PagedList<UserTokenEntity> GetPagingByFilter(QueryUserTokenRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<UserTokenEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryUserTokenRequestViewModel filter , string selectColumn = "*")
        {
            string sql = $@"SELECT {selectColumn} FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND ID in @IDs ";

            if (filter.UserIDs.IsNotEmptyNull())
                sql += " AND UserID in @UserIDs ";

            if (filter.Token.IsNotNullNotWhiteSpace())
                sql += " AND Token = [REDACTED] ";

            if (filter.IsNotExpiredOnly)
                sql += " AND ExpiredTime > GetUTCDATE() ";
            
            return sql;
        }

    }


}
