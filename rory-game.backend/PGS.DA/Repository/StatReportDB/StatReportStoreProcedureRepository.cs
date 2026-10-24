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
using Lck.Utility.Extensions.Helper;
using target.Utility;
using target.Models.Models.StoreProcedures;
using System.Data.SqlClient;
using System.Data;
using target.Models.Enums;

namespace target.DA.Repository
{

    /// <summary>
    /// 這個DB 很多都是透過DBA寫SP溝通的。
    /// </summary>
    public class StatReportStoreProcedureRepository
    {
        protected string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.StatReportDB); }
        }


        public List<MerchantDailySummaryResultModel> GetMerchantDailySummary(Stat_MerchantDailySummaryParameterModel filter)
        {
            var spName = "PrWeb_Stat_MerchantDailySummary";           

            using (var conn = new SqlConnection(ConnectionString))
            {
                var result = conn.Query<MerchantDailySummaryResultModel>(spName, filter, commandType: CommandType.StoredProcedure);
                return result.ToList();
            }
        }


        public List<HomePageSummaryViewModel> GetHomePageSummary(QueryHomePageReportRequestViewModel filter)
        {
            var spName = "PrWeb_Stat_HomePageSummary";
            var para = new
            {
                userID = filter.UserID,
                dateStart = filter.Begin,
                dateEnd = filter.End,
                merchantList = filter.MerchantIDs.IsEmptyOrNull() ? null : string.Join(",", filter.MerchantIDs)
            };

            using (var conn = new SqlConnection(ConnectionString))
            {
                var result = conn.Query<HomePageSummaryViewModel>(spName, para, commandType: CommandType.StoredProcedure);
                return result.ToList();
            }
        }

        /// <summary>
        /// 数据报表(商戶總報表?)
        /// </summary>
        public List<Stat_MerchantSummarySPResultModel> GetMerchantSummaryStatReport(DateTime begin, DateTime end, List<int> merchantIDs)
        {
            var para = new
            {
                dateStart = begin,
                dateEnd = end,
                merchantList = merchantIDs.IsEmptyOrNull() ? null : string.Join(",", merchantIDs),
            };

            using (var conn = new SqlConnection(ConnectionString))
            {
                var result = conn.Query<Stat_MerchantSummarySPResultModel>("PrWeb_Stat_MerchantSummary", para, commandType: CommandType.StoredProcedure);
                return result.ToList();
            }
        }



        /// <summary>
        /// 遊戲報表
        /// </summary>
        public List<Stat_GameSummarySPResultModel> GetGameSummaryStatReport(DateTime begin, DateTime end, List<int> merchantIDs, List<int> gameIDs)
        {
            var para = new
            {
                dateStart = begin,
                dateEnd = end,
                merchantList = merchantIDs.IsEmptyOrNull() ? null : string.Join("," , merchantIDs),
                gameList = gameIDs.IsEmptyOrNull() ? null : string.Join(",", gameIDs),
            };

            using (var conn = new SqlConnection(ConnectionString))
            {
                var result = conn.Query<Stat_GameSummarySPResultModel>("PrWeb_Stat_GameSummary", para, commandType: CommandType.StoredProcedure);
                return result.ToList();
            }
        }


        /// <summary>
        /// 圖表數據
        /// </summary>
        public List<Stat_ChartDataSPResultModel> GetStatReportChartData(DateTime begin, DateTime end, List<int> merchantIDs, ChartDataTypeEnum chartType)
        {
            var para = new
            {
                dateStart = begin,
                dateEnd = end,
                merchantList = merchantIDs.IsEmptyOrNull() ? null : string.Join(",", merchantIDs),
                chartType
            };

            using (var conn = new SqlConnection(ConnectionString))
            {
                var result = conn.Query<Stat_ChartDataSPResultModel>("PrWeb_Stat_MerchantSummary_Chart", para, commandType: CommandType.StoredProcedure);
                return result.ToList();
            }
        }

        private string GetConnectionStringByDbName(string dbName)
        {
            var conn = BaseSettingHelper.GetSettingValue($"ConnectionStrings:{dbName}");
            return conn;
        }

    }


}
