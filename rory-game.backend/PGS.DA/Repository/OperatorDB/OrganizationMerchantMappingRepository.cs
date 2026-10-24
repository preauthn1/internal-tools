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

    public class OrganizationMerchantMappingRepository : FormatIntegerPrimaryKeyRepository<OrganizationMerchantMappingEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.OperationDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[OrganizationMerchantMapping]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

       

        public void DeleteAndInsertNew(int orgID, List<int> merchantIDs)
        {
            merchantIDs = merchantIDs.Distinct().ToList();
            var mappingEntity = merchantIDs.Select(x => new OrganizationMerchantMappingEntity
                {
                    MerchantID = x,
                    OrganizationID = orgID
                })
                .ToList();

            DeleteAndInsertNewData("OrganizationID", orgID, mappingEntity);
        }


        public List<OrganizationMerchantMappingEntity> GetByOrgID(int orgIDs)
        {
            var filter = new QueryOrganizationMerchantMappingRequestViewModel
            {
                OrganizationIDs = orgIDs.ToListCollection(),
            };

            var query = GetByFilter(filter);
            return query;
        }


        public List<OrganizationMerchantMappingEntity> GetByFilter(QueryOrganizationMerchantMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<OrganizationMerchantMappingEntity>(sql, filter);
            return query;
        }


        public PagedList<OrganizationMerchantMappingEntity> GetPagingByFilter(QueryOrganizationMerchantMappingRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<OrganizationMerchantMappingEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryOrganizationMerchantMappingRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.OrganizationIDs.IsNotEmptyNull())
                sql += " AND OrganizationID in @OrganizationIDs ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            return sql;
        }

    }


}
