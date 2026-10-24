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
    public class MerchantWithdrawalRepository : IDapperRepository<MerchantWithdrawalEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[MerchantWithdrawal]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<MerchantWithdrawalEntity> GetByFilter(QueryMerchantWithdrawalRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<MerchantWithdrawalEntity>(sql, filter);
            return query;
        }


        public PagedList<MerchantWithdrawalEntity> GetPagingByFilter(QueryMerchantWithdrawalRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<MerchantWithdrawalEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryMerchantWithdrawalRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND ID in @IDs ";

            return sql;
        }

    }


}
