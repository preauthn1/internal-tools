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

    public class OrganizationRoleMappingRepository : FormatIntegerPrimaryKeyRepository<OrganizationRoleMappingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.OperationDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[OrganizationRoleMapping]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }


        public List<OrganizationRoleMappingEntity> GetByOrgID(int orgID)
        {
            var query = GetByFilter(new QueryOrganizationRoleMappingRequestViewModel
            {
                OrganizationIDs = orgID.ToListCollection()
            }) ;

            return query;
        }



        public void DeleteAndInsertNew(int orgID , List<int> roleIDs)
        {
            roleIDs = roleIDs.Distinct().ToList();
            var mappingEntity = roleIDs.Select(x => new OrganizationRoleMappingEntity
                {
                    RoleID = x,
                    OrganizationID = orgID
                })
                .ToList();

            DeleteAndInsertNewData("OrganizationID", orgID , mappingEntity);
        }


        public List<OrganizationRoleMappingEntity> GetByFilter(QueryOrganizationRoleMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<OrganizationRoleMappingEntity>(sql, filter);
            return query;
        }


        public PagedList<OrganizationRoleMappingEntity> GetPagingByFilter(QueryOrganizationRoleMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<OrganizationRoleMappingEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryOrganizationRoleMappingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.OrganizationIDs.IsNotEmptyNull())
                sql += " AND OrganizationID in @OrganizationIDs ";

            if (filter.RoleIDs.IsNotEmptyNull())
                sql += " AND RoleID in @RoleIDs ";

            return sql;
        }

    }


}
