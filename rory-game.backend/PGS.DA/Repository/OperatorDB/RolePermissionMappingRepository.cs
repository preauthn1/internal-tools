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

    public class RolePermissionMappingRepository : FormatIntegerPrimaryKeyRepository<RolePermissionMappingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.OperationDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[RolePermissionMapping]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }


        public void CopyPermissionFromRole(int roleID, int templateRoleID)
        {
            var filter = new QueryRolePermissionMappingRequestViewModel
            {
                RoleIDs = templateRoleID.ToListCollection()
            };

            var insertCopySettings = GetByFilter(filter);
            DeleteAndInsertNew(roleID, insertCopySettings.Select(x=> x.PermissionID).ToList());
        }


        public void DeleteAndInsertNew(int roleID, List<int> roleIDs)
        {
            roleIDs = roleIDs.Distinct().ToList();
            var models = roleIDs.Select(x => new RolePermissionMappingEntity
                {
                    PermissionID = x,
                    RoleID = roleID
                })
                .ToList();

            DeleteAndInsertNewData("RoleID", roleID, models);
        }



        public List<RolePermissionMappingEntity> GetByFilter(QueryRolePermissionMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<RolePermissionMappingEntity>(sql, filter);
            return query;
        }


        public PagedList<RolePermissionMappingEntity> GetPagingByFilter(QueryRolePermissionMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<RolePermissionMappingEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryRolePermissionMappingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.RoleIDs.IsNotEmptyNull())
                sql += " AND RoleID in @RoleIDs ";

            if (filter.PermissionIDs.IsNotEmptyNull())
                sql += " AND PermissionID in @PermissionIDs ";

            return sql;
        }

    }


}
