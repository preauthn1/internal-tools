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
    public class BetRecordMerchantNotifyQueueRepository : FormatLongPrimaryKeyRepository<BetRecordMerchantNotifyQueueEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.GameLogDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[BetRecordMerchantNotifyQueue]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public void UpdateNextTime(long id , DateTime nextTime)
        {
            // 更新下次時間
            // 更新retry次數
            string sql = @$" UPDATE {TableName} SET NextNotifyTime = @nextTime , RetryTimes = RetryTimes + 1 
                             WHERE Id = @id";

            ExecuteSqlByDapper(sql, new { id = id, nextTime = nextTime });
        }

        public void UpdateToSuccess(long id)
        {
            DateTime maxTime = new DateTime(2999, 1, 1);
            string sql = @$" UPDATE {TableName} SET NextNotifyTime = @nextTime , IsSuccess = 1
                             WHERE Id = @id";

            ExecuteSqlByDapper(sql, new { id = id, nextTime = maxTime });
        }


        public List<BetRecordMerchantNotifyQueueEntity> GetByFilter(QueryBetRecordMerchantNotifyQueueRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<BetRecordMerchantNotifyQueueEntity>(sql, filter);
            return query;
        }


        public PagedList<BetRecordMerchantNotifyQueueEntity> GetPagingByFilter(QueryBetRecordMerchantNotifyQueueRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<BetRecordMerchantNotifyQueueEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryBetRecordMerchantNotifyQueueRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if(filter.BetRecordNoPGS.IsNotNullEmptyText()) 
                sql += " AND BetRecordNoPGS = @BetRecordNoPGS ";

            if (filter.ListOfBetRecordNoPGS.IsNotEmptyNull())
                sql += " AND BetRecordNoPGS in @ListOfBetRecordNoPGS ";

            if (filter.IsMustNotify)            
                sql += " AND ( NextNotifyTime <= GETUTCDATE() )";
            
            return sql;
        }

    }


}
