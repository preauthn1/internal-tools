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

    public class GameRoomService : BaseGameService, IGameRoomService
    {

        public PlayerWalletRepository PlayerWalletRepository { get; set; }

        public PlayerGameSummaryRepository PlayerGameSummaryRepository { get; set; }

        public ServerListRepository ServerListRepository { get; set; }

        public PlayerCurrentServerRepository PlayerCurrentServerRepository { get; set; }

        public ForbiddenIPAddressRepository ForbiddenIPAddressRepository { get; set; }

        public ServerMerchantMappingRepository ServerMerchantMappingRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }

        public MerchantGameSettingRepository MerchantGameSettingRepository { get; set; }

        public SystemConfigSettingRepository SystemConfigSettingRepository { get; set; }

        public IPlayerWalletRefService PlayerWalletRefService { get; set; }



        public GameRoomEnterViewModel GameRoomEnter(int playerID , GameRoomEnterRequestModel model)
        {
            model.RuleFor(x => x.ServerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.IPAddress).Rule_NotEmptyText();

            var queryServer = ServerListRepository.GetById(model.ServerID);

            //檢查server、IP
            if (queryServer == null)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            if (queryServer.Status != NormalActiveStatus.Enable.ToInteger())
                throw new TerproException(ApiErrorCode.InvalidStatus_DisableNow);

            var queryDenyIP = ForbiddenIPAddressRepository.GetByIPAddress(model.IPAddress);

            if (queryDenyIP != null)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            
            var queryPlayer = PlayerRepository.GetSingleByCache(playerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            //內部測試用商戶 可以跳過系統維護的檢查
            if (!MerchantCustomizedSettingRepository.IsInternalTestMerchant(queryPlayer.MerchantID))
            {
                //檢查系統維護
                SystemConfigSettingRepository.CheckSystemMaintainAndThrowException();
            }

            var queryCurrentRoom = PlayerCurrentServerRepository.GetById(playerID);

            //檢查是否已在其他遊戲房，若已進入其他房間 (彈出錯誤訊息)
            if (queryCurrentRoom != null && queryCurrentRoom.ServerID != model.ServerID)            {

                CheckRoomLockErrorOrSkip(queryPlayer.MerchantID , queryCurrentRoom.ServerID);
            }

            bool isReconn = queryCurrentRoom == null ? false : true;

            string accessToken = [REDACTED];

            var parseToken = [REDACTED];

            var queryWallet = PlayerWalletRefService.GetWalletAmountByMerchantWalletMode(playerID , parseToken?.SessionTokenPGS);

            //檢查入場金幣
            if(queryWallet.Amount < queryServer.MinAvailableMoney)
            {
                string replaceMsg = MultiLanguageTextHelper.EnterGameRoom_AmountNotEnough(queryServer.MinAvailableMoney , queryWallet.Amount);
                throw new TerproException(ApiErrorCode.Forbidden_EnterGameError_WalletAmountNotEnough , replaceMsg);
            }

            //檢查某些商戶開關設置
            bool isGameEnable = MerchantGameSettingRepository.IsMerchantGameEnableInSystem(queryPlayer.MerchantID , queryServer.GameID);

            if (!isGameEnable)
                throw new TerproException(ApiErrorCode.InvalidStatus_GameStatus);

            bool isServerMerchantEnable = ServerMerchantMappingRepository.IsServerMerchantEnable(queryServer.ServerID, queryPlayer.MerchantID);
            
            if (!isServerMerchantEnable)
                throw new TerproException(ApiErrorCode.NotFoundAvailableServer_MerchantApiSelected);

            int HDType = 1;
            long tradeID = RandomHelper.GetTradeNo();

            //呼叫SP，相關紀錄的寫入邏輯
            PlayerCurrentServerRepository.Player_GameRoom_In(playerID , queryServer.ServerID , queryServer.GameID , tradeID, HDType , queryWallet.Amount , model.IPAddress);
            
            var queryPlayerGameSummary = PlayerGameSummaryRepository.GetByPlayerAndGameID(playerID , queryServer.GameID);

            //用途不明 timeSeq ? 
            return new GameRoomEnterViewModel
            {
                  IsReConnect = isReconn,
                  TradeID = tradeID,
                  PlayerInfo = new PlayerInfoModel
                  {                      
                      PlayerID = queryPlayer.PlayerID,
                      PlayerName = queryPlayer.PlayerName,
                      DisplayName = queryPlayer.DisplayName ?? queryPlayer.DisplayName,
                      PlayerType = queryPlayer.PlayerType,
                      Sex = queryPlayer.Sex,
                      DrawCount = queryPlayerGameSummary.DrawCount,
                      WinCount = queryPlayerGameSummary.WinCount,
                      LoseCount = queryPlayerGameSummary.LoseCount,
                      LastWinLose = queryPlayerGameSummary.LastWinLose,
                      FaceID = queryPlayer.FaceID,
                  },
            };

        }


        public GameRoomLeaveViewModel GameRoomLeave(GameRoomLeaveRequestModel model)
        {            
            model.RuleFor(x => x.PlayerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.ServerID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.LoginTime).Rule_NotEmptyDateTime();
            model.RuleFor(x => x.LogoutTime).Rule_NotEmptyDateTime();

            var queryPlayer = PlayerRepository.GetById(model.PlayerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            //機器人的話 可以不帶traderID
            if(queryPlayer.PlayerType == (int)PlayerTypeEnum.Normal)
            {
                model.RuleFor(x => x.TradeID).Rule_NotZeroNotMinus().Rule_IsValidTradeID();
            }

            var queryServer = ServerListRepository.GetById(model.ServerID);

            if (queryServer == null)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            //呼叫SP，相關紀錄的寫入邏輯
            var spPara = model.ReturnInjecteObject<PlayerRoomOutSPParameterModel>();
            spPara.ValidAmount = model.FlowAmount; //flowAmount(舊欄位) 視為 validAmount
            spPara.GameTimeSecond =(long)(model.LogoutTime - model.LoginTime).TotalSeconds;


            spPara.BetAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, spPara.BetAmount);
            spPara.PayoutAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, spPara.PayoutAmount);
            spPara.ValidAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, spPara.ValidAmount);
            spPara.UpdateAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, spPara.UpdateAmount);
            spPara.TaxAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, spPara.TaxAmount);
            spPara.FlowAmount = ConvertGamePointToMoney(queryPlayer.MerchantID, spPara.FlowAmount);

            PlayerCurrentServerRepository.Player_GameRoom_Out(spPara);

            return new GameRoomLeaveViewModel
            {

            };
        }


        private void CheckRoomLockErrorOrSkip(int merchantID , int previousServerID)
        {
            var queryPreviousServer = ServerListRepository.GetBriefByCache(previousServerID);

            bool isSlotsGame = queryPreviousServer.GameType == GameTypeEnum.Slot.ToInteger();

            if(isSlotsGame)
            {
                var querySetting = MerchantCustomizedSettingRepository.GetMerchantSettingByCache(merchantID, MerchantCustomizedSettingKeyEnum.GameSetting_IsSlotsGameRoomLock);

                //如果商戶配置有設定 可以跳過檢查
                if (querySetting != null && !querySetting.SettingValueOfBoolean)
                {
                    return;
                }
            }

            string gameName = DbMultiLanguageTextHelper.GetGameDisplayName(queryPreviousServer.GameID);
            string msg = MultiLanguageTextHelper.GameServerReconnectMsg(gameName);

            //卡房錯誤
            throw new TerproException(ApiErrorCode.Forbidden_DuplicateEnterGameServer, msg);
        }


    }
}
