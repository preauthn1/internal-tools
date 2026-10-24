using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class RobotListRequestModel
    {

        public int ServerID { get; set; }


        public int CurrentRobotCount { get; set; }


        public decimal MinMoneyAmount { get; set; }


        public decimal MaxMoneyAmount { get; set; }

    }


    public class RobotListResultModel
    {

        public int LogoutRobotCount { get; set; }

        public List<RobotListItemResultModel> Robots { get; set; }

    }

    /// <summary>
    /// SlotBonusRecord細項
    /// </summary>
    public class RobotListItemResultModel
    {
        /// <summary>
        /// 玩家ID
        /// </summary>
        public int PlayerID { get; set; }

        /// <summary>
        /// 玩家類型
        /// </summary>
        public int PlayerType { get; set; }

        /// <summary>
        /// 玩家暱稱
        /// </summary>
        public string NickName { get; set; }

        /// <summary>
        /// 頭像ID
        /// </summary>
        public int FaceID { get; set; }

        /// <summary>
        /// 性別
        /// </summary>
        public int Sex { get; set; }

        /// <summary>
        /// 玩家金額
        /// </summary>
        public decimal CashAmount { get; set; }

        /// <summary>
        /// WinCount
        /// </summary>
        public int WinCount { get; set; }

        /// <summary>
        /// LoseCount
        /// </summary>
        public int LoseCount { get; set; }

        /// <summary>
        /// DrawCount
        /// </summary>
        public int DrawCount { get; set; }

        /// <summary>
        /// 玩家等級
        /// </summary>
        public int LevelID { get; set; }
    }
}
