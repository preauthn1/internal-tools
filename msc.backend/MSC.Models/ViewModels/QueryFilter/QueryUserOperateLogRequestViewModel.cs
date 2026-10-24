using MSC.Models.Paging;
using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace MSC.Models.ViewModels
{

    public class QueryUserOperateLogApiRequestParameterModel : QueryUserOperateLogRequestViewModel
    {
        //建議新增一個class 是這樣的命名 並放到 RequestViewModel 同一個檔案中
        //建議新增一個class 是這樣的命名 並放到 RequestViewModel 同一個檔案中
        //建議新增一個class 是這樣的命名 並放到 RequestViewModel 同一個檔案中
        //建議新增一個class 是這樣的命名 並放到 RequestViewModel 同一個檔案中
        //建議新增一個class 是這樣的命名 並放到 RequestViewModel 同一個檔案中
        //建議新增一個class 是這樣的命名 並放到 RequestViewModel 同一個檔案中
    }


    public class QueryUserOperateLogRequestViewModel : PaginationParameter
    {

        public List<int> OperationTypes { get; set; }

        public List<int> UserIDs { get; set; }

        public string UserName { get; set; }

        public DateTime? CreateTimeBegin { get; set; }

        public DateTime? CreateTimeEnd { get; set; }
    }


}
