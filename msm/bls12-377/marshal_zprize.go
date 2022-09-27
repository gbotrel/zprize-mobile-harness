package bls12377

import (
	"encoding/binary"
	"errors"
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fp"
	"github.com/gbotrel/zprize-mobile-harness/msm/bls12-377/fr"
)

// SizeOfG1AffineCompressed represents the size in bytes that a G1Affine need in binary form, compressed
const SizeOfG1AffineCompressed = 48

// SizeOfG1AffineUncompressed represents the size in bytes that a G1Affine need in binary form, uncompressed
const SizeOfG1AffineUncompressed = SizeOfG1AffineCompressed * 2

// Note this follows arkworks little endian serialization format, NOT gnark original
// it provides util method to read and save test vectors as defined in celo test harness

func ReadScalars(path string) (scalars [][]fr.Element, err error) {
	fScalars, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var buf [fr.Bytes]byte

	for {
		_, err = io.ReadFull(fScalars, buf[:8])
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		l := binary.LittleEndian.Uint64(buf[:8])
		v := make([]fr.Element, l)
		for i := 0; i < int(l); i++ {
			_, err = io.ReadFull(fScalars, buf[:])
			if err != nil {
				return
			}
			v[i].ZSetBytes(buf[:])
		}
		scalars = append(scalars, v)
	}
}

func SerializeScalars(scalars [][]fr.Element) []byte {
	rSize := len(scalars) * 8 // reserve space for size of the vectors
	for _, v := range scalars {
		rSize += len(v) * fr.Bytes
	}
	r := make([]byte, rSize)

	var buf [fr.Bytes]byte
	at := 0
	for _, v := range scalars {
		binary.LittleEndian.PutUint64(r[at:at+8], uint64(len(v)))
		at += 8
		for _, s := range v {
			buf = s.ZBytes()
			// reverse(buf[:]) // to little endian
			copy(r[at:at+fr.Bytes], buf[:])
			at += fr.Bytes
		}
	}

	return r
}

func ReadResults(path string) (points []G1Affine, err error) {
	fPoints, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var buf [SizeOfG1AffineCompressed]byte

	for {
		_, err = io.ReadFull(fPoints, buf[:])
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		var p G1Affine
		_, err = p.ZSetBytes(buf[:])
		if err != nil {
			return
		}
		points = append(points, p)
	}
}

func SerializeResults(points []G1Affine) []byte {
	rSize := len(points) * SizeOfG1AffineCompressed // reserve space for size of the vectors
	r := make([]byte, rSize)

	var buf [SizeOfG1AffineCompressed]byte
	at := 0
	for _, p := range points {
		buf = p.ZBytes()
		// reverse(buf[:]) // to little endian
		copy(r[at:at+SizeOfG1AffineCompressed], buf[:])
		at += SizeOfG1AffineCompressed
	}

	return r
}

func ReadPoints(path string) (points [][]G1Affine, err error) {
	fPoints, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var buf [SizeOfG1AffineCompressed]byte

	for {
		_, err = io.ReadFull(fPoints, buf[:8])
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return
		}
		l := binary.LittleEndian.Uint64(buf[:8])
		v := make([]G1Affine, l)
		for i := 0; i < int(l); i++ {
			_, err = io.ReadFull(fPoints, buf[:])
			if err != nil {
				return
			}
			// to big endian
			// reverse(buf[:])
			_, err = v[i].ZSetBytes(buf[:])
			if err != nil {
				return
			}
		}
		points = append(points, v)
	}
}

func SerializePoints(points [][]G1Affine) []byte {
	rSize := len(points) * 8 // reserve space for size of the vectors
	for _, v := range points {
		rSize += len(v) * SizeOfG1AffineCompressed
	}
	r := make([]byte, rSize)

	var buf [SizeOfG1AffineCompressed]byte
	at := 0
	for _, v := range points {
		binary.LittleEndian.PutUint64(r[at:at+8], uint64(len(v)))
		at += 8
		for _, s := range v {
			buf = s.ZBytes()
			// reverse(buf[:]) // to little endian
			copy(r[at:at+SizeOfG1AffineCompressed], buf[:])
			at += SizeOfG1AffineCompressed
		}
	}

	return r
}

// ZBytes returns binary representation of p
// will store X coordinate in regular form and a parity bit
// we follow the BLS12-381 style encoding as specified in ZCash and now IETF
//
// The most significant bit, when set, indicates that the point is in compressed form. Otherwise, the point is in uncompressed form.
//
// The second-most significant bit indicates that the point is at infinity. If this bit is set, the remaining bits of the group element's encoding should be set to zero.
//
// The third-most significant bit is set if (and only if) this point is in compressed form and it is not the point at infinity and its y-coordinate is the lexicographically largest of the two associated with the encoded x-coordinate.
func (p *G1Affine) ZBytes() (res [SizeOfG1AffineCompressed]byte) {

	// check if p is infinity point
	if p.X.IsZero() && p.Y.IsZero() {
		res[len(res)-1] = zmInfinity
		// reverse(res[:])
		return
	}

	// tmp is used to convert from montgomery representation to regular
	var tmp fp.Element

	msbMask := byte(0)
	// compressed, we need to know if Y is lexicographically bigger than -Y
	// if p.Y ">" -p.Y
	if !p.Y.LexicographicallyLargest() {
		msbMask = zmPositiveY
	}

	// we store X  and mask the most significant word with our metadata mask
	tmp = p.X
	tmp.FromMont()
	binary.LittleEndian.PutUint64(res[0:8], tmp[0])
	binary.LittleEndian.PutUint64(res[8:16], tmp[1])
	binary.LittleEndian.PutUint64(res[16:24], tmp[2])
	binary.LittleEndian.PutUint64(res[24:32], tmp[3])
	binary.LittleEndian.PutUint64(res[32:40], tmp[4])
	binary.LittleEndian.PutUint64(res[40:48], tmp[5])

	res[len(res)-1] |= msbMask

	return
}

// To encode G1Affine and G2Affine points, we mask the most significant bits with these bits to specify without ambiguity
// metadata needed for point (de)compression
// we follow the BLS12-381 style encoding as specified in ZCash and now IETF
//
// The most significant bit, when set, indicates that the point is in compressed form. Otherwise, the point is in uncompressed form.
//
// The second-most significant bit indicates that the point is at infinity. If this bit is set, the remaining bits of the group element's encoding should be set to zero.
//
// The third-most significant bit is set if (and only if) this point is in compressed form and it is not the point at infinity and its y-coordinate is the lexicographically largest of the two associated with the encoded x-coordinate.
const (
	zmMask      byte = 0b11 << 6
	zmPositiveY byte = 0b1 << 7
	zmInfinity  byte = 0b1 << 6
)

// ZSetBytes sets p from binary representation in buf and returns number of consumed bytes
// this follow arkworks little endian and flags conventions
// https://docs.rs/ark-serialize/latest/src/ark_serialize/flags.rs.html#74-76
// https://github.com/arkworks-rs/algebra/blob/80857c9714c5a59068f8c20f1298e2138440a1d0/ff/src/fields/models/fp/mod.rs#L581
func (p *G1Affine) ZSetBytes(buf []byte) (int, error) {
	const subGroupCheck = false
	if len(buf) != SizeOfG1AffineCompressed {
		return 0, io.ErrShortBuffer
	}

	// reverse(buf) // to big endian
	// most significant byte
	mData := buf[len(buf)-1] & zmMask
	positiveY := (mData & zmPositiveY) == zmPositiveY
	isInfinity := (mData & zmInfinity) == zmInfinity

	if positiveY && isInfinity {
		return 0, errors.New("positiveY & isInfinity sets")
	}

	if isInfinity {
		p.X.SetZero()
		p.Y.SetZero()
		return SizeOfG1AffineCompressed, nil
	}

	buf[len(buf)-1] &= ^zmMask

	// read X coordinate
	p.X.ZSetBytes(buf)
	// p.X.SetBytes(buf[:fp.Bytes])

	var YSquared, Y fp.Element

	YSquared.Square(&p.X).Mul(&YSquared, &p.X)
	YSquared.Add(&YSquared, &bCurveCoeff)
	if Y.Sqrt(&YSquared) == nil {
		return 0, errors.New("invalid compressed coordinate: square root doesn't exist")
	}

	if Y.LexicographicallyLargest() {
		// Y ">" -Y
		if positiveY {
			Y.Neg(&Y)
		}
	} else {
		// Y "<=" -Y
		if !positiveY {
			Y.Neg(&Y)
		}
	}

	p.Y = Y

	// subgroup check
	if subGroupCheck && !p.IsInSubGroup() {
		return 0, errors.New("invalid point: subgroup check failed")
	}

	return SizeOfG1AffineCompressed, nil
}

// Execute process in parallel the work function
func Execute(nbIterations int, work func(int, int), maxCpus ...int) {

	nbTasks := runtime.NumCPU()
	if len(maxCpus) == 1 {
		nbTasks = maxCpus[0]
	}
	nbIterationsPerCpus := nbIterations / nbTasks

	// more CPUs than tasks: a CPU will work on exactly one iteration
	if nbIterationsPerCpus < 1 {
		nbIterationsPerCpus = 1
		nbTasks = nbIterations
	}

	var wg sync.WaitGroup

	extraTasks := nbIterations - (nbTasks * nbIterationsPerCpus)
	extraTasksOffset := 0

	for i := 0; i < nbTasks; i++ {
		wg.Add(1)
		_start := i*nbIterationsPerCpus + extraTasksOffset
		_end := _start + nbIterationsPerCpus
		if extraTasks > 0 {
			_end++
			extraTasks--
			extraTasksOffset++
		}
		go func() {
			work(_start, _end)
			wg.Done()
		}()
	}

	wg.Wait()
}
