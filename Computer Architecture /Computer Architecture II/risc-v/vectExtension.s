    .section .data
# Define two arrays of single-precision floating point numbers and reserve space for the result.
vecA:   .float 1.0, 2.0, 3.0, 4.0       # First input vector with 4 elements.
vecB:   .float 5.0, 6.0, 7.0, 8.0       # Second input vector with 4 elements.
vecC:   .space 16                      # Reserve space for 4 floats (4 x 4 bytes)

    .section .text
    .global main
main:
    # -----------------------------------------------------------
    # Set Vector Length:
    # Configure the vector unit to process 4 elements of 32-bit floats (e32)
    # with vector register group multiplier m1.
    # The 'vsetvli' instruction sets the vector length (VL) based on the requested number of elements.
    # -----------------------------------------------------------
    li      t0, 4                    # Load the number of elements (4) into t0.
    vsetvli t0, t0, e32, m1           # Set VL for 4 elements, each 32 bits wide.

    # -----------------------------------------------------------
    # Load Base Addresses:
    # Load the memory addresses of vecA, vecB, and vecC into registers for vector operations.
    # -----------------------------------------------------------
    la      t1, vecA                 # Load address of vecA into t1.
    la      t2, vecB                 # Load address of vecB into t2.
    la      t3, vecC                 # Load address of vecC into t3.

    # -----------------------------------------------------------
    # Load Vectors:
    # Load 4-element vectors from memory into vector registers using 'vle.v'.
    # -----------------------------------------------------------
    vle.v   v0, (t1)                 # Load vector from vecA into vector register v0.
    vle.v   v1, (t2)                 # Load vector from vecB into vector register v1.

    # -----------------------------------------------------------
    # Vector Floating Point Addition:
    # Perform element-wise addition of vectors in v0 and v1.
    # The result is stored in vector register v2.
    # Expected result: [1.0+5.0, 2.0+6.0, 3.0+7.0, 4.0+8.0] = [6.0, 8.0, 10.0, 12.0]
    # -----------------------------------------------------------
    vfadd.vv v2, v0, v1              # Add vectors v0 and v1 element-wise, store result in v2.

    # -----------------------------------------------------------
    # Vector Floating Point Multiplication (Optional):
    # Perform element-wise multiplication of vectors in v0 and v1.
    # The result is stored in vector register v3.
    # Expected result: [1.0*5.0, 2.0*6.0, 3.0*7.0, 4.0*8.0] = [5.0, 12.0, 21.0, 32.0]
    # -----------------------------------------------------------
    vfmul.vv v3, v0, v1              # Multiply vectors v0 and v1 element-wise, store result in v3.

    # -----------------------------------------------------------
    # Store Result:
    # Store the resulting vector from addition (v2) back into memory at vecC.
    # -----------------------------------------------------------
    vse.v   v2, (t3)                 # Store vector v2 (result of addition) into vecC.

    # -----------------------------------------------------------
    # Exit the Program:
    # Use the exit system call for RISC-V Linux to terminate the program.
    # -----------------------------------------------------------
    li      a7, 93                 # Load the syscall number for exit (93) into a7.
    li      a0, 0                  # Set the exit code to 0 (successful termination).
    ecall                          # Make the system call to exit.
