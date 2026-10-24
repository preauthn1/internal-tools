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
using target.DA.Repository;
using target.Utility;
using target.Utility.Helper;
using target.Models.ViewModels.Base;

namespace target.BL.Services
{

    public class PlayerJwtCheckService : BaseService, IPlayerJwtCheckService
    {
        
        public PlayerRepository PlayerRepository { get; set; }

        public LobbyTokenLogRepository LobbyTokenLogRepository { get; set; }


        /// <summary>
        /// 提供給 AuthorizeFilter呼叫，檢查DB紀錄的token是否已過期
        /// </summary>
        public bool CheckIsTokenExist(string jwtHash)
        {
            //測試環境先不需要後踢前
            if (SettingHelper.IsDev || SettingHelper.IsMSTest)
                return true;

            var queryPlayer = PlayerRepository.GetByToken(jwtHash);

            if (queryPlayer == null)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);

            return true;
        }


        /// <summary>
        /// 提供給 AuthorizeFilter呼叫，檢查DB紀錄的token是否已過期
        /// </summary>
        public bool CheckIsLobbyTokenExist(string tokenRawText)
        {
            CheckLobbyTokenExistAndReturnJwtTokenValue(tokenRawText);
            return true;
        }


        /// <summary>
        /// 解析lobbyToken
        /// </summary>
        public PlayerTokenJwtModel ParsingLobbyToken(string tokenRawText)
        {
            string tokenRawJwt = CheckLobbyTokenExistAndReturnJwtTokenValue(tokenRawText);

            var parseModel = MyJwtHelper.GetPlayerModelFromJwtString(tokenRawJwt);

            if (parseModel == null)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_NotValid);

            long nowUTC = DateTimeHelper.GetCurrentTimeTimestamp();

            if (nowUTC >= parseModel.TokenExpiredTime)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);

            return parseModel;
        }


        private string CheckLobbyTokenExistAndReturnJwtTokenValue(string tokenRawText)
        {
            if (string.IsNullOrWhiteSpace(tokenRawText))
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);

            LobbyTokenLogEntity queryToken = [REDACTED];

            if (tokenRawText.IsGuid())
            {
                Guid guid = Guid.Parse(tokenRawText);
                queryToken = [REDACTED];
            }
            else
            {
                queryToken = [REDACTED];
            }

            if (queryToken == null)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_NotValid);

            var now = DateTimeHelper.GetCurrentTime();

            if (now >= queryToken.ExpiredTime)
                throw new TerproException(ApiErrorCode.Authorize_AccessToken_Expired);

            return queryToken.RawJwtToken;
        }

    }
}
