using target.Models;
using System;
using System.Collections.Generic;
using System.Text;
using target.Utility.Extensions;
using target.Models.Models;
using target.Models.Enums;
using target.Models.Table;
using target.BL.Interface;
using target.Utility.IOC;

namespace target.BL.Services
{
    public class ServiceInstanceHelper
    {
        public static IGameVendorIntegrationService InstanceVendorService(GameVendorTypeEnum vendorID)
        {
            return InstanceVendorService(vendorID.ToInteger());
        }


        public static IGameVendorIntegrationService InstanceVendorService(int vendorID)
        {
            try
            {
                return IOCHelper.ResolveByKey<IGameVendorIntegrationService>(vendorID);
            }
            catch(Exception ex)
            {
                //如果這邊拋錯誤 請確認 IOCInitializeHelper.cs 有正確註冊 每一家遊戲商的service
                throw ex;
            }
        }


    }
}
