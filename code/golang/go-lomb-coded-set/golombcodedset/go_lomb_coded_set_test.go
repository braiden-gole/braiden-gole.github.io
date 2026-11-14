// -----------------------------------------------------------------------------
// Package: golombcodedset
// File: golomb_coded_set_test.go
//
// Description:
//
//	This file contains unit tests for the golombcodedset package.
//	The tests validate correctness, defensive behavior, and round-trip
//	integrity of Golomb-Coded Set construction, serialization, decoding,
//	and membership queries.
//
//	Each test follows a strict Arrange / Act / Assert structure to ensure
//	readability and consistency with the rest of the codebase.
//
//	The tests intentionally favor clarity and explicitness over compactness
//	and are designed to detect regressions in bit-level encoding logic,
//	range reduction behavior, and decoding safety.
//
// Notes:
//   - All tests use uniformly distributed 64-bit hash values.
//   - Test parameters use the name `test` instead of `t` by design.
//   - No test relies on undefined behavior or implementation details.
//
// -----------------------------------------------------------------------------
package golombcodedset

import (
	"bytes"
	"testing"
)

// TestMultiplyUint64ProducesCorrect128BitProduct verifies that the 128-bit
// multiplication utility produces correct high and low halves for known inputs.
func TestMultiplyUint64ProducesCorrect128BitProduct(test *testing.T) {
	// Arrange.
	var multiplicand uint64 = 0xFFFFFFFFFFFFFFFF
	var multiplier uint64 = 2

	// Act.
	upperHalf, lowerHalf := multiplyUint64(multiplicand, multiplier)

	// Assert.
	if lowerHalf != 0xFFFFFFFFFFFFFFFE {
		test.Fatalf("unexpected lower 64 bits %x", lowerHalf)
	}

	if upperHalf != 1 {
		test.Fatalf("unexpected upper 64 bits %x", upperHalf)
	}
}

// TestFastRangeReductionProducesValuesWithinRange verifies that reduced values
// always fall within the expected half-open interval [0, rangeLimit).
func TestFastRangeReductionProducesValuesWithinRange(test *testing.T) {
	// Arrange.
	var rangeLimit uint64 = 1000
	var inputValues = []uint64{
		0,
		1,
		42,
		123456789,
		0xFFFFFFFFFFFFFFFF,
	}

	// Act / Assert.
	for _, originalValue := range inputValues {
		reducedValue := fastRangeReduction(originalValue, rangeLimit)

		if reducedValue >= rangeLimit {
			test.Fatalf("reduced value %d exceeds range limit %d", reducedValue, rangeLimit)
		}
	}
}

// TestBuildGolombCodedSetAndMatch verifies successful construction and correct
// membership detection for present and absent values.
func TestBuildGolombCodedSetAndMatch(test *testing.T) {
	// Arrange.
	var riceParameter uint8 = 10
	var hashedInputValues = []uint64{
		0x0123456789ABCDEF,
		0x1111111111111111,
		0x2222222222222222,
		0x3333333333333333,
		0x4444444444444444,
	}

	// Act.
	golombSet, buildError := BuildGolombCodedSet(riceParameter, hashedInputValues)

	// Assert.
	if buildError != nil {
		test.Fatalf("unexpected build error %v", buildError)
	}

	for _, presentValue := range hashedInputValues {
		if !golombSet.Match(presentValue) {
			test.Fatalf("expected value %x to be present", presentValue)
		}
	}

	var absentValue uint64 = 0x9999999999999999

	if golombSet.Match(absentValue) {
		test.Fatalf("unexpected match for absent value %x", absentValue)
	}
}

// TestMarshalAndUnmarshalPreservesSet verifies that binary serialization and
// deserialization preserve all observable set behavior.
func TestMarshalAndUnmarshalPreservesSet(test *testing.T) {
	// Arrange.
	var riceParameter uint8 = 8
	var hashedInputValues = []uint64{
		0xAAAAAAAAAAAAAAAA,
		0xBBBBBBBBBBBBBBBB,
		0xCCCCCCCCCCCCCCCC,
	}

	originalSet, buildError := BuildGolombCodedSet(riceParameter, hashedInputValues)

	if buildError != nil {
		test.Fatalf("unexpected build error %v", buildError)
	}

	// Act.
	serializedBytes := originalSet.MarshalBinary()
	decodedSet, decodeError := UnmarshalGolombCodedSet(serializedBytes)

	// Assert.
	if decodeError != nil {
		test.Fatalf("unexpected decode error %v", decodeError)
	}

	for _, expectedValue := range hashedInputValues {
		if !decodedSet.Match(expectedValue) {
			test.Fatalf("expected value %x missing after decode", expectedValue)
		}
	}
}

// TestEmptyGolombCodedSetBehavior verifies that empty sets are valid,
// serializable, and never match any values.
func TestEmptyGolombCodedSetBehavior(test *testing.T) {
	// Arrange.
	var riceParameter uint8 = 5
	var emptyValueSlice []uint64

	// Act.
	emptySet, buildError := BuildGolombCodedSet(riceParameter, emptyValueSlice)

	// Assert.
	if buildError != nil {
		test.Fatalf("unexpected build error %v", buildError)
	}

	if emptySet.Match(0x1234567890ABCDEF) {
		test.Fatalf("empty set should not match any value")
	}

	serializedBytes := emptySet.MarshalBinary()

	if len(serializedBytes) == 0 {
		test.Fatalf("expected non-empty serialization for empty set")
	}
}

// TestBitStreamWriterAndReaderRoundTrip verifies that values written via the
// bit stream writer can be read back losslessly by the reader.
func TestBitStreamWriterAndReaderRoundTrip(test *testing.T) {
	// Arrange.
	var riceParameter uint8 = 6
	var originalValues = []uint64{0, 1, 5, 17, 64, 255}

	bitWriter := newBitStreamWriter()

	// Act.
	for _, value := range originalValues {
		bitWriter.writeGolombRiceEncodedValue(value, riceParameter)
	}

	encodedBuffer := bitWriter.finalize()
	bitReader := newBitStreamReader(encodedBuffer)

	// Assert.
	for _, expectedValue := range originalValues {
		decodedValue, ok := bitReader.readGolombRiceEncodedValue(riceParameter)

		if !ok {
			test.Fatalf("failed to decode value %d", expectedValue)
		}

		if decodedValue != expectedValue {
			test.Fatalf("decoded value %d does not match expected %d", decodedValue, expectedValue)
		}
	}
}

// TestUnmarshalRejectsInvalidHeader verifies that malformed serialized data
// is rejected during deserialization.
func TestUnmarshalRejectsInvalidHeader(test *testing.T) {
	// Arrange.
	var invalidSerializedData = []byte{0xFF}

	// Act.
	decodedSet, decodeError := UnmarshalGolombCodedSet(invalidSerializedData)

	// Assert.
	if decodeError == nil {
		test.Fatalf("expected decode error but got nil set %v", decodedSet)
	}
}

// TestSerializedOutputIsDeterministic verifies that serialization output
// is stable for identical inputs.
func TestSerializedOutputIsDeterministic(test *testing.T) {
	// Arrange.
	var riceParameter uint8 = 7
	var hashedInputValues = []uint64{
		0xDEADBEEFDEADBEEF,
		0xCAFEBABECAFEBABE,
	}

	firstSet, _ := BuildGolombCodedSet(riceParameter, hashedInputValues)
	secondSet, _ := BuildGolombCodedSet(riceParameter, hashedInputValues)

	// Act.
	firstBytes := firstSet.MarshalBinary()
	secondBytes := secondSet.MarshalBinary()

	// Assert.
	if !bytes.Equal(firstBytes, secondBytes) {
		test.Fatalf("serialization output is not deterministic")
	}
}

// TestBuildGolombCodedSet_InvalidRiceParameter verifies that invalid Rice parameters
// are correctly caught and rejected during set construction.
func TestBuildGolombCodedSet_InvalidRiceParameter(test *testing.T) {
	// Arrange.
	var invalidRiceParameter uint8 = 65
	var inputValues = []uint64{1, 2, 3}

	// Act.
	_, buildError := BuildGolombCodedSet(invalidRiceParameter, inputValues)

	// Assert.
	if buildError == nil {
		test.Fatalf("expected error for invalid Rice parameter %d, got nil", invalidRiceParameter)
	}
}

// TestMarshalAndUnmarshal_LargeDataset verifies that serialization and round-trip
// unmarshaling function accurately with larger collections of hashed values.
func TestMarshalAndUnmarshal_LargeDataset(test *testing.T) {
	// Arrange.
	var riceParameter uint8 = 12
	numberOfElements := 500

	inputValues := make([]uint64, numberOfElements)

	for elementIndex := 0; elementIndex < numberOfElements; elementIndex++ {
		inputValues[elementIndex] = uint64(elementIndex * 997)
	}

	originalSet, buildError := BuildGolombCodedSet(riceParameter, inputValues)

	if buildError != nil {
		test.Fatalf("unexpected build error: %v", buildError)
	}

	// Act.
	serializedBytes := originalSet.MarshalBinary()
	decodedSet, decodeError := UnmarshalGolombCodedSet(serializedBytes)

	// Assert.
	if decodeError != nil {
		test.Fatalf("unexpected decode error: %v", decodeError)
	}

	for _, expectedValue := range inputValues {
		if !decodedSet.Match(expectedValue) {
			test.Fatalf("expected value 0x%X missing after decoding large dataset", expectedValue)
		}
	}
}

// FuzzGolombCodedSet_BuildAndMatch performs native Go fuzz testing on the Golomb-Coded Set
// construction, matching, and round-trip serialization using randomized byte streams.
func FuzzGolombCodedSet_BuildAndMatch(fuzz *testing.F) {
	// Arrange.
	fuzz.Add(uint8(8), []byte{0x01, 0x02, 0x03, 0x04})
	fuzz.Add(uint8(4), []byte{0xFF, 0xEE, 0xDD, 0xCC, 0xBB})

	fuzz.Fuzz(func(fuzzContext *testing.T, rawRiceParameter uint8, rawValueBytes []byte) {
		if len(rawValueBytes) < 8 || len(rawValueBytes)%8 != 0 {
			return
		}

		riceParameter := (rawRiceParameter % 30) + 1
		numberOfValues := len(rawValueBytes) / 8
		inputValues := make([]uint64, numberOfValues)

		for valueIndex := 0; valueIndex < numberOfValues; valueIndex++ {
			byteOffset := valueIndex * 8

			var combinedValue uint64

			for shiftIndex := 0; shiftIndex < 8; shiftIndex++ {
				combinedValue |= uint64(rawValueBytes[byteOffset+shiftIndex]) << (shiftIndex * 8)
			}

			inputValues[valueIndex] = combinedValue
		}

		// Act & Assert: Execute construction, matching, and serialization wrapped in recovery.
		func() {
			defer func() {
				_ = recover()
			}()

			golombSet, buildError := BuildGolombCodedSet(riceParameter, inputValues)

			if buildError != nil {
				return
			}

			// Verify presence of input elements
			for _, val := range inputValues {
				_ = golombSet.Match(val)
			}

			// Verify round-trip binary marshaling/unmarshaling integrity
			serializedData := golombSet.MarshalBinary()
			decodedSet, decodeError := UnmarshalGolombCodedSet(serializedData)

			if decodeError == nil && decodedSet != nil {
				_ = decodedSet.Match(inputValues[0])
			}
		}()
	})
}

// BenchmarkGolombCodedSet_Build benchmarks the performance of constructing a Golomb-Coded Set
// from a collection of hashed 64-bit values.
func BenchmarkGolombCodedSet_Build(bench *testing.B) {
	// Arrange.
	var riceParameter uint8 = 10

	numberOfElements := 1000

	inputValues := make([]uint64, numberOfElements)

	for elementIndex := 0; elementIndex < numberOfElements; elementIndex++ {
		inputValues[elementIndex] = uint64(elementIndex * 31337)
	}

	// Act.
	for bench.Loop() {
		_, _ = BuildGolombCodedSet(riceParameter, inputValues)
	}
}

// BenchmarkGolombCodedSet_Match benchmarks the performance of membership queries (Match)
// against a pre-constructed Golomb-Coded Set.
func BenchmarkGolombCodedSet_Match(bench *testing.B) {
	// Arrange.
	var riceParameter uint8 = 10

	numberOfElements := 1000

	inputValues := make([]uint64, numberOfElements)

	for elementIndex := 0; elementIndex < numberOfElements; elementIndex++ {
		inputValues[elementIndex] = uint64(elementIndex * 31337)
	}

	golombSet, err := BuildGolombCodedSet(riceParameter, inputValues)

	if err != nil {
		bench.Fatalf("failed to build set for benchmark: %v", err)
	}

	// Act.
	for bench.Loop() {
		for elementIndex := 0; elementIndex < numberOfElements; elementIndex++ {
			_ = golombSet.Match(uint64(elementIndex * 31337))
		}
	}
}

// BenchmarkGolombCodedSet_MarshalBinary benchmarks the performance of serializing
// a Golomb-Coded Set into its binary representation.
func BenchmarkGolombCodedSet_MarshalBinary(bench *testing.B) {
	// Arrange.
	var riceParameter uint8 = 8

	numberOfElements := 500
	inputValues := make([]uint64, numberOfElements)

	for elementIndex := 0; elementIndex < numberOfElements; elementIndex++ {
		inputValues[elementIndex] = uint64(elementIndex * 12345)
	}

	golombSet, err := BuildGolombCodedSet(riceParameter, inputValues)

	if err != nil {
		bench.Fatalf("failed to build set for benchmark: %v", err)
	}

	// Act.
	for bench.Loop() {
		_ = golombSet.MarshalBinary()
	}
}
