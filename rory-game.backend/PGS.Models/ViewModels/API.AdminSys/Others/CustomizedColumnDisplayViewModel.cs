using target.Models.Table;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace target.Models.ViewModels
{
    public class CustomizedColumnDisplayRequestModel
    {

        /// <summary>
        /// 對應 DynamicQueryCommandRepository 這張表的ID
        /// </summary>
        public int QuerySqlID { get; set; }

        public int? PlayerID { get; set; }

        public string PlayerName { get; set; }

        public int? MerchantID { get; set; }

        public int? GameID { get; set; }

        public DateTime? BeginTime { get; set; }

        public DateTime? EndTime { get; set; }



    }


    public class CustomizedColumnDynamicQueryResultViewModel
    {
        public List<Dictionary<string , object>> Items { get; set; }

        public CustomizedColumnExtraInfoViewModel ExtraObject { get; set; }
    }


    public class CustomizedColumnExtraInfoViewModel
    {
        public List<CustomizedColumnDisplayColumnViewModel> Columns { get; set; }

        public string CopyExcel { get; set; }
    }

    public class CustomizedColumnDisplayColumnViewModel
    {

        public CustomizedColumnDisplayColumnViewModel()
        {

        }

        public CustomizedColumnDisplayColumnViewModel(string columnName)
        {
            this.ColumnName = columnName;
        }

        public string ColumnName { get; set; }
    }


    public class CustomizedColumnDisplayViewModel
    {
        public string Column1 { get; set; }
        public string Column2 { get; set; }
        public string Column3 { get; set; }
        public string Column4 { get; set; }
        public string Column5 { get; set; }
        public string Column6 { get; set; }
        public string Column7 { get; set; }
        public string Column8 { get; set; }
        public string Column9 { get; set; }
    }


}
