Aliyun ACR (Container Registry) Prometheus Collector
====================================================

WARNING: Not ready for production!

TODO
----
+ Collect cache
+ ACR EE support
+ More metrics
+ More configuration
+ etc

Metrics
-------

| Metric                                      | Value         | Description        |
| -------------                               | ------------- | -------            |
| aliyun_acr_namespace_info                   | 1             | 命名空间           |
| aliyun_acr_repo_info                        | 1             | 仓库信息           |
| aliyun_acr_repo_downloads                   | Gauge         | 仓库下载次数       |
| aliyun_acr_build_total                      | Gauge         | 总构建数量         |
| aliyun_acr_build_succeeded_duration_seconds | Gauge         | 成功构建的构建时长 |
| aliyun_acr_build_status                     | Gauge         | 构建状态           |

Configuration
-------------

| Name          | Description             |
| ------------- | -------------           |
| ALIYUN_REGION | 区域 like `cn-hangzhou` |
| ALIYUN_AK     | Access Key              |
| ALIYUN_SK     | Access Secret           |

