using System.Collections.Generic;

namespace target.Models.ViewModels
{
    public class ServerListAndSettingRequestModel
    {

    }



    public class ServerSettingViewModel
    {

        public int ServerID { get; set; }
        

        public string ServerName { get; set; }


        public int FlagLevel { get; set; }


        /// <summary>
        /// 已分配給哪幾家供應商使用
        /// </summary>
        public List<int> AvailableMerchantIDs { get; set; }

    }




}
