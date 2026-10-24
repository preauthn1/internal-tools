using System;
using System.Collections.Generic;
using target.Models.CoreLibs;
using target.Models.Paging;
using target.Models.Table;
using target.Models.ViewModels;
using target.Utility.Extensions;
using target.BL.Interface;
using target.Models.Enums;
using target.Models.Models;
using System.Linq;
using target.Utility.Helper;
using target.Models;
using target.DA.Repository;
using Lck.Utility.Extensions.Helper;
using target.Models.Attributes;
using target.Models.ViewModels.API.AdminSys;

namespace target.BL.Services
{

    public class MerchantService : BaseAdminSysService, IMerchantService
    {

        public MerchantRepository MerchantRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }

        public VendorConfigMerchantMappingRepository VendorConfigMerchantMappingRepository { get; set; }

        public MerchantGameSettingRepository MerchantGameSettingRepository { get; set; }

        public ServerMerchantMappingRepository ServerMerchantMappingRepository { get; set; }


        public PagedList<MerchantViewModel> GetByPagingFilter(QueryMerchantApiRequestParameterModel filter)
        {
            filter.OrganizationIDs = filter.OrganizationID.ToIntegerListCollection();

            var query = MerchantRepository.GetPagingByFilter(filter);

            var viewModels = ConvertToNewTypePagingObject<MerchantEntity, MerchantViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x);
                    return convert;
                })
                .ToList();

            SetExtraProperty(viewModels.Items);
            return viewModels;
        }


        public MerchantViewModel GetByID(int merchantID)
        {
            var queryMerchant = MerchantRepository.GetById(merchantID);

            if (queryMerchant == null)
                throw new TerproException(ApiErrorCode.NotFoundMerchant);

            var convert = ConvertToViewModels(queryMerchant);
            return convert;
        }

        public MerchantCopyInfoPreviewViewModel GetCopyPreviewInfo(int merchantID)
        {
            var queryMerchant = MerchantRepository.GetById(merchantID);

            if (queryMerchant == null)
                throw new TerproException(ApiErrorCode.NotFoundMerchant);

            var convert = ConvertToViewModels(queryMerchant);


            var queryEnableGames = MerchantGameSettingRepository.GetByFilter(new QueryMerchantGameSettingRequestViewModel
            {
                Status = NormalActiveStatus.Enable.ToInteger(),
                MerchantIDs = merchantID.ToListCollection()
            });

            var queryEnableServers = ServerMerchantMappingRepository.GetByFilter(new QueryServerMerchantRelationRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection(),
            });

            decimal queryConvertRate = MerchantCustomizedSettingRepository.GetMerchantGamePointConvertRateByCache(merchantID, MerchantCustomizedSettingKeyEnum.GameVendor_ConvertRate_PGS);

            bool isAllPGSGames = queryEnableGames.All(z => !AppHelper.IsExternalVendorGame(z.GameID));

            bool isSingleWallet = CommonLogicHelper.IsSingleWalletMode(queryMerchant.WalletMode);

            return new MerchantCopyInfoPreviewViewModel
            {
                Info = convert,
                IsEnableGameEmpty = queryEnableGames.Count() == 0,
                IsEnablePGSGameOnlyButServerEmpty = isAllPGSGames && queryEnableServers.Count() == 0,
                IsCodeNotSetup = isSingleWallet && string.IsNullOrWhiteSpace(queryMerchant.Code),
                EnablePgsGames = queryEnableGames.Where(z=> !AppHelper.IsExternalVendorGame(z.GameID)).Count(),
                EnableThirdPartyGames = queryEnableGames.Where(z => AppHelper.IsExternalVendorGame(z.GameID)).Count(),
                EnablePgsGameServers = queryEnableServers.Count(),
                PGSConvertRate = queryConvertRate
            };
        }


        public void Create(MerchantCreateRequestModel model)
        {
            model.RuleFor(x => x.Name).Rule_NotEmptyText().Rule_MaxLengthOK(50);
            model.RuleFor(x => x.Memo).Rule_MaxLengthOK(50);
            model.RuleFor(x => x.Currency).Rule_NotEmptyText().Rule_MaxLengthOK(10);
            model.RuleFor(x => x.DefaultLanguage).Rule_NotEmptyText().Rule_IsSupportLanguageCode();
            model.RuleFor(x => x.Status).Rule_NotMinus();
            model.RuleFor(x => x.WalletMode).Rule_NotZeroNotMinus();


            var newEntity = model.ReturnInjecteObject<MerchantEntity>();

            var currentMerchant = GetCurrentUser();
            newEntity.CreateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateBy = currentMerchant.LoginName;
            newEntity.CreateBy = currentMerchant.LoginName;
            newEntity.ChargeCurrency = model.Currency;

            // step 1. 新增基本資料
            var newMerchantID = MerchantRepository.Create(newEntity);

            // step 2. 遊戲商配置 塞預設值
            VendorConfigMerchantMappingRepository.CreateMerchantAllDefaultValue(newMerchantID);

            // step 3. 預設的商戶個性化配置
            CreateCustomizedSetting(newMerchantID, MerchantCustomizedSettingKeyEnum.DefaultWalletBalance, "0");
            CreateCustomizedSetting(newMerchantID, MerchantCustomizedSettingKeyEnum.Integration_SingleWalletApiPgsValidToken, "6f3a0ee2-eed7-4c95-b7ac-6e487f5d26b2");
        }



        public void Update(MerchantUpdateRequestModel model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.Name).Rule_NotEmptyText().Rule_MaxLengthOK(50);
            model.RuleFor(x => x.Memo).Rule_MaxLengthOK(50);
            model.RuleFor(x => x.Currency).Rule_NotEmptyText().Rule_MaxLengthOK(10);
            model.RuleFor(x => x.DefaultLanguage).Rule_NotEmptyText().Rule_IsSupportLanguageCode();
            model.RuleFor(x => x.Status).Rule_NotMinus();
            model.RuleFor(x => x.WalletMode).Rule_NotZeroNotMinus();

            var queryExist = MerchantRepository.GetById(model.MerchantID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundMerchant);

            queryExist.Name = model.Name;
            queryExist.Memo = model.Memo;
            queryExist.DefaultLanguage = model.DefaultLanguage;
            queryExist.Status = model.Status;
            queryExist.ApiSecretKey = model.ApiSecretKey;
            queryExist.UpdateTime = DateTimeHelper.GetCurrentTime();
            queryExist.UpdateBy = GetCurrentUser().LoginName;
            queryExist.WalletMode = model.WalletMode;
            queryExist.Type = model.Type;

            if (model.Code.IsNotNullEmptyText())
            {
                queryExist.Code = model.Code;
            }

            MerchantRepository.Update(queryExist);

            //更新Cache
            CacheRemoveHelper.RemoveCache_WhenUpdateMerchant(queryExist.MerchantID);

        }


        public void Delete(int merchantID)
        {
            var queryExist = MerchantRepository.GetById(merchantID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundMerchant);

            queryExist.Status = NormalActiveStatus.PermanentDelete.ToInteger();
            MerchantRepository.Update(queryExist);
        }



        private void CreateCustomizedSetting(int merchantID, MerchantCustomizedSettingKeyEnum setting, string value)
        {
            MerchantCustomizedSettingRepository.Create(new MerchantCustomizedSettingEntity
            {
                CustomizedSettingKey = setting.ToInteger(),
                MerchantID = merchantID,
                SettingValue = value,
                CreateTime = DateTimeHelper.GetCurrentTime(),
            });
        }


        private MerchantViewModel ConvertToViewModels(MerchantEntity queryMerchant)
        {
            var convert = queryMerchant.ReturnInjecteObject<MerchantViewModel>();
            convert.WalletModeText = MultiLanguageTextHelper.GetEnumsCultureText((MerchantWalletModeEnum)queryMerchant.WalletMode);
            convert.MerchantTypeName = MultiLanguageTextHelper.GetEnumsCultureText((MerchantTypeEnum)queryMerchant.Type);
           
            return convert;
        }


        private void SetExtraProperty(List<MerchantViewModel> viewModels)
        {
            List<int> merchantIDs = viewModels.Select(x => x.MerchantID).ToList();

            var querySetting = MerchantCustomizedSettingRepository.GetMultipleByMerchantSettingKey(merchantIDs, MerchantCustomizedSettingKeyEnum.DefaultWalletBalance);
                      
            var queryGameMapping = MerchantGameSettingRepository.GetByFilter(new QueryMerchantGameSettingRequestViewModel
            {
                MerchantIDs = merchantIDs,
                Status = NormalActiveStatus.Enable.ToInteger()
            });

            var queryServerMapping = ServerMerchantMappingRepository.GetByFilter(new QueryServerMerchantRelationRequestViewModel
            {
                MerchantIDs = merchantIDs,
            });

            foreach (var item in viewModels)
            {
                var matchSetting = querySetting.FirstOrDefault(x => item.MerchantID == x.MerchantID);

                if (matchSetting != null)
                {
                    item.DefaultBalance = matchSetting.SettingValueOfDecimal;
                }

                item.NumberOfGames = queryGameMapping.Where(z => z.MerchantID == item.MerchantID).Count();
                item.NumberOfServers = queryServerMapping.Where(z => z.MerchantID == item.MerchantID).Count();
            }

            return;
        }



    }
}
