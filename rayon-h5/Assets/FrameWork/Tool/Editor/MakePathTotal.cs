using UnityEditor;
using UnityEngine;


public static class MakePathTotal
{
    const string MethodName = "GameObject/Tool/By/Path Total";

    [MenuItem(MethodName, true, 1)]
    public static bool BuildScriptValidation()
    {
        return Selection.activeObject is GameObject;
    }

    [MenuItem(MethodName, false, 1)]
    public static void BuildScript()
    {
        MakePathMeshRenderer.BuildScript();
        MakePathSkinnedMeshRenderer.BuildScript();
        MakePathAnimator.BuildScript();
        MakePathParticleSystem.BuildScript();
    }
}
