# Exam Solutions

## Exercise 1: CPU Analysis with RV32IMAFD

## Code in Risc-V
```riscv
fdiv.s f1, f10, f2
fdiv.s f6, f1, f5
fadd.s f6, f2, f3
fadd.s f3, f8, f2
fidv.s f4, f6, f1
fmul.s f6, f10, f5
```

### 1. Dependency Graph
- **Task**: Draw the dependency graph for the given code.
- **Objective**: Deduce the minimum number of clock cycles needed, considering only the dependencies.
- **Solution**: 

```riscv
fdiv.s f1, f10, f2    
fdiv.s f6, f1, f5     # RAW hazard: f6 depends on f1
fadd.s f6, f2, f3     # WAW hazard: f6 is written here, overwriting its value from fdiv.s
fadd.s f3, f8, f2     # WAR hazard: f3 is read before being written in the next instruction
fidv.s f4, f6, f1     # RAW hazard: f4 depends on f6 and f1
fmul.s f6, f10, f5    # WAR/WAR hazard: f6 is written here, overwriting its previous value, and it could be done before the previous instruction
```
| Functional Unit         | Execution Time         | Reservation Stations |
| :-----------------------| :--------------------- | :--------------------|
| FADD                    | 2                      | 2                    |
| FMUL                    | 3                      | 2                    |
| FDIV                    | 4                      | 2                    |


### 2. Execution Dynamics with 1 CDB
- **Setup**:
  - Single CDB.
  - One IF and ID stage.
  - Two RS for each functional unit.
  - WB of the D unit takes priority in conflicts.
- **Task**: Illustrate the instruction scheduling and execution dynamics.

### 3. Optimizations
- **Possible Modifications**:
  - Increase RS to three.
  - Double/triple the CDB.
  - Add functional units with two RS each.
- **Task**: Identify the optimal architecture modification to minimize clock cycles. Justify and illustrate the new execution dynamics.

### 4. Register f6 Dynamics
- **Task**: Draw the timeline of the f6 register during execution.

---

## Exercise 2: Multicore System and Cache Analysis
Let A and B be arrays of 512 double-precision floating-point elements starting at addresses 0x0000_1000 and 0x0000_2000, respectively, and 
C and D be variables of the same type starting at addresses 0x0000_3000 and 0x0000_4000, respectively.

LSB: Least Significant Byte                    
MSB: Most Significant Byte
### 1. Memory Mapping
#### 1.1 Address Ranges
- **Task**: Specify the LSB and MSB for the first and last elements of vectors A, B, and variables C, D.
- **Solution**: Keep in mind that address are in **hexadecimal**
   - $512 \text{ elements} \times 64 \text{ bits per element} = 32,768 \text{ bits}$  
   - $32,768 / 8 = 4,096 \text{ bytes} = 4 \, \text{KiB}$ 
   - $4096_{10} = 1000_{16}$  

- A:
    - $\text{A}[0]$ goes from $0x0000\_1000$ (LSB) to $0x0000\_1000 + 8 -1 = 0x0000\_1007$ (MSB)
    - $\text{A}[511]$ ends at  $0x0000\_1000+1000-1=0x0000\_2000-1=0x0000\_1FFF$ (MSB) and starts in $0x0000\_1FFF-7=0x0000\_1FF7$ (LSB)
- B:
    - $\text{B}[0]$ goes from $0x0000\_2000$ (LSB) to $0x0000\_2000 + 8 -1 = 0x0000\_2007$ (MSB)
    - $\text{B}[511]$ ends at  $0x0000\_2000+1000-1=0x0000\_3000-1=0x0000\_2FFF$ (MSB) and starts in $0x0000\_2FFF-7=0x0000\_2FF7$ (LSB)
- C:
    - $\text{C}$ goes from $0x0000\_3000$ (LSB) to $0x0000\_3000 + 8 -1 = 0x0000\_3007$ (MSB)
- D:
     - $\text{D}$ goes from $0x0000\_4000$ (LSB) to $0x0000\_4000 + 8 -1 = 0x0000\_4007$ (MSB)

#### 1.2 Cache Tag and Set ID
- **Task**: Determine the number and value of TAG and Set ID bits for the above memory regions.
- **Solution**: 2-way Cache with 8KiB and a 32B cacheline, so:
    - $8\text{ KiB}=8192\text{ Byte}$
    - $8192\text{ B}/32\text{ B}=256\text{ Blocks}$
    - $256\text{ Blocks} / 2\text{ way} = 128\text{ Sets}$
    - $\text{\# Bit SetId}=\log_2 128=7$
    - $\text{\# Bit Offset}=\log_2 32=5$
    - $\text{\# Bit TAG}= 32 -7 -5 = 20$

How to write down the values of TAG and SetId:
- $\text{A}[0]= 0x0000\_1000=0000\ 0000\ 0000\ 0000\ 0001\ 0000\ 0000\ 0000$
    - From 0 to 4 we have the Offset: $00000$
    - From 5 to 11 we have the SetId: $0000000 = 0x0$
    - From 12 to 31 we have the TAG:  $0000000000000001=0x0000\ 1$
- $\text{A}[512]= 0x0000\_1FFF=0000\ 0000\ 0000\ 0000\ 0001\ 1111\ 1111\ 1111$
    - From 0 to 4 we have the Offset: $11111$
    - From 5 to 11 we have the SetId: $1111111 = 0x7F$
    - From 12 to 31 we have the TAG:  $0000000000000001=0x0000\ 1$


#### 1.3 Memory and Cache Fit
- **Task**: Compute the total memory size occupied by A, B, C, and D. Assess whether all can fit in Core 0's cache.

### 2. Cache Dynamics
#### 2.1 Cache and MESI States
- **Task**: Show the MESI states and LRU bit changes during the first loop execution.

#### 2.2 Cache Final State
- **Task**: Display the final cache contents, MESI states, and LRU bits.

#### 2.3 HIT/MISS Analysis
- **Task**: Calculate access count, HIT, MISS, HIT rate, MISS rate, and WB cycles for each core.

#### 2.4 Impact of Moving `D`
- **Task**: Describe qualitatively how the analysis changes if `D` is relocated to `0x0000_3010`.

---