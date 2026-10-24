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

    public class PlayerCurrentServerRepository : FormatIntegerPrimaryKeyRepository<PlayerCurrentServerEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[PlayerCurrentServer]"; }
        }

        protected override string PrimaryKeyColumnName
        {
            get { return "PlayerID"; }
        }


        /// <summary>
        /// 解鎖某台server全部玩家
        /// </summary>
        public void UnlockAllPlayerByServerID(int serverID)
        {
            string sql = $@" DELETE {TableName} WHERE ServerID = @serverID ";

            ExecuteSqlByDapper(sql, new { serverID  = serverID });
        }


        /// <summary>
        /// 玩家-遊戲進房
        /// </summary>
        public void Player_GameRoom_In(int playerID , int serverID , int gameID , long tradeID, int HDType , decimal currentAmount , string ipAddress)
        {
            string spName = "[PrWeb_GamePlay_PlayerRoomIn]";

            var para = new
            {
                playerID = playerID,
                serverID = serverID,
                gameID = gameID,
                tradeID = tradeID,
                HDType = HDType,
                currentAmount = currentAmount,
                ipAddress = ipAddress
            };

            ExecuteStoreProcedureByDapper(spName , para);
        }


        /// <summary>
        /// 玩家-遊戲退房
        /// </summary>
        public void Player_GameRoom_Out(PlayerRoomOutSPParameterModel para)
        {
            string spName = "[PrWeb_GamePlay_PlayerRoomOut]";
            ExecuteStoreProcedureByDapper(spName, para);
        }


        /// <summary>
        /// 玩家-遊戲退房
        /// </summary>
        public PlayerCurrentServerEntity GetByPlayerID(int playerID)
        {
            var query = GetByFilter(new QueryPlayerCurrentServerRequestViewModel
            {
                PlayerIDs = playerID.ToListCollection()
            });

            return query.FirstOrDefault();
        }


        public List<ServerConcurrentCountModel> GetOnlinePlayersCount(List<int> serverIDs)
        {
            serverIDs = serverIDs.Distinct().ToList();

            string sql = $@"SELECT ServerID , count(*) as  TotalCount
                            FROM [dbo].[PlayerCurrentServer]
                            WHERE ServerID in @serverIDs
                            Group by ServerID ";

            var query = QueryByDapper<ServerConcurrentCountModel>(sql, new { serverIDs = serverIDs });

            return serverIDs
                .Select(x => new ServerConcurrentCountModel
                {
                    ServerID = x,
                    TotalCount = query.Where(z => z.ServerID == x).Count()
                })
                .ToList();
        }


        public List<PlayerCurrentServerEntity> GetByFilter(QueryPlayerCurrentServerRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<PlayerCurrentServerEntity>(sql, filter);
            return query;
        }


        public PagedList<PlayerCurrentServerEntity> GetPagingByFilter(QueryPlayerCurrentServerRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);

            var query = GetPagingEntitiesBySQLScript<PlayerCurrentServerEntity>(sql, filter, filter);
            return query;
        }


        private string GenerateWhereSQL(QueryPlayerCurrentServerRequestViewModel filter)
        {
            string sql = $@"SELECT ps.* FROM {TableName} as ps
                           WHERE 1=1 ";

            if (filter.PlayerIDs.IsNotEmptyNull())
                sql += " AND ps.PlayerID in @PlayerIDs ";

            if (filter.ServerIDs.IsNotEmptyNull())
                sql += " AND ps.ServerID in @ServerIDs ";

            return sql;
        }

    }


}
