module studyGit

go 1.15

require (
	git.in.codoon.com/backend/common v0.0.0-20210702104524-7c7f4f773c49
	git.in.codoon.com/backend/framework v0.0.0-20210705055424-a9a890125198 // indirect
	git.in.codoon.com/backend/jsonrpc v0.0.0-20200423070224-aecd90a7ae72 // indirect
	git.in.codoon.com/backend/rpc v0.0.0-20200423070202-931c7a68cfbd // indirect
	git.in.codoon.com/backend/serverapi v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/check_dirty v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/codoon_race v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/codoon_tieba_v2_api v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/dataserver v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/equipment_api v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/feedserver_api v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/frequency_control v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/give_medal v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/gomore-api v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/gps_api v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/hot_zone_server v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/login_in v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/mallapi v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/member_api v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/mns_api v0.0.0-20210705054134-b7c250f1178a // indirect
	git.in.codoon.com/backend/serverapi/msg_api v0.0.0-20210705054134-b7c250f1178a // indirect
	git.in.codoon.com/backend/serverapi/pb_api v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/relation_chain/user_relation v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/share_center v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/sports_exam v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/sportslevel v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/training_plan_base_api v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/training_plan_v2_api v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/ucenter v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/userinfo_center v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/backend/serverapi/userprofile v0.0.0-20210705054134-b7c250f1178a
	git.in.codoon.com/codoon_mall/alert v0.0.0-20170707060403-1de75e52d6b6 // indirect
	git.in.codoon.com/codoon_mall/http-recover v0.0.0-20200421095720-5ac16d5b8448 // indirect
	git.in.codoon.com/codoon_mall/xcommon v0.0.0-20210419063838-0ea9c8c4a1c5
	git.in.codoon.com/lib/encrypt v0.0.0-20200402054539-4aec4e43ed67 // indirect
	git.in.codoon.com/third/aliyun-oss-go-sdk/oss v0.0.0-20200331102045-b4cf4109c8b5 // indirect
	git.in.codoon.com/third/beego/context v0.0.0-20200331102055-215ba36bc878 // indirect
	git.in.codoon.com/third/beego/orm v0.0.0-20200331102055-215ba36bc878 // indirect
	git.in.codoon.com/third/beego/session v0.0.0-20200331102055-215ba36bc878 // indirect
	git.in.codoon.com/third/beego/utils v0.0.0-20200331102055-215ba36bc878 // indirect
	git.in.codoon.com/third/context v0.0.0-20200401053909-4c3008b5cf79
	git.in.codoon.com/third/crc32 v0.0.0-20200401030510-4553cb28d416 // indirect
	git.in.codoon.com/third/etcd-client v0.0.0-20200909103706-1364dbdb1d5c // indirect
	git.in.codoon.com/third/etcd/pkg/pathutil v0.0.0-20200409060612-d4ce68b4f25d // indirect
	git.in.codoon.com/third/etcd/pkg/types v0.0.0-20200409060612-d4ce68b4f25d // indirect
	git.in.codoon.com/third/etcd_v3 v0.0.0-20200917100123-e4db3d835d93
	git.in.codoon.com/third/g2s v0.0.0-20200331102136-ab730e68af69 // indirect
	git.in.codoon.com/third/gin v0.0.0-20200331102141-b56734d01bd6
	git.in.codoon.com/third/gin/binding v0.0.0-20200331102141-b56734d01bd6 // indirect
	git.in.codoon.com/third/gin/render v0.0.0-20200331102141-b56734d01bd6 // indirect
	git.in.codoon.com/third/go-cache v0.0.0-20200331102148-2b3e3724e572
	git.in.codoon.com/third/go-codec/codec v0.0.0-20200331102655-c489e00a1ee4 // indirect
	git.in.codoon.com/third/go-colorable v0.0.0-20200331102151-42fd58b90add // indirect
	git.in.codoon.com/third/go-iap/appstore v0.0.0-20200331104311-c827044d82e0
	git.in.codoon.com/third/go-local v0.0.0-20200331102159-d88358be4dcf // indirect
	git.in.codoon.com/third/go-logging v0.0.0-20200331102200-9289b86becf2 // indirect
	git.in.codoon.com/third/go-metrics v0.0.0-20200331102202-685cf2d9773b // indirect
	git.in.codoon.com/third/go-resiliency/breaker v0.0.0-20200331102846-bcb70885e7e4 // indirect
	git.in.codoon.com/third/go-spew/spew v0.0.0-20200331102900-5bbd30c34478 // indirect
	git.in.codoon.com/third/go-xerial-snappy v0.0.0-20200331102213-46a4796787f3 // indirect
	git.in.codoon.com/third/go.uuid v1.2.0 // indirect
	git.in.codoon.com/third/golang.org/x/net/http2 v0.0.0-20200421021357-cf13fb474bce // indirect
	git.in.codoon.com/third/golang.org/x/net/http2/hpack v0.0.0-20200421021357-cf13fb474bce // indirect
	git.in.codoon.com/third/gorm v0.0.0-20210705083209-799f58fce370
	git.in.codoon.com/third/goroutine v0.0.0-20200331102230-c912e6db4d6a // indirect
	git.in.codoon.com/third/goroutineid v0.0.0-20200717063048-b6accc6b85f6 // indirect
	git.in.codoon.com/third/http_client_cluster v0.0.0-20200922064940-834b3323a368 // indirect
	git.in.codoon.com/third/httprouter v0.0.0-20200331102241-bd9396311bcb // indirect
	git.in.codoon.com/third/hystrix v0.0.0-20200331102243-29626b7c0e08 // indirect
	git.in.codoon.com/third/hystrix/metric_collector v0.0.0-20200331102243-29626b7c0e08 // indirect
	git.in.codoon.com/third/hystrix/rolling v0.0.0-20200331102243-29626b7c0e08 // indirect
	git.in.codoon.com/third/jwt-go v0.0.0-20200331102248-17c6c1d3d827 // indirect
	git.in.codoon.com/third/lz4 v0.0.0-20200331102251-62791713f25a // indirect
	git.in.codoon.com/third/pili-engineering/pili-sdk-go/pili v0.0.0-20200401100447-8348bdfa5f25 // indirect
	git.in.codoon.com/third/pq/hstore v0.0.0-20200331102306-1572c04087a8 // indirect
	git.in.codoon.com/third/qiniupkg.com/api.v7/api v0.0.0-20200407055923-d9e91162bab1 // indirect
	git.in.codoon.com/third/qiniupkg.com/api.v7/auth/qbox v0.0.0-20200407055923-d9e91162bab1 // indirect
	git.in.codoon.com/third/qiniupkg.com/api.v7/conf v0.0.0-20200407055923-d9e91162bab1
	git.in.codoon.com/third/qiniupkg.com/api.v7/kodo v0.0.0-20200407055923-d9e91162bab1
	git.in.codoon.com/third/qiniupkg.com/api.v7/kodocli v0.0.0-20200407055923-d9e91162bab1
	git.in.codoon.com/third/qiniupkg.com/x/bytes.v7 v0.0.0-20200407055923-d9e91162bab1 // indirect
	git.in.codoon.com/third/qiniupkg.com/x/bytes.v7/seekable v0.0.0-20200407055923-d9e91162bab1 // indirect
	git.in.codoon.com/third/qiniupkg.com/x/ctype.v7 v0.0.0-20200407055923-d9e91162bab1 // indirect
	git.in.codoon.com/third/qiniupkg.com/x/log.v7 v0.0.0-20200407055923-d9e91162bab1 // indirect
	git.in.codoon.com/third/qiniupkg.com/x/reqid.v7 v0.0.0-20200407055923-d9e91162bab1 // indirect
	git.in.codoon.com/third/qiniupkg.com/x/rpc.v7 v0.0.0-20200407055923-d9e91162bab1 // indirect
	git.in.codoon.com/third/qiniupkg.com/x/url.v7 v0.0.0-20200407055923-d9e91162bab1 // indirect
	git.in.codoon.com/third/qiniupkg.com/x/xlog.v7 v0.0.0-20200407055923-d9e91162bab1 // indirect
	git.in.codoon.com/third/queue v0.0.0-20200331102307-c24dd2e446c4 // indirect
	git.in.codoon.com/third/raven-go v0.0.0-20200331102309-4497cb87ba1b // indirect
	git.in.codoon.com/third/redigo/internal v0.0.0-20200401051426-418e88245970 // indirect
	git.in.codoon.com/third/redigo/redis v0.0.0-20200401051426-418e88245970
	git.in.codoon.com/third/sarama v0.0.0-20200401033751-9d1a38883991 // indirect
	git.in.codoon.com/third/sarama-cluster v0.0.0-20200331102322-700b5c6aa6b8 // indirect
	git.in.codoon.com/third/sarama_1.26/aescts v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/compress/fse v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/compress/huff0 v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/compress/snappy v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/compress/zstd v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/compress/zstd/internal/xxhash v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/dnsutils v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/go-metrics v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/go-resiliency/breaker v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/go-spew/spew v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/go-uuid v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/go-xerial-snappy v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gofork/encoding/asn1 v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gofork/x/crypto/pbkdf2 v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/asn1tools v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/client v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/config v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/credentials v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/crypto v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/crypto/common v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/crypto/etype v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/crypto/rfc3961 v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/crypto/rfc3962 v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/crypto/rfc4757 v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/crypto/rfc8009 v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/gssapi v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/iana v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/iana/addrtype v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/iana/adtype v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/iana/asnAppTag v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/iana/chksumtype v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/iana/errorcode v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/iana/etypeID v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/iana/flags v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/iana/keyusage v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/iana/msgtype v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/iana/nametype v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/iana/patype v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/kadmin v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/keytab v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/krberror v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/messages v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/pac v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/gokrb5/types v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/golang.org/x/crypto/md4 v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/golang.org/x/crypto/pbkdf2 v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/golang.org/x/net/internal/socks v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/golang.org/x/net/proxy v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/lz4 v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/lz4/internal/xxh32 v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/queue v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/rpc/mstypes v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/rpc/ndr v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/sarama v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sarama_1.26/snappy v0.0.0-20200604024659-dbe3dcaac7c3 // indirect
	git.in.codoon.com/third/sentinel-golang/api v0.0.0-20200615061810-caa09fc6babf // indirect
	git.in.codoon.com/third/sentinel-golang/core v0.0.0-20200615061810-caa09fc6babf // indirect
	git.in.codoon.com/third/sentinel-golang/logging v0.0.0-20200615061810-caa09fc6babf // indirect
	git.in.codoon.com/third/sentinel-golang/util v0.0.0-20200615061810-caa09fc6babf // indirect
	git.in.codoon.com/third/snappy v0.0.0-20200331102324-efa7634832d3 // indirect
	git.in.codoon.com/third/text/encoding v0.0.0-20200401053849-6eb37c02527e // indirect
	git.in.codoon.com/third/text/encoding/internal v0.0.0-20200401053849-6eb37c02527e // indirect
	git.in.codoon.com/third/text/encoding/internal/identifier v0.0.0-20200401053849-6eb37c02527e // indirect
	git.in.codoon.com/third/text/encoding/simplifiedchinese v0.0.0-20200401053849-6eb37c02527e
	git.in.codoon.com/third/text/transform v0.0.0-20200401053849-6eb37c02527e // indirect
	git.in.codoon.com/third/upyun v0.0.0-20200331102342-b9b47dbc8ae1 // indirect
	git.in.codoon.com/third/uuid v0.0.0-20200331102343-3dadab925980
	git.in.codoon.com/third/xlsx v0.0.0-20200331102352-f141cc04db59 // indirect
	git.in.codoon.com/third/xlsx_v2 v0.0.0-20200331102354-fec82e5bc549
	git.in.codoon.com/third/xxHash/xxHash32 v0.0.0-20200331102959-be4fc98e942f // indirect
)
