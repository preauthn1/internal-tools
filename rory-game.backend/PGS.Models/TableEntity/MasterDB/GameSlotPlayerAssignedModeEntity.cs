using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    public class GameSlotPlayerAssignedModeEntity
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
        /// 建立時間
        /// </summary>
        public DateTime CreateTime { get; set; }

        /// <summary>
        /// 特殊模式 (0: 預設, 1: 放分, 2: 殺分)
        /// </summary>
        public int Mode { get; set; }

        /// <summary>
        /// 過期時間
        /// </summary>
        public DateTime ExpireTime { get; set; }
    }
}
