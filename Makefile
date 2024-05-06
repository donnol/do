.PHONY:

install_letgo:
	go install ./cmd/letgo

test_letgo_sql2struct:install_letgo
	letgo sql2struct --pkg=sqlparser -f ./cmd/letgo/sqlparser/test.sql -o ./cmd/letgo/sqlparser/test.go

test_letgo_sql2struct_insert:install_letgo
	letgo sql2struct --pkg=sqlparser -f ./cmd/letgo/sqlparser/test.sql insert --amount=3

# 生成证书(参照：https://cloud.tencent.com/developer/article/1548350)
# .key: 私钥
# .csr: 证书请求文件，这个文件中会包含申请人的一些信息
# .crt: 自签名证书，CA机构用自己的私钥和证书申请文件生成自己签名的证书，俗称自签名证书，这里可以理解为根证书。
# .pem: 内容与crt文件一样。Privacy Enhanced Mail. [You may have seen digital certificate files with a variety of filename extensions, such as .crt, .cer, .pem, or .der. These extensions generally map to two major encoding schemes for X.509 certificates and keys: PEM (Base64 ASCII), and DER (binary).](https://www.ssl.com/guide/pem-der-crt-and-cer-x-509-encodings-and-conversions/)
# 
# 在生成 HTTPS 服务器端证书时，需要填写CN: Common Name (e.g. server FQDN or YOUR name), 即访问服务的域名信息，如果有很多子域名，可以用 * 代替，如 *.test.com。
certdir=testdata/cert/
certgen:
	cd $(certdir) && \
	openssl genrsa -out server.key 1024 && \
	openssl req -new -key server.key -out server.csr && \
	openssl genrsa -out ca.key 1024 && \
	openssl req -new -key ca.key -out ca.csr && \
	openssl x509 -req -in ca.csr -signkey ca.key -out ca.crt && \
	openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -in server.csr -out server.crt && \
	echo 'generate cert completed.'
