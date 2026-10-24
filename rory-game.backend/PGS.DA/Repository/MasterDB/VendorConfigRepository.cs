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

    public class VendorConfigRepository : FormatIntegerPrimaryKeyRepository<VendorConfigEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[ThirdParty].[VendorConfig]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "VendorConfigID"; }
        }

        public List<VendorConfigEntity> GetByFilter(QueryVendorConfigRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<VendorConfigEntity>(sql, filter);
            return query;
        }


        public PagedList<VendorConfigEntity> GetPagingByFilter(QueryVendorConfigRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<VendorConfigEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryVendorConfigRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.VendorConfigID.IsNotEmptyNull())
                sql += " AND VendorConfigID in @VendorConfigIDs ";

            if (filter.VendorIDs.IsNotEmptyNull())
                sql += " AND VendorID in @VendorIDs ";

            return sql;
        }

    }


}
