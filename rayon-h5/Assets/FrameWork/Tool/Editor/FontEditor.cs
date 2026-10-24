using UnityEngine;
using UnityEditor;
using example.invalid;

public class FontEditor
{
    public const string fontWords = "-1234567890.";      // 需要字体的文本内容，必须按照切图顺序排列

    const string CreateCustomFont = "Assets/生成图片字体";

    [MenuItem(CreateCustomFont, true, 62)]
    public static bool SetCurrentProject()
    {
        return Selection.activeObject is Texture2D;
    }

    [MenuItem(CreateCustomFont, false, 62)]
    static public void OpenBMFontMaker()
    {
        Object[] selectObjs = Selection.objects;

        foreach (var obj in selectObjs)
        {
            if (!(obj is Texture2D))
                continue;

            Texture2D mainTexture = obj as Texture2D;

            Object[] sprites = AssetDatabase.LoadAllAssetRepresentationsAtPath(AssetDatabase.GetAssetPath(obj));
            char[] fontChars = fontWords.ToCharArray();

            if (sprites.Length != fontChars.Length)
            {
                Debug.Log("图片数量和文字数量不相同");
                return;
            }

            CharacterInfo[] characterInfo = new CharacterInfo[sprites.Length];

            for (int i = 0; i < sprites.Length; i++)
            {
                Sprite sprite = sprites[i] as Sprite;

                Rect rect = sprite.rect;

                CharacterInfo info = new CharacterInfo();

                info.index = (int)fontChars[i];

                info.uvBottomLeft = new Vector2(rect.min.x / mainTexture.width, rect.min.y / mainTexture.height);
                info.uvBottomRight = new Vector2(rect.max.x / mainTexture.width, rect.min.y / mainTexture.height);
                info.uvTopLeft = new Vector2(rect.min.x / mainTexture.width, rect.max.y / mainTexture.height);
                info.uvTopRight = new Vector2(rect.max.x / mainTexture.width, rect.max.y / mainTexture.height);

                info.maxX = (int)rect.size.x;

                info.glyphHeight = info.glyphHeight * -1;
                info.minY = -1 * (int)rect.size.y;
                info.advance = (int)sprite.rect.width;
                //info.advance = advance;
                characterInfo[i] = info;
            }

            // Create custom font
            Font font = new Font();

            font.name = mainTexture.name;
            font.characterInfo = characterInfo;
            font.material = CreateFontMeterial(mainTexture);

            // Save font.
            string fontFile = GetExportFile(mainTexture) + ".fontsettings";
            //string tempFile = GetExportFile(mainTexture) + "_temp.fontsettings";

            //AssetDatabase.DeleteAsset(fontFile);
            AssetDatabase.CreateAsset(font, fontFile);
            //File.Copy(tempFile, fontFile, true);
            //AssetDatabase.DeleteAsset(tempFile);
        }
    }

    static private string GetExportFile(Texture2D texture)
    {
        string rootPath = Path.GetDirectoryName(AssetDatabase.GetAssetPath(texture));
        string exportPath = rootPath + "/" + Path.GetFileNameWithoutExtension(texture.name);
        return exportPath;
    }

    // Create material
    static private Material CreateFontMeterial(Texture2D fontTexture)
    {
        Shader shader = Shader.Find("UI/Default");
        Material material = new Material(shader);
        material.mainTexture = fontTexture;

        string materialFile = GetExportFile(fontTexture) + ".mat";
        AssetDatabase.DeleteAsset(materialFile);
        AssetDatabase.CreateAsset(material, GetExportFile(fontTexture) + ".mat");

        return material;
    }
}
