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

    public class ForbiddenIPAddressRepository : FormatIntegerPrimaryKeyRepository<ForbiddenIPAddressEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[ForbiddenIPAddress]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }


        public ForbiddenIPAddressEntity GetByIPAddress(string ipAddress)
        {
            var filter = new QueryForbiddenIPAddressRequestViewModel
            {
                IPAddress = ipAddress
            };

            var query = GetByFilter(filter);
            return query.FirstOrDefault();
        }


        public List<ForbiddenIPAddressEntity> GetByFilter(QueryForbiddenIPAddressRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<ForbiddenIPAddressEntity>(sql, filter);
            return query;
        }


        public PagedList<ForbiddenIPAddressEntity> GetPagingByFilter(QueryForbiddenIPAddressRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<ForbiddenIPAddressEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryForbiddenIPAddressRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND Id in @IDs ";

            if (!string.IsNullOrWhiteSpace(filter.IPAddress))
                sql += " AND IPAddress = @IPAddress ";

            return sql;
        }

    }


}
