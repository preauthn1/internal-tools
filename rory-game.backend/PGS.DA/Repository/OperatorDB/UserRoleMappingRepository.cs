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
using target.Models.Models;
using Lck.Utility.Extensions.Helper;

namespace target.DA.Repository
{

    public class UserRoleMappingRepository : FormatIntegerPrimaryKeyRepository<UserRoleMappingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.OperationDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[UserRoleMapping]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<UserPermissionModel> GetRolePermissionsByCache(List<int> roleIDs)
        {
            if (roleIDs.IsEmptyOrNull())
                return new List<UserPermissionModel>();

            roleIDs = roleIDs.Distinct().ToList();

            //多重角色的情況，cache不好處理
            bool isSingleOne = roleIDs.Count() == 1;

            string roleIDsText = string.Join("_", roleIDs.OrderBy(x => x));

            string cacheKey = CacheKey.AdminSys_RolePermission(roleIDsText);

            var queryCache = CacheHelper.GetCacheOrQuery(cacheKey, () =>
               {
                   var query = GetPermissionsByRoleIDs(roleIDs);

                   //為了避免麻煩 轉小寫判斷
                   query.ForEach(x => x.PermissionToken = [REDACTED];

                   return query;

               }, isSingleOne ? 5 : 1);

            return queryCache;
        }


        public List<UserPermissionModel> GetPermissionsByRoleIDs(List<int> roleIDs)
        {
            string sql = @" SELECT PermissionID, PermissionToken FROM Permission WHERE PermissionID in
                            (
	                            SELECT PermissionID FROM RolePermissionMapping where roleID in @roleIDs
                            )";

            using (var conn = InstanceSQLConn())
            {
                var query = conn.Query<UserPermissionModel>(sql, new { roleIDs = roleIDs });
                return query.ToList();
            }
        }


        public void DeleteAndInsertNew(int userID , List<int> roleIDs)
        {
            roleIDs = roleIDs.Distinct().ToList();
            var models = roleIDs.Select(x => new UserRoleMappingEntity
                {
                    UserID = userID,
                    RoleID = x
                })
                .ToList();

            DeleteAndInsertNewData("UserID" , userID , models);
        }


        public void DeleteByRoleID(int roleID)
        {
            string deleteSql = $"DELETE {TableName} WHERE RoleID = @RoleID ";

            using (var conn = InstanceSQLConn())
            {
                conn.Execute(deleteSql, new { RoleID = roleID });
            }
        }

        public List<UserRoleMappingEntity> GetByUserID(int userID)
        {
            var filter = new QueryUserRoleMappingRequestViewModel
            {
                UserIDs = userID.ToListCollection()
            };

            var query = GetByFilter(filter);
            return query;
        }



        public List<UserRoleMappingEntity> GetByFilter(QueryUserRoleMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<UserRoleMappingEntity>(sql, filter);
            return query;
        }


        public PagedList<UserRoleMappingEntity> GetPagingByFilter(QueryUserRoleMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<UserRoleMappingEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryUserRoleMappingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.UserIDs.IsNotEmptyNull())
                sql += " AND UserID in @UserIDs ";

            if (filter.RoleIDs.IsNotEmptyNull())
                sql += " AND RoleID in @RoleIDs ";

            return sql;
        }

    }


}
