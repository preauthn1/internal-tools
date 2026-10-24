using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class PlayerInfoViewModel
    {

        public int PlayerID { get; set; }


        public string PlayerName { get; set; }


        public int Status { get; set; }


        public int MerchantID { get; set; }


        public string MerchantName { get; set; }
        

        public decimal WalletAmount { get; set; }


        public int WalletMode { get; set; }


        public DateTime CreateTime { get; set; }


        public string CreateIPAddress { get; set; }


        public string LastIPAddress { get; set; }


        public DateTime? LastLoginTime { get; set; }

    }




    public class PlayerInfoSummaryRequestModel
    {
        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

    }


    public class PlayerInfoSummaryViewModel
    {
        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public int PlayerType { get; set; }

        public string PlayerTypeText { get; set; }

        public string DisplayName { get; set; }

        public int Status { get; set; }

        public int MerchantID { get; set; }

        public string MerchantName { get; set; }

        public int WalletMode { get; set; }

        public decimal WalletAmount { get; set; }

        public DateTime CreateTime { get; set; }

        public string CreateIPAddress { get; set; }

        public DateTime LastLoginTime { get; set; }

        public string PlayerCurrentGameGame { get; set; }

        public int PlayerCurrentServerID { get; set; }

        public bool IsTestAccountOrMerchant { get; set; }


        public PlayerTotalSummaryViewModel TotalSummary { get; set; }

        public List<PlayerGameSummaryViewModel> PlayerGameSummaries { get; set; }

        public PlayerGameSummaryViewModel PlayerGameSummaryRowTotal { get; set; }

    }


    public class PlayerTotalSummaryViewModel 
    {

        public int TotalLoginTimes { get; set; }

    }



    public class PlayerGameSummaryViewModel 
    {
        public string GameName { get; set; }
        public decimal BonusAmount { get; set; }
        public int PlayerID { get; set; }
        public int GameID { get; set; }
        public int WinCount { get; set; }
        public int LoseCount { get; set; }
        public int DrawCount { get; set; }
        public decimal Score { get; set; }
        public int LastWinLose { get; set; }
        public long GameTimeSeconds { get; set; }




        public decimal BetAmount { get; set; }
        public decimal ValidAmount { get; set; }
        public decimal PayoutAmount { get; set; }
        public decimal RTP { get; set; }








    }




}
