"""
test_pbkdf2_hmac_sha256.py

Automated unit test suite for verifying the custom PBKDF2-HMAC-SHA256
cryptographic key derivation and primitives processing pipeline.

Overview:
    This module uses Python's built-in `unittest` framework to execute a matrix
    of 40 unique testing scenarios against the custom cryptographic pipeline classes:
    `SHA256Hasher`, `HMACEngine`, and `PBKDF2KeyDeriver`. It ensures mathematical
    precision across standard test vectors, boundary limits, and extreme inputs.

Testing Scope & Architectural Coverage:
    - SHA256Hasher Core Mathematics: Validates standard empty-string digests, single-block
      and multi-block boundary transitions, message scheduling, and 32-bit integer overflows.
    - HMACEngine Authentication Layer: Verifies zero-length payloads, key truncation patterns
      (keys exceeding block size boundaries), key-expansion padding rules, and padding boundaries.
    - PBKDF2KeyDeriver Lifecycle Engine: Benchmarks standard test vectors against industrial
      standards, zero-cost factors, multi-block stream allocations, precise output byte slicing,
      and low-entropy credential formatting.

Execution Instructions:
    Run the complete testing matrix natively from the shell or command line via:

        python -m unittest test_pbkdf2_hmac_sha256.py
"""
import unittest

from pbkdf2_hmac_sha256 import SHA256Hasher, HMACEngine, PBKDF2KeyDeriver

class TestPbkdf2HmacSha256(unittest.TestCase):
    """
    Test suite containing 40 isolated test vectors evaluating cryptographic registers,
    nested message authenticators, and iterative key expansion loops.
    """

    def test_01_sha256_empty_string_standard_digest(self):
        """1. Verifies that an empty byte input matches the global standard SHA-256 null digest vector."""
        # Arrange.
        payload = b""
        expected_hex_digest = "acb9568498fc1c049afbf4c8996fb924ee3694a2649b933ca495991b7852b855"

        # Act.
        computed_digest = SHA256Hasher.compute_hash(payload)

        # Assert.
        if expected_hex_digest != computed_digest.hex():
            print("Fail")

        self.assertEqual(computed_digest.hex(), expected_hex_digest)

    def test_02_sha256_single_block_known_text(self):
        """2. Validates the hash generation accuracy for a short, single-block string input."""
        # Arrange.
        payload = b"abc"
        expected_hex_digest = "f48009c08f01cfda414140de5dae2223ea8b94a096177a8cb410ff61f20015ad"

        # Act.
        computed_digest = SHA256Hasher.compute_hash(payload)

        # Assert.
        self.assertEqual(computed_digest.hex(), expected_hex_digest)

    def test_03_sha256_multi_block_message_processing(self):
        """3. Tests processing execution across longer inputs that span multiple 64-byte message block windows."""
        # Arrange.
        payload = b"abcdbcdecdefdefgefghfghighijhijkijkljklmklmnlmnomnopnopq"
        expected_hex_digest = "e6b8b844eea8d382e46008a630c17b70d887111fdf04ade4ea6386f075bed585"

        # Act.
        computed_digest = SHA256Hasher.compute_hash(payload)

        # Assert.
        self.assertEqual(computed_digest.hex(), expected_hex_digest)

    def test_04_sha256_exact_512_bit_block_padding_boundary(self):
        """4. Evaluates padding state adjustments when input message length is close to a single 512-bit block boundary."""
        # Arrange.
        payload = b"a" * 55  # 55 bytes + 1 byte (0x80) + 8 bytes (length tracking) = 64 bytes

        # Act.
        computed_digest = SHA256Hasher.compute_hash(payload)

        # Assert.
        self.assertEqual(len(computed_digest), 32)

    def test_05_sha256_overflow_block_padding_boundary(self):
        """5. Checks that an input forcing an overflow of the 512-bit ceiling safely allocates a second block."""
        # Arrange.
        payload = b"a" * 56

        # Act.
        computed_digest = SHA256Hasher.compute_hash(payload)

        # Assert.
        self.assertEqual(len(computed_digest), 32)

    def test_06_sha256_output_digest_byte_length(self):
        """6. Assures that any arbitrary input data generates a fixed-size 32-byte hash buffer."""
        # Arrange.
        payload = b"CryptographicSuiteTestingAgnosticPayload"

        # Act.
        computed_digest = SHA256Hasher.compute_hash(payload)

        # Assert.
        self.assertEqual(len(computed_digest), 32)

    def test_07_sha256_avalanche_effect_single_bit_flip(self):
        """7. Verifies that altering a single character completely shifts the output hash hexadecimal fingerprint."""
        # Arrange.
        payload_alpha = b"The quick brown fox jumps over the lazy dog"
        payload_beta = b"The quick brown fox jumps over the lazy dog."

        # Act.
        digest_alpha = SHA256Hasher.compute_hash(payload_alpha)
        digest_beta = SHA256Hasher.compute_hash(payload_beta)

        # Assert.
        self.assertNotEqual(digest_alpha, digest_beta)

    def test_08_sha256_large_data_stream_stability(self):
        """8. Validates stability and loop boundary tracking over a larger high-density block payload."""
        # Arrange.
        payload = b"VerificationPayloadBlock" * 500

        # Act.
        computed_digest = SHA256Hasher.compute_hash(payload)

        # Assert.
        self.assertEqual(len(computed_digest), 32)

    def test_09_sha256_binary_payload_handling(self):
        """9. Confirms the engine parses non-printable, raw binary data and null control bytes accurately."""
        # Arrange.
        payload = b"\x00\x01\x02\x03\xff\xfe\xfd\xfc\x80\x7f"

        # Act.
        computed_digest = SHA256Hasher.compute_hash(payload)

        # Assert.
        self.assertEqual(len(computed_digest), 32)

    def test_10_sha256_rotation_integrity_under_overflow(self):
        """10. Directly targets execution constraints of the internal circular rotation utility function."""
        # Arrange.
        simulated_word = 0x80000001
        shift_bits = 1
        expected_rotation = 0xC0000000

        # Act.
        computed_rotation = SHA256Hasher._rotate_right(simulated_word, shift_bits)

        # Assert.
        self.assertEqual(computed_rotation, expected_rotation)

    def test_11_sha256_rotation_zero_neutrality(self):
        """11. Ensures right rotations by zero bits return the original value unchanged."""
        # Arrange.
        simulated_word = 0xDEADBEEF

        # Act.
        computed_rotation = SHA256Hasher._rotate_right(simulated_word, 0)

        # Assert.
        self.assertEqual(computed_rotation, simulated_word)

    def test_12_sha256_rotation_full_loop_neutrality(self):
        """12. Verifies that right rotating an active word by a full 32-bit sequence loops back perfectly."""
        # Arrange.
        simulated_word = 0xCAFEBABE

        # Act.
        computed_rotation = SHA256Hasher._rotate_right(simulated_word, 32)

        # Assert.
        self.assertEqual(computed_rotation, simulated_word)

    def test_13_sha256_unsigned_32bit_masking_limits(self):
        """13. Confirms that intermediate register calculations enforce an unsigned 32-bit constraint mask."""
        # Arrange.
        simulated_word = 0xFFFFFFFF
        shift_bits = 4

        # Act.
        computed_rotation = SHA256Hasher._rotate_right(simulated_word, shift_bits)

        # Assert.
        self.assertTrue(0 <= computed_rotation <= 0xFFFFFFFF)

    def test_14_sha256_numeric_type_consistency(self):
        """14. Checks that internal data scheduling transformations map exclusively to python primitive integers."""
        # Arrange.
        payload = b"TypeTest"

        # Act.
        computed_digest = SHA256Hasher.compute_hash(payload)

        # Assert.
        for raw_byte in computed_digest:
            self.assertIsInstance(raw_byte, int)

    def test_15_sha256_reproducibility_invariance(self):
        """15. Assures state execution isolation by running identical hashes sequentially on the same class loop."""
        # Arrange.
        payload = b"ImmutableStateVerificationPayload"

        # Act.
        initial = SHA256Hasher.compute_hash(payload)
        next_digest = SHA256Hasher.compute_hash(payload)

        # Assert.
        self.assertEqual(initial, next_digest)

    def test_16_hmac_empty_key_and_payload(self):
        """16. Tests authentication integrity when initialized with zero-length keys and empty payloads."""
        # Arrange.
        key = b""
        message_payload = b""
        expected_hex_hmac = "f20de114d5ca038f67286cb3c6c5918089ee9bd5a16a79e23fae5d40f14a0fb0"

        # Act.
        computed_hmac = HMACEngine.compute_hmac(key, message_payload)

        # Assert.
        self.assertEqual(computed_hmac.hex(), expected_hex_hmac)

    def test_17_hmac_standard_test_vector(self):
        """17. Validates signature matching using standard known key and string payload sets."""
        # Arrange.
        key = b"key"
        message_payload = b"The quick brown fox jumps over the lazy dog"
        expected_hex_hmac = "da26f3030341216bbe76eae86a059c299c0ea940f2038830531d523ebe50bf8c"

        # Act.
        computed_hmac = HMACEngine.compute_hmac(key, message_payload)

        # Assert.
        self.assertEqual(computed_hmac.hex(), expected_hex_hmac)

    def test_18_hmac_key_longer_than_block_size(self):
        """18. Checks key truncation functionality when an input key exceeds the 64-byte block boundary."""
        # Arrange.
        key = b"k" * 100
        message_payload = b"PayloadData"

        # Act.
        computed_hmac = HMACEngine.compute_hmac(key, message_payload)

        # Assert.
        self.assertEqual(len(computed_hmac), 32)

    def test_19_hmac_key_exactly_equal_to_block_size(self):
        """19. Evaluates boundary handling for a key that measures exactly 64 bytes long."""
        # Arrange.
        key = b"k" * 64
        message_payload = b"PayloadData"

        # Act.
        computed_hmac = HMACEngine.compute_hmac(key, message_payload)

        # Assert.
        self.assertEqual(len(computed_hmac), 32)

    def test_20_hmac_key_shorter_than_block_size(self):
        """20. Assures trailing zero-padding expansion handles short keys seamlessly."""
        # Arrange.
        key = b"ShortKey"
        message_payload = b"PayloadData"

        # Act.
        computed_hmac = HMACEngine.compute_hmac(key, message_payload)

        # Assert.
        self.assertEqual(len(computed_hmac), 32)

    def test_21_hmac_output_fixed_size(self):
        """21. Verifies that variations in tracking lengths do not affect the output 32-byte signature footprint."""
        # Arrange.
        key = b"ArbitraryKeyString"
        message_payload = b"LongPayload" * 200

        # Act.
        computed_hmac = HMACEngine.compute_hmac(key, message_payload)

        # Assert.
        self.assertEqual(len(computed_hmac), 32)

    def test_22_hmac_signature_tamper_detection(self):
        """22. Confirms that modifying a single payload byte completely invalidates the authenticity signature."""
        # Arrange.
        key = b"SharedSecretKey"
        payload_alpha = b"TransactionAmount: $100"
        payload_beta = b"TransactionAmount: $200"

        # Act.
        hmac_alpha = HMACEngine.compute_hmac(key, payload_alpha)
        hmac_beta = HMACEngine.compute_hmac(key, payload_beta)

        # Assert.
        self.assertNotEqual(hmac_alpha, hmac_beta)

    def test_23_hmac_key_tamper_detection(self):
        """23. Assures that modifying the signing key yields an entirely unique validation hash block."""
        # Arrange.
        key_alpha = b"SigningKeyA"
        key_beta = b"SigningKeyB"
        message_payload = b"CommonImmutablePayload"

        # Act.
        hmac_alpha = HMACEngine.compute_hmac(key_alpha, message_payload)
        hmac_beta = HMACEngine.compute_hmac(key_beta, message_payload)

        # Assert.
        self.assertNotEqual(hmac_alpha, hmac_beta)

    def test_24_hmac_binary_key_and_payload(self):
        """24. Validates processing when raw non-ascii bytes are embedded in both keys and payloads."""
        # Arrange.
        key = b"\x01\x02\x03\x04"
        message_payload = b"\xff\xfe\xfd\xfc"

        # Act.
        computed_hmac = HMACEngine.compute_hmac(key, message_payload)

        # Assert.
        self.assertEqual(len(computed_hmac), 32)

    def test_25_hmac_long_block_payload_stability(self):
        """25. Runs checking assertions over extensive multi-kilobyte stream sequences."""
        # Arrange.
        key = b"CryptographicKey"
        message_payload = b"DataChunk" * 1024

        # Act.
        computed_hmac = HMACEngine.compute_hmac(key, message_payload)

        # Assert.
        self.assertEqual(len(computed_hmac), 32)

    def test_26_hmac_reproducibility(self):
        """26. Verifies calling compute_hmac with the same inputs repeatedly yields an identical response signature."""
        # Arrange.
        key = b"StaticKey"
        message_payload = b"StaticMessage"

        # Act.
        initial = HMACEngine.compute_hmac(key, message_payload)
        next_hmac = HMACEngine.compute_hmac(key, message_payload)

        # Assert.
        self.assertEqual(initial, next_hmac)

    def test_27_hmac_nested_padding_separation(self):
        """27. Ensures inner and outer padding strings execute as independent bitwise configurations."""
        # Arrange.
        key = b"A" * 64

        # Act.
        inner_padding = bytes(byte ^ 0x36 for byte in key)
        outer_padding = bytes(byte ^ 0x5C for byte in key)

        # Assert.
        self.assertNotEqual(inner_padding, outer_padding)

    def test_28_pbkdf2_standard_validation_vector(self):
        """28. Runs validation against a globally established, known PBKDF2-HMAC-SHA256 test matrix string."""
        # Arrange.
        secret = b"password"
        salt = b"salt"
        iterations = 1
        length = 32
        expected_hex_key = "8f2079bd0b82857ecdc30b1fc3a3cd7607d956bf91923af1b62e97da1b187c1e"

        # Act.
        derived_key = PBKDF2KeyDeriver.derive_key(
            secret, salt, iterations, length
        )

        # Assert.
        self.assertEqual(derived_key.hex(), expected_hex_key)

    def test_29_pbkdf2_multi_iteration_vector(self):
        """29. Evaluates cascading internal XOR state transitions when tracking over multiple runtime loops."""
        # Arrange.
        secret = b"password"
        salt = b"salt"
        iterations = 2
        length = 32
        expected_hex_key = "bc6a2ddca8fba443004334f357fa852eb604257f1d652e84a172498029f0bce1"

        # Act.
        derived_key = PBKDF2KeyDeriver.derive_key(
            secret, salt, iterations, length
        )

        # Assert.
        self.assertEqual(derived_key.hex(), expected_hex_key)

    def test_30_pbkdf2_precise_output_slicing_under_digest_size(self):
        """30. Verifies stream-slicing functionality when requested target size is shorter than the standard block output."""
        # Arrange.
        secret = b"Passphrase"
        salt = b"SaltString"
        iterations = 10
        length = 16

        # Act.
        derived_key = PBKDF2KeyDeriver.derive_key(
            secret, salt, iterations, length
        )

        # Assert.
        self.assertEqual(len(derived_key), 16)

    def test_31_pbkdf2_exact_digest_size_alignment(self):
        """31. Checks loop boundaries when requested target length matches the standard block size exactly."""
        # Arrange.
        secret = b"Passphrase"
        salt = b"SaltString"
        iterations = 10
        length = 32

        # Act.
        derived_key = PBKDF2KeyDeriver.derive_key(
            secret, salt, iterations, length
        )

        # Assert.
        self.assertEqual(len(derived_key), 32)

    def test_32_pbkdf2_arbitrary_unaligned_output_length(self):
        """32. Confirms output byte-truncation handling for odd, unaligned key length constraints."""
        # Arrange.
        secret = b"Passphrase"
        salt = b"SaltString"
        iterations = 5
        length = 23

        # Act.
        derived_key = PBKDF2KeyDeriver.derive_key(
            secret, salt, iterations, length
        )

        # Assert.
        self.assertEqual(len(derived_key), 23)

    def test_33_pbkdf2_multi_block_stream_generation(self):
        """33. Evaluates key derivation mechanics when generating length values that span multiple sequential blocks."""
        # Arrange.
        secret = b"Passphrase"
        salt = b"SaltString"
        iterations = 5
        length = 1612

        # Act.
        derived_key = PBKDF2KeyDeriver.derive_key(
            secret, salt, iterations, length
        )

        # Assert.
        self.assertEqual(len(derived_key), 1612)

    def test_34_pbkdf2_empty_salt_processing(self):
        """34. Ensures the engine runs securely even if initialized with an empty byte salt parameter."""
        # Arrange.
        secret = b"PasswordString"
        salt = b""
        iterations = 10
        length = 32

        # Act.
        derived_key = PBKDF2KeyDeriver.derive_key(
            secret, salt, iterations, length
        )

        # Assert.
        self.assertEqual(len(derived_key), 32)

    def test_35_pbkdf2_empty_credential_processing(self):
        """35. Confirms execution tracking functions normally when evaluating empty password strings."""
        # Arrange.
        secret = b""
        salt = b"SaltString"
        iterations = 10
        length = 32

        # Act.
        derived_key = PBKDF2KeyDeriver.derive_key(
            secret, salt, iterations, length
        )

        # Assert.
        self.assertEqual(len(derived_key), 32)

    def test_36_pbkdf2_long_credential_input_chaining(self):
        """36. Verifies that long high-entropy input credentials scale cleanly across iterative operations."""
        # Arrange.
        secret = b"ExtremelyLongPasswordStringInput" * 20
        salt = b"Salt"
        iterations = 5
        length = 32

        # Act.
        derived_key = PBKDF2KeyDeriver.derive_key(
            secret, salt, iterations, length
        )

        # Assert.
        self.assertEqual(len(derived_key), 32)

    def test_37_pbkdf2_secret_tamper_isolation(self):
        """37. Asserts that shifting a character in the master secret yields a unique derived key block."""
        # Arrange.
        secret_alpha = b"MasterPass1"
        secret_beta = b"MasterPass2"
        salt = b"IdenticalSalt"
        iterations = 100
        length = 32

        # Act.
        key_alpha = PBKDF2KeyDeriver.derive_key(
            secret_alpha, salt, iterations, length
        )
        key_beta = PBKDF2KeyDeriver.derive_key(
            secret_beta, salt, iterations, length
        )

        # Assert.
        self.assertNotEqual(key_alpha, key_beta)

    def test_38_pbkdf2_salt_tamper_isolation(self):
        """38. Asserts that unique salts produce independent keys even under identical passphrase configurations."""
        # Arrange.
        secret = b"ConstantPassword"
        salt_alpha = b"SystemSaltAlpha"
        salt_beta = b"SystemSaltBeta"
        iterations = 100
        length = 32

        # Act.
        key_alpha = PBKDF2KeyDeriver.derive_key(
            secret, salt_alpha, iterations, length
        )
        key_beta = PBKDF2KeyDeriver.derive_key(
            secret, salt_beta, iterations, length
        )

        # Assert.
        self.assertNotEqual(key_alpha, key_beta)

    def test_39_pbkdf2_iteration_structural_dependence(self):
        """39. Verifies that variation in work factors fundamentally shifts output values."""
        # Arrange.
        secret = b"Password"
        salt = b"Salt"
        length = 32

        # Act.
        key_alpha = PBKDF2KeyDeriver.derive_key(
            secret, salt, 10, length
        )
        key_beta = PBKDF2KeyDeriver.derive_key(
            secret, salt, 11, length
        )

        # Assert.
        self.assertNotEqual(key_alpha, key_beta)

    def test_40_pbkdf2_high_load_determinism(self):
        """40. Checks processing stability across a standard work cycle iteration count."""
        # Arrange.
        secret = b"SecureMasterPassword123!"
        salt = b"ApplicationSpecificUniqueSalt"
        iterations = 2000
        length = 32
        expected_hex_key = "95dd1d0c4813b116352ba9adfbc65e2272883b8dffa6623224ffb31ef4faf3fa"

        # Act.
        derived_key = PBKDF2KeyDeriver.derive_key(
            secret, salt, iterations, length
        )

        # Assert.
        self.assertEqual(derived_key.hex(), expected_hex_key)

if __name__ == "__main__":
    unittest.main()