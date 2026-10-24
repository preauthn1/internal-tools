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
using target.Models.Enums;

namespace target.DA.Repository
{

    public class PlayerLoginLogRepository : FormatLongPrimaryKeyRepository<PlayerLoginLogEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.TransactionDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[PlayerLoginLog]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }


        public void UpdateStatus(long id , NormalSuccessStatusEnum status)
        {
            string sql = $" UPDATE {TableName} SET status = @status WHERE Id = @id";

            var para = new
            {
                id = id,
                status = status.ToInteger()
            };

            ExecuteSqlByDapper(sql, para);
        }


        public List<PlayerLoginLogEntity> GetByFilter(QueryPlayerLoginLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<PlayerLoginLogEntity>(sql, filter);
            return query;
        }


        public PagedList<PlayerLoginLogEntity> GetPagingByFilter(QueryPlayerLoginLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<PlayerLoginLogEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryPlayerLoginLogRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} WITH (NOLOCK)
                            WHERE 1=1 ";

            if (filter.PlayerIDs.IsNotEmptyNull())
                sql += " AND PlayerID in @PlayerIDs ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.IsThirdPartyGame)
                sql += " AND GameID >= 100000";

            if (filter.VendorID.HasValue)
                sql += " AND (GameID >= @VendorID AND GameID < (@VendorID + 4999))";

            if (filter.PlayerName.IsNotNullEmptyText())
                sql += " AND PlayerName = @PlayerName";

            if (filter.CreateTimeBegin.HasValue)
                sql += " AND CreateTime >= @CreateTimeBegin";

            if (filter.CreateTimeEnd.HasValue)
                sql += " AND CreateTime < @CreateTimeEnd";

            return sql;
        }

    }


}
