using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class PerGameSlotHistoryDetailInfoViewModel
    {
        public long PerGameLogID { get; set; }

        public long PerGameDetailLogID { get; set; }

        public long SlotGameHistoryLogID { get; set; }

        public string PerGameTransactionNo { get; set; }

        public int GameID { get; set; }

        public string GameName { get; set; }

        public int PlayerID { get; set; }

        public string PlayerName { get; set; }

        public decimal BetAmount { get; set; }

        public decimal TaxAmount { get; set; }

        public decimal UpdateAmount { get; set; }

        public decimal BeforeAmount { get; set; }

        public decimal AfterAmount { get; set; }


        public DateTime GameTime { get; set; }


        public bool IsSingleWalletMode { get; set; }


        /// <summary>
        /// 有觸發特殊盤面
        /// </summary>
        public bool isMultipleSlotsResult { get; set; }


        public List<List<SlotHistoryResultIconViewModel>> LineResult { get; set; }


        public List<SlotHistoryResultIconViewModel> DistinctIcons { get; set; }

    }

    public class SlotHistoryResultIconViewModel
    {
        public string Code { get; set; }

        public string ImageUrl { get; set; }
    }




    public class MatchMultipleHistoryResultIconViewModel
    {
        public DateTime GameTime { get; set; }

        public long SlotGameHistoryLogID { get; set; }

        /// <summary>
        /// 盤面結果
        /// </summary>
        public List<List<string>> SlotsResult { get; set; }
    }


}
