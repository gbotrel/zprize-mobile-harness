# ZPrize submission

This is based on the [provided test harness](https://github.com/celo-org/zprize-mobile-harness).

The submission is written in Go, and derived from the [`gnark-crypto`](https://github.com/consensys/gnark-crypto) project. Note that the `bls12377` package was not audited (but the **very similar** packages `bn254`, `bls12381`, were audited).

## Results

On the target device, for a random bls12377 G1 MSM (n == 2**16), we measure

* ~2309ms for the reference benchmark
* ~509ms for our submission

(-77.9%).

## Pre-requisites

* [Install Go 1.19](https://go.dev/doc/install) (latest version)
* /!\ important: for simplicity as the repo is private until the deadline; clone this repository inside the $GOPATH;

```bash
echo $GOPATH # be sure that your Go install is OK
mkdir -p $GOPATH/src/github.com/gbotrel/
cd $GOPATH/src/github.com/gbotrel/
git clone git@github.com:gbotrel/zprize-mobile-harness.git
```

## Getting started

For simplicity, our code follows closely the `benchmarkMSMRandomMultipleVecs` and `benchmarkMSMFile` provided.

We use the provided test harness; any command (test vectors or random instances) is ran both with the reference arkworks,
and with our submission, and both results are displayed. We also check result consistency here --

"PRESS TO RUN FROM TEST VECTOR FILE" and "PRESS TO RUN USING RANDOM ELEMENTS" are both implemented; the result is displayed under the reference one (prefixed by "[gnark]").

Additionally, when running the test vector file, we compare the `result.txt` file with our output `gnark_result.txt`.

We also append in the same format the results to the `resulttimes.txt` file.

### Quick start

Simply run the test harness as described in the original [README](https://github.com/celo-org/zprize-mobile-harness).

### Build the app from sractch

```bash
make install
```

This build our code, copies it to the Android application folder. Then, launch the Android app.

### Misc

See the Makefile for options to test locally (`make local`) or on the device without building an Android app (`make run`).

## Optimizations

We experimented several approaches; here is a description of the key findings for the final one.

1. It uses `gnark-crypto/bls12377` package, which out of the box performs very well (> `arkworks`). The MSM algorithm is described in the Go doc of the corresponding method. Additionaly, we introduce a bls12-377 algorithmic optimization; the "bucket/pippenger" method now uses an optimized twisted edwards extended cordinate system, resulting in a significant performance improvement. You can read more about this in [the attached note](). 

2. We perform a static build targetting a 64bit arm linux architecture, which allows without a complicated build procress to run 64bit code on the target device. We copy the output in the armv7 (32bit) destination folder; in a production deployment, Java calling code must at runtime check for the actual CPU architecture and switch to a fallback if it's 32bit (outside of the scope of the challenge). Note that while the submission spawn a process at each msm call, other ways may turn out more efficient (allocate the verifying key on the stack, communicate with the process with unix sockets, ...).

3. We hand tuned the field arithmetic for the Multiplication targetting the `arm64` architecture. Our pure-go version performed better than the arm assembly one, and resulted in a ~20% speed up on some platforms compared to existing version in `gnark-crypto`.

4. We implemented and optimized a dedicated Squaring algorithm (rather than calling the Multiplication as in `gnark-crypto`) following our previous work https://hackmd.io/@gnark/modular_multiplication , which resulted in significant perf improvement on the target device. This is not used in the twisted edwards extended MSM, only in the parameterized Jacobian version which uses Affine points as input (branch: TODO, performance: ~620ms for 2**14).

5. For the target (arm64) we add ~40lines of arm assembly for a small function (`fp.Butterfly(a, b) -> a = a + b; b = a - b`). The perf impact is ~5%, as it speeds up a bit the `UnifiedMixedAdd` point addition in the buckets (msm). The rest of the submission is compiled from pure Go code;

6. On this device, our GPU experimentations were not promising.

7. We raised [an issue](https://github.com/golang/go/issues/54607) to the Golang team. Once the fix is merged into the latest Golang compiler release, we might squeeze an extra 5-10% perf improvement.

---------

Our code also includes some serialization helpers and modification to be compatible with arkworks format. Once the competition is over, the `Mul` and `Square` optimizations for `arm64` will land in `gnark-crypto`.  [Get in touch](gnark@consensys.net) if you have any questions.

The new MSM with optimized twisted edwards extended cordinate system will also be supported in `gnark-crypto`, for the curves that allows it.

## License

© 2022 ConsenSys [gnark@consensys.net].

This project is licensed under either of

* Apache License, Version 2.0 (LICENSE-APACHE)
* MIT license (LICENSE-MIT)

at your option.