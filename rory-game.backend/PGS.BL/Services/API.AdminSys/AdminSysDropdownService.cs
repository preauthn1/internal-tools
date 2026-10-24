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
using target.Models.ViewModels.API.AdminSys;
using Microsoft.Extensions.Caching.Memory;
using target.Models.Attributes;

namespace target.BL.Services
{

    public class AdminSysDropdownService : BaseAdminSysService, IAdminSysDropdownService
    {

        public MerchantGameSettingRepository MerchantGameSettingRepository { get; set; }

        public SystemConfigSettingRepository SystemConfigSettingRepository { get; set; }

        public OrganizationRepository OrganizationRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public MerchantCustomizedSettingRepository MerchantCustomizedSettingRepository { get; set; }

        public RoleRepository RoleRepository { get; set; }

        public OrganizationRoleMappingRepository OrganizationRoleMappingRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }

        public ServerListRepository ServerListRepository { get; set; }

        public PlayerRepository PlayerRepository { get; set; }

        public ServerMerchantMappingRepository ServerMerchantMappingRepository { get; set; }

        public DynamicQueryCommandRepository DynamicQueryCommandRepository { get; set; }


        public List<DropdownValueViewModel> GetGames(int userID)
        {
            var currentUser = GetCurrentOwnPermissionByCache(userID);

            //已經被標記為 999(刪除) 的遊戲不顯示了
            var queryGames = GameListRepository.GetByFilter(new QueryGameListRequestViewModel
            {
                IsNotDelete = true
            });

            //如果不是超級管理員 僅顯示有開通的遊戲
            if (!currentUser.IsSuperAdminOrg)
            {
                var enableGames = MerchantGameSettingRepository.GetByFilter(new QueryMerchantGameSettingRequestViewModel
                    {
                        MerchantIDs = currentUser.AvailableMerchantIDs,
                        Status = NormalActiveStatus.Enable.ToInteger()
                    })
                    .Select(z => z.GameID).Distinct().ToList();

                queryGames = queryGames.Where(z => enableGames.Contains(z.GameID)).ToList();
            }

            var result = queryGames
                .Where(z=> z.GameID > 0)
                .Select(x =>
                {
                    var text = DbMultiLanguageTextHelper.GetGameDisplayName(x.GameID);
                    text = $"{x.GameID}｜{text}";
                    var model = new DropdownValueViewModel
                    {
                        Text_CN = text,
                        Text_EN = text,
                        Value = x.GameID.ToString()
                    };

                    return model;
                })
                .ToList();


            return result;
        }


        public List<DropdownValueViewModel> GetServers(int userID , QueryServerListApiRequestParameterModel filter)
        {
            var currentUser = GetCurrentOwnPermissionByCache(userID);

            filter = filter ?? new QueryServerListApiRequestParameterModel();
            filter.GameIDs = filter.GameID.ToIntegerListCollection();

            if(!currentUser.IsSuperAdminOrg)
            {
                List<int> merchantIDs = currentUser.AvailableMerchantIDs;
                filter.ServerIDs = ServerMerchantMappingRepository.GetAvailableServersByMerchantIDs(merchantIDs);
            }

            var query = ServerListRepository.GetByFilter(filter);

            var result = query.Select(x => new DropdownValueViewModel
                {
                    Text_CN = x.ServerName,
                    Text_EN = x.ServerName,
                    Value = x.ServerID.ToString(),
                    ExtraProperty = new
                    {
                        GameID = x.GameID,
                        Currency = x.Currency
                    }
                })
                .ToList();

            return result;
        }


        public List<DropdownValueViewModel> GetTimeZoneAreas()
        {
            var allEnums = ExtensionOfEnum.GetEnumAllValue<TimeZoneAreaEnum>();

                var result = allEnums
                    .OrderBy(x=> DataModelConvertHelper.GetTimeZoneByAreaID(x.ToInteger()))
                    .Select(x =>
                    {
                        var displayName = MultiLanguageTextHelper.GetTimeZoneDisplayName(x.ToInteger());
                        var convert = new DropdownValueViewModel
                        {
                            Text_CN = displayName,
                            Text_EN = displayName,
                            Value = x.ToIntegerString(),
                            ExtraProperty = new
                            {
                                TimeZone = DataModelConvertHelper.GetTimeZoneByAreaID(x.ToInteger())
                            }
                        };

                        return convert;
                    })
                    .ToList();
            
            return result;
        }



        public List<DropdownValueViewModel> GetAvailableRoles(int userID)
        {
            var currentUser = GetCurrentUser();

            var queryRoles = RoleRepository.GetAll();

            //如果非超級管理員 篩選roles
            if (!currentUser.IsSuperAdminOrg)
            {
                var queryAvailableRoles = OrganizationRoleMappingRepository.GetByFilter(new QueryOrganizationRoleMappingRequestViewModel
                {
                    OrganizationIDs = currentUser.OrganizationID.ToListCollection()
                })
                .Select(z => z.RoleID)
                .ToList();

                queryRoles = queryRoles
                    .Where(x => queryAvailableRoles.Contains(x.RoleID))
                    .ToList();
            }

            var queryUserOwnPermission = GetCurrentOwnPermissionByCache(userID);

            return queryRoles.Select(x => new DropdownValueViewModel
                {
                    Text_CN = x.RoleName,
                    Text_EN = x.RoleName,
                    Value = x.RoleID.ToString()
                })
                .ToList();
        }
        

        public List<DropdownValueViewModel> GetAvailableMerchants(int userID , QueryDropdownMerchantRequestViewModel filter)
        {
            var querySettings = SystemConfigSettingRepository.GetParseModelByCache(SystemConfigSettingEnum.InternalTest_Favorit_MerchantIDs);

            List<int> favoritMerchantIDs = new List<int>();

            if(querySettings != null)
            {
                favoritMerchantIDs = querySettings.SettingValueOfListInt;
            }

            var queryUserOwnPermission = GetCurrentOwnPermissionByCache(userID);

            var query = MerchantRepository.GetByFilter(new QueryMerchantRequestViewModel
            {
                OrganizationIDs = filter.OrganizationID.HasValue ? filter.OrganizationID.Value.ToListCollection() : new List<int>(),
                MerchantIDs = queryUserOwnPermission.SuperAdminOrAvailableMerchantIDs
            });

            var currentUser = GetCurrentOwnPermissionByCache(userID);

            return query
                .OrderBy(x => favoritMerchantIDs.Contains(x.MerchantID) ? favoritMerchantIDs.IndexOf(x.MerchantID) : 9999999)
                .ThenBy(x => x.Name.Contains("測試") ? 999 : 1)
                .ThenBy(x => x.Name.Contains("机器人") ? 888 : 1)
                .ThenBy(x => x.MerchantID)
                .Select(x => new DropdownValueViewModel
                {
                    Text_CN = currentUser.IsSuperAdminOrg ? $"{x.MerchantID}｜{x.Name}" : x.Name, //超級管理組織會顯示merchantID
                    Text_EN = x.Name,
                    Value = x.MerchantID.ToString()
                })         
                .ToList();
        }


        public List<DropdownValueViewModel> GetInternalMerchants(int userID)
        {
            var querySettings = MerchantCustomizedSettingRepository.GetByFilter(new QueryMerchantCustomizedSettingRequestViewModel
            {
                CustomizedSettingKeyIDs = MerchantCustomizedSettingKeyEnum.IsInternalTestMerchant.ToInteger().ToListCollection(),
                SettingValueMatch = "1"
            });

            if (querySettings.IsEmptyOrNull())
                return new List<DropdownValueViewModel>();

            List<int> merchantIDs = querySettings.Select(z => z.MerchantID).ToList();

            var queryMerchants = MerchantRepository.GetByIds(merchantIDs);

            return queryMerchants.Select(x => new DropdownValueViewModel
            {
                Text_CN = $"{x.MerchantID}｜{x.Name}",
                Text_EN = x.Name,
                Value = x.MerchantID.ToString()
            })
            .ToList();
        }


        public List<DropdownValueViewModel> GetAvailableOrgs(int userID)
        {
            var queryUserOwnPermission = GetCurrentOwnPermissionByCache(userID);

            var filter = new QueryOrganizationRequestViewModel
            {
                OrganizationIDs = queryUserOwnPermission.IsSuperAdminOrg ? new List<int>() : queryUserOwnPermission.OrganizationID.ToListCollection()
            };

            var query = OrganizationRepository.GetByFilter(filter);
            return query.Select(x => new DropdownValueViewModel
            {
                Text_CN = x.OrganizationName,
                Text_EN = x.OrganizationName,
                Value = x.OrganizationID.ToString(),
                ExtraProperty = new
                {
                    OrganizationCode = x.OrganizationCode
                }
            })
            .ToList();
        }


        public List<DropdownValueViewModel> GetMerchantTypes()
        {
            var queryAll = ExtensionOfEnum.GetEnumAllValue<MerchantTypeEnum>();

            return queryAll.Select(x => new DropdownValueViewModel
            {
                Text_CN = MultiLanguageTextHelper.GetEnumsCultureText(x),
                Text_EN = MultiLanguageTextHelper.GetEnumsCultureText(x),
                Value = x.ToIntegerString(),
            })
            .OrderBy(z=> z.Text_CN)
            .ToList();
        }


        public List<DropdownValueViewModel> GetRoomLevels()
        {
            var queryAll = ExtensionOfEnum.GetEnumAllValue<GameRoomFlagLevelEnum>();

            return queryAll.Select(x => new DropdownValueViewModel
            {
                Text_CN = MultiLanguageTextHelper.GetEnumsCultureText(x),
                Text_EN = MultiLanguageTextHelper.GetEnumsCultureText(x),
                Value = x.ToIntegerString(),
            })
            .ToList();
        }

        public List<DropdownValueViewModel> GetTransactionTypes()
        {
            var queryAll = ExtensionOfEnum.GetEnumAllValue<TransactionTypeEnum>();

            return queryAll.Select(x => new DropdownValueViewModel
            {
                Text_CN = MultiLanguageTextHelper.GetEnumsCultureText(x),
                Text_EN = MultiLanguageTextHelper.GetEnumsCultureText(x),
                Value = x.ToIntegerString(),
            })
            .ToList();
        }

        public List<DropdownValueViewModel> GetThirdPartyVendors()
        {
            var queryAll = ExtensionOfEnum.GetEnumAllValue<GameVendorTypeEnum>();

            //大於99999 才視為三方遊戲商
            return queryAll.Where(x=> x.ToInteger() > 99999)
                .Select(x => new DropdownValueViewModel
                {
                    Text_CN = x.ToString(),
                    Text_EN = x.ToString(),
                    Value = x.ToIntegerString(),
                })
                .ToList();
        }


        public List<DropdownValueViewModel> GetDynamicQueryItems()
        {
            var query = DynamicQueryCommandRepository.GetAll();

            return query
                .Select(x =>
                {                   
                    var convert = new DropdownValueViewModel
                    {
                        Text_CN = x.DisplayName,
                        Value = x.Id.ToString(),
                    };

                    return convert;
                })
                .ToList();
        }


        public List<DropdownValueViewModel> GetInternalTestPlayerAccounts()
        {
            var query = SystemConfigSettingRepository.GetParseModelByEnum(SystemConfigSettingEnum.InternalTestAccountPlayerIDs);

            if (query == null || query.SettingValueOfListInt.IsEmptyOrNull())
                return new List<DropdownValueViewModel>();

            var queryPlayers = PlayerRepository.GetByIds(query.SettingValueOfListInt);

            return queryPlayers
                .OrderBy(z => z.MerchantID)
                .Select(x =>
                {
                    string merchantName = MerchantRepository.GetSingleByCache(x.MerchantID)?.Name;
                    string displayName = $"{x.PlayerID}｜{x.PlayerName}｜{merchantName}";
                    var convert = new DropdownValueViewModel
                    {
                        Text_CN = displayName,
                        Value = x.PlayerID.ToString(),
                    };

                    return convert;
                })
                .ToList();
        }

        public List<DropdownValueViewModel> GetBetRecordStatus()
        {
            var queryAll = ExtensionOfEnum.GetEnumAllValue<GameVendorBetRecordStatusEnum>();

            return queryAll.Select(x => new DropdownValueViewModel
            {
                Text_CN = MultiLanguageTextHelper.GetEnumsCultureText(x),
                Text_EN = MultiLanguageTextHelper.GetEnumsCultureText(x),
                Value = x.ToIntegerString(),
            })
            .ToList();
        }


        public List<DropdownValueViewModel> GetGameTypes()
        {
            var queryAll = ExtensionOfEnum.GetEnumAllValue<GameTypeEnum>();

            return queryAll.Select(x => new DropdownValueViewModel
            {
                Text_CN = MultiLanguageTextHelper.GetEnumsCultureText(x),
                Text_EN = MultiLanguageTextHelper.GetEnumsCultureText(x),
                Value = x.ToIntegerString(),
            })
            .ToList();
        }


        public List<DropdownValueViewModel> GetGameVendorTypes()
        {
            var queryAll = ExtensionOfEnum.GetEnumAllValue<GameVendorTypeEnum>();

            return queryAll.Select(x => new DropdownValueViewModel
            {
                Text_CN = MultiLanguageTextHelper.GetEnumsCultureText(x),
                Text_EN = MultiLanguageTextHelper.GetEnumsCultureText(x),
                Value = x.ToIntegerString(),
            })
            .ToList();
        }



        private string GetValueOrDefaultValue(string expectValue , string defaultValue)
        {
            if (string.IsNullOrWhiteSpace(expectValue))
                return defaultValue;

            return expectValue;
        }
        
    }
}
