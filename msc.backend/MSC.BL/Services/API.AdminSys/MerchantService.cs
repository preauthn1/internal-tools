using System;
using System.Collections.Generic;
using MSC.Models.CoreLibs;
using MSC.Models.Paging;
using MSC.Models.Table;
using MSC.Models.ViewModels;
using MSC.Utility.Extensions;
using MSC.BL.Interface;
using MSC.Models.Enums;
using MSC.Models.Models;
using System.Linq;
using MSC.Utility.Helper;
using MSC.Models;
using MSC.DA.Repository;
using MSC.Models.Attributes;
using target.Utility.Helper;
using Lck.Utility.Extensions.Helper;

namespace MSC.BL.Services
{

    public class MerchantService : BaseAdminService, IMerchantService
    {
        
        public MerchantRepository MerchantRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }

        public MerchantWalletRepository MerchantWalletRepository { get; set; }

        public UserMerchantMappingRepository UserMerchantMappingRepository { get; set; }


        public PagedList<MerchantViewModel> GetByPagingFilter(int userID, QueryMerchantApiRequestParameterModel filter)
        {
            filter.MerchantIDs = filter.MerchantID.ToIntegerListCollection();
            var query = MerchantRepository.GetPagingByFilter(filter);

            Dictionary<int, List<MerchantCustomizedSettingConvertModel>> settings = new Dictionary<int, List<MerchantCustomizedSettingConvertModel>>();

            if(query.Items.IsNotEmptyNull())
            {
                List<int> merchantIDs = query.Items.Select(x => x.MerchantID).ToList();
                settings = MerchantCustomizedSettingRepository.GetMultipleMerchants(merchantIDs , new List<int>
                {
                    MerchantCustomizedSettingKeyEnum.MaxTimesOfRequestLinkPerMember.ToInteger(),
                    MerchantCustomizedSettingKeyEnum.Basic_AvailableAmounts.ToInteger(),
                    MerchantCustomizedSettingKeyEnum.Basic_CardInputPage_DomainUrl.ToInteger(),
                    MerchantCustomizedSettingKeyEnum.DefaultTransactionFeeRate.ToInteger(),
                });
            }
            
            var viewModels = ConvertToNewTypePagingObject<MerchantEntity, MerchantViewModel>(query);
            viewModels.Items = query.Items.Select(x =>
                {
                    var convert = ConvertToViewModels(x , settings);
                    return convert;
                })
                .ToList();

            SetExtraAttribute(viewModels.Items);

            return viewModels;
        }



        public SystemConfigSettingPageViewModel GetMerchantSettings(int userID , int merchantID)
        {
            var querySettings = MerchantCustomizedSettingRepository.GetByFilter(new QueryMerchantCustomizedSettingRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection()
            });

            var allEnums = ExtensionOfEnum.GetEnumAllValue<MerchantCustomizedSettingKeyEnum>().Select(x =>
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
                var dbSetting = querySettings.FirstOrDefault(x => x.CustomizedSettingKey == item.SettingKey);

                if (dbSetting != null)
                {
                    //把db設值帶入
                    item.Value = dbSetting.SettingValue;
                }
            }

            return convertResult;
        }


        public void UpdateMerchantSettings(int userID, int merchantID , List<SystemConfigSettingUpdateRequestItemModel> models)
        {
            var querySettings = MerchantCustomizedSettingRepository.GetByFilter(new QueryMerchantCustomizedSettingRequestViewModel
            {
                MerchantIDs = merchantID.ToListCollection(),
                CustomizedSettingKeyIDs = models.Select(z=> z.SettingKey).ToList()
            });

            foreach(var formItem in models)
            {
                //防呆格式檢查
                if(formItem.SettingKey == (int)MerchantCustomizedSettingKeyEnum.Basic_AvailableAmounts)
                {
                    formItem.Value = CheckAmountsArrayIsValid(formItem.Value );
                }

                var matchExist = querySettings.FirstOrDefault(x => x.CustomizedSettingKey == formItem.SettingKey);

                if(matchExist == null)
                {
                    MerchantCustomizedSettingRepository.Create(new MerchantCustomizedSettingEntity
                    {
                        CustomizedSettingKey = formItem.SettingKey,
                        SettingValue = formItem.Value,
                        CreateTime = DateTimeHelper.GetCurrentTime(),
                        MerchantID = merchantID,
                    });
                }
                else
                {
                    matchExist.SettingValue = formItem.Value;
                    MerchantCustomizedSettingRepository.Update(matchExist);
                }
            }

            CacheRemoveHelper.WhenUpdateMerchant(merchantID);

        }


        public void Create(int userID, MerchantCreateRequestModel model)
        {
            model.RuleFor(x => x.Status).Rule_NotMinus();
            model.RuleFor(x => x.DisplayName).Rule_NotEmptyText();
            model.RuleFor(x => x.Currency).Rule_NotEmptyText();

            string userName = GetUserNameByUserID(userID);

            var newEntity = model.ReturnInjecteObject<MerchantEntity>();
            newEntity.ApiSecretKey = Guid.NewGuid().ToString().Replace("-" , "");
            newEntity.CreateTime = DateTimeHelper.GetCurrentTime();
            newEntity.UpdateTime = DateTimeHelper.GetCurrentTime();
            newEntity.CreateBy = userName;
            newEntity.UpdateBy = userName;
            newEntity.Currency = model.Currency;


            // Step 1. 建立商戶資料
            int resultMerchantID = (int)MerchantRepository.Create(newEntity);

            // Step 2. 建立錢包資料
            MerchantWalletRepository.Create(new MerchantWalletEntity
            {
                MerchantID = resultMerchantID,
                WalletType = MerchantWalletTypeEnum.Main.ToInteger(),
                Amount = 0
            });

            // Step 3. 塞預設的setting資料
            UpdateMerchantSettings(resultMerchantID, MerchantCustomizedSettingKeyEnum.Basic_AvailableAmounts, "");
            UpdateMerchantSettings(resultMerchantID, MerchantCustomizedSettingKeyEnum.Basic_CardInputPage_DomainUrl, "");
            UpdateMerchantSettings(resultMerchantID, MerchantCustomizedSettingKeyEnum.MaxTimesOfRequestLinkPerMember, "50");
            UpdateMerchantSettings(resultMerchantID, MerchantCustomizedSettingKeyEnum.DefaultTransactionFeeRate, "0.15");


            // Step 4. 將新增的merchant分配給當前使用者
            UserMerchantMappingRepository.Create(new UserMerchantMappingEntity
            {
                UserID = userID,
                MerchantID = resultMerchantID
            });


            // Step 4. 清除Cache
            CacheRemoveHelper.WhenCreateMerchant(userID);
        }


        public void Update(int userID, MerchantUpdateRequestModel model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.DisplayName).Rule_NotEmptyText();
            model.RuleFor(x => x.Status).Rule_NotMinus();

            var queryExist = MerchantRepository.GetById(model.MerchantID);

            if (queryExist == null)
                throw new MSCException(ErrorCode.NotFoundMerchant);

            string userName = GetUserNameByUserID(userID);

            queryExist.DisplayName = model.DisplayName;
            queryExist.Memo = model.Memo;
            queryExist.Status = model.Status;
            queryExist.CreateTime = DateTimeHelper.GetCurrentTime();
            queryExist.UpdateTime = DateTimeHelper.GetCurrentTime();
            queryExist.CreateBy = userName;
            queryExist.UpdateBy = userName;

            MerchantRepository.Update(queryExist);

            //清除cache
            CacheRemoveHelper.WhenUpdateMerchant(queryExist.MerchantID);
        }


        public void InitializeWalletData()
        {
            //針對沒有wallet的商戶 建立預設值          

            var queryAllMerchants = MerchantRepository.GetAll();

            var queryAllWallet = MerchantWalletRepository.GetByFilter(new QueryMerchantWalletRequestViewModel
            {
                WalletTypes = MerchantWalletTypeEnum.Main.ToInteger().ToListCollection()
            });

            List<int> walletExistMerchants = queryAllWallet.Select(z => z.MerchantID).ToList();

            //找出沒有建wallet資料的商戶
            var withoutWalletMerchants = queryAllMerchants
                .Where(z => !walletExistMerchants.Contains(z.MerchantID))
                .ToList();


            var insertData = new List<MerchantWalletEntity>();

            foreach(var item in withoutWalletMerchants)
            {
                insertData.Add(new MerchantWalletEntity
                {
                    MerchantID = item.MerchantID,
                    WalletType = MerchantWalletTypeEnum.Main.ToInteger(),
                    Amount = 0
                });
            }

            MerchantWalletRepository.CreateMany(insertData);
        }


        private string CheckAmountsArrayIsValid(string availableAmounts)
        {
            if (string.IsNullOrWhiteSpace(availableAmounts))
                return string.Empty;

            var split = availableAmounts.Split(',').Select(x => x.Trim()).Where(z=> z.IsNumeric()).ToList();
            return string.Join(",", split);
        }


        private void SetExtraAttribute(List<MerchantViewModel> models)
        {
            List<int> merchantIDs = models.Select(z => z.MerchantID).ToList();

            var queryWallet = MerchantWalletRepository.GetByFilter(new QueryMerchantWalletRequestViewModel
            {
                MerchantIDs = merchantIDs
            });

            foreach (var item in models)
            {
                var match = queryWallet.FirstOrDefault(x => x.MerchantID == item.MerchantID);

                if (match != null)
                {
                    item.WalletAmount = match.Amount;
                }
            }
        }


        private void UpdateMerchantSettings(int merchantID , MerchantCustomizedSettingKeyEnum settingKey, string value)
        {
            var queryExist = MerchantCustomizedSettingRepository.GetByFilter(new QueryMerchantCustomizedSettingRequestViewModel
                {
                    CustomizedSettingKeyIDs = settingKey.ToInteger().ToListCollection(),
                    MerchantIDs = merchantID.ToListCollection()
                })
                .FirstOrDefault();

            if(queryExist == null)
            {
                MerchantCustomizedSettingRepository.Create(new MerchantCustomizedSettingEntity
                {
                    CreateTime = DateTimeHelper.GetCurrentTime(),
                    MerchantID = merchantID,
                    SettingValue = value,
                    CustomizedSettingKey = settingKey.ToInteger(),                    
                });
            }
            else
            {
                queryExist.SettingValue = value;
                MerchantCustomizedSettingRepository.Update(queryExist);
            }

        }



        private MerchantViewModel ConvertToViewModels(MerchantEntity model , Dictionary<int, List<MerchantCustomizedSettingConvertModel>> settingModels)
        {
            var convert = model.ReturnInjecteObject<MerchantViewModel>();

            var match = settingModels.GetDictionaryOrDefaultNull(model.MerchantID);

            if(match != null)
            {
                convert.CardInputPage_DomainUrl = GetByKey(match , MerchantCustomizedSettingKeyEnum.Basic_CardInputPage_DomainUrl)?.SettingValue;
                convert.AvailableAmounts = GetByKey(match, MerchantCustomizedSettingKeyEnum.Basic_AvailableAmounts)?.SettingValue;
                convert.MaxTimesOfRequestLinkPerMember  = GetByKey(match, MerchantCustomizedSettingKeyEnum.MaxTimesOfRequestLinkPerMember)?.SettingValueOfInt ?? 0;
                convert.FeeRate = GetByKey(match, MerchantCustomizedSettingKeyEnum.DefaultTransactionFeeRate)?.SettingValueOfDecimal ?? 0;
            }
            convert.StatusText = MultiLanguageTextHelper.GetEnums_i18nText((NormalActiveStatus)model.Status);
            return convert;
        }

        private MerchantCustomizedSettingConvertModel GetByKey(List<MerchantCustomizedSettingConvertModel> settingModels , MerchantCustomizedSettingKeyEnum enumValue)
        {          
            return settingModels.FirstOrDefault(x=> x.CustomizedSettingKey == enumValue.ToInteger());
        }



    }
}
