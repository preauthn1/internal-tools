using Dapper;
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
    public class RoleRepository : IDapperRepository<RoleEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[AdminSys].[Role]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
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

            var query = GetPageingEntitiesBySQLScript<RoleEntity>(sql, filter, filter);
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
