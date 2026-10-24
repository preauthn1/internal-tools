using Dapper;
using MSC.Core.Repository;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using System;
using System.Collections.Generic;
using System.Linq;

namespace MSC.DA.Repository
{
    public class MerchantCardStatisticReportRepository : IDapperRepository<MerchantCardStatisticReportEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[Statistic].[MerchantCardStatisticReport]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<MerchantCardStatisticReportEntity> GetByFilter(QueryMerchantCardStatisticReportRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<MerchantCardStatisticReportEntity>(sql, filter);
            return query;
        }


        public PagedList<MerchantCardStatisticReportEntity> GetPagingByFilter(QueryMerchantCardStatisticReportRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<MerchantCardStatisticReportEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryMerchantCardStatisticReportRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.BeginTime.HasValue)
                sql += " AND StatisticTime >= @BeginTime ";

            if (filter.EndTime.HasValue)
                sql += " AND StatisticTime < @BeginTime ";

            return sql;
        }

    }


}
