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
using target.DA.Repository;
using target.Utility.Helper;
using target.Utility;

namespace target.BL.Services
{

    public class GameServerSettingService : BaseService, IGameServerSettingService
    {

        public ServerListRepository ServerListRepository { get; set; }

        public ServerMerchantMappingRepository ServerMerchantMappingRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }

        public RobotListRepository RobotListRepository { get; set; }

        public PlayerCurrentServerRepository PlayerCurrentServerRepository { get; set; }

        public ServerRobotConfigRepository ServerRobotConfigRepository { get; set; }


        public RobotListResultModel GetRobotList(RobotListRequestModel model)
        {
            model.RuleFor(x => x.ServerID).Rule_NotZero().Rule_NotMinus();

            var queryServer = ServerListRepository.GetBriefByCache(model.ServerID);

            if (queryServer == null)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            //Gs 定期去重抓需要的機器人數量 或 重新洗牌機器人帳號 (每分鐘一次的頻率是正常的)
            var result = RobotListRepository.RobotLogon(model);
            return result;
        }


        public List<ServerSettingViewModel> GetServerListAndSetting(ServerListAndSettingRequestModel model)
        {
            var queryAllServer = ServerListRepository.GetAll();

            var queryAllRelationSetting = ServerMerchantMappingRepository.GetAll();

            var result = queryAllServer.Select(x => new ServerSettingViewModel
                {
                    ServerID = x.ServerID,
                    ServerName = x.ServerName,
                    FlagLevel = x.FlagLevel,
                    AvailableMerchantIDs = queryAllRelationSetting
                        .Where(z => z.ServerID == x.ServerID)
                        .Select(y => y.MerchantID)
                        .ToList()
                })
                .ToList();

            return result;
        }



        public ServerStartProcessViewModel RunServerStartProcess(ServerStartProcessRequestModel model)
        {
            model.RuleFor(x => x.ServerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.FlagLevel).Rule_NotMinus();
            model.RuleFor(x => x.MinAvailableMoney).Rule_NotMinus();

            var queryServer = ServerListRepository.GetById(model.ServerID);

            if (queryServer == null)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            if (queryServer.GameID != queryServer.GameID)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            //狀態設為delete的 禁止
            if (queryServer.Status == ServerStatusEnum.Delete.ToInteger())
                throw new TerproException(ApiErrorCode.InvalidStatus);

            var queryGame = GameListRepository.GetSingleByCache(model.GameID);

            if (queryGame == null)
                throw new TerproException(ApiErrorCode.NotFoundGame);


            queryServer.ClientAddr = model.ClientAddr;
            queryServer.HttpAddr = model.HttpAddr;
            queryServer.NotifyAddr = model.NotifyAddr;
            queryServer.MinAvailableMoney = model.MinAvailableMoney;
            queryServer.FlagLevel = model.FlagLevel;
            queryServer.UpdateBy = "api_server_start";
            queryServer.UpdateTime = DateTimeHelper.GetCurrentTime();

            //如果server狀態原本是 暫時關閉，改為開啟
            if (queryServer.Status == ServerStatusEnum.ServerDown.ToInteger())
                queryServer.Status = ServerStatusEnum.Enable.ToInteger();


            //step 1. 更新服務器
            ServerListRepository.Update(queryServer);

            //step 2. 解锁玩家房间锁
            PlayerCurrentServerRepository.UnlockAllPlayerByServerID(queryServer.ServerID);

            //step 3. 机器人解锁
            RobotListRepository.RemoveServerAssignRobots(queryServer.ServerID);


            //參考 PrGs_ServerStart (CenterDB)
            //--更新金币源
            //--清除房间网络消息记录

            return new ServerStartProcessViewModel
            {
                GameID = queryGame.GameID,
                GameStatus = queryGame.Status
            };
        }


        public void RunServerCloseProcess(ServerCloseProcessRequestModel model)
        {
            model.RuleFor(x => x.ServerID).Rule_NotZeroNotMinus();

            var query = ServerListRepository.GetById(model.ServerID);

            if (query == null)
                throw new TerproException(ApiErrorCode.NotFoundServer);
            
            query.Status = ServerStatusEnum.ServerDown.ToInteger();
            query.UpdateBy = "api_server_close";
            query.UpdateTime = DateTimeHelper.GetCurrentTime();

            //step 1. 更新服務器
            ServerListRepository.Update(query);

            //step 2. 解锁玩家房间锁
            PlayerCurrentServerRepository.UnlockAllPlayerByServerID(query.ServerID);

            //step 3. 机器人解锁
            RobotListRepository.RemoveServerAssignRobots(query.ServerID);
        }



        /// <summary>
        /// 這機制(應該)是用來檢查伺服器健康狀況的
        /// </summary>
        /// <param name="model"></param>
        /// <returns></returns>
        public void ServerReRegister(SystemReregisterRequestModel model)
        {
            ServerListRepository.UpdateLastTime(model.ServerID);
        }


    }
}
