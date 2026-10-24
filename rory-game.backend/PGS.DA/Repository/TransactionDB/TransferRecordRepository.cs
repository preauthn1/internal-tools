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

    public class TransferRecordRepository : FormatLongPrimaryKeyRepository<TransferRecordEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.TransactionDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[TransferRecord]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "TransferRecordID"; }
        }

        public List<TransferRecordEntity> GetByFilter(QueryTransferRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<TransferRecordEntity>(sql, filter);
            return query;
        }


        public PagedList<TransferRecordEntity> GetPagingByFilter(QueryTransferRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<TransferRecordEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryTransferRecordRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} WITH (NOLOCK)
                           WHERE 1=1 ";

            if (filter.TransferRecordIDs.IsNotEmptyNull())
                sql += " AND TransferRecordID in @TransferRecordIDs ";

            if (filter.PlayerIDs.IsNotEmptyNull())
                sql += " AND PlayerID in @PlayerIDs ";

            if(filter.TransferType.HasValue)
                sql += " AND TransferType = @TransferType ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (!string.IsNullOrWhiteSpace(filter.OrderNo))
                sql += " AND OrderNo = @OrderNo ";

            if (!string.IsNullOrWhiteSpace(filter.MerchantOrderNo))
                sql += " AND MerchantOrderNo = @MerchantOrderNo ";

            if (!string.IsNullOrWhiteSpace(filter.PlayerName))
                sql += " AND PlayerName = @PlayerName ";

            if (filter.CreateTimeBegin.HasValue)
                sql += " AND CreateTime >= @CreateTimeBegin ";

            if (filter.CreateTimeEnd.HasValue)
                sql += " AND CreateTime <= @CreateTimeEnd ";

            return sql;
        }

    }


}
