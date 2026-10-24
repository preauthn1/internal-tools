using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class GamePlayBetRecordRequestModel
    {

        /// <summary>
        /// 每筆投注紀錄
        /// </summary>
        public List<BetRecordItemRequestModel> BetRecords { get; set; }

    }

    public class BetRecordItemRequestModel
    {
        public int PlayerID { get; set; }

        public decimal Amount { get; set; }

        public int Status { get; set; }
    }


    public class GamePlayBetRecordViewModel
    {



    }

}
