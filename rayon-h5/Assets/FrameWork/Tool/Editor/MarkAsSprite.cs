using UnityEditor;
using UnityEngine;

public class MarkAsSprite : MonoBehaviour
{
    [MenuItem("Tools/选中图片加sprite_前缀", false, 201)]
    static private void AddSuffix()
    {
        Object[] selectObjs = Selection.objects;

        foreach (var obj in selectObjs)
        {
            if (!(obj is Texture2D) && obj.name.Contains("sprite_"))
                continue;

            string path = AssetDatabase.GetAssetPath(obj);
            AssetDatabase.RenameAsset(path, "sprite_" + obj.name);
        }
    }
}