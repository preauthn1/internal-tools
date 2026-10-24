using Dapper;
using MSC.Core.Repository;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using System;
using System.Collections.Generic;
using MSC.Utility.Extensions;
using System.Linq;
using System.Data.SqlClient;
using System.Data;
using MSC.Models.Models;

namespace MSC.DA.Repository
{

    public class CardConsumeRecordRepository : IDapperRepository<CardConsumeRecordEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[CardConsumeRecord]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        /// <summary>
        /// 以Transaction交易的方式 寫入一筆紀錄
        /// </summary>
        /// <param name="model"></param>
        /// <returns></returns>
        public CommonSPResultModel TransactionWriteConsumeRecord(CardConsumeRecordEntity model)
        {
            string spName = "USP_WriteCardConsumeAndUpdateStatus";
          
            var para = InstanceDynamicParameterWithReturnValue();
            para.Add("@cardNo", model.CardNo);
            para.Add("@merchantID", model.MerchantID);
            para.Add("@orderNo", model.OrderNo);
            para.Add("@point", model.Point);
            para.Add("@memberCode", model.MemberCode);
            para.Add("@merchantRequestID", model.MerchantRequestID);
            para.Add("@ipAddress", model.IPAddress);
            para.Add("@device", model.DeviceID);

            using (var conn = new SqlConnection(ConnectionString))
            {
                var result = conn.Execute(spName, param: para, commandType: CommandType.StoredProcedure);
                int returnCode = para.GetDynamicReturnCode();

                return new CommonSPResultModel
                {
                    Code = (ErrorCode)returnCode,
                };
            }
        }



        public CardConsumeRecordEntity GetByMerchantIDAndRequestID(int merchantID, string merchantRequestID)
        {
            var filter = new QueryCardConsumeRecordRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection(),
                MerchantRequestID = merchantRequestID
            };

            var query = GetByFilter(filter);
            return query.FirstOrDefault();
        }


        public List<CardConsumeRecordEntity> GetByFilter(QueryCardConsumeRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<CardConsumeRecordEntity>(sql, filter);
            return query;
        }


        public PagedList<CardConsumeRecordEntity> GetPagingByFilter(QueryCardConsumeRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<CardConsumeRecordEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryCardConsumeRecordRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND ID in @IDs ";

            if (filter.ListOfCardNo.IsNotEmptyNull())
                sql += " AND CardNo in @ListOfCardNo ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.OrderNo.IsNotNullNotWhiteSpace())
                sql += " AND OrderNo = @OrderNo ";

            if (filter.MemberCode.IsNotNullNotWhiteSpace())
                sql += " AND MemberCode = @MemberCode ";

            if (filter.CardNo.IsNotNullNotWhiteSpace())
                sql += " AND CardNo = @CardNo ";

            if (filter.MerchantRequestID.IsNotNullNotWhiteSpace())
                sql += " AND MerchantRequestID = @MerchantRequestID ";

            if (filter.CreateTimeBegin.HasValue)
                sql += " AND CreateTime >= @CreateTimeBegin ";

            if (filter.CreateTimeEnd.HasValue)
                sql += " AND CreateTime < @CreateTimeEnd ";

            return sql;
        }

    }


}
