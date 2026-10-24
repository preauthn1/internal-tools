    using System.Collections;
    using System.Collections.Generic;
    using example.invalid;
    using UnityEditor;
    using UnityEditor.SceneManagement;
    using UnityEngine;
    using UnityEngine.UI;

public class FindAllSpriteReferences
{
    [MenuItem("Assets/资源优化/查找选中图片引用", priority = 22)]
    static private void ExcuteMenu()
    {
        UnityEngine.Object[] selectObjs = Selection.objects;

        for (int m = 0; m < selectObjs.Length; m++)
        {
            if (!(selectObjs[m] is Texture2D))
                return;

            string assetPath = AssetDatabase.GetAssetPath(selectObjs[m]);
            Sprite obj = AssetDatabase.LoadAssetAtPath<Sprite>(assetPath);

            try
            {
                string[] paths = Directory.GetFiles(AppDefine.CurrentProjectPath, "*.prefab", SearchOption.AllDirectories);

                for (int i = 0; i < paths.Length; i++)
                {
                    GameObject prefabObj = AssetDatabase.LoadAssetAtPath(EditorFileUtility.GetRelativeAssetPath(paths[i]), typeof(GameObject)) as GameObject;

                    EditorUtility.DisplayProgressBar("查找sprite引用", string.Format("查找{0}在预制体{1}中的引用...", obj.name, prefabObj.name), 1f * i / paths.Length);

                    // ------------------Image-------------------------
                    List<Image> imageList = new List<Image>();

                    (prefabObj as GameObject).GetComponentsInChildren<Image>(true, imageList);  // 获取所有Text组件

                    for (int j = 0; j < imageList.Count; j++)
                    {
                        if (imageList[j].sprite == obj)
                            Debug.Log(string.Format("<color=#3FFF20>预制体</color> {0} 的 {1} 引用了{2}。", prefabObj.name, imageList[j].name, obj.name));
                    }

                    // ------------------Sprite------------------------
                    List<SpriteRenderer> spriteList = new List<SpriteRenderer>();

                    (prefabObj as GameObject).GetComponentsInChildren<SpriteRenderer>(true, spriteList);  // 获取所有Text组件

                    for (int j = 0; j < spriteList.Count; j++)
                    {
                        if (spriteList[j].sprite == obj)
                            Debug.Log(string.Format("<color=#3FFF20>预制体</color> {0} 的 {1} 引用了{2}。", prefabObj.name, spriteList[j].name, obj.name));
                    }
                }

                string[] scenePaths = Directory.GetFiles(AppDefine.CurrentProjectPath, "*.unity", SearchOption.AllDirectories);

                for (int i = 0; i < scenePaths.Length; i++)
                {
                    string senceName = EditorFileUtility.GetFileName(scenePaths[i].Replace(@"\", @"/"));

                    EditorUtility.DisplayProgressBar("查找sprite引用", string.Format("查找{0}在场景{1}中的引用...", obj.name, senceName), 1f * i / scenePaths.Length);

                    //打开场景
                    EditorSceneManager.OpenScene(scenePaths[i]);

                    //获取场景中的所有游戏对象
                    GameObject[] gos = (GameObject[])GameObject.FindObjectsOfType(typeof(GameObject));

                    foreach (GameObject go in gos)
                    {
                        //判断GameObject是否为一个Prefab的引用
                        if (PrefabUtility.GetPrefabType(go) != PrefabType.PrefabInstance)
                        {
                            // ------------------Image-------------------------
                            List<Image> imageList = new List<Image>();

                            (go as GameObject).GetComponentsInChildren<Image>(true, imageList);  // 获取所有Text组件

                            for (int j = 0; j < imageList.Count; j++)
                            {
                                if (imageList[j].sprite == obj)
                                    Debug.Log(string.Format("<color=#ff0000>场景</color> {0} 的 {1} 引用了{2}。", senceName, imageList[j].name, obj.name));
                            }

                            // ------------------Sprite------------------------
                            List<SpriteRenderer> spriteList = new List<SpriteRenderer>();

                            (go as GameObject).GetComponentsInChildren<SpriteRenderer>(true, spriteList);  // 获取所有Text组件

                            for (int j = 0; j < spriteList.Count; j++)
                            {
                                if (spriteList[j].sprite == obj)
                                    Debug.Log(string.Format("<color=#ff0000>场景</color> {0} 的 {1} 引用了{2}。", senceName, spriteList[j].name, obj.name));
                            }
                        }
                    }
                }
            }
            catch (System.Exception ex)
            {
                throw ex;
            }
            finally
            {
                EditorUtility.ClearProgressBar();
            }
        }
    }

    [MenuItem("Assets/资源优化/抖抖抖色RGBA4444", priority = 23)]
    static private void Ditherrgba4444Menu()
    {
        UnityEngine.Object[] selectObjs = Selection.objects;

        for (int m = 0; m < selectObjs.Length; m++)
        {
            if (!(selectObjs[m] is Texture2D))
                return;

            string assetPath = AssetDatabase.GetAssetPath(selectObjs[m]);

            TextureImporter textureImporter = (TextureImporter)TextureImporter.GetAtPath(assetPath);
#if UNITY_ANDROID
            TextureImporterPlatformSettings AndroidPlatformSettings = textureImporter.GetPlatformTextureSettings("Android");
            AndroidPlatformSettings.overridden = true;
            AndroidPlatformSettings.format = TextureImporterFormat.RGBA32;
            textureImporter.SetPlatformTextureSettings(AndroidPlatformSettings);
#endif
#if UNITY_IOS
                   TextureImporterPlatformSettings iOSPlatformSettings = textureImporter.GetPlatformTextureSettings("iPhone");
                   iOSPlatformSettings.overridden = true;
                   iOSPlatformSettings.format = TextureImporterFormat.RGBA32;
                   textureImporter.SetPlatformTextureSettings(iOSPlatformSettings);
#endif
            AssetDatabase.ImportAsset(assetPath, ImportAssetOptions.ForceUpdate);

            Texture2D texture = AssetDatabase.LoadAssetAtPath<Texture2D>(assetPath);

            try
            {
                var texw = texture.width;
                var texh = texture.height;

                var pixels = texture.GetPixels();
                var offs = 0;

                var k1Per15 = 1.0f / 15.0f;
                var k1Per16 = 1.0f / 16.0f;
                var k3Per16 = 3.0f / 16.0f;
                var k5Per16 = 5.0f / 16.0f;
                var k7Per16 = 7.0f / 16.0f;

                for (var y = 0; y < texh; y++)
                {
                    for (var x = 0; x < texw; x++)
                    {
                        float a = pixels[offs].a;
                        float r = pixels[offs].r;
                        float g = pixels[offs].g;
                        float b = pixels[offs].b;

                        var a2 = Mathf.Clamp01(Mathf.Floor(a * 16) * k1Per15);
                        var r2 = Mathf.Clamp01(Mathf.Floor(r * 16) * k1Per15);
                        var g2 = Mathf.Clamp01(Mathf.Floor(g * 16) * k1Per15);
                        var b2 = Mathf.Clamp01(Mathf.Floor(b * 16) * k1Per15);

                        var ae = a - a2;
                        var re = r - r2;
                        var ge = g - g2;
                        var be = b - b2;

                        pixels[offs].a = a2;
                        pixels[offs].r = r2;
                        pixels[offs].g = g2;
                        pixels[offs].b = b2;

                        var n1 = offs + 1;   // (x+1,y)
                        var n2 = offs + texw - 1; // (x-1 , y+1)
                        var n3 = offs + texw;  // (x, y+1)
                        var n4 = offs + texw + 1; // (x+1 , y+1)

                        if (x < texw - 1)
                        {
                            pixels[n1].a += ae * k7Per16;
                            pixels[n1].r += re * k7Per16;
                            pixels[n1].g += ge * k7Per16;
                            pixels[n1].b += be * k7Per16;
                        }

                        if (y < texh - 1)
                        {
                            pixels[n3].a += ae * k5Per16;
                            pixels[n3].r += re * k5Per16;
                            pixels[n3].g += ge * k5Per16;
                            pixels[n3].b += be * k5Per16;

                            if (x > 0)
                            {
                                pixels[n2].a += ae * k3Per16;
                                pixels[n2].r += re * k3Per16;
                                pixels[n2].g += ge * k3Per16;
                                pixels[n2].b += be * k3Per16;
                            }

                            if (x < texw - 1)
                            {
                                pixels[n4].a += ae * k1Per16;
                                pixels[n4].r += re * k1Per16;
                                pixels[n4].g += ge * k1Per16;
                                pixels[n4].b += be * k1Per16;
                            }
                        }

                        offs++;
                    }
                }

                texture.SetPixels(pixels);
                byte[] _bytes = texture.EncodeToPNG();
                example.invalid.File.WriteAllBytes(assetPath, _bytes);
                Debug.Log(_bytes.Length / 1024 + "Kb was saved as: " + assetPath);


            }
            catch (System.Exception ex)
            {
                throw ex;
            }
            finally
            {
                EditorUtility.ClearProgressBar();
            }
        }
    }
    [MenuItem("Assets/资源优化/抖抖抖色RGB565", priority = 24)]
    static private void Dithergba565Menu()
    {
        UnityEngine.Object[] selectObjs = Selection.objects;

        for (int m = 0; m < selectObjs.Length; m++)
        {
            if (!(selectObjs[m] is Texture2D))
                return;

            string assetPath = AssetDatabase.GetAssetPath(selectObjs[m]);

            TextureImporter textureImporter = (TextureImporter)TextureImporter.GetAtPath(assetPath);
#if UNITY_ANDROID
            TextureImporterPlatformSettings AndroidPlatformSettings = textureImporter.GetPlatformTextureSettings("Android");
            AndroidPlatformSettings.overridden = true;
            AndroidPlatformSettings.format = TextureImporterFormat.RGBA32;
            textureImporter.SetPlatformTextureSettings(AndroidPlatformSettings);
#endif
#if UNITY_IOS
                   TextureImporterPlatformSettings iOSPlatformSettings = textureImporter.GetPlatformTextureSettings("iPhone");
                   iOSPlatformSettings.overridden = true;
                   iOSPlatformSettings.format = TextureImporterFormat.RGBA32;
                   textureImporter.SetPlatformTextureSettings(iOSPlatformSettings);
#endif
            AssetDatabase.ImportAsset(assetPath, ImportAssetOptions.ForceUpdate);
            Texture2D texture = AssetDatabase.LoadAssetAtPath<Texture2D>(assetPath);
            try
            {
                var texw = texture.width;
                var texh = texture.height;

                var pixels = texture.GetPixels();
                var offs = 0;

                var k1Per31 = 1.0f / 31.0f;

                var k1Per32 = 1.0f / 32.0f;
                var k5Per32 = 5.0f / 32.0f;
                var k11Per32 = 11.0f / 32.0f;
                var k15Per32 = 15.0f / 32.0f;

                var k1Per63 = 1.0f / 63.0f;

                var k3Per64 = 3.0f / 64.0f;
                var k11Per64 = 11.0f / 64.0f;
                var k21Per64 = 21.0f / 64.0f;
                var k29Per64 = 29.0f / 64.0f;

                var k_r = 32; //R&B压缩到5位，所以取2的5次方
                var k_g = 64; //G压缩到6位，所以取2的6次方

                for (var y = 0; y < texh; y++)
                {
                    for (var x = 0; x < texw; x++)
                    {
                        float r = pixels[offs].r;
                        float g = pixels[offs].g;
                        float b = pixels[offs].b;

                        var r2 = Mathf.Clamp01(Mathf.Floor(r * k_r) * k1Per31);
                        var g2 = Mathf.Clamp01(Mathf.Floor(g * k_g) * k1Per63);
                        var b2 = Mathf.Clamp01(Mathf.Floor(b * k_r) * k1Per31);

                        var re = r - r2;
                        var ge = g - g2;
                        var be = b - b2;

                        var n1 = offs + 1;
                        var n2 = offs + texw - 1;
                        var n3 = offs + texw;
                        var n4 = offs + texw + 1;

                        if (x < texw - 1)
                        {
                            pixels[n1].r += re * k15Per32;
                            pixels[n1].g += ge * k29Per64;
                            pixels[n1].b += be * k15Per32;
                        }

                        if (y < texh - 1)
                        {
                            pixels[n3].r += re * k11Per32;
                            pixels[n3].g += ge * k21Per64;
                            pixels[n3].b += be * k11Per32;

                            if (x > 0)
                            {
                                pixels[n2].r += re * k5Per32;
                                pixels[n2].g += ge * k11Per64;
                                pixels[n2].b += be * k5Per32;
                            }

                            if (x < texw - 1)
                            {
                                pixels[n4].r += re * k1Per32;
                                pixels[n4].g += ge * k3Per64;
                                pixels[n4].b += be * k1Per32;
                            }
                        }

                        pixels[offs].r = r2;
                        pixels[offs].g = g2;
                        pixels[offs].b = b2;

                        offs++;
                    }
                }

                texture.SetPixels(pixels);
                EditorUtility.CompressTexture(texture, TextureFormat.RGBA32,(int)UnityEngine.TextureCompressionQuality.Best);
                byte[] _bytes = texture.EncodeToPNG();
                example.invalid.File.WriteAllBytes(assetPath, _bytes);
                Debug.Log(_bytes.Length / 1024 + "Kb was saved as: " + assetPath);


            }
            catch (System.Exception ex)
            {
                throw ex;
            }
            finally
            {
                EditorUtility.ClearProgressBar();
            }
        }
    }

  
    }
