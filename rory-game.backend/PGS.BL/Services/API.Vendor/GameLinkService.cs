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

namespace target.BL.Services
{

    public class GameLinkService : BaseVendorService, IGameLinkService
    {

        public IGameLaunchService GameLaunchService { get; set; }


        public GameLinkGenerateResultViewModel GenerateSingleWalletLaunchLink(GameLaunchVendorApiSingleWalletModeRequestModel model)
        {
            model.RuleFor(z => z.PlayerUniqueID);

            var convert = model.ReturnInjecteObject<GameLaunchRequestModel>();
            convert.PlayerName = model.PlayerUniqueID;

            var result = GenerateGameLink(convert);
            return result;
        }


        public GameLinkGenerateResultViewModel GenerateTransferModeLaunchLink(GameLaunchVendorApiTransferModeRequestModel model)
        {
            model.RuleFor(x => x.PlayerName).Rule_NotEmptyText();

            var convert = model.ReturnInjecteObject<GameLaunchRequestModel>();

            var result = GenerateGameLink(convert);
            return result;
        }


        private GameLinkGenerateResultViewModel GenerateGameLink(GameLaunchRequestModel model)
        {
            this.GetMerchantAndCheckIsValid(model.MerchantID, model.MerchantSecretKey);

            string url = GameLaunchService.LaunchGameAndGetGamePageUrl(model);

            return new GameLinkGenerateResultViewModel
            {
                Link = url
            };
        }





    }
}
