using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IGameServerSettingService
    {
        RobotListResultModel GetRobotList(RobotListRequestModel model);

        List<ServerSettingViewModel> GetServerListAndSetting(ServerListAndSettingRequestModel model);

        ServerStartProcessViewModel RunServerStartProcess(ServerStartProcessRequestModel model);

        void RunServerCloseProcess(ServerCloseProcessRequestModel model);

        void ServerReRegister(SystemReregisterRequestModel model);

    }
}
