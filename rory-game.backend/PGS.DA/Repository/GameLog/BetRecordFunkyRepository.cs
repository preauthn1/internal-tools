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
using target.Models.ViewModels.GameVendor;
using target.Models.Enums;
using System.Transactions;
using target.Models.Models;
using target.Utility.Helper;

namespace target.DA.Repository
{

    public class BetRecordFunkyRepository : FormatLongPrimaryKeyRepository<BetRecordFunkyEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.GameLogDB); }
        }

        protected override string TableName
        {
            get { return "[ThirdParty].[BetRecordFunky]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }


        public BetRecordSummaryReportModel GetTotalSummaryReport(DateTime beginUTC , DateTime endUTC)
        {
            string sql = $@"SELECT count(*) as TotalBetCount , 
                                sum(BetAmount) as TotalBetAmount , 
                                sum(WinLossAmount) as TotalWinLoss , 
                                sum(PayoutAmount) as TotalPayoutAmount,
                                sum(ValidAmount) as TotalValidBetAmount
                           FROM {TableName} WITH (NOLOCK)
                           WHERE Status = @status AND SettlementTime >= @beginUTC AND SettlementTime < @endUTC";

            var para = new 
            { 
                status = GameVendorBetRecordStatusEnum.SettleComplete.ToInteger(),
                beginUTC, 
                endUTC 
            };

            var query = QueryByDapper<BetRecordSummaryReportModel>(sql, para);
            return query.FirstOrDefault();
        }


        public BetRecordFunkyEntity GetByBetRecordNoPGS(string orderNoPGS)
        {
            if (string.IsNullOrWhiteSpace(orderNoPGS))
                return null;

            var query = GetByFilter(new QueryBetRecordGeneralRequestViewModel
            {
                BetRecordNoPGS = orderNoPGS
            });

            return query.FirstOrDefault();
        }


        public BetRecordFunkyEntity GetByBetRecordNo(string orderNo , string gameCode = "")
        {
            if (string.IsNullOrWhiteSpace(orderNo))
                return null;

            var query = GetByFilter(new QueryBetRecordGeneralRequestViewModel
            {
                ThirdPartyGameCode = gameCode,
                BetRecordNo = orderNo
            });

            return query.FirstOrDefault();
        }


        public List<BetRecordFunkyEntity> GetByFilter(QueryBetRecordGeneralRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<BetRecordFunkyEntity>(sql, filter);
            return query;
        }
        
        public PagedList<BetRecordFunkyEntity> GetPagingByFilter(QueryBetRecordGeneralRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<BetRecordFunkyEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryBetRecordGeneralRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.BetRecordNo.IsNotNullEmptyText())
                sql += " AND BetRecordNo = @BetRecordNo ";

            if (filter.BetRecordNoPGS.IsNotNullEmptyText())
                sql += " AND BetRecordNoPGS = @BetRecordNoPGS ";

            return sql;
        }

    }


}
