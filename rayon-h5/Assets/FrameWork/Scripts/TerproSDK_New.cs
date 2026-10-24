using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class TerproSDK_New : MonoBehaviour {

#if UNITY_EDITOR

#elif UNITY_ANDROID
    private static AndroidJavaClass JavaClass;
    private static AndroidJavaObject Activity;
#endif

    public static bool s_IsInit = false;

    public static void Init()
    {
        s_IsInit = true;

#if UNITY_EDITOR

#elif UNITY_ANDROID
        Debug.Log("*************** snail");

        JavaClass = new AndroidJavaClass("com.snail.sdk");
        //JavaClass = new AndroidJavaClass("com.target.sdk");
        if (JavaClass == null)
            Debug.Log("JavaClass is null");
        else
            Debug.Log("JavaClass is not null");

        //AndroidJavaClass cls = new AndroidJavaClass("com.onevcat.uniwebview.AndroidPlugin");
        AndroidJavaClass cls = new AndroidJavaClass("com.unity3d.player.UnityPlayer");
        if (cls == null)
        {
            Debug.Log("JavaClass is null");
        }
        else
        {
            Debug.Log("JavaClass is not null");

            Activity = cls.GetStatic<AndroidJavaObject>("currentActivity");
            if (Activity == null)
                Debug.Log("Activity is null");
            else
                Debug.Log("Activity is not null");
        }
#endif
    }

    private static string Call(string Name, bool IsActivity)
    {
        if (s_IsInit == false)
            Init();

        string result = "";
#if UNITY_EDITOR

#elif UNITY_ANDROID
        if (JavaClass != null)
        {
            if (IsActivity == false)
                result = JavaClass.CallStatic<string>(Name);
            else
            {
                if (IsActivity != null)
                    result = JavaClass.CallStatic<string>(Name, Activity);
                else
                    Debug.Log("IsActivity is null");
            }
        }
        else
        {
            Debug.Log("JavaClass is null");
        }
#endif

        return result;
    }

    public static string Check()
    {
        Debug.Log("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~1111111111111111");

        string Text = Call("Check", true);
        return Text;
    }

    public static string Check1()
    {
        Debug.Log("~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~2222222222");

        string Text = Call("Check1", true);
        return Text;
    }
}