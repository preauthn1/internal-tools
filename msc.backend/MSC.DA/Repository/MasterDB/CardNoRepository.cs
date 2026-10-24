using Dapper;
using MSC.Core.Repository;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using System;
using System.Collections.Generic;
using MSC.Utility.Extensions;
using System.Linq;
using System.Data;
using System.Data.SqlClient;

namespace MSC.DA.Repository
{

    public class CardNoRepository : IDapperRepository<CardNoEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[CardNo]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "CardNo"; }
        }

        public List<CheckCardNoExist_SPResult> GetExistCardNo(List<string> listOfCardNo)
        {
            var spName = "USP_CheckCardNoExist";
            var para = new
            {
                @listOfCardNo = string.Join("," , listOfCardNo),
            };

            using (var conn = new SqlConnection(ConnectionString))
            {
                var result = conn.Query<CheckCardNoExist_SPResult>(spName, para, commandType: CommandType.StoredProcedure);
                return result.ToList();
            }
        }


        public List<CardNoEntity> GetByFilter(QueryCardNoRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<CardNoEntity>(sql, filter);
            return query;
        }


        public PagedList<CardNoEntity> GetPagingByFilter(QueryCardNoRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<CardNoEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryCardNoRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.ListOfCardNo.IsNotEmptyNull())
                sql += " AND CardNo in @ListOfCardNo ";

            return sql;
        }

    }


}
