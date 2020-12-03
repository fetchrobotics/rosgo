package ros

import (
	"io"
	"testing"

	"google3/third_party/golang/cmp/cmp"
)

func TestReader(t *testing.T) {
	type nextOp struct {
		n    int
		want []byte
	}
	type readOp struct {
		buffer  []byte
		want    []byte
		wantN   int
		wantErr error
	}
	type op struct {
		read *readOp
		next *nextOp
	}

	for _, tc := range []struct {
		name string
		data []byte
		ops  []op
	}{
		{
			name: "ReadAll",
			data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24},
			ops: []op{
				{read: &readOp{
					buffer: make([]byte, 8),
					want:   []byte{1, 2, 3, 4, 5, 6, 7, 8},
					wantN:  8,
				}},
				{read: &readOp{
					buffer: make([]byte, 8),
					want:   []byte{9, 10, 11, 12, 13, 14, 15, 16},
					wantN:  8,
				}},
				{read: &readOp{
					buffer: make([]byte, 16),
					want:   []byte{17, 18, 19, 20, 21, 22, 23, 24, 0, 0, 0, 0, 0, 0, 0, 0},
					wantN:  8,
				}},
				{read: &readOp{
					buffer:  make([]byte, 1),
					want:    []byte{0},
					wantN:   0,
					wantErr: io.EOF,
				}},
			},
		},
		{
			name: "NextAll",
			data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24},
			ops: []op{
				{next: &nextOp{
					n:    8,
					want: []byte{1, 2, 3, 4, 5, 6, 7, 8},
				}},
				{next: &nextOp{
					n:    8,
					want: []byte{9, 10, 11, 12, 13, 14, 15, 16},
				}},
				{next: &nextOp{
					n:    16,
					want: []byte{17, 18, 19, 20, 21, 22, 23, 24},
				}},
				{next: &nextOp{
					n:    8,
					want: []byte{},
				}},
			},
		},
		{
			name: "ReadAndNextAll",
			data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24},
			ops: []op{
				{read: &readOp{
					buffer: make([]byte, 8),
					want:   []byte{1, 2, 3, 4, 5, 6, 7, 8},
					wantN:  8,
				}},
				{next: &nextOp{
					n:    8,
					want: []byte{9, 10, 11, 12, 13, 14, 15, 16},
				}},
				{read: &readOp{
					buffer: make([]byte, 16),
					want:   []byte{17, 18, 19, 20, 21, 22, 23, 24, 0, 0, 0, 0, 0, 0, 0, 0},
					wantN:  8,
				}},
				{next: &nextOp{
					n:    8,
					want: []byte{},
				}},
			},
		},
		{
			name: "ReadEmtpyBuffer",
			data: []byte{},
			ops: []op{
				{read: &readOp{
					buffer:  make([]byte, 1),
					want:    []byte{0},
					wantN:   0,
					wantErr: io.EOF,
				}},
			},
		},
		{
			name: "ReadEmtpyBufferEmptyRead",
			data: []byte{},
			ops: []op{
				{read: &readOp{
					buffer:  make([]byte, 0),
					want:    []byte{},
					wantN:   0,
					wantErr: io.EOF,
				}},
			},
		},
		{
			name: "NextEmtpyBuffer",
			data: []byte{},
			ops: []op{
				{next: &nextOp{
					n:    1,
					want: []byte{},
				}},
			},
		},
		{
			name: "NextEmtpyBufferZeroBytes",
			data: []byte{},
			ops: []op{
				{next: &nextOp{
					n:    0,
					want: []byte{},
				}},
			},
		},
	} {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := NewReader(tc.data)
			for _, op := range tc.ops {
				if op.read != nil {
					n, err := r.Read(op.read.buffer)

					if diff := cmp.Diff(op.read.want, op.read.buffer); diff != "" {
						t.Errorf("Read(%v) read unexpected bytes: diff (-want +got): %v\nreader: %+v", op.read.buffer, diff, r)
					}

					if op.read.wantN != n {
						t.Errorf("Read(%v) read unexpected number of bytes: got %v, want %v\nreader: %+v", op.read.buffer, n, op.read.wantN, r)
					}

					if op.read.wantErr != err {
						t.Errorf("Read(%v) returned unexpected error: got %v, want %v\nreader: %+v", op.read.buffer, err, op.read.wantErr, r)
					}
				} else if op.next != nil {
					got := r.Next(op.next.n)
					if diff := cmp.Diff(op.next.want, got); diff != "" {
						t.Errorf("Next(%v) returned unexpected bytes: diff (-want +got): %v\nreader: %+v", op.next.n, diff, r)
					}
				}
			}
		})
	}
}
