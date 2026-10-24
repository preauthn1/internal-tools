using target.Models.Enums;
using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{


    public class JobSettingControlViewModel
    {

        public bool IsActive { get; set; }

        public DateTime? LastExecuteTimes { get; set; }

    }

    public class JobExecuteRecordModel
    {
        public JobActionTypeEnum JobType { get; set; }

        public string Memo { get; set; }

        public DateTime StartTime { get; set; }

        public DateTime EndTime { get; set; }

    }




}
