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
using target.Utility.Helper.Dependency;

namespace target.BL.Services
{

    public partial class GameLaunchService : BaseVendorService, IGameLaunchService
    {

        public LobbyGenerateLinkResultModel GenerateGameLinkForLobby(LobbyApiGenerateLinkParameterModel model)
        {
            model.RuleFor(z => z.GameID).Rule_NotZeroNotMinus();
            model.RuleFor(z => z.Ts).Rule_NotZeroNotMinus();
            model.RuleFor(z => z.PlayerName).Rule_NotEmptyText();
            model.RuleFor(z => z.SignKey).Rule_NotEmptyText();

            SignKeyRuleHelper.CheckIsValidPublicApiLobbySignkey(model.Ts, model.PlayerName, model.SignKey);

            string token = [REDACTED];

            var parseModel = GetLobbyTokenAndCheckValid(token);

            if (parseModel.LoginName != model.PlayerName)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_NotValid);

            var linkUrl = LaunchGameAndGetGamePageUrl(new GameLaunchRequestModel
            {
                GameID = model.GameID,
                MerchantID = parseModel.MerchantID,
                PlayerName = parseModel.LoginName,
                Lang = parseModel.Lang,
                SessionToken = [REDACTED],
                IsAllowDuplicateSessionToken = [REDACTED], //使用同一組sessionToken
            });

            return new LobbyGenerateLinkResultModel
            {
                Url = linkUrl
            };
        }


    }
}
