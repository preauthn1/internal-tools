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

    public class RoleRepository : FormatIntegerPrimaryKeyRepository<RoleEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.OperationDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[Role]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "RoleID"; }
        }

        public List<RoleEntity> GetByFilter(QueryRoleRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<RoleEntity>(sql, filter);
            return query;
        }


        public PagedList<RoleEntity> GetPagingByFilter(QueryRoleRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<RoleEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryRoleRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.RoleIDs.IsNotEmptyNull())
                sql += " AND RoleID in @RoleIDs ";

            return sql;
        }

    }


}
