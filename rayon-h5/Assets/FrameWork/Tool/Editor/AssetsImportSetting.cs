using UnityEngine;
using UnityEditor;
using example.invalid;

public class ImportTextureSetting : AssetPostprocessor
{
    // PSD导入到UI的地方已经处理过格式问题，其实不需要在导入图片的时候再强制设置一次，这里视情况决定是否开放
    void OnPreprocessTexture()
    { 
        if (assetPath.Contains("/HeadIcon/"))
            return;

        if (assetPath.Contains("/UI/"))
        {
            string projectName = string.Empty;

            if (assetPath.Contains("Assets/Hall"))
                projectName = "Hall";

            if (assetPath.Contains("Assets/Game"))
                projectName = assetPath.Split('/')[2];

            string AtlasName = new DirectoryInfo(Path.GetDirectoryName(assetPath)).Name;
            TextureImporter textureImporter = assetImporter as TextureImporter;

            if (textureImporter.textureType == TextureImporterType.Sprite)
            {
                textureImporter.textureCompression = TextureImporterCompression.Compressed;
                textureImporter.spritePackingTag = projectName.ToLower() + "_" + AtlasName;

                if (textureImporter.spriteImportMode == SpriteImportMode.None)
                    textureImporter.spriteImportMode = SpriteImportMode.Single;

                textureImporter.wrapMode = TextureWrapMode.Clamp;
                textureImporter.mipmapEnabled = false;
                textureImporter.alphaIsTransparency = true;
            }

            if (textureImporter.textureType == TextureImporterType.Default)
            {
                textureImporter.textureCompression = TextureImporterCompression.Compressed;
                textureImporter.mipmapEnabled = false;
                // textureImporter.npotScale = TextureImporterNPOTScale.ToNearest;
            }

            /*
            TextureImporterPlatformSettings iOSPlatformSettings = textureImporter.GetPlatformTextureSettings("iPhone");
            iOSPlatformSettings.overridden = true;
            //iOSPlatformSettings.format = TextureImporterFormat.PVRTC_RGBA2;
            iOSPlatformSettings.format = TextureImporterFormat.RGBA32;
            textureImporter.SetPlatformTextureSettings(iOSPlatformSettings);
            */

			TextureImporterPlatformSettings AndroidPlatformSettings = textureImporter.GetPlatformTextureSettings("Android");
			AndroidPlatformSettings.overridden = true;
			AndroidPlatformSettings.format = TextureImporterFormat.ETC2_RGBA8;
			textureImporter.SetPlatformTextureSettings(AndroidPlatformSettings);
		}
    }

    private static void OnPostprocessAllAssets(string[] importedAssets, string[] deletedAssets, string[] movedAssets, string[] movedFromPath)
    {
        foreach (string move in movedAssets)
            AssetDatabase.ImportAsset(move);
    }
}