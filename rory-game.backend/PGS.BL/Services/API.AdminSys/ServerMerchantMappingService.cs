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

namespace target.BL.Services
{

    public class ServerMerchantMappingService : BaseAdminSysService, IServerMerchantMappingService
    {
        public ServerListRepository ServerListRepository { get; set; }

        public ServerMerchantMappingRepository ServerMerchantMappingRepository { get; set; }

        public MerchantRepository MerchantRepository { get; set; }

        public PlayerCurrentServerRepository PlayerCurrentServerRepository { get; set; }

        public GameListRepository GameListRepository { get; set; }



        public AllServerMerchantMappingViewModel GetAllServerMerchantMappings()
        {
            var queryAllServers = ServerListRepository.GetAll();

            var queryAllMappings = ServerMerchantMappingRepository.GetAll();
            

            var queryAllMerchants = MerchantRepository.GetByFilter(new QueryMerchantRequestViewModel
                {
                    IsIncludeDeleteData = false
                })
                .OrderBy(x=> x.Currency)
                .ThenBy(x => x.MerchantID)
                .ToList();


            List<string> existCurrency = queryAllMerchants.Select(x => x.Currency).Distinct().ToList();

            var currencyGroups = new List<CurrencyCategoryItemModel>();


            List<int> GameIDs = queryAllServers.Select(x => x.GameID).Distinct().ToList();

            var queryGames = GameListRepository.GetByFilter(new QueryGameListRequestViewModel
            {
                GameIDs = GameIDs
            });
            

            //吐給前端是每個幣別一個分組
            foreach(string currency in existCurrency)
            {
                var eachCurrency = new CurrencyCategoryItemModel
                {
                    Merchants = queryAllMerchants
                        .Where(x=> x.Currency == currency)
                        .Select(x => new BatchEditMappingMerchantItemViewModel
                        {
                            MerchantID = x.MerchantID,
                            MerchantName = x.Name,
                            Currency = x.Currency
                        })
                        .ToList(),
                    Servers = new List<BatchEditMappingServerItemViewModel>(),
                    Currency = currency
                };

                List<BatchEditMappingGameItemViewModel> gamesViewModels = queryGames.Select(x => new BatchEditMappingGameItemViewModel
                {
                    GameID = x.GameID,
                    GameName = x.GameName,
                })
                .ToList();


                foreach (var item in gamesViewModels)
                {
                    var matchServersThisGame = queryAllServers
                        .Where(x => x.GameID == item.GameID)
                        .Where(x=> x.Currency == currency)
                        .ToList();

                    if (matchServersThisGame.IsEmptyOrNull())
                        continue;

                    var servers = matchServersThisGame.Select(x => new BatchEditMappingServerItemViewModel
                        {
                            GameID = item.GameID,
                            GameName = item.GameName,
                            ServerID = x.ServerID,
                            ServerName = x.ServerName,
                            Currency = x.Currency
                        })
                        .ToList();

                    foreach (var eachServer in servers)
                    {
                        eachServer.MerchantEnableFlags = new List<BatchEditMappingMerchantFlagItemViewModel>();

                        foreach (var eachMerchant in eachCurrency.Merchants)
                        {
                            var matchSetting = queryAllMappings.FirstOrDefault(x => x.MerchantID == eachMerchant.MerchantID && x.ServerID == eachServer.ServerID);
                            bool isEnable = matchSetting != null ? true : false;

                            eachServer.MerchantEnableFlags.Add(new BatchEditMappingMerchantFlagItemViewModel
                            {
                                IsEnable = isEnable,
                                MerchantName = eachMerchant.MerchantName,
                                MerchantID = eachMerchant.MerchantID,
                                MerchantCurrency = eachMerchant.Currency,
                            });
                        }
                    }

                    item.Servers = servers;
                    eachCurrency.Servers.AddRange(servers);
                }

                //eachCurrency.Games = gamesViewModels;

                currencyGroups.Add(eachCurrency);
            }

            return new AllServerMerchantMappingViewModel
            {
                AvailableCurrency = currencyGroups.Where(z=> z.Servers.IsNotEmptyNull()).Select(z=> z.Currency).ToList(),
                CurrencyGroups = currencyGroups.OrderBy(x => x.Servers.IsEmptyOrNull() ? 999 : 1).ThenBy(x => x.Currency).ToList()
            };
        }


        public ReversServerMerchantResultViewModel ReversMerchantMappings(ReversServerMerchantRequestViewModel model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            model.RuleFor(x => x.ServerID).Rule_NotZeroNotMinus();

            var queryServerMapping = ServerMerchantMappingRepository.GetByFilter(new QueryServerMerchantRelationRequestViewModel
                {
                    ServerIDs = model.ServerID.ToListCollection(),
                    MerchantIDs = model.MerchantID.ToListCollection()
                })
               .FirstOrDefault();

            bool beforeValue = queryServerMapping == null ? false : true;
            bool afetrValue = !beforeValue;

            this.UpdateServerMerchantMappings(model.ServerID, new ServerAvailableMerchantUpdateRequestModel
            {
                IsEnable = afetrValue,
                MerchantID = model.MerchantID
            });

            return new ReversServerMerchantResultViewModel
            {
                AfterValue = afetrValue
            };
        }


        public ServerMerchantMappingViewModel GetServerMerchantMappings(int serverID)
        {
            var queryExist = ServerListRepository.GetById(serverID);

            if (queryExist == null)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            var queryServerMapping = ServerMerchantMappingRepository.GetByServerID(serverID);

            var queryAllMerchants = MerchantRepository.GetByFilter(new QueryMerchantRequestViewModel
            {
                IsIncludeDeleteData = false
            });

            return new ServerMerchantMappingViewModel
            {
                AvailableSettings = queryAllMerchants.Select(x => new ServerMerchantMappingSettingVIewModel
                {
                    MerchantID = x.MerchantID,
                    MerchantName = x.Name,
                    IsEnable = queryServerMapping.Any(z => z.MerchantID == x.MerchantID)
                })
                .OrderBy(x => x.MerchantID)
                .ToList()
            };
        }


        public void UpdateServerMerchantMappings(int serverID, ServerAvailableMerchantUpdateRequestModel model)
        {
            model.RuleFor(x => x.MerchantID).Rule_NotZeroNotMinus();
            
            var queryServer = ServerListRepository.GetById(serverID);

            if (queryServer == null)
                throw new TerproException(ApiErrorCode.NotFoundServer);

            var queryMercahnt = MerchantRepository.GetById(model.MerchantID);
            
            if (queryMercahnt == null)
                throw new TerproException(ApiErrorCode.NotFoundMerchant);

            if(queryServer.Currency != queryMercahnt.Currency)
                throw new TerproException(ApiErrorCode.BusinessLogicDeny_CurrencyNotMatch);

            var queryServerMapping = ServerMerchantMappingRepository.GetByFilter(new QueryServerMerchantRelationRequestViewModel
                {
                    ServerIDs = serverID.ToListCollection(),
                    MerchantIDs = model.MerchantID.ToListCollection()
                })
                .FirstOrDefault();

            if(queryServerMapping == null)
            {
                ServerMerchantMappingRepository.Create(new ServerMerchantMappingEntity
                {
                    MerchantID = model.MerchantID,
                    ServerID = serverID,
                });
            }
            else
            {
                if(model.IsEnable)
                {

                }
                else
                {                    
                    ServerMerchantMappingRepository.Delete(queryServerMapping.Id);
                }
            }
        }



    }
}
