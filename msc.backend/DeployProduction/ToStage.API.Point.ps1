$now = Get-Date
$dateStr = $now.ToString("yyyyMMdd_HHmm")
echo $dateStr

$currentPath = Split-Path $script:MyInvocation.MyCommand.Path
echo $currentPath

$SHA1 = git rev-parse HEAD
echo $SHA1


cd $currentPath
cd ..
cd MSC.API.Point


$text_content =  -join($dateStr , "____" ,$SHA1);
Set-Content -Path "./AppData/AppVersion.txt" -Value $text_content -Force
& "C:\Program Files (x86)\Microsoft Visual Studio\2019\Enterprise\MSBuild\Current\Bin\MSBuild.exe" MSC.API.Point.csproj /p:Platform="AnyCPU" /p:DeployOnBuild=true /p:AllowUntrustedCertificate=True /p:PublishProfile="ogcwebapipoint-ogcWebAPIPointStage - Web Deploy.pubxml" /p:Password=[REDACTED]"