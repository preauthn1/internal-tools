using System;
using System.Collections.Generic;
using game.Entity.Game;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IPlayerInfoService
    {
        PagedList<PlayerInfoViewModel> GetByPagingFilter(int userID, QueryPlayerApiRequestParameterModel filter);

        int GetPlayerIdByPlayerNameOrID(int userID, string playerNameOrID);

        PlayerInfoSummaryViewModel GetPlayerSummaryInfo(int playerID);

        PlayerManualDepositResultViewModel InternalTestAccountWalletDepositAmount(int userID, InternalTestAccountPlusWalletRequestModel model);


        void CreateInernalTestAccount(int userID, CreateInternalTestAccountRequestModel model);


        void UpdatePlayerStatus(int userID, PlayerUpdateStatusRequestModel model);


        void UpdateMultiplePlayerStatus(int userID, PlayerMultipleUpdateStatusRequestModel model);


        void RemovePlayerRoomLock(int userID, int playerID);

        /// <summary>
        /// 手動幫玩家補錢 (僅針對轉帳模式)
        /// </summary>
        PlayerManualDepositResultViewModel ManualDepositPlayerWallet(int userID, PlayerManualDepositRequestModel model);
    }
}
