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
using Lck.Utility.Extensions.Helper;

namespace target.DA.Repository
{

    public class UserRepository : FormatIntegerPrimaryKeyRepository<UserEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.OperationDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[User]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "UserID"; }
        }


        public UserEntity GetSingleByCache(int userID)
        {
            string cacheKey = CacheKey.AdminSys_UserInfo(userID);

            var queryCache = CacheHelper.GetCacheOrQuery(cacheKey, () =>
               {
                   var q = GetById(userID);
                   return q;

               }, expiredMinute: 15);

            return queryCache;
        }
        

        public void UpdateLastLoginInfo(int userID , string ipAddress)
        {
            string sql = $@"UPDATE {TableName}
                            SET [LastLogonTime] = GETUTCDATE() , 
                                [LastLogonIPAddress] = @LastLogonIPAddress
                            WHERE UserID = @UserID";

            var para = new
            {
                LastLogonIPAddress = ipAddress,
                UserID = userID
            };

            ExecuteSqlByDapper(sql , para);
        }


        public UserEntity GetByAccount(string account)
        {
            var filter = new QueryUserRequestViewModel
            {
                IsIncludeDelete = true,  //已刪除的(999) 也包含進去好了
                Account = account
            };

            var query = GetByFilter(filter);
            return query.FirstOrDefault();
        }
        

        public UserEntity GetByAccountPassword(string account , string password)
        {
            var filter = new QueryUserRequestViewModel
            {
                Account = account,
                Password = [REDACTED]
            };

            var query = GetByFilter(filter);
            return query.FirstOrDefault();
        }


        public List<OrganizationEntity> SelectSpecificColumnsByFilter(QueryUserRequestViewModel filter , string selectColumns = "*")
        {
            string sql = GenerateWhereSQL(filter , selectColumns);
            var query = QueryByDapper<OrganizationEntity>(sql, filter);
            return query;
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

            var query = GetPagingEntitiesBySQLScript<UserEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryUserRequestViewModel filter, string selectColumns = "*")
        {
            string sql = $@"SELECT {selectColumns} FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.UserIDs.IsNotEmptyNull())
                sql += " AND UserID in @UserIDs ";

            if (filter.OrganizationIDs.IsNotEmptyNull())
                sql += " AND OrganizationID in @OrganizationIDs ";

            if (!string.IsNullOrWhiteSpace(filter.Account))
                sql += " AND UserName = @Account ";

            if (!string.IsNullOrWhiteSpace(filter.Password))
                sql += " AND Password = [REDACTED] ";

            if (!filter.IsIncludeDelete)
                sql += $" AND Status != {NormalActiveStatus.PermanentDelete.ToInteger()}";

            return sql;
        }

    }


}
