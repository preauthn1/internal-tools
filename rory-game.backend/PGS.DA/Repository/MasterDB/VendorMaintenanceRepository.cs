using Dapper;
using game.Entity;
using Lck.Utility.Extensions.Helper;
using target.Core.Repository;
using target.Models.Enums;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using System;
using System.Collections.Generic;
using System.Linq;

namespace target.DA.Repository
{
    public class VendorMaintenanceRepository : FormatIntegerPrimaryKeyRepository<VendorMaintenanceEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[ThirdParty].[VendorMaintenance]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public VendorMaintenanceEntity GetVendorMaintenanceByCache(int vendorID)
        {
            string cacheKey = CacheKey.DbEntity_VendorMaintain_VendorID(vendorID);

            var queryCache = CacheHelper.GetCacheOrQuery(cacheKey, () =>
               {
                   var query = GetByFilter(new QueryVendorMaintenanceRequestViewModel
                   {
                       VendorIDs = vendorID.ToListCollection(),
                       Status = NormalActiveStatus.Enable.ToInteger(),
                       IsNotExpired = true
                   });

                   return query.FirstOrDefault();

               }, 1);

            return queryCache;            
        }


        public VendorMaintenanceEntity GetVendorMaintenance(int vendorID)
        {
            var query = GetByFilter(new QueryVendorMaintenanceRequestViewModel
            {
                VendorIDs = vendorID.ToListCollection(),
                Status = NormalActiveStatus.Enable.ToInteger(),
                IsNotExpired = true
            });

            return query.FirstOrDefault();
        }


        public List<VendorMaintenanceEntity> GetByFilter(QueryVendorMaintenanceRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<VendorMaintenanceEntity>(sql, filter);
            return query;

        }

        public PagedList<VendorMaintenanceEntity> GetPagingByFilter(QueryVendorMaintenanceRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<VendorMaintenanceEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryVendorMaintenanceRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.VendorIDs.IsNotEmptyNull())
                sql += " AND VendorID in @VendorIDs ";

            if (!filter.IsIncludeDelete)
                sql += $" AND Status != {NormalActiveStatus.PermanentDelete.ToInteger()} ";

            if (filter.IsNotExpired)
                sql += $" AND (GETUTCDATE() >= StartTime AND EndTime > GETUTCDATE()) ";

            if (filter.Status.HasValue)
                sql += $" AND Status = @Status ";

            return sql;
        }

    }


}
