using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    /// <summary>
    /// 取得彩池金額 (請求參數)
    /// </summary>
    public class GameSlotJackpotRequestModel
    {
        public bool IsLink { get; set; }

        public int FlagID { get; set; }

        public int LevelID { get; set; }

        public int TableNo { get; set; }
    }


    /// <summary>
    /// 取得彩池金額 (回應參數)
    /// </summary>
    public class GameJackpotViewModel
    {
        public List<GameJackpotItemViewModel> Jackpots { get; set; }
    }


    public class GameJackpotItemViewModel
    {
        /// <summary>
        /// 某些遊戲 彩池有分成 大獎、中獎、小獎
        /// </summary>
        public int GameID { get; set; }
        public int TableNo { get; set; }
        public int FlagID { get; set; }
        public int LevelID { get; set; }
        public decimal Contribution { get; set; }
        public decimal TotalContribution { get; set; }
        public decimal TotalAward { get; set; }
        public bool Type { get; set; }
    }
}
