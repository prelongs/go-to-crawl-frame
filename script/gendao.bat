REM root path run this bat
REM -c支持选config*.yaml bug多，因此只支持去改./hack/config.yaml中的gen.dao节点
REM 保持统一性，只用gf.exe，Linux下的gf找不到对应版本了，除非整体升最新稳定版
./gf.exe gen dao