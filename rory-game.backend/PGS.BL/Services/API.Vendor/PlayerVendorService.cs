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

namespace target.BL.Services
{

    public class PlayerVendorService : BaseVendorService, IPlayerVendorService
    {
        public PlayerRepository PlayerRepository { get; set; }

        public IPlayerCreateRefService PlayerCreateRefService { get; set; }


        public PlayerVendorViewModel IsAccountExist(QueryPlayerAccountExistRequestModel filter)
        {
            this.GetMerchantAndCheckIsValid(filter.MerchantID, filter.MerchantSecretKey);

            filter.RuleFor(x => x.PlayerName).Rule_NotEmptyText();
            
            var queryPlayerVendor = PlayerRepository.GetByPlayerNameAndMerchantID(filter.PlayerName , filter.MerchantID);

            return new PlayerVendorViewModel
            {
                IsExist = queryPlayerVendor != null,
                CreateTime = queryPlayerVendor?.CreateTime
            };
        }



        public PlayerCreateResultVenodrApiViewModel CreateNewPlayer(CreatePlayerAccountExistRequestModel model)
        {
            var queryMerchant =  this.GetMerchantAndCheckIsValid(model.MerchantID, model.MerchantSecretKey);

            model.RuleFor(x => x.PlayerName).Rule_NotEmptyText();
            
            var resultPlayerID = PlayerCreateRefService.CreatePlayerForMerchantRequest(new PlayerCreateParameterModel
            {
                Currency = queryMerchant.Currency,
                MerchantID = model.MerchantID,
                PlayerName = model.PlayerName
            });

            return new PlayerCreateResultVenodrApiViewModel
            {

            };
        }



    }
}
