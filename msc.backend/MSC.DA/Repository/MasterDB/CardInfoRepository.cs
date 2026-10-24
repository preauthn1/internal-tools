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
using MSC.Models.Enums;
using MSC.Models;
using MSC.Models.Models;

namespace MSC.DA.Repository
{

    public class CardInfoRepository : IDapperRepository<CardInfoEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[CardInfo]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "CardNo"; }
        }


        /// <summary>
        /// 卡號被使用 主表 會透過以下欄位記錄是哪個商戶使用
        /// </summary>
        public void UpdateConsumeMerchantID(List<string> listOfCards, int merchantID)
        {
            string sql = @$" UPDATE {TableName} 
                             SET ConsumeMerchantID = @merchantID
                             WHERE CardNo in @listOfCards AND ConsumeMerchantID is null ";

            ExecuteSqlByDapper(sql, new { listOfCards, merchantID });
        }


        /// <summary>
        /// 卡號被使用 主表 會透過以下欄位記錄是哪個商戶使用
        /// </summary>
        public void UpdateConsumeMerchantID(string cardNo, int merchantID)
        {
            string sql = @$" UPDATE {TableName} 
                             SET ConsumeMerchantID = @merchantID
                             WHERE CardNo = @cardNo ";

            ExecuteSqlByDapper(sql, new { cardNo, merchantID });
        }

        
        public CommonSPResultModel CheckIsValidCard(string cardNo , int merchantID)
        {
            string spName = "USP_CheckIsValidCard";

            var para = InstanceDynamicParameterWithReturnValue();
            para.Add("@cardNo", cardNo);
            para.Add("@merchantID", merchantID);

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

        public CardInfoEntity GetByCardNo(string cardNo)
        {
            var query = GetByFilter(new QueryCardInfoRequestViewModel
            {
                ListOfCardNo = cardNo.ToListCollection()
            });

            return query.FirstOrDefault();
        }

        public List<CardInfoEntity> GetByFilter(QueryCardInfoRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<CardInfoEntity>(sql, filter);
            return query;
        }


        public PagedList<CardInfoEntity> GetPagingByFilter(QueryCardInfoRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<CardInfoEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryCardInfoRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} {GenerateNoLockQuery(filter.IsLockQuery)}
                           WHERE 1=1 ";

            if (filter.CardNo.IsNotNullNotWhiteSpace())
                sql += " AND CardNo = @CardNo ";

            if (filter.Memo.IsNotNullNotWhiteSpace())
                sql += " AND Memo = @Memo ";

            if (filter.Currency.IsNotNullNotWhiteSpace())
                sql += " AND Currency = @Currency ";

            if (filter.ListOfCardNo.IsNotEmptyNull())
                sql += " AND CardNo in @ListOfCardNo ";

            if (filter.Point.HasValue)
                sql += " AND Point = @Point ";

            if (filter.CreateTimeBegin.HasValue)
                sql += " AND CreateTime >= @CreateTimeBegin ";

            if (filter.CreateTimeEnd.HasValue)
                sql += " AND CreateTime <= @CreateTimeEnd ";

            if(filter.CardStatus.HasValue)
            {
                if(filter.CardStatus.Value == CardInfoStatusEnum.Consumed)                
                    sql += " AND CardNo in (SELECT CardNo from [dbo].[CardConsumeRecord]) ";                

                if (filter.CardStatus.Value == CardInfoStatusEnum.Pending)                
                    sql += "  AND CardNo not in (SELECT CardNo from [dbo].[CardConsumeRecord])  ";                

                if (filter.CardStatus.Value == CardInfoStatusEnum.Destroy)
                {
                    int magicMerchant = SystemConstant.MagicConsumeMerchantID;
                    sql += $" AND ConsumeMerchantID = {magicMerchant} ";
                }
            }

            return sql;
        }

    }


}
