using MSC.DA.Repository;
using MSC.Models.Table;
using MSC.Utility;
using MSC.Utility.Helper;
using MSC.Utility.Utility;
using MSC.Utility.Extensions;
using System;
using System.Collections.Generic;
using Lck.Utility.Extensions.Helper;
using MSC.Models.Models;
using System.Linq;
using MSC.Models.Enums;

namespace MSC.BL.Services
{

    public class RequestForbiddenRefService : BaseService, IRequestForbiddenRefService
    {

        public void AddForbiddenCounterTimes(RequestForbiddenTypeEnum type, string forbiddenKeyOrIP)
        {
            long currentTime = DateTimeHelper.GetCurrentTimeTimestamp();

            string cacheKey = GetRequestCounterRecord_CacheKey(type , forbiddenKeyOrIP);
            List<long> cacheValues = GetIPAddressRequestRecords(cacheKey);

            if (cacheValues == null)
            {
                cacheValues = new List<long>
                {
                    currentTime
                };
            }
            else
            {
                cacheValues.Add(currentTime);
            }

            //如果時間分鐘內 又再一次請求進來 cache存活時間延長 (cache存活時間為鎖定時間的兩倍)
            int expiredMinutes = GetFobbiddenLockMinutesByType(type) * 2;
            CacheHelper.SetCacheValue(cacheKey, cacheValues, new TimeSpan(0 , expiredMinutes , 0));
        }


        public IPForbiddenCheckResultModel CheckForbiddenCounterTooManyTimes(RequestForbiddenTypeEnum type , string forbiddenKeyOrIP , int maxTimes)
        {
            string myIP = IPHelper.GetIP();

            string cacheKeyOfForbidden = GetRequestForbidden_CacheKey(type , forbiddenKeyOrIP);

            var queryForbiddenCache = CacheHelper.GetCacheValue<string>(cacheKeyOfForbidden);

            //如果Forbidden cache已經存在，直接告知呼叫端拋錯誤
            if (queryForbiddenCache != null)
            {
                return new IPForbiddenCheckResultModel
                {
                    IsFobidden = true,
                    ForbiddenCacheKey = cacheKeyOfForbidden
                };
            }

            string cacheKey = GetRequestCounterRecord_CacheKey(type ,forbiddenKeyOrIP);
            List<long> cacheValues = GetIPAddressRequestRecords(cacheKey);

            int requestTimes = cacheValues.Count();

            //次數太少 不做後面的邏輯檢查 (減少運算量)
            if (cacheValues.Count() < maxTimes)
            {
                return new IPForbiddenCheckResultModel
                {
                    IsFobidden = false
                };
            }

            long previous_hours = DateTimeHelper.GetCurrentTime().AddHours(-1).CovertToUnixTimestamp();

            //搜尋cache紀錄中 時間在一小時內的紀錄 (這cache裡面存的是ts的陣列)
            cacheValues = cacheValues.Where(z => z >= previous_hours).ToList();

            if (cacheValues.Count() < maxTimes)
            {
                return new IPForbiddenCheckResultModel
                {
                    IsFobidden = false
                };
            }

            int expiredMinutes = GetFobbiddenLockMinutesByType(type);
            CacheHelper.SetCacheValue(cacheKeyOfForbidden, "true", new TimeSpan(0, expiredMinutes, 0));


            //如果Forbidden cache存在 後面的實作可能是 拋錯誤/寫紀錄 ...
            return new IPForbiddenCheckResultModel
            {
                IsFobidden = true,
                IsNewForbiddenEvent = true,
                ForbiddenCacheKey = cacheKeyOfForbidden
            };
        }

        /// <summary>
        /// 根據type定義 每一種 請求鎖定(黑名單) 的時間長短
        /// </summary>
        /// <param name="type"></param>
        /// <returns></returns>
        private int GetFobbiddenLockMinutesByType(RequestForbiddenTypeEnum type)
        {
            switch (type)
            {
                case RequestForbiddenTypeEnum.Request_CardInput_CheckCardNumber_IPAddress:
                    return 60;

                case RequestForbiddenTypeEnum.Request_CardInput_CheckCardNumber_TokenHash:
                    return 240;

                default:
                    return 60;
            }
        }



        /// <summary>
        /// CacheKey - 計數器
        /// </summary>
        private string GetRequestCounterRecord_CacheKey(RequestForbiddenTypeEnum type, string forbiddenKeyOrIP)
        {
            string cacheKey = $"{type.ToString()}_{forbiddenKeyOrIP}";
            return cacheKey;
        }


        /// <summary>
        /// CacheKey - 請求封鎖 (ip、Token...etc)
        /// </summary>
        private string GetRequestForbidden_CacheKey(RequestForbiddenTypeEnum type , string forbiddenKeyOrIP)
        {
            string cacheKey = $"{type.ToString()}_Forbidden_{forbiddenKeyOrIP}";
            return cacheKey;
        }



        private List<long> GetIPAddressRequestRecords(string cacheKey)
        {
            var cacheValues = CacheHelper.GetCacheValue<List<long>>(cacheKey);
            return cacheValues ?? new List<long>();
        }

    }
}
