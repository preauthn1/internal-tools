using System;
using System.Collections.Generic;
using XLua;

/// <summary>
/// 网页插件 xlua自定义导出
/// </summary>
public static class XLuaWebViewExport
{

    /// <summary>
    /// dotween的扩展方法在lua中调用
    /// </summary>
    [LuaCallCSharp]
    [ReflectionUse]
    public static List<Type> webView_lua_call_cs_list = new List<Type>()
    { 
#if UNITY_IOS || UNITY_ANDROID || UNITY_WP8
        typeof(UniWebView),
        typeof(UniWebViewMessage),
        typeof(UniWebViewOrientation),
        typeof(UniWebViewEdgeInsets),
        typeof(UniWebViewHelper),
        typeof(Dictionary<string, string>)          
#endif
    };

    //C#静态调用Lua的配置（包括事件的原型），仅可以配delegate，interface
    [CSharpCallLua]
    public static List<Type> webView_cs_call_lua_List = new List<Type>()
    { 
#if UNITY_IOS || UNITY_ANDROID || UNITY_WP8 
        typeof(UniWebView.ReceivedMessageDelegate),
        typeof(UniWebView.LoadCompleteDelegate),
        //typeof(UniWebView.EvalJavaScriptFinishedDelegate),
        //typeof(UniWebView.WebViewShouldCloseDelegate)      
#endif
    };
}