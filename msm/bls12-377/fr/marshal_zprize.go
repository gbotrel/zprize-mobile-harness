package fr

import "encoding/binary"

// ZBytes returns the value of z as a litte-endian byte array
func (z *Element) ZBytes() (res [Limbs * 8]byte) {
	_z := z.ToRegular()
	binary.LittleEndian.PutUint64(res[0:8], _z[0])
	binary.LittleEndian.PutUint64(res[8:16], _z[1])
	binary.LittleEndian.PutUint64(res[16:24], _z[2])
	binary.LittleEndian.PutUint64(res[24:32], _z[3])

	return
}

// ZSetBytes arkworks helper
func (z *Element) ZSetBytes(res []byte) {
	z[0] = binary.LittleEndian.Uint64(res[0:8])
	z[1] = binary.LittleEndian.Uint64(res[8:16])
	z[2] = binary.LittleEndian.Uint64(res[16:24])
	z[3] = binary.LittleEndian.Uint64(res[24:32])

	z.ToMont()
}
