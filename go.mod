module github.com/threefoldtech/tfexplorer

go 1.14

require (
	github.com/asaskevich/govalidator v0.0.0-20200907205600-7a23bdc65eef
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/containerd/containerd v1.4.0-beta.1.0.20200615192441-ae2f3fdfd1a4 // indirect
	github.com/dave/jennifer v1.3.0
	github.com/emicklei/dot v0.10.1
	github.com/fzipp/gocyclo v0.3.1 // indirect
	github.com/gordonklaus/ineffassign v0.0.0-20210522101830-0589229737b2 // indirect
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/iancoleman/strcase v0.0.0-20190422225806-e506e3ef7365
	github.com/jbenet/go-base58 v0.0.0-20150317085156-6237cf65f3a6
	github.com/mitchellh/go-homedir v1.1.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.5.1
	github.com/rakyll/statik v0.1.7
	github.com/rs/zerolog v1.19.0
	github.com/rusart/muxprom v0.0.0-20200323164249-36ea051efbe6
	github.com/stellar/go v0.0.0-20201005172357-947b63b1099f
	github.com/stretchr/testify v1.6.1
	github.com/threefoldtech/zos v0.4.0-rc9-b.0.20200918140104-b46553b0c680
	github.com/tyler-smith/go-bip39 v1.0.2
	github.com/urfave/cli v1.22.4
	github.com/whs/nacl-sealed-box v0.0.0-20180930164530-92b9ba845d8d // indirect
	github.com/xdg/scram v0.0.0-20180814205039-7eeb5667e42c // indirect
	github.com/zaibon/httpsig v0.0.0-20210219100301-931cc471f406
	go.mongodb.org/mongo-driver v1.3.2
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20191219145116-fa6499c8e75f
	google.golang.org/grpc v1.29.1 // indirect
	gopkg.in/yaml.v2 v2.2.7
	gotest.tools v2.2.0+incompatible
)

replace github.com/docker/distribution v2.7.1+incompatible => github.com/docker/distribution v2.7.1-0.20190205005809-0d3efadf0154+incompatible
