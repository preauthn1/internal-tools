using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;

namespace target.Models.Table
{

    /// <summary>
    /// 這張表 是為了加快perGameDetailLog查詢而設計
    /// 如果透過jackpotAmount的欄位查資料 需要額外建index
    /// </summary>
    public class PerGameDetailLogTagMappingEntity
    {
        [Key]
        public long Id { get; set; }

        public long PerGameDetailLogID { get; set; }

        public int TagID { get; set; }

    }
}
