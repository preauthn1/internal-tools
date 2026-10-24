using target.BL.Interface;
using target.BL.Services;
using target.DA.Repository;
using target.Models.Attributes;
using target.Models.CoreLibs;
using target.Models.Enums;
using target.Models.Models;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.Utility.Helper;
using System;
using System.Collections.Generic;
using System.Data;
using System.Data.SqlClient;
using System.Linq;

namespace target.BL.Services
{
    public class VendorConfigService : BaseAdminSysService, IVendorConfigService
    {

        public VendorConfigRepository VendorConfigRepository { get; set; }

        public VendorConfigSettingRepository VendorConfigSettingRepository { get; set; }

        public VendorConfigMerchantMappingRepository VendorConfigMerchantMappingRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }


        public PagedList<VendorConfigViewModel> GetByPagingFilter(int userID, QueryVendorConfigApiRequestParameterModel filter)
        {
            var cureentUser = GetCurrentOwnPermissionByCache(userID);
          
            var query = VendorConfigRepository.GetPagingByFilter(filter);

            var queryMappings = VendorConfigMerchantMappingRepository.GetByFilter(new QueryVendorConfigMerchantMappingRequestViewModel
            {
                IsIncludeDeleteMerchant = false
            });

            var viewModels = ConvertToNewTypePagingObject<VendorConfigEntity, VendorConfigViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModel(x , queryMappings);
                    return convert;
                })
                .ToList();

            return viewModels;
        }


        public SystemConfigSettingPageViewModel GetVendorConfigSettins(int vendorConfigID)
        {
            var queryAllSettings = VendorConfigSettingRepository.GetByFilter(new QueryVendorConfigSettingRequestViewModel
            {
                VendorConfigIDs = vendorConfigID.ToListCollection()
            });

            var allEnums = ExtensionOfEnum.GetEnumAllValue<VendorConfigSettingEnum>().Select(x =>
                new SettingEnumStandardFormatConvertModel
                {
                    Group = x.GetAttributeFromEnum<DataCategoryAttribute>()?.Value ?? DataCategoryEnum.Common,
                    SettingKey = x.ToInteger(),
                    SettingKeyEnum = x
                })
                .ToList();

            var convertResult = SettingKeyEnumLogicHelper.ConvertToEditPageViewModel(allEnums);


            foreach (var item in convertResult.Groups.SelectMany(z => z.Items))
            {
                var dbSetting = queryAllSettings.FirstOrDefault(x => x.SettingKey == item.SettingKey);

                if (dbSetting != null)
                {
                    //把db設值帶入
                    item.Value = dbSetting.Value;
                }
            }


            return convertResult;
        }


        public void SyncDefaultVendorConfigMapping(int vendorConfigID)
        {
            var queryVendorConfig = VendorConfigRepository.GetById(vendorConfigID);

            if (queryVendorConfig == null)
                throw new TerproException(ApiErrorCode.NotFoundData);

            //找出針對這個遊戲商 還沒配置的商戶
            var queryMerchants = VendorConfigMerchantMappingRepository.GetNotSettingMerchants(queryVendorConfig.VendorID);


            List<VendorConfigMerchantMappingEntity> isnertItems = new List<VendorConfigMerchantMappingEntity>();
            foreach(var merchantID in queryMerchants)
            {
                isnertItems.Add(new VendorConfigMerchantMappingEntity
                {
                    MerchantID = merchantID,
                    VendorConfigID = queryVendorConfig.VendorConfigID,
                    VendorID = queryVendorConfig.VendorID
                });
            }

            //把這些商戶塞預設值
            VendorConfigMerchantMappingRepository.CreateMany(isnertItems);
        }


        public void DeleteVendorConfigMerchantMapping(VendorConfigDeleteRequestModel model)
        {
            model.RuleFor(z => z.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(z => z.VendorConfigID).Rule_NotZeroNotMinus();

            var query = VendorConfigMerchantMappingRepository.GetByFilter(new QueryVendorConfigMerchantMappingRequestViewModel
                {
                    MerchantIDs = model.MerchantID.ToListCollection(),
                    VendorConfigIDs = model.VendorConfigID.ToListCollection()
                })
                .FirstOrDefault();

            if (query == null)
                return;

            var deleteKey = new
            {
                MerchantID = model.MerchantID,
                VendorConfigID = model.VendorConfigID
            };

            VendorConfigMerchantMappingRepository.DeleteByModel(deleteKey);
        }


        private VendorConfigViewModel ConvertToViewModel(VendorConfigEntity query , List<VendorConfigMerchantMappingEntity> mappings)
        {
            //將model轉成viewModels
            var convert = query.ReturnInjecteObject<VendorConfigViewModel>();
            convert.VendorName = MultiLanguageTextHelper.GetGameVendorDisplayName(query.VendorID);

            var myMerchants = mappings.Where(z => z.VendorConfigID == query.VendorConfigID).ToList();

            convert.NumberOfMerchants = myMerchants.Count();
            convert.MappingMerchants = myMerchants.OrderBy(z=> z.MerchantID).Select(z =>
                {
                    string name = MerchantRepository.GetSingleByCache(z.MerchantID)?.Name;
                    return new VendorConfigMappingMerchantViewModel
                    {
                        MerchantID = z.MerchantID,
                        MerchantName = $"{z.MerchantID}｜{name}"
                    };
                })
                .ToList();

            return convert;
        }



    }

}
