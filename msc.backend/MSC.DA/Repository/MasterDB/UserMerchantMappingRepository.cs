using Dapper;
using MSC.Core.Repository;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using System;
using System.Collections.Generic;
using MSC.Utility.Extensions;
using System.Linq;

namespace MSC.DA.Repository
{

    public class UserMerchantMappingRepository : IDapperRepository<UserMerchantMappingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[AdminSys].[UserMerchantMapping]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "MappingID"; }
        }


        public void DeleteByUserID(int userID)
        {
            string deleteSql = $"DELETE {TableName} WHERE userID = @userID ";

            using (var conn = InstanceSQLConn())
            {
                conn.Execute(deleteSql, new { userID = userID });
            }
        }

        public bool IsMerchantPermission(int userID , int merchantID)
        {
            var query = GetByFilter(new QueryUserMerchantMappingRequestViewModel
            {
                UserIDs = userID.ToListCollection(),
                MerchantIDs = merchantID.ToListCollection()
            });

            return query.Any();
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

            var query = GetPageingEntitiesBySQLScript<UserMerchantMappingEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryUserMerchantMappingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.UserIDs.IsNotEmptyNull())
                sql += " AND UserID in @UserIDs ";

            return sql;
        }

    }


}
