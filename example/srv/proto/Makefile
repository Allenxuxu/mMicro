.PHONY: proto
proto:
	protoc --proto_path=. --go_out=. --micro_out=. hello/hello.proto
	ls hello/*.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
