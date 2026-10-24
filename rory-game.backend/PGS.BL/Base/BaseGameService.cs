using target.Core.Repository;
using target.Models.Paging;
using target.Utility.IOC;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Reflection;
using target.Utility.Extensions;
using target.Utility.Helper;
using target.DA.Repository;
using target.Models.Models;
using target.Models.Enums;

namespace target.BL.Services
{
    public class BaseGameService : BaseService
    {      
        public PlayerRepository PlayerRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }


        protected PlayerMerchantInfoModel GetValidPlayerAndMerchant(int playerID)
        {
            var queryPlayer = PlayerRepository.GetSingleByCache(playerID);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.NotFoundPlayer);

            if (queryPlayer.Status != NormalActiveStatus.Enable.ToInteger())
                throw new TerproException(ApiErrorCode.InvalidStatus_AccountStatus);

            var queryMerchant = MerchantRepository.GetSingleByCache(queryPlayer.MerchantID);

            if (queryMerchant == null)
                throw new TerproException(ApiErrorCode.NotFoundMerchant);

            if (queryMerchant.Status != NormalActiveStatus.Enable.ToInteger())
                throw new TerproException(ApiErrorCode.InvalidStatus_DisableNow);


            return new PlayerMerchantInfoModel
            {
                Currency = queryMerchant.Currency,
                MerchantID = queryMerchant.MerchantID,
                PlayerID = queryPlayer.PlayerID,
                PlayerName = queryPlayer.PlayerName,
                WalletMode = (MerchantWalletModeEnum)queryMerchant.WalletMode
            };

        }

        protected decimal GetMerchantConvertRate_PGSGame(int merchantID)
        {           
            var rep = new MerchantCustomizedSettingRepository();
            decimal rate = rep.GetMerchantGamePointConvertRateByCache(merchantID, MerchantCustomizedSettingKeyEnum.GameVendor_ConvertRate_PGS);
            return rate;
        }

        protected decimal ConvertMoneyToGamePoint(int merchantID, decimal money)
        {
            decimal convertRate = GetMerchantConvertRate_PGSGame(merchantID);
            return money * convertRate;
        }

        protected decimal ConvertGamePointToMoney(int merchantID, decimal point)
        {
            decimal convertRate = GetMerchantConvertRate_PGSGame(merchantID);
            return point / convertRate;
        }
    }
}
