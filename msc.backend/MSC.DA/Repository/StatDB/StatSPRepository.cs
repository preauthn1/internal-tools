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
using MSC.Utility;
using MSC.Models.Models;
using System.Data.SqlClient;
using System.Data;

namespace MSC.DA.Repository
{

    public class StatSPRepository 
    {
        protected string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }


        public void ExecuteCalculateStatisticMerchantReport()
        {
            var spName = "PrJob_CalculateMerchantStatisticReport";

            using (var conn = new SqlConnection(ConnectionString))
            {
                conn.Execute(spName, null, commandType: CommandType.StoredProcedure);
            }
        }


        public List<MerchantDailySummarySPResult> GetMerchantDailySummary(MerchantDailySummarySPRequestParameter filter)
        {
            var spName = "USP_GetStatReport_MerchantDailySummary";

            using (var conn = new SqlConnection(ConnectionString))
            {
                var result = conn.Query<MerchantDailySummarySPResult>(spName, filter, commandType: CommandType.StoredProcedure);
                return result.ToList();
            }
        }

        public List<MerchantDailySummarySPResultV2> GetMerchantDailySummaryV2(MerchantDailySummarySPRequestParameter filter)
        {
            var spName = "USP_GetStatReport_MerchantDailySummary_V2";

            using (var conn = new SqlConnection(ConnectionString))
            {
                var result = conn.Query<MerchantDailySummarySPResultV2>(spName, filter, commandType: CommandType.StoredProcedure);
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
