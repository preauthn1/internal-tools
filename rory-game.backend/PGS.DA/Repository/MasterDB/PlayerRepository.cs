using Dapper;
using game.Entity.Game;
using target.Core.Repository;
using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Data;
using System.Data.Common;
using System.Linq;
using System.Reflection;
using target.Utility.Extensions;
using target.Models.CoreLibs;
using target.Core.Helper;
using target.Models.Enums;
using target.Models.ViewModels;
using target.Models.Paging;
using Lck.Utility.Extensions.Helper;

namespace target.DA.Repository
{
    public class PlayerRepository : FormatIntegerPrimaryKeyRepository<PlayerEntity>
    {
        protected override string ConnectionString
        {
            get { return GetConnectionStringByDbName(DBConnType.MasterDB); }
        }

        protected override string TableName
        {
            get { return "[dbo].[Player]"; }
        }

        protected override string PrimaryKeyColumnName 
        {
            get { return "PlayerID"; }
        }


        public void TestTesttttttttttttttttttt()
        {
            ExecuteStoreProcedureByDapper("[dbo].[SP_TTTTTTEST_Delay]", null);
            return;
        }



        public void TestDebug_TransactionIn()
        {
            using (var connection = InstanceSQLConn())
            {
                connection.Open();

                using (var transaction = connection.BeginTransaction())
                {
                    string sql1 = @" update [dbo].[Player] set UpdateTime = GETDATE() where PlayerID = 1000000 ";
                    string sql2 = @" update [dbo].[PlayerWallet] set Amount = Amount + 1 where PlayerID = 1000000 ";

                    connection.Execute(sql1, null, transaction);
                    connection.Execute(sql2, null, transaction);
                    transaction.Commit();
                }
            }           
        }


        public void UpdatePlayerWalletAfterGameSettle(WriteGameResulSPParameterModel model)
        {
            string spName = "PrWeb_GamePlay_UpdatePlayerWalletAfterGameSettle";
            ExecuteStoreProcedureByDapper(spName , model);
        }


        public CRUDResultModel CreatePlayerAndWallet(PlayerEntity model , decimal defaultWalletAmount)
        {
            string insertPlayerSQL = SQLHelper.GetInsertCommand<PlayerEntity>(TableName);
            string insertPlayerWalletSQL = SQLHelper.GetInsertCommand<PlayerWalletEntity>("[dbo].[PlayerWallet]");
            
            //step 1. 建立player
            //step 2. 建立playerWallet
            //step 3. 建立PlayerTotalSummary

            string sql = $@" DECLARE @newID int;
                            {insertPlayerSQL}
                            SET @newID = (SELECT @@IDENTITY)     
                            INSERT INTO [dbo].[PlayerWallet] (PlayerID , WalletType , Amount) 
                            Values (@newID , {WalletTypeEnum.GameWallet.ToInteger()} , {defaultWalletAmount}) 
                            INSERT INTO [dbo].[PlayerTotalSummary] (PlayerID ,SumTax ,SumLoginTimes ,SumFlowAmount) 
                            VALUES (@newID , 0 , 0 , 0) 
                            SELECT @newID ";

            //包交易transaction
            using (var connection = InstanceSQLConn())
            {
                connection.Open();

                using (var transaction = connection.BeginTransaction())
                {
                    var resultId = connection.Query<int>(sql, model , transaction).FirstOrDefault();
                    transaction.Commit();
                    return new CRUDResultModel
                    {
                        IsSuccess = true,
                        ResultId = resultId
                    };
                }
            }
        }

        public PlayerEntity GetSingleByCache(int playerID)
        {
            string cacheKey = CacheKey.DbEntity_Player(playerID);

            var queryCachae = CacheHelper.GetCacheOrQuery(cacheKey, () =>
            {
                var querySingle = GetById(playerID);
                return querySingle;

            } , expiredMinute: 5);

            return queryCachae;
        }


        public PlayerBriefModel GetBriefByCache(int playerID)
        {
            if (playerID <= 0)
                return null;

            string cacheKey = CacheKey.AdminSys_PlayerBriefInfo(playerID);

            var queryCache = CacheHelper.GetCacheOrQuery(cacheKey, () =>
               {
                   var q = GetBriefByID(playerID);
                   return q;

               }, expiredMinute: 15);

            return queryCache;
        }


        /// <summary>
        /// 查詢精簡版欄位
        /// </summary>
        public PlayerBriefModel GetBriefByID(int playerID)
        {
            string sql = $@" SELECT PlayerID , MerchantID , PlayerName , Status , Currency
                             FROM {TableName} WITH (NOLOCK) 
                             WHERE playerID = @playerID ";

            var query = QueryByDapper<PlayerBriefModel>(sql, new { playerID = playerID });
            return query.FirstOrDefault();
        }
        

        public PlayerEntity GetByPlayerNameAndMerchantID(string playerName, int merchantID)
        {
            var filter = new QueryPlayerRequestViewModel
            {
                PlayerNames = new List<string> { playerName },
                MerchantIDs = merchantID.ToListCollection()
            };

            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<PlayerEntity>(sql, filter);
            return query.FirstOrDefault();
        }


        public PlayerEntity GetByToken(string token)
        {
            var filter = new QueryPlayerRequestViewModel
            {
                IsNoLockQuery = true,
                Token = [REDACTED]
            };

            var query = GetByFilter(filter);
            return query.FirstOrDefault();
        }


        public void UpdateStatus(int playerID, int targetStatus)
        {
            string sql = $@"UPDATE {TableName} 
                            SET Status = @targetStatus , UpdateTime = GETUTCDATE()
                            WHERE PlayerID = @PlayerID";

            var para = new
            {
                targetStatus = targetStatus,
                PlayerID = playerID
            };

            ExecuteSqlByDapper(sql, para);
        }



        public void UpdateTokenAndLastTime(int playerID, string newTokenHash , string lastIPAddress)
        {
            string sql = $@"UPDATE {TableName} SET 
                                LastLoginToken = [REDACTED] ,
                                LastLoginTime = GETUTCDATE() ,
                                LastIPAddress = @lastIPAddress
                            WHERE PlayerID = @PlayerID";
            
            var para = new
            {
                LastLoginToken = [REDACTED],
                PlayerID = playerID,
                lastIPAddress = lastIPAddress
            };

            ExecuteSqlByDapper(sql, para);
        }


        /// <summary>
        /// 多商戶架構下，可能不是唯一值
        /// </summary>
        public List<PlayerEntity> GetByPlayerName(string playerName , int? merchantID = null)
        {
            var filter = new QueryPlayerRequestViewModel
            {
                PlayerNames = new List<string> { playerName },
                MerchantIDs = merchantID.HasValue ? merchantID.Value.ToListCollection() : new List<int>()
            };

            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<PlayerEntity>(sql, filter);
            return query;
        }

        public List<int> GetPlayerIDsByPlayerName(string playerName, int? merchantID = null)
        {
            var query = GetByPlayerName(playerName, merchantID);
            return query.Select(x=> x.PlayerID).ToList();
        }
        

        public List<PlayerEntity> GetByFilter(QueryPlayerRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = QueryByDapper<PlayerEntity>(sql, filter);
            return query;
        }


        public PagedList<PlayerEntity> GetPagingByFilter(QueryPlayerRequestViewModel filter)
        {
            string sql = GenerateWhereSQL(filter);
            var query = GetPagingEntitiesBySQLScript<PlayerEntity>(sql, filter, filter);
            return query;
        }



        private string GenerateWhereSQL(QueryPlayerRequestViewModel filter)
        {
            string noLock = GetNoLockSQL(filter.IsNoLockQuery);
            string sql = $@"SELECT * FROM {TableName} {noLock}
                           WHERE 1=1 ";

            if (filter.PlayerIDs.IsNotEmptyNull())
                sql += " AND PlayerID in @PlayerIDs ";

            if (filter.PlayerNames.IsNotEmptyNull())
                sql += " AND PlayerName in @PlayerNames ";

            if (filter.StatusValues.IsNotEmptyNull())
                sql += " AND Status in @StatusValues ";

            if (filter.MerchantIDs.IsNotEmptyNull())
                sql += " AND MerchantID in @MerchantIDs ";

            if (!string.IsNullOrWhiteSpace(filter.Token))
                sql += " AND LastLoginToken = [REDACTED] ";

            if(!string.IsNullOrWhiteSpace(filter.PlayerName))
                sql += $" AND {GenerateLikeQuery("PlayerName", "@PlayerName")} ";

            if (!string.IsNullOrWhiteSpace(filter.LastIPAddress))
                sql += $" AND LastIPAddress = @LastIPAddress ";

            if (!string.IsNullOrWhiteSpace(filter.CreateIPAddress))
                sql += $" AND CreateIPAddress = @CreateIPAddress ";

            if (filter.CreateTimeBegin.HasValue)
                sql += " AND CreateTime >= @CreateTimeBegin ";

            if (filter.CreateTimeEnd.HasValue)
                sql += " AND CreateTime <= @CreateTimeEnd ";

            if (filter.LastLoginTimeBegin.HasValue)
                sql += " AND LastLoginTime >= @LastLoginTimeBegin ";

            if (filter.LastLoginTimeEnd.HasValue)
                sql += " AND LastLoginTime <= @LastLoginTimeEnd ";

            return sql;
        }

    }


}
