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

    public class PlayerCreateRefService : BaseService, IPlayerCreateRefService
    {
        public PlayerRepository PlayerRepository { get; set; }


        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }


        public int CreatePlayerOrReturnPlayerID(PlayerCreateParameterModel model)
        {
            var queryPlayer = PlayerRepository.GetByPlayerNameAndMerchantID(model.PlayerName, model.MerchantID);

            //資料已存在 回應playerID
            if (queryPlayer != null)
                return queryPlayer.PlayerID;

            //資料不存在 呼叫 create
            var result = CreatePlayer(model);
            return result.PlayerID;
        }



        public int CreatePlayerForMerchantRequest(PlayerCreateParameterModel model)
        {
            var queryPlayer = PlayerRepository.GetByPlayerNameAndMerchantID(model.PlayerName, model.MerchantID);

            //資料已存在 回應playerID
            if (queryPlayer != null)
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_DuplicateDataExist);

            var result = CreatePlayer(model);
            return result.PlayerID;
        }



        public PlayerCreateResultModel CreatePlayer(PlayerCreateParameterModel model)
        {
            DateTime now = DateTimeHelper.GetCurrentTime();

            var newPlayer = new PlayerEntity
            {
                PlayerName = model.PlayerName,
                DisplayName = model.PlayerName,
                MerchantID = model.MerchantID,
                CreateTime = now,
                UpdateTime = now,
                PlayerType = PlayerTypeEnum.Normal.ToInteger(),
                Currency = model.Currency,
                Status = NormalActiveStatus.Enable.ToInteger(),
                LastLoginToken = "",
                CreateIPAddress = IPHelper.GetIP(),
                FaceID = RandomHelper.GetRandIntegerNumbers(1, 16), //隨機頭像
            };

            decimal defaultAmount = 0;


            //查出是否有設置 預設金幣(基本上是for Dev、QA環境用的)
            var queryDefaultBalanceSetting = MerchantCustomizedSettingRepository.GetByMerchantSettingKey(model.MerchantID, MerchantCustomizedSettingKeyEnum.DefaultWalletBalance);

            if (queryDefaultBalanceSetting != null)
            {
                defaultAmount = queryDefaultBalanceSetting.SettingValueOfDecimal;
            }
            
            var result = PlayerRepository.CreatePlayerAndWallet(newPlayer , defaultAmount);
            return new PlayerCreateResultModel
            {
                PlayerID = result.ResultId,
                DefaultBalance = defaultAmount
            };
        }


    }
}
