using Dapper;
using game.Entity;
using Lck.Utility.Extensions.Helper;
using target.Core.Repository;
using target.Models.Enums;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.Utility.Utility;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.DA.Repository
{
    public class DynamicQueryCommandRepository : FormatIntegerPrimaryKeyRepository<DynamicQueryCommandEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[DynamicQueryCommand]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<Dictionary<string , object>> QueryByTargetSql(string sql , object filter)
        {
            using (var conn = InstanceSQLConn())
            {
                var query = conn.Query(sql , filter);
                return query
                    .Cast<IDictionary<string, object>>()
                    .Select(it => it.ToDictionary(it => it.Key, it => it.Value))
                    .ToList();
            }
        }


        public List<DynamicQueryCommandEntity> GetByFilter(QueryDynamicQueryCommandRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<DynamicQueryCommandEntity>(sql, filter);
            return query;
        }


        public PagedList<DynamicQueryCommandEntity> GetPagingByFilter(QueryDynamicQueryCommandRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<DynamicQueryCommandEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryDynamicQueryCommandRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX in @IDs ";

            return sql;
        }

    }


}
