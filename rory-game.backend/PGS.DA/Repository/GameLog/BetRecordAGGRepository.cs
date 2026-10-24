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
    public class BetRecordAGGRepository : FormatLongPrimaryKeyRepository<BetRecordAGGEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.GameLogDB); }
        }

        protected override string TableName
        {
            get { return "[ThirdParty].[BetRecordAGG]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

      
        public BetRecordAGGEntity GetByBetRecordNoPGS(string orderNoPGS)
        {
            if (string.IsNullOrWhiteSpace(orderNoPGS))
                return null;

            var query = GetByFilter(new QueryBetRecordGeneralRequestViewModel
            {
                BetRecordNoPGS = orderNoPGS
            });

            return query.FirstOrDefault();
        }


        public BetRecordAGGEntity GetByBetRecordNo(string orderNo, string thirdPartyGameCode = "")
        {
            if (string.IsNullOrWhiteSpace(orderNo))
                return null;

            var query = GetByFilter(new QueryBetRecordGeneralRequestViewModel
            {
                ThirdPartyGameCode = thirdPartyGameCode,
                BetRecordNo = orderNo
            });

            return query.FirstOrDefault();
        }


        public List<BetRecordAGGEntity> GetByFilter(QueryBetRecordGeneralRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<BetRecordAGGEntity>(sql, filter);
            return query;
        }


        public PagedList<BetRecordAGGEntity> GetPagingByFilter(QueryBetRecordGeneralRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<BetRecordAGGEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryBetRecordGeneralRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} WITH (NOLOCK) 
                           WHERE 1=1 ";

            if (filter.BetRecordNo.IsNotNullEmptyText())
                sql += " AND BetRecordNo = @BetRecordNo ";

            if (filter.BetRecordNoPGS.IsNotNullEmptyText())
                sql += " AND BetRecordNoPGS = @BetRecordNoPGS ";

            if (filter.ThirdPartyGameCode.IsNotNullEmptyText())
                sql += " AND ThirdPartyGameCode = @ThirdPartyGameCode ";

            if (filter.PlayerIDs.IsNotEmptyNull())
                sql += " AND PlayerID in @PlayerIDs ";

            return sql;
        }

    }


}
