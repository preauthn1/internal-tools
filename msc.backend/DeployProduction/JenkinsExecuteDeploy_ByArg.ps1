$projectName = $args[0]; #參數請傳入專案名稱 e.g. MSC.Web.Point
$webDeployFile = $args[1]; #參數請傳入發布設定檔
$webDeployPassword = [REDACTED]; #參數請傳入發布密碼

if ($projectName -eq $null -Or $webDeployFile -eq $null -Or $webDeployPassword -eq $null) 
{
    echo 'null!';
    exit;
}

$now = Get-Date
$dateStr = $now.ToString("yyyyMMdd_HHmm")
echo $dateStr

$currentPath = Split-Path $script:MyInvocation.MyCommand.Path
echo $currentPath

$SHA1 = git rev-parse HEAD
$SHA1_s = $SHA1.substring(0, 8)

echo $SHA1

cd $currentPath
cd ..
cd $projectName

$text_content =  -join($dateStr , "____" ,$SHA1);
Set-Content -Path "./AppData/AppVersion.txt" -Value $text_content -Force
& "C:\Program Files (x86)\Microsoft Visual Studio\2019\Community\MSBuild\Current\Bin\MSBuild.exe" "$projectName.csproj" /p:Platform="AnyCPU" /p:DeployOnBuild=true /p:AllowUntrustedCertificate=True /p:PublishProfile="$webDeployFile" /p:Password=[REDACTED]"

[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
$postParams = "token=[REDACTED] ``$SHA1_s``"
Invoke-WebRequest -Uri https://example.invalid/api/chat.postMessage -Method POST -Body $postParams 