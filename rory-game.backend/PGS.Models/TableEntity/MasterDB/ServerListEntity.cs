using target.Models.CoreLibs;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{
    
    [IsNotAutoKeyTable]
    public class ServerListEntity
    {
        [Key]
        public int ServerID { get; set; }
        public int GameID { get; set; }
        public string ServerName { get; set; }
        public int FlagLevel { get; set; }
        public string Currency { get; set; }

        /// <summary>
        /// 對應 ServerStatusEnum
        /// </summary>
        public int Status { get; set; }
        public string ClientAddr { get; set; }
        public string NotifyAddr { get; set; }
        public string HttpAddr { get; set; }
        public DateTime CreateTime { get; set; }
        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }
        public decimal MinAvailableMoney { get; set; }
        public int RoomFlagLevel { get; set; }



    }
}
