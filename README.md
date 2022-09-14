# ZPrize submission

This is based on the [provided test harness](https://github.com/celo-org/zprize-mobile-harness).

The submission is written in Go, and derived from the [`gnark-crypto`](https://github.com/consensys/gnark-crypto) project. Note that the `bls12377` package was not audited (but the **very similar** packages `bn254`, `bls12381`, were audited).

## Results

On the target device, for a random bls12377 G1 MSM (n == 2**16), we measure

* ~2309ms for the reference benchmark
* ~650ms for our submission

(-71.8%).

## Getting started

For simplicity, our code follows closely the `benchmarkMSMRandomMultipleVecs` and `benchmarkMSMFile` provided.

### Pre-requisites

* [Install Go 1.19](https://go.dev/doc/install) (latest version)

### Build the app

```bash
make install
```

This build our code, copies it to the Android application folder. Then simply run the test harness as described in the original [README](https://github.com/celo-org/zprize-mobile-harness). 

"PRESS TO RUN FROM TEST VECTOR FILE" and "PRESS TO RUN USING RANDOM ELEMENTS" are both implemented; the result is displayed under the reference one (prefixed by "[gnark]").

Additionally, when running the test vector file, we compare the `result.txt` file with our output `gnark_result.txt`.

We also append in the same format the results to the `resulttimes.txt` file.

### Misc

See the Makefile for options to test locally (`make local`) or on the device without building an Android app (`make run`).

## Optimizations

We experimented several approaches; here is a description of the key findings for the final one.

1. It uses `gnark-crypto/bls12377` package, which out of the box performs very well (> `arkworks`). The MSM algorithm is described in the Go doc of the corresponding method. Essentially, it's the "bucket" method using parameterized Jacobian coordinates (x=X/ZZ, y=Y/ZZZ, ZZ³=ZZZ²).
2. We perform a static build targetting a 64bit arm linux architecture, which allows without a complicated build procress to run 64bit code on the target device. We copy the output in the armv7 (32bit) destination folder; in a production deployment, Java calling code must at runtime check for the actual CPU architecture and switch to a fallback if it's 32bit (outside of the scope of the challenge). Note that while the submission spawn a process at each msm call, other ways may turn out more efficient (allocate the verifying key on the stack, communicate with the process with unix sockets, ..., couple of ideas, contact us for more info).
3. To increase performance, we empirically modified the C parameter of the MSM (the window size). Essentialy, for 2**16 points, we aim for C-1 (vs default value proposed by `gnark-crypto`) as it seems the memory pressure is higher in this environment. That is, we have less buckets than on non-mobile env to allocate on the heap.

4. We hand tuned the field arithmetic for the Multiplication targetting the `arm64` architecture. Our pure-go version performed better than the arm assembly one, and resulted in a ~20% speed up on some platforms compared to existing version in `gnark-crypto`.

5. We implemented and optimized a dedicated Squaring algorithm (rather than calling the Multiplication as in `gnark-crypto`) following our previous work https://hackmd.io/@gnark/modular_multiplication , which resulted in significant perf improvement on the target device.
6. Our GPU experimentations were not promising (on this device).

---------

Our code also includes some serialization helpers and modification to be compatible with arkworks format. Once the competition is over, the `Mul` and `Square` optimizations for `arm64` will land in `gnark-crypto`.  [Get in touch](gnark@consensys.net) if you have any questions.

## License

© 2022 ConsenSys [gnark@consensys.net].

This project is licensed under either of

* Apache License, Version 2.0 (LICENSE-APACHE)
* MIT license (LICENSE-MIT)

at your option.