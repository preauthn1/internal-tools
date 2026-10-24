using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class CustomDefine : MonoBehaviour
{
    public int ServerOption;
    public int PlatformID;
    public int StaticSiteID;
}

public class JiuLiaoInfo
{
    public int Index;
    public string JiuLiaoId = "";
    public string JiuLiaoSecrety = "";
    public string JiuLiaoName = "";
    public string DownloadUrl = "";
}

public class BuildInfo
{
    public int Index;
    public string CompanyName = "";
    public string ProductName = "";
    public string BundleID = "";
    public string AppName = "";
    public string AppleID = "";
    public int ServerOption;
    public int Branch = 0;
    public int PlatformID;
    public int SiteID;
    public int StaticSiteID;
    public bool FlagShare = true;
    public string IconName = "";
    public int Version = 0;
    public int Agent = 0;
    public int NotExport = 0;

    public BuildInfo Copy()
    {
        BuildInfo info = new BuildInfo();
        info.Index = Index;
        info.CompanyName = CompanyName;
        info.ProductName = ProductName;
        info.BundleID = BundleID;
        info.AppName = AppName;
        info.AppleID = AppleID;
        info.ServerOption = ServerOption;
        info.Branch = Branch;
        info.PlatformID = PlatformID;
        info.SiteID = SiteID;
        info.StaticSiteID = StaticSiteID;
        info.FlagShare = FlagShare;
        info.IconName = IconName;

        return info;
    }

    public string GetAppKey()
    {
        string _key = BundleID.Replace(".", "_");
        _key += StaticSiteID;
        return _key;
    }

    public string GetKeyWord()
    {
        string _key = Index.ToString("00") + " " + StaticSiteID.ToString("0000") + " " + ProductName;

        return _key;
    }
}

