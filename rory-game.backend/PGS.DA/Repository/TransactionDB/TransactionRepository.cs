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
using System.Data.SqlClient;
using System.Data;
using target.Models.Models.StoreProcedures;

namespace target.DA.Repository
{

    public class TransactionRepository : FormatLongPrimaryKeyRepository<TransactionEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.TransactionDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[Transaction]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "TransactionID"; }
        }


        public List<TransactionEntity> GetByFilter(QueryTransactionRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<TransactionEntity>(sql, filter);
            return query;
        }


        public PagedList<TransactionEntity> GetPagingByFilter(QueryTransactionRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<TransactionEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryTransactionRecordRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} WITH (NOLOCK) " +
                           "WHERE 1=1 ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += "AND MerchantID in @MerchantIDs ";

            if (filter.PlayerID.HasValue)
                sql += "AND PlayerID = @PlayerID ";

            if (filter.PlayerName.IsNotEmptyNull())
                sql += "AND PlayerName = @PlayerName ";

            if (filter.GameID.HasValue)
                sql += "AND GameID = @GameID ";

            if (filter.TransactionType.HasValue)
                sql += "AND TransactionType = @TransactionType ";

            if (filter.IPAddress.IsNotEmptyNull())
                sql += "AND IPAddress = @IPAddress ";

            if (filter.CreateTimeBegin.HasValue)
                sql += "AND CreateTime >= @CreateTimeBegin ";

            if (filter.CreateTimeEnd.HasValue)
                sql += "AND CreateTime <= @CreateTimeEnd ";

            return sql;
        }
    }
}
