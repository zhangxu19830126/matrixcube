// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: timestamp.proto

package hlcpb

import (
	fmt "fmt"
	io "io"
	math "math"
	math_bits "math/bits"

	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

// Timestamp is a HLC time value. All its field should never be accessed
// directly by its users.
type Timestamp struct {
	// PhysicalTime is the physical component of the HLC, it is read from a node's
	// wall clock time as Unix epoch time in nanoseconds. HLC requires this field
	// to be monotonically increase on each node.
	PhysicalTime int64 `protobuf:"varint,1,opt,name=physical_time,json=physicalTime,proto3" json:"physical_time,omitempty"`
	// LogicalTime is the logical component of the HLC, its value is maintained
	// according to the HLC algorithm. The HLC paper further establishes that its
	// value will not overflow in a real production environment.
	LogicalTime          uint32   `protobuf:"varint,2,opt,name=logical_time,json=logicalTime,proto3" json:"logical_time,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Timestamp) Reset()         { *m = Timestamp{} }
func (m *Timestamp) String() string { return proto.CompactTextString(m) }
func (*Timestamp) ProtoMessage()    {}
func (*Timestamp) Descriptor() ([]byte, []int) {
	return fileDescriptor_edac929d8ae1e24f, []int{0}
}
func (m *Timestamp) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Timestamp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Timestamp.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalTo(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Timestamp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Timestamp.Merge(m, src)
}
func (m *Timestamp) XXX_Size() int {
	return m.Size()
}
func (m *Timestamp) XXX_DiscardUnknown() {
	xxx_messageInfo_Timestamp.DiscardUnknown(m)
}

var xxx_messageInfo_Timestamp proto.InternalMessageInfo

func (m *Timestamp) GetPhysicalTime() int64 {
	if m != nil {
		return m.PhysicalTime
	}
	return 0
}

func (m *Timestamp) GetLogicalTime() uint32 {
	if m != nil {
		return m.LogicalTime
	}
	return 0
}

func init() {
	proto.RegisterType((*Timestamp)(nil), "hlcpb.Timestamp")
}

func init() { proto.RegisterFile("timestamp.proto", fileDescriptor_edac929d8ae1e24f) }

var fileDescriptor_edac929d8ae1e24f = []byte{
	// 152 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0xc9, 0xcc, 0x4d,
	0x2d, 0x2e, 0x49, 0xcc, 0x2d, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0xcd, 0xc8, 0x49,
	0x2e, 0x48, 0x92, 0xd2, 0x4d, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x4b, 0xce, 0xcf, 0xd5, 0x4f,
	0xcf, 0x4f, 0xcf, 0xd7, 0x07, 0xcb, 0x26, 0x95, 0xa6, 0x81, 0x79, 0x60, 0x0e, 0x98, 0x05, 0xd1,
	0xa5, 0x14, 0xcc, 0xc5, 0x19, 0x02, 0x33, 0x48, 0x48, 0x99, 0x8b, 0xb7, 0x20, 0xa3, 0xb2, 0x38,
	0x33, 0x39, 0x31, 0x27, 0x1e, 0x64, 0xbc, 0x04, 0xa3, 0x02, 0xa3, 0x06, 0x73, 0x10, 0x0f, 0x4c,
	0x10, 0xa4, 0x52, 0x48, 0x91, 0x8b, 0x27, 0x27, 0x3f, 0x1d, 0xa1, 0x86, 0x49, 0x81, 0x51, 0x83,
	0x37, 0x88, 0x1b, 0x2a, 0x06, 0x52, 0xe2, 0xc4, 0x73, 0xe2, 0x91, 0x1c, 0xe3, 0x85, 0x47, 0x72,
	0x8c, 0x0f, 0x1e, 0xc9, 0x31, 0x26, 0xb1, 0x81, 0x6d, 0x32, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff,
	0x98, 0x91, 0xa3, 0x7a, 0xb2, 0x00, 0x00, 0x00,
}

func (m *Timestamp) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalTo(dAtA)
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Timestamp) MarshalTo(dAtA []byte) (int, error) {
	var i int
	_ = i
	var l int
	_ = l
	if m.PhysicalTime != 0 {
		dAtA[i] = 0x8
		i++
		i = encodeVarintTimestamp(dAtA, i, uint64(m.PhysicalTime))
	}
	if m.LogicalTime != 0 {
		dAtA[i] = 0x10
		i++
		i = encodeVarintTimestamp(dAtA, i, uint64(m.LogicalTime))
	}
	if m.XXX_unrecognized != nil {
		i += copy(dAtA[i:], m.XXX_unrecognized)
	}
	return i, nil
}

func encodeVarintTimestamp(dAtA []byte, offset int, v uint64) int {
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return offset + 1
}
func (m *Timestamp) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.PhysicalTime != 0 {
		n += 1 + sovTimestamp(uint64(m.PhysicalTime))
	}
	if m.LogicalTime != 0 {
		n += 1 + sovTimestamp(uint64(m.LogicalTime))
	}
	if m.XXX_unrecognized != nil {
		n += len(m.XXX_unrecognized)
	}
	return n
}

func sovTimestamp(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTimestamp(x uint64) (n int) {
	return sovTimestamp(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Timestamp) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTimestamp
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Timestamp: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Timestamp: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field PhysicalTime", wireType)
			}
			m.PhysicalTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTimestamp
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.PhysicalTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field LogicalTime", wireType)
			}
			m.LogicalTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTimestamp
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.LogicalTime |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipTimestamp(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthTimestamp
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthTimestamp
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipTimestamp(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTimestamp
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTimestamp
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowTimestamp
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthTimestamp
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthTimestamp
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowTimestamp
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipTimestamp(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthTimestamp
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthTimestamp = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTimestamp   = fmt.Errorf("proto: integer overflow")
)
