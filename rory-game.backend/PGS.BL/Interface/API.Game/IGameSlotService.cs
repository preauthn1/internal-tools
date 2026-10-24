using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{
    public interface IGameSlotService
    {
        #region SlotGameLog
        PagedList<GameSlotHistoryLogViewModel> GetSlotGameLog(GameSlotLogRequestModel model);
        void AddSlotGameLog(GameSlotHistoryLogRequestModel model);
        #endregion

        #region jackpot
        GameJackpotViewModel UpdateJackpot(int gameID, GameSlotJackpotUpdateRequestModel model);
        #endregion

        #region Slot Probability Record
        GameSlotProbabilitySettingViewModel GetSlotProbabilitySetting(GameSlotRequestModel model);
        void UpdateSlotProbRecord(GameSlotProbRequestModel model);
        void DeleteSlotProbRecord(GameSlotRequestModel model);
        #endregion

        #region Slot Bonus Record
        PlayerSlotBonusRecordViewModel GetSlotBonusRecord(GameSlotRequestModel model);
        void UpdateSlotBonusRecord(GameSlotBonusRecordRequestModel model);
        void ResetSlotBonusRecord(GameSlotBonusRequestModel model);
        #endregion

        #region Slot Free Spin Record
        PlayerSlotFreeSpinRecordViewModel GetSlotFreeSpinRecord(GameSlotRequestModel model);
        void UpdateSlotFreeSpinRecord(GameSlotFreeSpinRequestModel model);
        #endregion

        #region Slot Lightning Record
        PlayerSlotLightingRecordViewModel GetSlotLightningRecord(GameSlotRequestModel model);
        void UpdateSlotLightningRecord(GameSlotLightningRequestModel model);
        #endregion

        #region Deferred Award Record
        GameSlotDeferredAwardViewModel GetSlotDeferredAward(GameSlotDeferredAwardRequestModel model);
        void AddSlotDeferredAward(GameSlotDeferredAwardUpdateRequestModel model);
        GameSlotDeferredAwardItemViewModel DeleteSlotDeferredAward(GameSlotDeferredAwardRequestModel model);
        #endregion

        #region Player Assigned Mode
        GameSlotPlayerAssignedModeViewModel GetSlotAssignedMode(GameSlotRequestModel model);
        void AddSlotAssignedMode(GameSlotPlayerAssignedRequestModel model);
        void DeleteSlotAssignedMode(GameSlotRequestModel model);
        #endregion

    }
}
