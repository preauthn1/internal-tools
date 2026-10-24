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

    public class MerchantGameSettingRepository : FormatIntegerPrimaryKeyRepository<MerchantGameSettingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[MerchantGameSetting]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        

        /// <summary>
        /// 檢查 1. 商戶是否開啟此遊戲
        /// 檢查 2. 系統是否開啟此遊戲
        /// </summary>
        public bool IsMerchantGameEnableInSystem(int merchantID , int gameID)
        {
            string sql = @"SELECT Id FROM [dbo].[MerchantGameSetting]
                           WHERE MerchantID = @merchantID AND[Status] = 1 AND GameID in
                            (
                                SELECT GameID FROM GameList
                                WHERE Status = 1 AND GameID = @gameID
                            ) ";

            using (var conn = InstanceSQLConn())
            {
                var query = conn.Query<int>(sql, new { merchantID, gameID }).FirstOrDefault();
                return query > 0;
            }
        }
        

        public List<MerchantGameSettingEntity> GetByMerchantID(int merchantID , NormalActiveStatus? status = null)
        {
            var filter = new QueryMerchantGameSettingRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection()
            };

            if (status.HasValue)
                filter.Status = status.Value.ToInteger();

            var query = GetByFilter(filter);
            return query;
        }


        public MerchantGameSettingEntity GetByMerchantIDAndGameID(int merchantID , int gameID)
        {
            var filter = new QueryMerchantGameSettingRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection(),
                GameIDs = gameID.ToListCollection()
            };

            var query = GetByFilter(filter);
            return query.FirstOrDefault();
        }


        public List<MerchantGameSettingEntity> GetByFilter(QueryMerchantGameSettingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<MerchantGameSettingEntity>(sql, filter);
            return query;
        }


        public PagedList<MerchantGameSettingEntity> GetPagingByFilter(QueryMerchantGameSettingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<MerchantGameSettingEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryMerchantGameSettingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.GameIDs.IsNotEmptyNull())
                sql += " AND GameID in @GameIDs ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.Status.HasValue)
                sql += " AND Status = @Status ";

            if (filter.EnableGameOnly)
                sql += " AND GameID in (SELECT GameID FROM [dbo].[GameList] WHERE Status = 1) ";

            return sql;
        }

    }


}
