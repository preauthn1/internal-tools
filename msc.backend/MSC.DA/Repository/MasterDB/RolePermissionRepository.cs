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
    public class RolePermissionRepository : IDapperRepository<RolePermissionEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[AdminSys].[RolePermission]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<RolePermissionEntity> GetRolePermissionByCache(int roleID)
        {
            string cacheKey = CacheKey.DB_Data_RolePermissions(roleID);

            var queryCache = CacheHelper.GetCacheOrQuery(cacheKey, () =>
               {
                   var query = GetByFilter(new QueryRolePermissionRequestViewModel
                   {
                       RoleIDs = roleID.ToListCollection()
                   });

                   return query;
               });

            return queryCache;
        }

        public List<RolePermissionEntity> GetByFilter(QueryRolePermissionRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<RolePermissionEntity>(sql, filter);
            return query;
        }


        public PagedList<RolePermissionEntity> GetPagingByFilter(QueryRolePermissionRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<RolePermissionEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryRolePermissionRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND ID in @IDs ";

            if (filter.RoleIDs.IsNotEmptyNull())
                sql += " AND RoleID in @RoleIDs ";

            return sql;
        }

    }


}
