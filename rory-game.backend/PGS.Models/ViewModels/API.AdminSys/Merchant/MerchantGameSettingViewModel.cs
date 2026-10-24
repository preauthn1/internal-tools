using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{

    public class MerchantGameSettingViewModel
    {

        public List<MerchantGameSettingItemViewModel> GameSettingItems { get; set; }

    }



    public class MerchantGameSettingItemViewModel
    {
        public int GameID { get; set; }

        public string GameName { get; set; }

        public int Status { get; set; }

    }



}
