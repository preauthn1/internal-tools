using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class ServerListViewModel
    {

        public int ServerID { get; set; }
        public int GameID { get; set; }
        public string ServerName { get; set; }
        public string GameName { get; set; }
        public string Currency { get; set; }
        public int FlagLevel { get; set; }
        public int Status { get; set; }
        public string ClientAddr { get; set; }
        public string NotifyAddr { get; set; }
        public string HttpAddr { get; set; }
        public DateTime CreateTime { get; set; }
        public DateTime UpdateTime { get; set; }
        public string UpdateBy { get; set; }
        public decimal MinAvailableMoney { get; set; }
        public int RoomFlagLevel { get; set; }
        public string RoomFlagLevelText { get; set; }



        /// <summary>
        /// 已分配給哪幾家商戶
        /// </summary>
        public int AssignMerchantTotalCount { get; set; }

    }



}
