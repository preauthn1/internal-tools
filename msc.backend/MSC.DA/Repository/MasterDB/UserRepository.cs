using Dapper;
using MSC.Core.Repository;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using System;
using System.Collections.Generic;
using MSC.Utility.Extensions;
using System.Linq;
using Lck.Utility.Extensions.Helper;

namespace MSC.DA.Repository
{

    public class UserRepository : IDapperRepository<UserEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[AdminSys].[User]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "UserID"; }
        }


        public UserEntity GetSingleByCache(int userID)
        {
            string cacheKey = CacheKey.DB_Data_User(userID);

            var queryCache = CacheHelper.GetCacheOrQuery(cacheKey, () =>
               {
                   var query = GetById(userID);
                   return query;

               }, 5);

            return queryCache;
        }


        public string GetUserNameByUserID(int userID)
        {
            var filter = new QueryUserRequestViewModel
            {
                UserIDs = userID.ToListCollection()
            };

            string sql = GenerateWhereSQL(filter, " UserName ");
            var query = QueryByDapper<UserEntity>(sql, filter);
            return query.FirstOrDefault()?.UserName ?? "";
        }

        public List<UserEntity> GetByFilter(QueryUserRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<UserEntity>(sql, filter);
            return query;
        }


        public PagedList<UserEntity> GetPagingByFilter(QueryUserRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<UserEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryUserRequestViewModel filter , string columns = "*")
        {
            string sql = $@"SELECT {columns} FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.UserIDs.IsNotEmptyNull())
                sql += " AND UserID in @UserIDs ";

            if (filter.UserName.IsNotNullNotWhiteSpace())
                sql += " AND UserName = @UserName ";

            if (filter.PasswordHash.IsNotNullNotWhiteSpace())
                sql += " AND Password = [REDACTED] ";

            return sql;
        }

    }


}
