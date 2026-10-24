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

    public class ThirdPartySlot88MD5RecordRepository : FormatIntegerPrimaryKeyRepository<ThirdPartySlot88MD5RecordEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[ThirdPartySlot88MD5Record]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<ThirdPartySlot88MD5RecordEntity> GetByFilter(QueryThirdPartySlot88MD5RecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<ThirdPartySlot88MD5RecordEntity>(sql, filter);
            return query;
        }


        public PagedList<ThirdPartySlot88MD5RecordEntity> GetPagingByFilter(QueryThirdPartySlot88MD5RecordRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<ThirdPartySlot88MD5RecordEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryThirdPartySlot88MD5RecordRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (!string.IsNullOrWhiteSpace(filter.MD5Text))
                sql += " AND MD5 = @MD5Text ";

            if (filter.PlayerIDs.IsNotEmptyNull())
                sql += " AND PlayerID in @PlayerIDs ";

            return sql;
        }

    }


}
