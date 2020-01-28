module github.com/kytra-app/users-srv

go 1.12

require (
	dmitri.shuralyov.com/app/changes v0.0.0-20181114035150-5af16e21babb // indirect
	dmitri.shuralyov.com/service/change v0.0.0-20190203025214-430bf650e55a // indirect
	github.com/99designs/gqlgen v0.7.2 // indirect
	github.com/RoaringBitmap/roaring v0.4.16 // indirect
	github.com/anacrolix/tagflag v0.0.0-20180803105420-3a8ff5428f76 // indirect
	github.com/asaskevich/govalidator v0.0.0-20190424111038-f61b66f89f4a // indirect
	github.com/fatih/structs v1.1.0
	github.com/glycerine/go-unsnap-stream v0.0.0-20181221182339-f9677308dec2 // indirect
	github.com/go-ozzo/ozzo-validation v3.5.0+incompatible
	github.com/golang/protobuf v1.3.2
	github.com/gregjones/httpcache v0.0.0-20190203031600-7a902570cb17 // indirect
	github.com/jinzhu/gorm v1.9.10
	github.com/kytra-app/helpers/microgorm v1.0.0
	github.com/kytra-app/helpers/passwordhasher v1.1.1
	github.com/lucas-clemente/quic-go v0.11.2 // indirect
	github.com/micro/go-api v0.6.0 // indirect
	github.com/micro/go-micro v1.8.1
	github.com/micro/go-plugins v1.1.1
	github.com/micro/kubernetes v0.7.0 // indirect
	github.com/microcosm-cc/bluemonday v1.0.2 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/russross/blackfriday v2.0.0+incompatible // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/shurcooL/go v0.0.0-20190121191506-3fef8c783dec // indirect
	github.com/shurcooL/gofontwoff v0.0.0-20181114050219-180f79e6909d // indirect
	github.com/shurcooL/highlight_diff v0.0.0-20181222201841-111da2e7d480 // indirect
	github.com/shurcooL/highlight_go v0.0.0-20181215221002-9d8641ddf2e1 // indirect
	github.com/shurcooL/home v0.0.0-20190204141146-5c8ae21d4240 // indirect
	github.com/shurcooL/htmlg v0.0.0-20190120222857-1e8a37b806f3 // indirect
	github.com/shurcooL/httpfs v0.0.0-20181222201310-74dc9339e414 // indirect
	github.com/shurcooL/issues v0.0.0-20190120000219-08d8dadf8acb // indirect
	github.com/shurcooL/issuesapp v0.0.0-20181229001453-b8198a402c58 // indirect
	github.com/shurcooL/notifications v0.0.0-20181111060504-bcc2b3082a7a // indirect
	github.com/shurcooL/octicon v0.0.0-20181222203144-9ff1a4cf27f4 // indirect
	github.com/shurcooL/reactions v0.0.0-20181222204718-145cd5e7f3d1 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/shurcooL/webdavfs v0.0.0-20181215192745-5988b2d638f6 // indirect
	go4.org v0.0.0-20181109185143-00e24f1b2599 // indirect
	golang.org/x/perf v0.0.0-20190124201629-844a5f5b46f4 // indirect
	sourcegraph.com/sqs/pbtypes v1.0.0 // indirect
)

replace github.com/kytra-app/helpers/passwordhasher => ../helpers/passwordhasher

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1

replace github.com/kytra-app/helpers/microgorm => ../helpers/microgorm
