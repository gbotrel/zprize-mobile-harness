.PHONY: install run local clean updatefield

install:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o gnark . 
	cp gnark ./android/ZPrize/app/src/main/jniLibs/armeabi-v7a/lib_gnark_.so
	cp gnark ./android/ZPrize/app/src/main/jniLibs/arm64-v8a/lib_gnark_.so
	rm gnark

local:
	go build -ldflags="-s -w"  -o gnark .
	./gnark -n 14
	
run:
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o gnark . 
	adb push gnark /data/local/tmp/
	adb shell "/data/local/tmp/gnark -n 8 -wd /data/local/tmp/"

clean:
	rm gnark
	rm *.txt

updatefield:
	cp ../../consensys/gnark-crypto/ecc/bls12-377/fp/* ./msm/bls12-377/fp/
	cd ./msm/bls12-377/fp/ && rm *amd64* && rm element_mul.go
	cp ../../consensys/gnark-crypto/ecc/bls12-377/fr/*.go ./msm/bls12-377/fr/
	cp ../../consensys/gnark-crypto/ecc/bls12-377/fr/*.s ./msm/bls12-377/fr/
	cd ./msm/bls12-377/fr/ && rm *amd64* && rm element_mul.go

updatemsm:
	cp ../../consensys/gnark-crypto/ecc/bls12-377/g1.go ./msm/bls12-377/
	cp ../../consensys/gnark-crypto/ecc/bls12-377/multiexp.go ./msm/bls12-377/
	cp ../../consensys/gnark-crypto/ecc/bls12-377/multiexp_test.go ./msm/bls12-377/
	cp ../../consensys/gnark-crypto/ecc/bls12-377/g1_test.go ./msm/bls12-377/