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
using target.Models.ViewModels.GameVendor;
using target.Models.Enums;

namespace target.DA.Repository
{

    public class MerchantGameSummaryDataRepository : FormatIntegerPrimaryKeyRepository<MerchantGameSummaryDataEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.StatReportDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[MerchantGameSummaryData]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "StatTime"; }
        }


        /// <summary>
        /// 遊戲商對帳用
        /// </summary>
        public BetRecordSummaryReportModel GetSummaryForVendorReportDataCheck(List<int> merchantIDs, List<int> gameIDs, DateTime begin, DateTime end)
        {
            string extraWhere = merchantIDs.IsEmptyOrNull() ? "" : " AND MerchantID in @merchantIDs";
            string extraWhereGameID = gameIDs.IsEmptyOrNull() ? "" : " AND GameID in @gameIDs";
            string sql = $@" SELECT SUM(BetAmount) as TotalBetAmount ,
                                    SUM(WinLossAmount) as TotalWinLoss  ,
                                    SUM(ValidAmount) as TotalValidBetAmount  ,
                                    SUM(PayoutAmount) as TotalPayoutAmount 
                            FROM [dbo].[MerchantGameSummaryData] WITH (NOLOCK)
                            WHERE StatTime >= @begin AND StatTime < @end {extraWhere} {extraWhereGameID} ";

            var filter = new
            {
                gameIDs = gameIDs,
                merchantIDs = merchantIDs,
                begin = begin,
                end = end
            };

            var query = QueryByDapper<BetRecordSummaryReportModel>(sql, filter);
            var result = query.FirstOrDefault();
            result.AvailableColumns = new List<GameVendorReportApiAvailableColumnsEnum>
            {
                GameVendorReportApiAvailableColumnsEnum.TotalPayoutAmount,
                GameVendorReportApiAvailableColumnsEnum.TotalValidBetAmount,
                GameVendorReportApiAvailableColumnsEnum.TotalBetAmount,
                GameVendorReportApiAvailableColumnsEnum.TotalWinLoss
            };

            return result;
        }


        public List<MerchantGameSummaryDataEntity> GetEachGameTotalSumOnDuration(List<int> merchantIDs, List<int> gameIDs, DateTime begin, DateTime end)
        {
            string extraWhere = merchantIDs.IsEmptyOrNull() ? "" : " AND MerchantID in @merchantIDs";
            string extraWhereGameID = gameIDs.IsEmptyOrNull() ? "" : " AND GameID in @gameIDs";
            string sql = $@" SELECT GameID , 

                                SUM(TotalRound) as TotalRound, 
                                SUM([BetPlayer]) as [BetPlayer] , 
                                SUM(WinAmount) as WinAmount , 
                                SUM(LoseAmount) as LoseAmount , 
                                SUM(TaxAmount) as TaxAmount ,
                                SUM(FlowAmount) as FlowAmount , 
                                SUM(JackpotAmount) as JackpotAmount ,
                                SUM(BetAmount) as BetAmount ,
                                SUM(PayoutAmount) as PayoutAmount  ,
                                SUM(ValidAmount) as ValidAmount  ,
                                SUM(WinLossAmount) as WinLossAmount  

                            FROM [dbo].[MerchantGameSummaryData] WITH (NOLOCK)
                            WHERE StatTime >= @begin AND StatTime < @end {extraWhere} {extraWhereGameID} 
                            GROUP BY GameID";

            var filter = new
            {
                gameIDs = gameIDs,
                merchantIDs = merchantIDs,
                begin = begin,
                end = end
            };

            var query = QueryByDapper<MerchantGameSummaryDataEntity>(sql, filter);
            return query;
        }


        public List<MerchantGameSummaryDataEntity> GetByFilter(QueryMerchantGameSummaryDataRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<MerchantGameSummaryDataEntity>(sql, filter);
            return query;
        }


        public PagedList<MerchantGameSummaryDataEntity> GetPagingByFilter(QueryMerchantGameSummaryDataRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<MerchantGameSummaryDataEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryMerchantGameSummaryDataRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.GameIDs.IsNotEmptyNull())
                sql += " AND GameID in @GameIDs ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.StartTimeBegin.HasValue)
                sql += " AND StartTime >= @StartTimeBegin ";

            if (filter.StartTimeEnd.HasValue)
                sql += " AND StartTime < @StartTimeEnd ";

            return sql;
        }

    }


}
