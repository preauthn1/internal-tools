using Dapper;
using game.Entity;
using Lck.Utility.Extensions.Helper;
using target.Core.Repository;
using target.Models.Enums;
using target.Models.Models;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Transactions;

namespace target.DA.Repository
{
    public class BetRecordRepository : FormatLongPrimaryKeyRepository<BetRecordEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.GameLogDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[BetRecord]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }


        private BetRecordEntity GetByBetRecordNoPGS(string transactionNoPGS)
        {
            var query = GetByFilter(new QueryBetRecordRequestViewModel
            {
                BetRecordNoPGS = transactionNoPGS
            });

            return query.FirstOrDefault();
        }

        public PagedList<UnionBetRecordQueryModel> GetUnionBetRecordsByPagingFilter(QueryUnionBetRecordRequestViewModel filter)
        {
            filter.OrderByColumn = " CreateTime ";
            string commonExtraSQL = "  0 as type , ";
            List<string> allSQLs = new List<string>
                {
                    GenerateUnionQuerySQL(GameVendorTypeEnum.target , " Type as type , " , "[dbo].[BetRecord]"),
                    GenerateUnionQuerySQL(GameVendorTypeEnum.AGG , commonExtraSQL , "[ThirdParty].[BetRecordAGG]"),
                    GenerateUnionQuerySQL(GameVendorTypeEnum.CrazyFishing , commonExtraSQL , "[ThirdParty].[BetRecordFish]"),

                    //GenerateUnionQuerySQL(GameVendorTypeEnum.Funky , commonExtraSQL , "[ThirdParty].[BetRecordFunky]"),
                };

            string unionSQL = string.Join(" UNION " , allSQLs);       
            
            string final_sql = $" SELECT * FROM ({unionSQL}) as u WHERE 1=1 ";

            if (filter.TransactionNoPGS.IsNotNullEmptyText())
                final_sql += " AND BetRecordNoPGS = @TransactionNoPGS ";

            if (filter.TransactionNo.IsNotNullEmptyText())
                final_sql += " AND BetRecordNo = @TransactionNo ";

            if (filter.PlayerName.IsNotNullEmptyText())
                final_sql += " AND PlayerName = @PlayerName ";

            if (filter.Status.HasValue)
                final_sql += " AND Status = @Status ";

            if (filter.PlayerID.HasValue)
                final_sql += " AND PlayerID = @PlayerID ";

            if (filter.MerchantID.HasValue)
                final_sql += " AND MerchantID = @MerchantID ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                final_sql += " AND MerchantID in @MerchantIDs ";

            if (filter.GameID.HasValue)
                final_sql += " AND GameID = @GameID ";

            if (filter.CreateTimeBegin.HasValue)
                final_sql += " AND CreateTime >= @CreateTimeBegin ";

            if (filter.CreateTimeEnd.HasValue)
                final_sql += " AND CreateTime < @CreateTimeEnd ";

            var query = GetPagingEntitiesBySQLScript<UnionBetRecordQueryModel>(final_sql, filter, filter);
            return query;
        }


        private string GenerateUnionQuerySQL(GameVendorTypeEnum gameVendor , string extraColumns , string tableName)
        {
            string common_column = @" BetRecordNo , BetRecordNoPGS , MerchantID, GameID, PlayerID, PlayerName, BetAmount, WinLossAmount , PayoutAmount , Status, CreateTime, SettlementTime ";
            string resultSQL = @$" SELECT {(int)gameVendor} as GameVendor , {extraColumns}  {common_column} FROM {tableName} WITH (NOLOCK) ";
            return resultSQL;
        }


        public void UpdateStatus(string betRecordNoPGS, GameVendorBetRecordStatusEnum status , string memo = "")
        {
            memo = memo.Truncate(500);
            string sql = $@"UPDATE {TableName} SET 
                                Status = @status , 
                                Memo = @memo ,
                                UpdateTime = @updateTime
                            WHERE BetRecordNoPGS = @betRecordNoPGS";

            var para = new
            {
                status = status.ToInteger(),
                memo = memo,
                updateTime = DateTimeHelper.GetCurrentTime(),
                betRecordNoPGS = betRecordNoPGS,
            };

            ExecuteSqlByDapper(sql , para);
        }


        public void UpdateMemo(string betRecordNoPGS, string memo)
        {
            string sql = $@"UPDATE {TableName} SET 
                                Memo = @memo ,
                                UpdateTime = @updateTime
                            WHERE BetRecordNoPGS = @betRecordNoPGS";

            var para = new
            {
                memo = memo,
                updateTime = DateTimeHelper.GetCurrentTime(),
                betRecordNoPGS = betRecordNoPGS,
            };

            ExecuteSqlByDapper(sql, para);
        }




        public BetRecordEntity GetByTransactionNo(string betRecordNoPGS)
        {
            var query = GetByFilter(new QueryBetRecordRequestViewModel
            {
                BetRecordNoPGS = betRecordNoPGS
            });

            return query.FirstOrDefault();
        }


        public List<BetRecordEntity> GetByFilter(QueryBetRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<BetRecordEntity>(sql, filter);
            return query;
        }


        public PagedList<BetRecordEntity> GetPagingByFilter(QueryBetRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<BetRecordEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryBetRecordRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.BetRecordNoPGS.IsNotNullEmptyText())
                sql += " AND BetRecordNoPGS = @BetRecordNoPGS ";

            return sql;
        }

    }


}
