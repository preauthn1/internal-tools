using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Utility.Utility.ThirdParty
{
    public enum SlackChannelEnum
    {

        [Description("https://example.invalid/services/T011TNRPVP0/B01951PUP3J/hMESIwm31JyltaYFXYj7sdPW")]
        notify_msc_api_adminsys,

        [Description("https://example.invalid/services/T011TNRPVP0/B01951LD3FA/PGOF1e6Ayz7KJjbmFl0r7SAS")]
        notify_msc_errors,            

        [Description("https://example.invalid/services/T011TNRPVP0/B0180DEMFFH/ihDwZQCHVWEsKXczVv52ZKaY")]
        notify_msc_api_points,

        [Description("https://example.invalid/services/T011TNRPVP0/B018DQCLBAR/u44BW1vH1Gbx6LkjGEGSqxQ5")]
        notify_msc_api_vendor,

        [Description("https://example.invalid/services/T011TNRPVP0/B018DQDA77F/Xvf9e7VOqc7ONXzv0ZWe8fUF")]
        notify_msc_api_web,

        [Description("https://example.invalid/services/T011TNRPVP0/B01J98RAJHZ/iQ5vdMJrgJW6wCgnxkuzY3MN")]
        notify_msc_api_job,                

    }


    /// <summary>
    /// 純粹拿來標記 是屬於哪個組織下的頻道 
    /// 跟程式碼運作無關
    /// </summary>
    public class SlackTeam : Attribute
    {
        public SlackTeam(string teamName)
        {

        }
    }
    

}
