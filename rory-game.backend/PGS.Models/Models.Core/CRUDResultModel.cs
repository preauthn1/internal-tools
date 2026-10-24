namespace target.Models.CoreLibs
{
    public class CRUDResultModel
    {
        public CRUDResultModel()
        {

        }
        
        public CRUDResultModel (string msg)
        {
            Message = msg;
        }

        public CRUDResultModel(bool isOK , string msg = "储存成功")
        {
            IsSuccess = isOK;
            Message = msg;
        }



        public bool IsSuccess { get; set; }

        public string Message { get; set; }

        public int ResultId { get; set; }

        public object ResultObject { get; set; }
    }
}
