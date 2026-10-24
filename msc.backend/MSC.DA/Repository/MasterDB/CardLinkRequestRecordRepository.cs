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

namespace MSC.DA.Repository
{

    public class CardLinkRequestRecordRepository : IDapperRepository<CardLinkRequestRecordEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[CardLinkRequestRecord]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public void UpdateToDisabled(long id)
        {
            var para = new 
            {
                IsDisabled = true
            };

            Update(id, para);
        }

        
        /// <summary>
        /// 取得某一段時間 該商戶 + memberCode 建了幾筆單
        /// </summary>
        public int GetTotalCountForTheDuration(int merchantID , string memberCode , DateTime begin , DateTime end )
        {
            string sql = @" SELECT count(*) as totalCount 
                            FROM [dbo].[CardLinkRequestRecord] WITH (NOLOCK)
                            WHERE MerchantID = @MerchantID AND MemberCode = @memberCode AND CreateTime between @begin AND  @end";


            var para = new
            {
                MerchantID = merchantID,
                memberCode = memberCode,
                begin = begin,
                end = end,
            };

            using(var conn = new SqlConnection(ConnectionString))
            {
                var query = conn.Query<int>(sql, para);
                return query.FirstOrDefault();
            }
        }


        public CardLinkRequestRecordEntity GetByTokenHash(string hashValue)
        {
            var filter = new QueryCardLinkRequestRecordRequestViewModel
            {
                TokenHash = hashValue
            };

            var query = GetByFilter(filter);
            return query.FirstOrDefault();
        }


        /// <summary>
        /// 基本上是給MSTest使用的
        /// </summary>
        public void MSTest_ResetIsDisabledFlag(int merchantID , string merchantRequestID)
        {
            string sql = @$"UPDATE {TableName} SET IsDisabled = 0 
                            WHERE MerchantID = @merchantID AND MerchantRequestID = @merchantRequestID";

            ExecuteSqlByDapper(sql, new { merchantID , merchantRequestID });
        }


        public bool IsTokenHashExist(string hashValue)
        {
            string sql = $"SELECT Id FROM {TableName} WITH (NOLOCK) WHERE TokenHashValue = @hashValue AND IsDisabled = 0";

            using (var conn = InstanceSQLConn())
            {
                var resultId = conn.Query<int>(sql, new { hashValue = hashValue }).FirstOrDefault();

                if (resultId != 0)
                    return true;

                return false;
            }
        }



        public CardLinkRequestRecordEntity GetByMerchantIDAndRequestID(int merchantID , string merchantRequestID)
        {
            var filter = new QueryCardLinkRequestRecordRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection(),
                MerchantRequestID = merchantRequestID
            };

            var query = GetByFilter(filter);
            return query.FirstOrDefault();
        }



        public List<CardLinkRequestRecordEntity> GetByFilter(QueryCardLinkRequestRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<CardLinkRequestRecordEntity>(sql, filter);
            return query;
        }


        public PagedList<CardLinkRequestRecordEntity> GetPagingByFilter(QueryCardLinkRequestRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<CardLinkRequestRecordEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryCardLinkRequestRecordRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.MerchantRequestID.IsNotNullNotWhiteSpace())
                sql += " AND MerchantRequestID = @MerchantRequestID ";

            if (filter.SystemRequestID.IsNotNullNotWhiteSpace())
                sql += " AND SystemRequestID = @SystemRequestID ";

            if (filter.TokenHash.IsNotNullNotWhiteSpace())
                sql += " AND TokenHashValue = @TokenHash ";

            if (filter.MemberCode.IsNotNullNotWhiteSpace())
                sql += " AND MemberCode = @MemberCode ";

            if (filter.CreateTimeBegin.HasValue)
                sql += " AND CreateTime >= @CreateTimeBegin ";

            if (filter.CreateTimeEnd.HasValue)
                sql += " AND CreateTime < @CreateTimeEnd ";

            return sql;
        }

    }


}
