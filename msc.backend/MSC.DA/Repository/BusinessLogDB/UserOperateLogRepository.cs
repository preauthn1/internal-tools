using Dapper;
using MSC.Core.Repository;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using System;
using System.Collections.Generic;
using MSC.Utility.Extensions;
using System.Linq;
using MSC.Models.Enums;
using MSC.Utility.Helper;
using MSC.Utility;

namespace MSC.DA.Repository
{

    public class UserOperateLogRepository : IDapperRepository<UserOperateLogEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.BusinessLogDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[UserOperateLog]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "Id"; }
        }


        public void CreateByParameter(int userID , UserOperateTypeEnum type , long? refID = null , object debugInfo = null , string memo = "")
        {
            string userName = RepositoryHelper.GetUserNameByUserID(userID);
            var insert = new UserOperateLogEntity
            {
                CreateTime = DateTimeHelper.GetCurrentTime(),
                DebugInfo = debugInfo == null ? "" : debugInfo.ConvertToJsonString(),
                UserID = userID,
                UserName = userName,
                IPAddress = IPHelper.GetIP(),
                OperationType = type.ToInteger(),
                DeviceID = HeaderValueHelper.GetDeviceID(),
                Memo = memo,
                ReferenceID = refID,
            };

            this.Create(insert);
            return;
        }

        
        public List<UserOperateLogEntity> GetByFilter(QueryUserOperateLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<UserOperateLogEntity>(sql, filter);
            return query;
        }


        public PagedList<UserOperateLogEntity> GetPagingByFilter(QueryUserOperateLogRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPageingEntitiesBySQLScript<UserOperateLogEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryUserOperateLogRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.UserIDs.IsNotEmptyNull())
                sql += " AND UserID in @UserIDs ";

            if (filter.UserName.IsNotNullNotWhiteSpace())
                sql += " AND UserName = @UserName ";

            if (filter.OperationTypes.IsNotEmptyNull())
                sql += " AND OperationType in @OperationTypes ";

            if (filter.CreateTimeBegin.HasValue)
                sql += " AND CreateTime >= @CreateTimeBegin ";

            if (filter.CreateTimeEnd.HasValue)
                sql += " AND CreateTime < @CreateTimeEnd ";

            return sql;
        }

    }


}
