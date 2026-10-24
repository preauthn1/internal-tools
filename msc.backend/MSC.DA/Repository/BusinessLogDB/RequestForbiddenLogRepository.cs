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

    public class RequestForbiddenLogRepository : IDapperRepository<RequestForbiddenLogEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.BusinessLogDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[RequestForbiddenLog]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }

        public List<RequestForbiddenLogEntity> GetByFilter(QueryRequestForbiddenLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<RequestForbiddenLogEntity>(sql, filter);
            return query;
        }


        public PagedList<RequestForbiddenLogEntity> GetPagingByFilter(QueryRequestForbiddenLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<RequestForbiddenLogEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryRequestForbiddenLogRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.IDs.IsNotEmptyNull())
                sql += " AND Id in @IDs ";

            if (filter.Types.IsNotEmptyNull())
                sql += " AND Type in @Types ";

            if (filter.MemberCode.IsNotNullNotWhiteSpace())
                sql += " AND MemberCode = @MemberCode ";

            if (filter.CreateTimeBegin.HasValue)
                sql += " AND CreateTime >= @CreateTimeBegin ";

            if (filter.CreateTimeEnd.HasValue)
                sql += " AND CreateTime < @CreateTimeEnd ";


            return sql;
        }

    }


}
