using Dapper;
using MSC.Core.Repository;
using MSC.Models.Enums;
using MSC.Models.Models;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using System;
using System.Collections.Generic;
using System.Data;
using System.Data.SqlClient;
using System.Linq;

namespace MSC.DA.Repository
{
    public class MerchantWalletRepository : IDapperRepository<MerchantWalletEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[MerchantWallet]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "MerchantID"; }
        }


        public CommonSPResultModel UpdateWalletAndWriteLog(UpdateMerchantWalletSPRequestParameter model)
        {
            string spName = "USP_UpdateMerchantWalletAndLog";

            var para = InstanceDynamicParameterWithReturnValue();
            para.Add("@merchantID", model.MerchantID);
            para.Add("@sourceAmount", model.SourceAmount);
            para.Add("@feeRate", model.FeeRate);
            para.Add("@transactionType", model.TransactionType.ToInteger());
            para.Add("@refNo", model.RefNo);

            using (var conn = new SqlConnection(ConnectionString))
            {
                var result = conn.Execute(spName, param: para, commandType: CommandType.StoredProcedure);
                int returnCode = para.GetDynamicReturnCode();

                return new CommonSPResultModel
                {
                    Code = (ErrorCode)returnCode,
                };
            }
        }


        public MerchantWalletEntity GetByMerchantID(int merchantID)
        {
            var query = GetByFilter(new QueryMerchantWalletRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection(),
                WalletTypes = MerchantWalletTypeEnum.Main.ToInteger().ToListCollection()
            });

            return query.FirstOrDefault();
        }



        public List<MerchantWalletEntity> GetByFilter(QueryMerchantWalletRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<MerchantWalletEntity>(sql, filter);
            return query;
        }


        public PagedList<MerchantWalletEntity> GetPagingByFilter(QueryMerchantWalletRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<MerchantWalletEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryMerchantWalletRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (filter.WalletTypes.IsNotEmptyNull())
                sql += " AND WalletType in @WalletTypes ";

            return sql;
        }

    }


}
