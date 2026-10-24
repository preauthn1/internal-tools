/**********************************************************
// FileName : DebugFilter.cs
// Reference: http://example.invalid/2015/09/unity3d-debuglog.html
**********************************************************/
#if NOLOG

using System.Diagnostics;
using UnityEngine;

public static class Debug
{
    [Conditional("NOLOG")]
    public static void Break() { }
    [Conditional("NOLOG")]
    public static void ClearDeveloperConsole() { }
    [Conditional("NOLOG")]
    public static void DebugBreak() { }

    [Conditional("NOLOG")]
    public static void DrawLine(Vector3 start, Vector3 end) { }
    [Conditional("NOLOG")]
    public static void DrawLine(Vector3 start, Vector3 end, Color color) { }
    [Conditional("NOLOG")]
    public static void DrawLine(Vector3 start, Vector3 end, Color color, float duration) { }
    [Conditional("NOLOG")]
    public static void DrawLine(Vector3 start, Vector3 end, Color color, float duration, bool depthTest) { }

    [Conditional("NOLOG")]
    public static void DrawRay(Vector3 start, Vector3 dir) { }
    [Conditional("NOLOG")]
    public static void DrawRay(Vector3 start, Vector3 dir, Color color) { }
    [Conditional("NOLOG")]
    public static void DrawRay(Vector3 start, Vector3 dir, Color color, float duration) { }
    [Conditional("NOLOG")]
    public static void DrawRay(Vector3 start, Vector3 dir, Color color, float duration, bool depthTest) { }

    [Conditional("NOLOG")]
    public static void Log(object message) { }
    [Conditional("NOLOG")]
    public static void Log(object message, Object context) { }
    [Conditional("NOLOG")]
    public static void LogFormat(string format, params object[] args) { }
    [Conditional("NOLOG")]
    public static void LogFormat(Object context, string format, params object[] args) { }

    [Conditional("NOLOG")]
    public static void LogWarning(object message) { }
    [Conditional("NOLOG")]
    public static void LogWarning(object message, Object context) { }
    [Conditional("NOLOG")]
    public static void LogWarningFormat(string format, params object[] args) { }
    [Conditional("NOLOG")]
    public static void LogWarningFormat(Object context, string format, params object[] args) { }

    [Conditional("NOLOG")]
    public static void LogError(object message) { }
    [Conditional("NOLOG")]
    public static void LogError(object message, Object context) { }
    [Conditional("NOLOG")]
    public static void LogErrorFormat(string format, params object[] args) { }
    [Conditional("NOLOG")]
    public static void LogErrorFormat(Object context, string format, params object[] args) { }

    [Conditional("NOLOG")]
    public static void LogException(System.Exception exception) { }
    [Conditional("NOLOG")]
    public static void LogException(System.Exception exception, Object context) { }
}
#endif