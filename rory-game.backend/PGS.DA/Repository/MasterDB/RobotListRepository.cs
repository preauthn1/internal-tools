using Dapper;
using game.Entity;
using target.Core.Repository;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using System;
using System.Data;
using System.Collections.Generic;
using target.Utility.Extensions;
using System.Linq;

namespace target.DA.Repository
{

    public class RobotListRepository : FormatIntegerPrimaryKeyRepository<RobotListEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[RobotList]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "PlayerID"; }
        }


        /// <summary>
        /// 清除某台server所有機器人指派 (已被指派的 不會再被機器人分配選中)
        /// </summary>
        public void RemoveServerAssignRobots(int serverID)
        {
            string sql = $@" UPDATE {TableName}
                             SET AssignServerID = 0	
                             WHERE AssignServerID = @serverID";

            ExecuteSqlByDapper(sql , new { serverID = serverID });
        }

        public RobotListResultModel RobotLogon(RobotListRequestModel model)
        {
            RobotListResultModel result = new RobotListResultModel();

            string spName = "[PrGs_RobotLogon]";
            var para = new DynamicParameters();
            para.Add("intServerID", model.ServerID);
            para.Add("intCurrentRobotCount", model.CurrentRobotCount);
            para.Add("lngMinMoneyAmount", model.MinMoneyAmount);
            para.Add("lngMaxMoneyAmount", model.MaxMoneyAmount);
            para.Add("intLogoutRobotCount", dbType: DbType.Int32, direction: ParameterDirection.Output);

            result.Robots = ExecuteStoreProcedureResultByDapper<RobotListItemResultModel>(spName, para);
            result.LogoutRobotCount = para.Get<int>("intLogoutRobotCount");
            return result;
        }


        public List<RobotListEntity> GetByFilter(QueryRobotListRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<RobotListEntity>(sql, filter);
            return query;
        }

        public PagedList<RobotListEntity> GetPagingByFilter(QueryRobotListRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<RobotListEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryRobotListRequestViewModel filter)
        {
            string sql = $@"SELECT * FROM {TableName} 
                           WHERE 1=1 ";

            if (filter.PlayerIDs.IsNotEmptyNull())
                sql += " AND PlayerIDs in @PlayerIDs ";

            if (filter.LevelIDs.IsNotEmptyNull())
                sql += " AND LevelID in @LevelIDs ";

            return sql;
        }

    }


}
