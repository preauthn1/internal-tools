using System.Collections.Generic;

namespace target.Models.ViewModels.Base
{
    public class AGGVendorApiGetDetailResponseModel
    {
        public List<AGGRecord> records { get; set; }
    }

    public class AGGJackpotSetting
    {
        public List<int> minBetList { get; set; }
        public List<int> maxBetList { get; set; }
    }

    public class AGGOption
    {
        public AGGJackpotSetting jackpotSetting { get; set; }
        public int winLimit { get; set; }
        public int rtp { get; set; }
    }

    public class AGGJackpotSetting2
    {
        public List<int> maxBetList { get; set; }
        public List<int> minBetList { get; set; }
    }

    public class AGGInit
    {
        public AGGJackpotSetting2 jackpotSetting { get; set; }
        public string strip { get; set; }
        public List<int> scatterMultiple { get; set; }
        public List<int> scatterMultipleWeight { get; set; }
        public string payTable { get; set; }
        public List<int> jackpotBetMultiple { get; set; }
    }

    public class AGGMeta
    {
        public AGGOption option { get; set; }
        public AGGInit init { get; set; }
        public int roundId { get; set; }
    }

    public class AGGWager
    {
        public string wagerId { get; set; }
        public int roundId { get; set; }
    }

    public class AGGDetail
    {
        public int bet { get; set; }
        public AGGWager wager { get; set; }
        public int clientBetUnit { get; set; }
        public int state { get; set; }
        public int betUnit { get; set; }
        public int clientBet { get; set; }
        public List<int> reelsScatter { get; set; }
        public int nextState { get; set; }
        public int win { get; set; }
        public List<int> jackpot { get; set; }
        public List<int> reels { get; set; }
    }

    public class AGGCurrency
    {
        public string code { get; set; }
        public bool isIso4217 { get; set; }
    }

    public class AGGBalance
    {
        public AGGCurrency currency { get; set; }
        public string amount { get; set; }
    }

    public class AGGRecord
    {
        public int recovery { get; set; }
        public string eventType { get; set; }
        public int roundId { get; set; }
        public string gameSessionId { get; set; }
        public bool? isLast { get; set; }
        public AGGMeta meta { get; set; }
        public bool? offlineSpin { get; set; }
        public AGGDetail detail { get; set; }
        public int? gameWin { get; set; }
        public int? bet { get; set; }
        public int? lose { get; set; }
        public string feature { get; set; }
        public int? betUnit { get; set; }
        public List<AGGBalance> balances { get; set; }
        public int? mid { get; set; }
        public int? jackpotWin { get; set; }
    }


}
