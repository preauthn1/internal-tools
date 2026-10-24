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

    public class UserMerchantMappingRepository : FormatIntegerPrimaryKeyRepository<UserMerchantMappingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.OperationDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[UserMerchantMapping]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public void DeleteAndInsertNew(int userID, List<int> merchantIDs)
        {
            merchantIDs = merchantIDs.Distinct().ToList();
            var models = merchantIDs.Select(x => new UserMerchantMappingEntity
                {
                    UserID = userID,
                    MerchantID = x
                })
                .ToList();

            DeleteAndInsertNewData("UserID", userID, models);
        }
        

        public List<UserMerchantMappingEntity> GetByUserID(int userID)
        {
            var filter = new QueryUserMerchantMappingRequestViewModel
            {
                UserIDs = userID.ToListCollection()
            };

            var query = GetByFilter(filter);
            return query;
        }


        public List<UserMerchantMappingEntity> GetByFilter(QueryUserMerchantMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<UserMerchantMappingEntity>(sql, filter);
            return query;
        }
        
        public PagedList<UserMerchantMappingEntity> GetPagingByFilter(QueryUserMerchantMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<UserMerchantMappingEntity>(sql, filter, filter);
            return query;
        }
        
        private string GenerateWhereSQL(QueryUserMerchantMappingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.UserIDs.IsNotEmptyNull())
                sql += " AND UserID in @UserIDs ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            return sql;
        }

    }


}
