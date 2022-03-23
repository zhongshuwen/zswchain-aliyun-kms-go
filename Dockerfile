FROM scratch
COPY zswchain-aliyun-kms-go /usr/bin/zswchain-aliyun-kms-go
ENTRYPOINT ["/usr/bin/zswchain-aliyun-kms-go"]
