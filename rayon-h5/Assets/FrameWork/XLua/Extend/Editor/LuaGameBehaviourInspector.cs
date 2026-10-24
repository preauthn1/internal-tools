using UnityEditor;
using UnityEngine;

[CustomEditor(typeof(LuaGameBehaviour))]
public class LuaGameBehaviourInspector : Editor
{
    override public void OnInspectorGUI()
    {
        LuaGameBehaviour _src = target as LuaGameBehaviour;

        base.OnInspectorGUI();


        GUILayout.BeginHorizontal();

        GUILayout.FlexibleSpace();
        if (GUILayout.Button("重設", GUILayout.Width(230)))
            _src.InitSceneResolution();
        GUILayout.FlexibleSpace();

        GUILayout.EndHorizontal();
    }
}
