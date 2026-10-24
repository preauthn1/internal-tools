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

    public class PermissionRepository : FormatIntegerPrimaryKeyRepository<PermissionEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.OperationDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[Permission]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "PermissionID"; }
        }

        public List<PermissionEntity> GetByFilter(QueryPermissionRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<PermissionEntity>(sql, filter);
            return query;
        }


        public PagedList<PermissionEntity> GetPagingByFilter(QueryPermissionRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<PermissionEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryPermissionRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.PermissionIDs.IsNotEmptyNull())
                sql += " AND PermissionID in @PermissionIDs ";

            return sql;
        }

    }


}
