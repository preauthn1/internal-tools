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

    public class ServerMerchantMappingRepository : FormatIntegerPrimaryKeyRepository<ServerMerchantMappingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[ServerMerchantMapping]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }


        public bool IsServerMerchantEnable(int serverID , int merchantID)
        {
            var filter = new QueryServerMerchantRelationRequestViewModel
            {
                ServerIDs = serverID.ToListCollection(),
                MerchantIDs = merchantID.ToListCollection()
            };

            string sql = GenerateWhereSQL(filter , "Id");
            var query = QueryByDapper<ServerMerchantMappingEntity>(sql, filter);
            return query.IsNotEmptyNull();
        }


        public List<ServerMerchantMappingEntity> GetByServerID(int serverID)
        {
            var filter = new QueryServerMerchantRelationRequestViewModel
            {
                ServerIDs = serverID.ToListCollection()
            };

            var query = GetByFilter(filter);
            return query;
        }

        public List<ServerMerchantMappingEntity> GetByMerchantID(int merchantID)
        {
            var filter = new QueryServerMerchantRelationRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection()
            };

            var query = GetByFilter(filter);
            return query;
        }

        public List<int> GetAvailableServersByMerchantIDs(List<int> merchantIDs)
        {
            var filter = new QueryServerMerchantRelationRequestViewModel
            {
                MerchantIDs = merchantIDs
            };

            var query = GetByFilter(filter);
            return query.Select(x=> x.ServerID).ToList();
        }


        public List<ServerMerchantMappingEntity> GetByFilter(QueryServerMerchantRelationRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<ServerMerchantMappingEntity>(sql, filter);
            return query;
        }


        public PagedList<ServerMerchantMappingEntity> GetPagingByFilter(QueryServerMerchantRelationRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<ServerMerchantMappingEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryServerMerchantRelationRequestViewModel filter , string columns = "*")
        {
            string sql = $@"SELECT {columns} FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.ServerIDs.IsNotEmptyNull())
                sql += " AND ServerID in @ServerIDs ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            return sql;
        }

    }


}
