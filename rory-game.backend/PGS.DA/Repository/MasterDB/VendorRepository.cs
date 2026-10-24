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
using Lck.Utility.Extensions.Helper;

namespace target.DA.Repository
{

    public class VendorRepository : FormatIntegerPrimaryKeyRepository<VendorEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[ThirdParty].[Vendor]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "VendorID"; }
        }

        public List<VendorEntity> GetByFilter(QueryVendorRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<VendorEntity>(sql, filter);
            return query;
        }


        public PagedList<VendorEntity> GetPagingByFilter(QueryVendorRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<VendorEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryVendorRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.VendorIDs.IsNotEmptyNull())
                sql += " AND VendorID in @VendorIDs ";

            return sql;
        }

    }


}
