using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using System;


public static class Protocol
{
    public static bool showProtocolName = true;

    // 封包類型，封包值
    public static Dictionary<string, Array> s_ProtocolEnum = null;

    private static void Init()
    {
        if (s_ProtocolEnum != null)
            return;

        s_ProtocolEnum = new Dictionary<string, Array>();

        s_ProtocolEnum.Add(((int)MessageBClassID.Common).ToString(),        Enum.GetValues(typeof(PlatformCommonClassID)));
        s_ProtocolEnum.Add(((int)MessageBClassID.Hall).ToString(),          Enum.GetValues(typeof(HallMsgClassID)));
        s_ProtocolEnum.Add(((int)MessageBClassID.ServerMgm).ToString(),     null);
        s_ProtocolEnum.Add(((int)MessageBClassID.DBServer).ToString(),      Enum.GetValues(typeof(DBServerClassID)));
        s_ProtocolEnum.Add(((int)MessageBClassID.GameRoom).ToString(),      Enum.GetValues(typeof(GameRoomClassID)));
        s_ProtocolEnum.Add(((int)MessageBClassID.Game).ToString(),          Enum.GetValues(typeof(GameMessageClassID)));
        s_ProtocolEnum.Add(((int)MessageBClassID.NotifyServer).ToString(),  Enum.GetValues(typeof(NotifyServerClassID)));


        // Type 10001
        // 各個小遊戲
        s_ProtocolEnum.Add("bjl",       Enum.GetValues(typeof(BJL_GameMessageClassID)));
        s_ProtocolEnum.Add("by",        Enum.GetValues(typeof(BY_GameMessageClassID)));
        s_ProtocolEnum.Add("ddz",       Enum.GetValues(typeof(DDZ_GameMessageClassID)));
        s_ProtocolEnum.Add("dfdc",      Enum.GetValues(typeof(DFDC_GameMessageClassID)));
        s_ProtocolEnum.Add("dzpk",      Enum.GetValues(typeof(DZPK_GameMessageClassID)));
        s_ProtocolEnum.Add("lhd",       Enum.GetValues(typeof(LHD_GameMessageClassID)));
        s_ProtocolEnum.Add("lhj",       Enum.GetValues(typeof(LHJ_GameMessageClassID)));
        s_ProtocolEnum.Add("qzesyd",    Enum.GetValues(typeof(QZ21_GameMessageClassID)));
        s_ProtocolEnum.Add("qznn",      Enum.GetValues(typeof(QZNN_GameMessageClassID)));
        s_ProtocolEnum.Add("sangong",   Enum.GetValues(typeof(SG_GameMessageClassID)));
        s_ProtocolEnum.Add("scmj",      Enum.GetValues(typeof(SCMJ_GameMessageClassID)));
        s_ProtocolEnum.Add("sgj",       Enum.GetValues(typeof(SGJ_GameMessageClassID)));
        s_ProtocolEnum.Add("srnn",      Enum.GetValues(typeof(SRNN_GameMessageClassID)));
        s_ProtocolEnum.Add("sss",       Enum.GetValues(typeof(SSS_GameMessageClassID)));
        s_ProtocolEnum.Add("tbnn",      Enum.GetValues(typeof(TBNN_GameMessageClassID)));
        s_ProtocolEnum.Add("ttz",       Enum.GetValues(typeof(TTZ_GameMessageClassID)));
        s_ProtocolEnum.Add("wzp",       Enum.GetValues(typeof(WZP_GameMessageClassID)));
        s_ProtocolEnum.Add("xyx",       Enum.GetValues(typeof(XYX_GameMessageClassID)));
        s_ProtocolEnum.Add("zjh",       Enum.GetValues(typeof(ZJH_GameMessageClassID)));
        s_ProtocolEnum.Add("brnn",      Enum.GetValues(typeof(BRNN_GameMessageClassID)));
        s_ProtocolEnum.Add("hhdz",      Enum.GetValues(typeof(HHDZ_GameMessageClassID)));
        s_ProtocolEnum.Add("hhnn",      Enum.GetValues(typeof(HHNN_GameMessageClassID)));
        s_ProtocolEnum.Add("ermj",      Enum.GetValues(typeof(ERMJ_GameMessageClassID)));
        s_ProtocolEnum.Add("codyqzesyd",Enum.GetValues(typeof(CD21_GameMessageClassID)));
        s_ProtocolEnum.Add("fkeb",      Enum.GetValues(typeof(FKEB_GameMessageClassID)));
        s_ProtocolEnum.Add("tgw",       Enum.GetValues(typeof(TGW_GameMessageClassID)));
    }


    public static string GetTypeString(int Type)
    {
        if (showProtocolName == false)
            return "";

        foreach (MessageBClassID val in Enum.GetValues(typeof(MessageBClassID)))
            if (Type == (int)val)
                return val.ToString();

        return "";
    }

    public static string GetProtocolString(int Type, int Protocol)
    {
        if (showProtocolName == false)
            return "";

        Init();

        // 如果是遊戲封包，先找各遊戲中的封包名稱
        if (Type == (int)MessageBClassID.Game)
        {
            string _GameValue = GetByString(GameManager.Instance.GetCurrentEnterGameName(), Protocol);

            // 有找到封包
            if (string.IsNullOrEmpty(_GameValue) == false)
                return _GameValue;
        }
       
        string _Value = GetByString(Type.ToString(), Protocol);

        return _Value;
    }

    private static string GetByString(string Type, int Protocol)
    {
        if (s_ProtocolEnum.ContainsKey(Type) == false)
            return "";

        Array _List = s_ProtocolEnum[Type];
        if (_List == null ||
            _List.Length == 0)
        {
            return "";
        }

        foreach (var val in _List)
        {
            if (Protocol == (int)val)
                return val.ToString();
        }
        return "";
    }
}
