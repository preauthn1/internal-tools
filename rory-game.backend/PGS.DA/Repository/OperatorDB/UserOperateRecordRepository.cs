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

namespace target.DA.Repository
{

    public class UserOperateRecordRepository : FormatIntegerPrimaryKeyRepository<UserOperateRecordEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.OperationDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[UserOperateRecord]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<UserOperateRecordEntity> GetByFilter(QueryUserOperateRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<UserOperateRecordEntity>(sql, filter);
            return query;
        }


        public PagedList<UserOperateRecordEntity> GetPagingByFilter(QueryUserOperateRecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<UserOperateRecordEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryUserOperateRecordRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.OperateTypes.IsNotEmptyNull())
                sql += " AND OperateType in @OperateTypes ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.TargetIDs.IsNotEmptyNull())
                sql += " AND TargetID in @TargetIDs ";

            if(filter.CreateTimeBegin.HasValue)
                sql += " AND CreateTime >= @CreateTimeBegin ";

            if (filter.CreateTimeEnd.HasValue)
                sql += " AND CreateTime <= @CreateTimeEnd ";

            if (!string.IsNullOrWhiteSpace(filter.Memo))
                sql += $" AND ( {GenerateLikeQuery("Memo", "@Memo")} )";

            return sql;
        }

    }


}
