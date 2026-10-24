using example.invalid;
using UnityEditor;
using UnityEngine;

[CustomEditor(typeof(GameConfig))]
public class GameConfigInspector : Editor
{
    private bool m_ChangeSkin = false;
    public override void OnInspectorGUI()
    {
        GameConfig m_Config = (GameConfig)target;
           
        GUI.Label(new Rect(0, 50, 200, 30), "版本号"); 
        m_Config.Version = GUI.TextField(new Rect(60, 50, 300, 20), m_Config.Version);


        GUI.Label(new Rect(0, 80, 200, 30), "样版");
        if (m_ChangeSkin == true)
            m_Config.Skin = GUI.TextField(new Rect(60, 80, 300, 20), m_Config.Skin);

        if (GUI.changed)
        {
            m_ChangeSkin = false;
            EditorUtility.SetDirty(m_Config);

            serializedObject.Update();
            serializedObject.ApplyModifiedProperties();
        }

        if (m_ChangeSkin == false && GUI.Button(new Rect(60, 80, 300, 20), m_Config.Skin))
            m_ChangeSkin = true;
    }
} 