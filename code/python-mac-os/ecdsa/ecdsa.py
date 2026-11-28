"""
ecdsa.py

An implementation of the Elliptic Curve Digital Signature Algorithm (ECDSA)
configured for standard cryptographic signing and verification pipelines.

This module provides non-parametric asymmetric cryptographic workflows for key
generation, message hashing, signature synthesis, and verification topologies. It
explicitly combines elliptic curve scalar mathematics over prime fields with
one-way hashing primitives to achieve deterministic authentication criteria.

Design Note:
    This implementation is optimized for educational evaluation and structural reference.
    It natively handles large-integer math using pure Python primitives and lacks
    built-in side-channel mitigation techniques (such as constant-time execution)
    necessary for production-grade security architectures.
"""
from __future__ import annotations

import hashlib
import secrets


class EllipticCurve:
    """
    Represents an Elliptic Curve over a prime field (F_p).
    Equation: y^2 = x^3 + a*x + b (mod prime_modulus)
    """

    def __init__(self, prime_modulus: int, coefficient_linear: int, coefficient_constant: int):
        """
        Initializes an independent elliptic curve configuration.

        Args:
            prime_modulus (int): The large prime modulus defining the finite field boundaries.
            coefficient_linear (int): The 'a' parameter multiplying the linear term.
            coefficient_constant (int): The 'b' parameter defining the constant term.
        """
        # Set the prime modulus defining the field order.
        self.prime_modulus = prime_modulus

        # Assign the linear coefficient value for the curve equation.
        self.coefficient_linear = coefficient_linear

        # Assign the constant coefficient value for the curve equation.
        self.coefficient_constant = coefficient_constant


class Point:
    """
    Represents a specific coordinate pair (x, y) resting on an active elliptic curve,
    or an abstract representation of the identity element (the Point at Infinity).
    """

    def __init__(self, curve: EllipticCurve, x: int = None, y: int = None):
        """
        Initializes an independent coordinate point or an identity state at infinity.

        Args:
            curve (EllipticCurve): The mathematical curve domain this point belongs to.
            x (int, optional): The horizontal field coordinate. None implies infinity.
            y (int, optional): The vertical field coordinate. None implies infinity.
        """
        # Store the reference to the underlying mathematical curve.
        self.curve = curve

        # Store the horizontal coordinate value.
        self.x = x

        # Store the vertical coordinate value.
        self.y = y

        # Determine whether the coordinate represents the point at infinity.
        self.is_infinity = x is None or y is None

    @property
    def coordinate_x(self):
        """Exposes the strict integer horizontal coordinate, raising an exception if it is undefined."""
        # Ensure the point is not at infinity and the x-coordinate exists.
        if self.is_infinity or self.x is None:
            raise ValueError("The point at Infinity does not posses a finite X-coordinate.")

        # Return the stored horizontal coordinate integer.
        return self.x

    @property
    def coordinate_y(self):
        """Exposes the strict integer vertical coordinate, raising an exception if it is undefined."""
        # Ensure the point is not at infinity and the y-coordinate exists.
        if self.is_infinity or self.y is None:
            raise ValueError("The point at Infinity does not posses a finite Y-coordinate.")

        # Return the stored vertical coordinate integer.
        return self.y

    def __eq__(self, other: Point) -> bool:
        """Evaluates structural equivalence between two separate point nodes."""
        # Check if both points represent the point at infinity.
        if self.is_infinity and other.is_infinity:
            return True

        # Check if only one of the points is at infinity.
        if self.is_infinity != other.is_infinity:
            return False

        # Compare both horizontal and vertical coordinate values.
        return (self.x == other.x and self.y == other.y)

    def __add__(self, other: Point) -> Point:
        """
        Executes point addition over a finite field to synthesize a combined coordinate point.

        Mathematical Formula:
            $$\\lambda = \\frac{y_2 - y_1}{x_2 - x_1} \\pmod p \\quad \\text{or} \\quad \\lambda = \\frac{3x_1^2 + a}{2y_1} \\pmod p$$
            $$x_3 = \\lambda^2 - x_1 - x_2 \\pmod p, \\quad y_3 = \\lambda(x_1 - x_3) - y_1 \\pmod p$$

        Args:
            other (Point): The secondary point configuration being added to the current point.

        Returns:
            Point: A newly instantiated point representing the chord-tangent field sum.
        """
        # Return the secondary point if the current point is the identity element.
        if self.is_infinity:
            return other

        # Return the current point if the secondary point is the identity element.
        if other.is_infinity:
            return self

        # Ensure both coordinates are fully instantiated finite values.
        assert self.x is not None and self.y is not None
        assert other.x is not None and other.y is not None

        # Extract the prime modulus from the associated curve context.
        prime = self.curve.prime_modulus

        # Return the point at infinity if adding a point to its vertical reflection.
        if self.x == other.x and (self.y + other.y) % prime == 0:
            return Point(self.curve)

        # Handle distinct points using standard chord-intercept addition formulas.
        if self != other:
            # Calculate the slope numerator using vertical differences.
            numerator = (other.y - self.y) % prime

            # Calculate the slope denominator using horizontal differences.
            denominator = (other.x - self.x) % prime

            # Compute the field slope using modular multiplicative inverse.
            slope = (numerator * pow(denominator, prime - 2, prime))
        else:
            # Handle identical points using point doubling / tangent addition formulas.
            # Calculate the tangent slope numerator using curve coefficients.
            numerator = (3 * pow(self.coordinate_x, 2, prime) + self.curve.coefficient_linear) % prime

            # Calculate the tangent slope denominator using the vertical coordinate.
            denominator = (2 * self.y) % prime

            # Compute the tangent field slope using modular multiplicative inverse.
            slope = (numerator * pow(denominator, prime - 2, prime))

        # Compute the resulting horizontal coordinate on the curve.
        new_x = (pow(slope, 2, prime) - self.x - other.x) % prime

        # Compute the resulting vertical coordinate on the curve.
        new_y = (slope * (self.x - new_x) - self.y) % prime

        # Return the newly constructed point on the elliptic curve.
        return Point(self.curve, new_x, new_y)

    def __rmul__(self, scalar: int) -> Point:
        """
        Multiplies an existing point by an integer scalar via the Double-and-Add heuristic.

        Args:
            scalar (int): The integer multiplier scaling the target curve coordinate.

        Returns:
            Point: The terminal point node produced by cyclic addition loops.
        """
        # Initialize the result accumulator as the point at infinity.
        result = Point(self.curve)

        # Set the initial addend point to the current point instance.
        addend = self

        # Copy the scalar multiplier value for bitwise iteration.
        bits = scalar

        # Iterate through all active bits in the scalar multiplier.
        while bits > 0:
            # Accumulate the current addend if the lowest bit is set.
            if bits & 1:
                result = result + addend

            # Double the addend point for the next bit position.
            addend = addend + addend

            # Shift the bitwise tracker right by one position.
            bits >>= 1

        # Return the final scaled point resulting from the multiplication.
        return result


class ECDSA:
    """
    An engine responsible for managing cryptographic keys, generating verifiable message
    signatures, and executing mathematical verification tests over a defined curve.
    """

    def __init__(self):
        """
        Initializes the cryptographic pipeline using standard secp256k1 base constraints.
        """
        # Define the secp256k1 prime modulus constant.
        prime_modulus = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEFFFFFC2F

        # Define the linear curve coefficient (zero for secp256k1).
        coefficient_linear = 0

        # Define the constant curve coefficient (seven for secp256k1).
        coefficient_constant = 7

        # Instantiate the elliptic curve specification.
        self.curve = EllipticCurve(prime_modulus, coefficient_linear, coefficient_constant)

        # Define the X coordinate for the base generator point.
        generator_x = 0x79BE667EF9DCBBAC55A06295CE870B07029BFCDB2DCE28D959F2815B16F81798

        # Define the Y coordinate for the base generator point.
        generator_y = 0x483ADA7726A3C4655DA4FBFC0E1108A8FD17B448A68554199C47D08FFB10D4B8

        # Instantiate the base generator point on the curve.
        self.generator_point = Point(self.curve, generator_x, generator_y)

        # Define the order of the generator subgroup.
        self.order = 0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141

    def _hash_message(self, message: bytes) -> int:
        """
        Transforms raw byte input vectors into uniform scalar integers using a SHA-256 digest.

        Args:
            message (bytes): The raw sequence of data bytes to evaluate.

        Returns:
            int: A large integer representation derived directly from the computed hash.
        """
        # Compute the cryptographic SHA-256 hash digest of the input message bytes.
        message_hash_bytes = hashlib.sha256(message).digest()

        # Convert the raw hash digest bytes into a large integer value.
        return int.from_bytes(message_hash_bytes, byteorder='big')

    def generate_key_pair(self) -> tuple[int, Point]:
        """
        Generates an independent private integer scalar and its corresponding public point.

        Returns:
            tuple[int, Point]: A pair containing the private key and public key point.
        """
        # Generate a cryptographically secure random integer within valid curve range.
        private_key = secrets.randbelow(self.order - 1) + 1

        # Derive the corresponding public key point via scalar multiplication.
        public_key = private_key * self.generator_point

        # Return the generated private scalar and public coordinate point pair.
        return private_key, public_key

    def sign_message(self, message: bytes, private_key: int) -> tuple[int, int]:
        """
        Computes a non-invertible digital signature (r, s) utilizing the private key state.

        Optimization Objective:
            $$r = (k \\cdot G)_x \\pmod n$$
            $$s = k^{-1} \\cdot (\\text{hash} + r \\cdot \\text{private\\_key}) \\pmod n$$

        Args:
            message (bytes): The raw data matrix passing through the signature pipeline.
            private_key (int): The secret key integer controlling identity verification.

        Returns:
            tuple[int, int]: The cryptographic signature components (r_value, s_value).
        """
        # Compute the numerical message digest integer.
        message_digest = self._hash_message(message)

        # Initialize signature proof assertion component (r).
        signature_proof_assertion = 0

        # Initialize signature witness component (s).
        signature_witness = 0

        # Loop until valid non-zero signature parameters are successfully generated.
        while signature_proof_assertion == 0 or signature_witness == 0:
            # Generate a secure random ephemeral nonce scalar.
            nonce = secrets.randbelow(self.order - 1) + 1

            # Compute the ephemeral point on the curve using the nonce.
            ephemeral_point = nonce * self.generator_point

            # Extract the x-coordinate modulo curve order as the r component.
            signature_proof_assertion = ephemeral_point.x % self.order

            # Retry if the r component evaluates to zero.
            if signature_proof_assertion == 0:
                continue

            # Compute the modular inverse of the ephemeral nonce.
            nonce_inverse = pow(nonce, self.order - 2, self.order)

            # Calculate the s component using the message digest, private key, and r component.
            signature_witness = ((nonce_inverse * (message_digest + signature_proof_assertion * private_key))
                                 % self.order)

        # Return the finalized (r, s) cryptographic signature pair.
        return signature_proof_assertion, signature_witness

    def verify_signature(self, message: bytes, signature: tuple[int, int], public_key: Point) -> bool:
        """
        Evaluates an external signature pair against a known public key validation path.

        Verification Criteria:
            $$u_1 = \\text{hash} \\cdot s^{-1} \\pmod n, \\quad u_2 = r \\cdot s^{-1} \\pmod n$$
            $$\\text{Verify that: } (u_1 \\cdot G + u_2 \\cdot Q)_x \\equiv r \\pmod n$$

        Args:
            message (bytes): The data vector claiming validation authority.
            signature (tuple[int, int]): A cluster comprising the (r, s) values under review.
            public_key (Point): The unencrypted validation target point.

        Returns:
            bool: True if structural verification metrics match; False if the payload was altered.
        """
        # Unpack the signature components into r and s variables.
        signature_proof_assertion, signature_witness = signature

        # Validate that the signature components lie within valid field boundaries.
        if not (0 < signature_proof_assertion < self.order) or not (0 < signature_witness < self.order):
            return False

        # Compute the message digest integer for verification.
        message_digest = self._hash_message(message)

        # Compute the modular inverse of the signature witness component (s).
        witness_inverse = pow(signature_witness, self.order - 2, self.order)

        # Calculate the scalar multiplier for the base generator point.
        generator_scalar_multiplier = (message_digest * witness_inverse) % self.order

        # Calculate the scalar multiplier for the public key point.
        public_key_scalar_multiplier = (signature_proof_assertion * witness_inverse) % self.order

        # Compute the linear combination point using both scalar multipliers.
        point = (generator_scalar_multiplier * self.generator_point) + (public_key_scalar_multiplier * public_key)

        # Return false if the resulting coordinate evaluates to the point at infinity.
        if point.is_infinity:
            return False

        # Return true if the x-coordinate modulo order matches the signature r component.
        return (point.x % self.order) == signature_proof_assertion