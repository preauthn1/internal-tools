using Dapper;
using game.Entity;
using Lck.Utility.Extensions.Helper;
using target.Core.Repository;
using target.Models.Enums;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.DA.Repository
{
    public class BetRecordTransactionNoMappingRepository : FormatIntegerPrimaryKeyRepository<BetRecordTransactionNoMappingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.GameLogDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[BetRecordTransactionNoMapping]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "TransactionNo"; }
        }

        public BetRecordTransactionNoMappingEntity GetByTransactionNo(string transactionNo)
        {
            var query = GetByFilter(new QueryBetRecordTransactionNoMappingRequestViewModel
            {
                BetRecordNoPGS = transactionNo
            });

            return query.FirstOrDefault();
        }


        public List<BetRecordTransactionNoMappingEntity> GetByFilter(QueryBetRecordTransactionNoMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<BetRecordTransactionNoMappingEntity>(sql, filter);
            return query;
        }


        public PagedList<BetRecordTransactionNoMappingEntity> GetPagingByFilter(QueryBetRecordTransactionNoMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<BetRecordTransactionNoMappingEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryBetRecordTransactionNoMappingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.BetRecordNoPGS.IsNotNullEmptyText())
                sql += " AND BetRecordNoPGS = @BetRecordNoPGS ";

            return sql;
        }

    }


}
