using Lck.Utility.Extensions.Helper;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using MSC.Utility.Extensions;
using Lck.Utility.Extensions.Cache;

namespace Lck.Utility.Extensions.Helper
{
   
    /// <summary>
    /// 純粹針對 MemoryCache
    /// </summary>
    public enum CacheTypeEnum
    {
        MemoryCache = 1,

        RedisCache = 2,      

    }
}
