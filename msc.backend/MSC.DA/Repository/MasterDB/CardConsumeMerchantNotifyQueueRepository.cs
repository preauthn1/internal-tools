using Dapper;
using MSC.Core.Repository;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using System;
using System.Collections.Generic;
using MSC.Utility.Extensions;
using System.Linq;
using MSC.Models.Enums;
using MSC.Utility.Helper;

namespace MSC.DA.Repository
{

    public class CardConsumeMerchantNotifyQueueRepository : IDapperRepository<CardConsumeMerchantNotifyQueueEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[CardConsumeMerchantNotifyQueue]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }


        public void UpdateStatus(long id, NotifyQueueStatusEnum status)
        {           
            var para = new
            {
                Status = status.ToInteger(),
                UpdateTime = DateTimeHelper.GetCurrentTime()
            };

            Update(id, para);
            return;
        }


        public void UpdateStatusAndLastExecutedInfo(long id , string lastTimeResponse , string requestUUID)
        {
            string sql = $@"UPDATE {TableName} SET 
                                LastTimeResponse = @lastTimeResponse,
                                LastTimeRequestUUID = @requestUUID,
                                RetryTimes = RetryTimes + 1,
                                UpdateTime = GETUTCDATE(),
                                Status = @status
                           WHERE id = @id";

            var para = new
            {
                id = id,
                lastTimeResponse = lastTimeResponse,
                requestUUID = requestUUID,
                status = NotifyQueueStatusEnum.Must_Retry
            };

            ExecuteSqlByDapper(sql, para);
            return;
        }

        public CardConsumeMerchantNotifyQueueEntity GetByMerchantIDAndRequestID(int merchantID, string merchantRequestID)
        {
            var filter = new QueryCardConsumeMerchantNotifyQueueRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection(),
                MerchantRequestID = merchantRequestID
            };

            var query = GetByFilter(filter);
            return query.FirstOrDefault();
        }


        public List<CardConsumeMerchantNotifyQueueEntity> GetByFilter(QueryCardConsumeMerchantNotifyQueueRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<CardConsumeMerchantNotifyQueueEntity>(sql, filter);
            return query;
        }


        public PagedList<CardConsumeMerchantNotifyQueueEntity> GetPagingByFilter(QueryCardConsumeMerchantNotifyQueueRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<CardConsumeMerchantNotifyQueueEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryCardConsumeMerchantNotifyQueueRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.Status.HasValue)
                sql += " AND Status = @Status ";

            if (filter.CardNo.IsNotNullNotWhiteSpace())
                sql += " AND CardNo = @CardNo ";

            if (filter.ReferenceOrderNo.IsNotNullNotWhiteSpace())
                sql += " AND ReferenceOrderNo = @ReferenceOrderNo ";

            if (filter.MerchantRequestID.IsNotNullNotWhiteSpace())
                sql += " AND MerchantRequestID = @MerchantRequestID ";

            if (filter.MerchantRequestIDs.IsNotEmptyNull())
                sql += " AND MerchantRequestID in @MerchantRequestIDs ";

            if (filter.CreateTimeBegin.HasValue)
                sql += " AND CreateTime >= @CreateTimeBegin ";

            if (filter.CreateTimeEnd.HasValue)
                sql += " AND CreateTime < @CreateTimeEnd ";

            return sql;
        }

    }


}
