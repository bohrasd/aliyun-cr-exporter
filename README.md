Aliyun ACR (Container Registry) Prometheus Collector
====================================================

WARNING: Not ready for production!

USAGE
-----
```
kubectl create secret generic aliyun-cr-secret --from-literal ALIYUN_REGION=cn-hangzhou --from-literal ALIYUN_AK=ABC --from-literal ALIYUN_SK=DEF -n your_namespace

kubectl apply -f ./k8s/ -n your_namespace
```

Metrics
-------

| Metric                                      | Value         | Description        |
| -------------                               | ------------- | -------            |
| aliyun_acr_namespace_info                   | 1             | 命名空间           |
| aliyun_acr_repo_info                        | 1             | 仓库信息           |
| aliyun_acr_repo_downloads                   | Gauge         | 仓库下载次数       |
| aliyun_acr_build_total                      | Gauge         | 总构建数量         |
| aliyun_acr_build_succeeded_duration_seconds | Gauge         | 成功构建的构建时长 |
| aliyun_acr_build_status                     | 1         | 构建状态           |

Configuration
-------------

| Name          | Description             |
| ------------- | -------------           |
| ALIYUN_REGION | 区域 like `cn-hangzhou` |
| ALIYUN_AK     | Access Key              |
| ALIYUN_SK     | Access Secret           |

TODO
----
+ ACR EE support
+ Multi tenent exporting
+ More metrics
+ More configuration
+ etc
