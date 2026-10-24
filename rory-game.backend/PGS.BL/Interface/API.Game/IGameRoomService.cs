using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;

namespace target.BL.Interface
{

    public interface IGameRoomService
    {

        GameRoomEnterViewModel GameRoomEnter(int playerID, GameRoomEnterRequestModel model);
        
        GameRoomLeaveViewModel GameRoomLeave(GameRoomLeaveRequestModel model);

    }
}
