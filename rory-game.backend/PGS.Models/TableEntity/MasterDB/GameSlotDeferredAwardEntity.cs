using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class GameSlotDeferredAwardEntity
    {
        [Key]
        public int ID { get; set; }

        /// <summary>
        /// 玩家ID
        /// </summary>
        public int PlayerID { get; set; }

        /// <summary>
        /// 遊戲ID
        /// </summary>
        public int GameID { get; set; }

        /// <summary>
        /// 房間類型ID
        /// </summary>
        public int FlagID { get; set; }

        /// <summary>
        /// 獎勵類型
        /// </summary>
        public string AwardType { get; set; }

        /// <summary>
        /// 建立時間
        /// </summary>
        public DateTime CreateTime { get; set; }

        /// <summary>
        /// 規則ID
        /// </summary>
        public int OriginID { get; set; }

        /// <summary>
        /// 遊戲資料(自行做字串處理)
        /// </summary>
        public string Data { get; set; }

        /// <summary>
        /// 此遊戲獲得的來源 Pattern
        /// </summary>
        public int PatternID { get; set; }
    }
}
