# 对象存储与数据存储裸露（OSS / S3 / MinIO / MySQL / MongoDB）

- 来源：FJNB5/fujiang-pojia@cf086fc `plugin/knowledge/gray-market/object-storage-datastore-exposure.md`
- 分类：gray-market ｜ 严重度：high ｜ CVSS：8.6 ｜ CVE：-
- 产品：阿里云 OSS / 腾讯云 COS / AWS S3 / MinIO / MySQL / MongoDB
- 影响版本：任意
- 端口：80,443,9000,3306,27017
- 关键词：oss, s3, bucket, minio, mongodb, mysql, 弱口令, 列举, kyc
- 触发点：-

## 定位要点

- 组件指纹：按 `阿里云 OSS / 腾讯云 COS / AWS S3 / MinIO / MySQL / MongoDB` 的版本声明判断候选面（任意），端口 80,443,9000,3306,27017 优先。
- 信号来源：oss、s3、bucket、minio、mongodb、mysql、弱口令、列举。
- 触发条件：-。
- 判断口径：本条为**未实测**的索引条目，命中只代表候选，需按实际回显与版本复核后升级为结论。

## 原始资料

https://github.com/FJNB5/fujiang-pojia/blob/HEAD/plugin/knowledge/gray-market/object-storage-datastore-exposure.md
