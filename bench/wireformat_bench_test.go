package bench

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"text/tabwriter"

	"github.com/KingrogKDR/omni/bench/pb"
	"google.golang.org/protobuf/proto"
)

type JSONRaftCmd struct {
	Term      uint64 `json:"term"`
	Index     uint64 `json:"index"`
	CmdType   uint64 `json:"cmd_type"`
	Key       []byte `json:"key"`
	Value     []byte `json:"value"`
	ClientID  uint64 `json:"client_id"`
	RequestID uint64 `json:"request_id"`
}

func sampleJSONCmd() *JSONRaftCmd {
	return &JSONRaftCmd{
		Term:      42,
		Index:     1337,
		CmdType:   1, // Put
		Key:       []byte("user:1000:profile"),
		Value:     []byte(`{"name":"foo-bar","plan":"free","region":"ap-south-1"}`),
		ClientID:  9001,
		RequestID: 555111,
	}
}

func samplePBCmd() *pb.RaftCmd {
	return &pb.RaftCmd{
		Term:      42,
		Index:     1337,
		CmdType:   1,
		Key:       []byte("user:1000:profile"),
		Value:     []byte(`{"name":"foo-bar","plan":"free","region":"ap-south-1"}`),
		ClientId:  9001,
		RequestId: 555111,
	}
}

func BenchmarkJSON_Marshal(b *testing.B) {
	cmd := sampleJSONCmd()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := json.Marshal(cmd); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkJSON_Unmarshal(b *testing.B) {
	data, err := json.Marshal(sampleJSONCmd())
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		var out JSONRaftCmd
		if err := json.Unmarshal(data, &out); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPB_Marshal(b *testing.B) {
	cmd := samplePBCmd()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := proto.Marshal(cmd); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPB_Unmarshal(b *testing.B) {
	data, err := proto.Marshal(samplePBCmd())
	if err != nil {
		b.Fatal(err)
	}
	b.ReportAllocs()
	for b.Loop() {
		var out pb.RaftCmd
		if err := proto.Unmarshal(data, &out); err != nil {
			b.Fatal(err)
		}
	}
}

// TestWireSize just prints the on-wire byte size,
// which matters as much as speed for a network layer decision.
func TestWireSize(t *testing.T) {
	jsonData, err := json.Marshal(sampleJSONCmd())
	if err != nil {
		t.Fatal(err)
	}
	pbData, err := proto.Marshal(samplePBCmd())
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("JSON size:     %d bytes", len(jsonData))
	t.Logf("Protobuf size: %d bytes", len(pbData))
}

func TestBenchmarkComparison(t *testing.T) {
	jsonCmd := sampleJSONCmd()
	pbCmd := samplePBCmd()

	jsonData, err := json.Marshal(jsonCmd)
	if err != nil {
		t.Fatal(err)
	}

	pbData, err := proto.Marshal(pbCmd)
	if err != nil {
		t.Fatal(err)
	}

	jsonMarshal := testing.Benchmark(BenchmarkJSON_Marshal)

	pbMarshal := testing.Benchmark(BenchmarkPB_Marshal)

	jsonUnmarshal := testing.Benchmark(BenchmarkJSON_Unmarshal)

	pbUnmarshal := testing.Benchmark(BenchmarkPB_Unmarshal)

	printComparison(
		jsonMarshal,
		pbMarshal,
		jsonUnmarshal,
		pbUnmarshal,
		len(jsonData),
		len(pbData),
	)
}

func printComparison(
	jsonMarshal, pbMarshal,
	jsonUnmarshal, pbUnmarshal testing.BenchmarkResult,
	jsonSize, pbSize int,
) {
	marshalSpeedup :=
		float64(jsonMarshal.NsPerOp()) /
			float64(pbMarshal.NsPerOp())

	unmarshalSpeedup :=
		float64(jsonUnmarshal.NsPerOp()) /
			float64(pbUnmarshal.NsPerOp())

	marshalMemoryRatio :=
		float64(jsonMarshal.AllocedBytesPerOp()) /
			float64(pbMarshal.AllocedBytesPerOp())

	unmarshalMemoryRatio :=
		float64(jsonUnmarshal.AllocedBytesPerOp()) /
			float64(pbUnmarshal.AllocedBytesPerOp())

	unmarshalAllocRatio :=
		float64(jsonUnmarshal.AllocsPerOp()) /
			float64(pbUnmarshal.AllocsPerOp())

	wireRatio :=
		float64(jsonSize) /
			float64(pbSize)

	fmt.Println()
	fmt.Println("JSON vs Protobuf Benchmark")
	fmt.Println("==========================")
	fmt.Println()

	w := tabwriter.NewWriter(
		os.Stdout,
		0,
		4,
		2,
		' ',
		0,
	)

	fmt.Fprintln(w, "Operation\tJSON\tProtobuf\tImprovement")
	fmt.Fprintln(w, "------\t----\t--------\t-----------")

	fmt.Fprintf(
		w,
		"Marshal\t%d ns/op\t%d ns/op\t%.2fx faster\n",
		jsonMarshal.NsPerOp(),
		pbMarshal.NsPerOp(),
		marshalSpeedup,
	)

	fmt.Fprintf(
		w,
		"Unmarshal\t%d ns/op\t%d ns/op\t%.2fx faster\n",
		jsonUnmarshal.NsPerOp(),
		pbUnmarshal.NsPerOp(),
		unmarshalSpeedup,
	)

	fmt.Fprintf(
		w,
		"Marshal memory\t%d B\t%d B\t%.2fx less\n",
		jsonMarshal.AllocedBytesPerOp(),
		pbMarshal.AllocedBytesPerOp(),
		marshalMemoryRatio,
	)

	fmt.Fprintf(
		w,
		"Unmarshal memory\t%d B\t%d B\t%.2fx less\n",
		jsonUnmarshal.AllocedBytesPerOp(),
		pbUnmarshal.AllocedBytesPerOp(),
		unmarshalMemoryRatio,
	)

	fmt.Fprintf(
		w,
		"Unmarshal allocations\t%d\t%d\t%.2fx fewer\n",
		jsonUnmarshal.AllocsPerOp(),
		pbUnmarshal.AllocsPerOp(),
		unmarshalAllocRatio,
	)

	fmt.Fprintf(
		w,
		"Wire size\t%d B\t%d B\t%.2fx smaller\n",
		jsonSize,
		pbSize,
		wireRatio,
	)

	w.Flush()
}
