using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.BL.Interface;
using target.Models.Enums;
using target.Models.Models;
using System.Linq;
using target.Utility.Helper;
using target.Models;
using target.DA.Repository;
using target.Utility;

namespace target.BL.Services
{

    public class GameServerDistributeRefService : BaseService, IGameServerDistributeRefService
    {

        public ServerListRepository ServerListRepository { get; set; }

        public PlayerWalletRepository PlayerWalletRepository { get; set; }

        public PlayerCurrentServerRepository PlayerCurrentServerRepository { get; set; }

        public IPlayerWalletRefService PlayerWalletRefService { get; set; }

        public MerchantRepository MerchantRepository { get; set; }



        public PlayerLaunchGameServerChooseResultModel ChooseAvailableGameServer(int playerID , int merchantID , int gameID , ServerDistributeOptionalParameter optional_para)
        {

            //進房再檢查即可
            //var queryCurrentServer = PlayerCurrentServerRepository.GetById(playerID);

            ////禁止重複入房
            //if (queryCurrentServer != null)
            //    throw new TerproException(ApiErrorCode.Forbidden_DuplicateEnterGameServer);


            var queryMerchant = MerchantRepository.GetById(merchantID);

            List<ServerListEntity> queryGameServers = ServerListRepository.GetByFilter(new QueryServerListRequestViewModel
            {
                GameIDs = gameID.ToListCollection(),
                MerchantIDs = merchantID.ToListCollection(),
                Status = ServerStatusEnum.Enable.ToInteger()
            });

            if (queryGameServers.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.NotFoundServer);

            queryGameServers = queryGameServers
                .Where(x => x.Currency == queryMerchant.Currency)
                .ToList();

            if (queryGameServers.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.NotFoundServer);

            if (queryGameServers.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.NotFoundAvailableServer_MerchantApiSelected);

            //如果有指定参数 游戏厅等级
            if (optional_para.RoomFlagLevelID.HasValue)
            {
                queryGameServers = queryGameServers
                    .Where(x => x.RoomFlagLevel == optional_para.RoomFlagLevelID.Value)
                    .ToList();
            }

            //如果有指定serverID (內部使用的隱藏參數)
            if (optional_para.SecretParameter_ServerID.HasValue)
            {
                queryGameServers = queryGameServers
                    .Where(x => x.ServerID == optional_para.SecretParameter_ServerID.Value)
                    .ToList();
            }

            if (queryGameServers.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.NotFoundAvailableServer_FlagLevelNotMatch);


            List<int> serverIDs = queryGameServers.Select(x => x.ServerID).ToList();

            List<ServerConcurrentCountModel> queryOnlinePlayers = GetServerOnlinePlayers(serverIDs);

            //檢查最大玩家數
            int maxPlayerCount = SettingHelper.MaxServerPlayerCount;

            List<int> sizeAvailableServers = queryOnlinePlayers
                .Where(x => x.TotalCount <= maxPlayerCount)
                .Select(x => x.ServerID)
                .ToList();

            if (sizeAvailableServers.IsEmptyOrNull())
                throw new TerproException(ApiErrorCode.ServerPlayerOutOfRange);

            queryGameServers = queryGameServers
                .Where(x => sizeAvailableServers.Contains(x.ServerID))
                .ToList();

            //如果有多個符合的servers 會在後續的邏輯導到選廳頁面
            bool isMultiple = queryGameServers.Count() > 1;
            bool isOnlyOne = !isMultiple;

            var matchFirst = queryGameServers.Where(x => queryOnlinePlayers.Select(z => z.ServerID).Contains(x.ServerID)).FirstOrDefault();

            var result = new PlayerLaunchGameServerChooseResultModel
            {
                IsRedirectToHall = isMultiple, //如果有限制入場金幣的話 不要自動導去唯一廳
                ServerID = matchFirst.ServerID,
                ServerIPAddress = matchFirst.ClientAddr,
            };


            //2020-08-25 企劃要求的邏輯 如果有限制入場金幣的話 不要自動導去遊戲頁
            if (isOnlyOne)
            {
                bool isRequiredMinAmount = queryGameServers.FirstOrDefault().MinAvailableMoney > 0;

                if (isRequiredMinAmount)
                {
                    result.IsRedirectToHall = true;
                }
            }

            return result;
        }


        /// <summary>
        /// 為了方便寫MSTest 寫成virtual
        /// </summary>
        protected virtual List<ServerConcurrentCountModel> GetServerOnlinePlayers(List<int> serverIDs)
        {
            return PlayerCurrentServerRepository.GetOnlinePlayersCount(serverIDs);
        }
    }
}
