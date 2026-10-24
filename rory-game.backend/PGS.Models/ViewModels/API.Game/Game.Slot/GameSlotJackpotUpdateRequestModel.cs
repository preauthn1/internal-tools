namespace target.Models.ViewModels
{
    /// <summary>
    /// 更新彩池金額
    /// </summary>
    public class GameSlotJackpotUpdateRequestModel
    {
        public int FlagID { get; set; }

        public int TableNo { get; set; }

        public int LevelID { get; set; }

        public decimal Contribution { get; set; }

        public decimal TotalContribution { get; set; }

        public decimal TotalAward { get; set; }

        public bool IsLink { get; set; }
    }
}
