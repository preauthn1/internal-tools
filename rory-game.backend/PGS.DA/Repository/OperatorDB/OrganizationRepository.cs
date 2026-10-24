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
using Lck.Utility.Extensions.Helper;
using target.Models.Enums;

namespace target.DA.Repository
{

    public class OrganizationRepository : FormatIntegerPrimaryKeyRepository<OrganizationEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.OperationDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[Organization]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "OrganizationID"; }
        }

        /// <summary>
        /// 從Cache取資料
        /// </summary>
        public OrganizationEntity GetSingleOrgByCache(int orgID)
        {
            string cacheKey = CacheKey.AdminSys_OrganizationsList();

            var queryCache = CacheHelper.GetCacheOrQuery(cacheKey, () =>
               {
                   var queryAll = GetAll();
                   return queryAll;

               } , expiredMinute: 5);

            var query = queryCache.FirstOrDefault(x => x.OrganizationID == orgID);

            if (query != null)
                return query;

            //如果沒查到可能是cache未刷新
            var queryAgain = GetById(orgID);
            return queryAgain;
        }



        public List<OrganizationEntity> GetByFilter(QueryOrganizationRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<OrganizationEntity>(sql, filter);
            return query;
        }


        public PagedList<OrganizationEntity> GetPagingByFilter(QueryOrganizationRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<OrganizationEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryOrganizationRequestViewModel filter , string selectColumns = "*")
        {
            string sql = $@"SELECT {selectColumns} FROM {TableName} 
                           WHERE 1=1 ";

            if (!filter.IsIncludeDeleteData)
                sql += $" AND Status != {NormalActiveStatus.PermanentDelete.ToInteger()} ";

            if (filter.OrganizationIDs.IsNotEmptyNull())
                sql += " AND OrganizationID in @OrganizationIDs ";

            if(!string.IsNullOrWhiteSpace(filter.OrganizationCode))
                sql += " AND OrganizationCode = @OrganizationCode ";

            return sql;
        }

    }


}
