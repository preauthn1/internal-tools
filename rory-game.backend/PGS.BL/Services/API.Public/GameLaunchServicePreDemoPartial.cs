using System;
using System.Collections.Generic;
using target.BL.Interface;
using target.BL.Services;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.Models.Models;
using System.Linq;
using target.DA.Repository;
using target.Models;
using target.Utility;
using target.Models.Enums;
using target.Models.ViewModels.Base;

namespace target.BL.Services
{

    public partial class GameLaunchService : BaseVendorService, IGameLaunchService
    {
       

        public GameServerViewModel GetGameServers(GameServerRequestModel model)
        {
            model.RuleFor(x => x.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.Ts).Rule_NotZeroNotMinus();

            string token = [REDACTED];

            var parseModel = GetTokenAndCheckValid(token);

            var queryGame = GameListRepository.GetSingleByCacheAndCheckValid(model.GameID);

            string tokenHash = token.MD5Hash();

            //檢查DB是否存在
            PlayerJwtCheckService.CheckIsTokenExist(tokenHash);


            List<ServerListEntity> queryGameServers = ServerListRepository.GetByFilter(new QueryServerListRequestViewModel
            {
                GameIDs = model.GameID.ToListCollection(),
                MerchantIDs = parseModel.MerchantID.ToListCollection(),
                Status = ServerStatusEnum.Enable.ToInteger()
            });
            
            string gameDomain = GetGameLinkDomain(parseModel.MerchantID);

            decimal balanceAmount = PlayerWalletRefService.GetWalletAmountByMerchantWalletMode(parseModel.LoginUserID , parseModel.SessionTokenPGS).Amount;
            decimal balanceGamePoint = ConvertMoneyToGamePoint(parseModel.MerchantID , balanceAmount);


            //部分設置是存在 CustomizedSetting 這邊
            List<MerchantCustomizedSettingModel> customizedSettings = GetMerchantCustomizedSettings(parseModel.MerchantID);
          

            return new GameServerViewModel
            {
                PlayerID = parseModel.LoginUserID,
                PlayerName = parseModel.LoginName,
                BalanceAmount = balanceGamePoint,
                GameID = queryGame.GameID,
                GameName = DbMultiLanguageTextHelper.GetGameDisplayName(queryGame.GameID),
                Servers = queryGameServers.Select(x=>
                    {
                        var gameLaunchResult = new GameLaunchViewModel
                        {
                            IsMultipleServers = false,
                            Info = new GameLaunchRequiredInfoViewModel
                            {
                                ServerID = x.ServerID,
                                ServerConnection = x.ClientAddr,
                                GameID = x.GameID,
                                Lang = parseModel.Lang,
                                PlayerID = parseModel.LoginUserID,
                                PlayerName = parseModel.LoginName,
                                GameCode = queryGame.GameCode?.ToLower(),
                                Access_Token = [REDACTED],
                            }                           
                        };

                        //客製customized設置
                        AssignCustomizedSettingsValueToInfo(customizedSettings, example.invalid);

                        //外部三方的遊戲 還有額外參數要設置
                        SetSlot88Information(gameLaunchResult, queryGame);

                        var convert = new GameServerItemViewModel
                        {
                            FlagID = x.FlagLevel,
                            RoomLevel = x.RoomFlagLevel,
                            RoomLevelName = MultiLanguageTextHelper.GetEnumsCultureText((GameRoomFlagLevelEnum)x.RoomFlagLevel),
                            MinAvailableMoney = x.MinAvailableMoney,
                            Url = GenerateRedirectUrl(gameDomain , gameLaunchResult)
                        };

                        string replaceText = $"({convert.RoomLevelName.GetInnerTextBetweenTwoWord("(", ")")})";
                        convert.RoomLevelName = convert.RoomLevelName.Replace(replaceText, "");

                        return convert;
                    })
                    .ToList()
            };
        }


        private PlayerTokenJwtModel GetTokenAndCheckValid(string jwtToken)
        {
            var parseModel = MyJwtHelper.GetPlayerModelFromJwtString(jwtToken);

            if (parseModel == null)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_NotValid);

            long nowUTC = DateTimeHelper.GetCurrentTimeTimestamp();

            if (nowUTC >= parseModel.TokenExpiredTime)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);

            return parseModel;
        }


        private PlayerTokenJwtModel GetLobbyTokenAndCheckValid(string guidToken)
        {
            var parseModel = PlayerJwtCheckService.ParsingLobbyToken(guidToken);
            return parseModel;
        }


        private decimal GetMerchantConvertRate_PGSGame(int merchantID)
        {
            decimal rate = MerchantCustomizedSettingRepository.GetMerchantGamePointConvertRateByCache(merchantID, MerchantCustomizedSettingKeyEnum.GameVendor_ConvertRate_PGS);
            return rate;
        }


        private decimal ConvertMoneyToGamePoint(int merchantID, decimal money)
        {
            decimal convertRate = GetMerchantConvertRate_PGSGame(merchantID);
            return money * convertRate;
        }


    }
}
