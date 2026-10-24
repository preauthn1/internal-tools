namespace target.Models.ViewModels
{
    /// <summary>
    /// 遊戲進房
    /// </summary>
    public class GameRoomEnterRequestModel
    {

        public int ServerID { get; set; }


        /// <summary>
        /// 是否為捕魚王 (邏輯匯特殊點)
        /// </summary>
        public bool IsFishHunter { get; set; }

        /// <summary>
        /// 用戶IP
        /// </summary>
        public string IPAddress { get; set; }



        [SwaggerIgnore]
        public string AccessToken { get; set; }


    }


    /// <summary>
    /// 遊戲進房
    /// </summary>
    public class GameRoomEnterViewModel
    {

        /// <summary>
        /// 是否斷線重連
        /// </summary>
        public bool IsReConnect { get; set; }


        public long TradeID { get; set; }


        public ServerReconnectViewModel ReConnectInfo { get; set; }


        public PlayerInfoModel PlayerInfo { get; set; }

    }






    public class PlayerInfoModel
    {
        public string PlayerName { get; set; }

        public string DisplayName { get; set; }

        public int PlayerID { get; set; }

        public int PlayerType { get; set; }

        public int FaceID { get; set; }

        public int Sex { get; set; }


        public int WinCount { get; set; }

        public int LoseCount { get; set; }

        public int DrawCount { get; set; }

        public int LastWinLose { get; set; }


    }





    public class ServerReconnectViewModel
    {

        public int ServerID { get; set; }


        public int GameID { get; set; }


        public string ClientAddr { get; set; }


        public string ReConnectMsg { get; set; }


    }







}
