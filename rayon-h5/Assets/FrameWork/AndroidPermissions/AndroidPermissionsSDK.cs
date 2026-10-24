using System.Collections;
using UnityEngine;

public class AndroidPermissionsSDK : MonoBehaviour
{
    public static bool s_IsCanUse = true;

    public const string STORAGE_PERMISSION = "android.permission.READ_EXTERNAL_STORAGE";
    public const string READ_PHONE_STATE = "android.permission.READ_PHONE_STATE";
    public const string ACCESS_CONTENT_PROVIDERS_EXTERNALLY = "android.permission.ACCESS_CONTENT_PROVIDERS_EXTERNALLY";
    public const string ACCESS_COARSE_LOCATION = "android.permission.ACCESS_COARSE_LOCATION";

    public void CheckPermissions1()
    {
        CheckPermissions(STORAGE_PERMISSION);
    }
    public void CheckPermissions2()
    {
        CheckPermissions(READ_PHONE_STATE);
    }
    public void CheckPermissions3()
    {
        CheckPermissions(ACCESS_CONTENT_PROVIDERS_EXTERNALLY);
    }
    public void CheckPermissions4()
    {
        CheckPermissions(ACCESS_COARSE_LOCATION);
    }

    public void GrantButtonPress1()
    {
        OnGrantButtonPress(STORAGE_PERMISSION, Result);
    }
    public void GrantButtonPress2()
    {
        OnGrantButtonPress(READ_PHONE_STATE, Result);
    }
    public void GrantButtonPress3()
    {
        OnGrantButtonPress(ACCESS_CONTENT_PROVIDERS_EXTERNALLY, Result);
    }
    public void GrantButtonPress4()
    {
        OnGrantButtonPress(ACCESS_COARSE_LOCATION, Result);
    }
    public void Result(int result)
    {

    }

    public static bool CheckPermissions(string Permission)
    {
        if (Application.platform != RuntimePlatform.Android)
        {
            return true;
        }

        bool Result = AndroidPermissionsManager.IsPermissionGranted(Permission);
        Debug.Log("### CS 確認權限 " + Permission + " " + Result);

        return Result;
    }

    public static void OnGrantButtonPress(string Permission, System.Action<int> Result)
    {
        Debug.Log("### CS 請求權限 " + Permission);
        AndroidPermissionsManager.RequestPermission(new[] { Permission }, new AndroidPermissionCallback(
            grantedPermission =>
            {
                // AndroidManifest 中必須有填要確認的permission
                // The permission was successfully granted, restart the change avatar routine
                if (!CheckPermissions(Permission))
                {
                    Debug.LogWarning("### CS AndroidManifest 檔中缺少 " + Permission);
                    Result(1);
                }
                else
                {
                    Debug.Log("### CS 請求權限成功...");
                    Result(0);
                }
            },
            // 許可被拒絕了
            deniedPermission =>
            {
                Debug.Log("### CS 請求權限被拒...");
                Result(2);
            })
        );
    }
}
