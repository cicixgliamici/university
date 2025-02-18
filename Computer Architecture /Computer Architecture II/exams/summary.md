# Computer Architecture 2 Summary

This document provides a concise overview of key concepts related to modern CPU pipelines, cache organization, and coherence protocols.

---

## 1. Dependencies: RAW, WAW, and WAR

1. **RAW (Read After Write)**
   - **Definition**: Instruction $I_n$ needs to read a register/memory location **after** it has been written by instruction $I_m$.
   - **Impact**: The pipeline must ensure $I_m$ completes its write **before** $I_n$ can read the updated value.
   - **Also Known As**: True dependency or flow dependency.

2. **WAW (Write After Write)**
   - **Definition**: Instruction $I_n$ writes to a register/memory location that is also written by an earlier instruction $I_m$.
   - **Impact**: We must preserve the order of the two writes to ensure the final state is correct (the last writer “wins”).
   - **Common in Out-of-Order Execution**: In simple in-order pipelines, WAW typically does not occur, but it does arise when instructions are reordered.

3. **WAR (Write After Read)**
   - **Definition**: Instruction $I_n$ writes to a register/memory location that is read by a subsequent instruction $I_m$.
   - **Impact**: If execution is re-ordered, $I_n$ might overwrite the value before $I_m$ reads it.  
   - **Often Avoided**: In most in-order pipelines, WAR hazards are not present. In out-of-order pipelines (like Tomasulo), WAR hazards can be a concern.

---

## 2. Tomasulo's Algorithm

**Goal**  
Tomasulo’s algorithm is designed to enable out-of-order execution of instructions in a pipeline while ensuring correctness (no hazards).

### Key Components

- **Reservation Stations (RS)**  
  Each functional unit (e.g., adder, multiplier) has multiple RS entries. An instruction will occupy a reservation station until its operands are ready and the functional unit is available.

- **Common Data Bus (CDB)**  
  The CDB is used to broadcast results from the functional units back to the reservation stations and to the registers. Only one result at a time can typically be placed on a single CDB.

- **Register Renaming**  
  Tomasulo uses reservation stations to rename registers, thereby avoiding WAW and WAR hazards at the architectural register level. Internally, each reservation station can hold the data or the “tag” of where the data will come from.  

  - **How Register Renaming Works in Tomasulo**  
    1. When an instruction is **issued**, it is allocated a reservation station that keeps track of where its operands come from (either from a register or from another RS that will produce the data).  
    2. If a source register is not yet available, the reservation station records the “tag” of the unit that will produce that register’s value.  
    3. Meanwhile, the instruction’s destination register is marked as “pending” or renamed so that future instructions will also wait for the correct tag rather than the old architectural register.  
    4. When the producing instruction **writes to the CDB**, all reservation stations watching for that tag latch the value. This mechanism resolves potential WAW and WAR hazards, because multiple instructions can be “in flight” without overwriting each other’s results in the architectural registers.

### Execution Steps

1. **Issue**  
   - The instruction is decoded, and an available reservation station is allocated. If none are free, the pipeline stalls.
2. **Execute**  
   - Once all operands are available (either from the register file or from the CDB), the instruction executes on its functional unit.  
   - During execution, the reservation station is busy until completion.
3. **Write Result (WB)**  
   - Upon completing execution, the functional unit places the result on the CDB. All reservation stations that are waiting for that result will capture it and become ready to execute when their time comes.
   - The register file is also updated with the final value (if this is the latest instruction writing to that register).

---

## 3. Cache Addressing: Offset, Index (Set ID), and Tag

In a typical cache organization, a memory address is divided into three parts:

1. **Offset (Block Offset)**  
   - Number of bits: $\log_2(\text{Cache line size})$ 
   - Identifies the byte within a single cache line. For example, if the cache line size is 64 bytes, you need $\log_2(64)$ bits for the offset.

2. **Index (Set Index)**  
   - Number of bits: $\log_2(\text{Number of sets})$
   - Selects which set in the cache will be used. For instance, if there are 128 sets, $\log_2(128)$ bits are needed.

3. **Tag**  
   - Number of bits:  
     $
       \text{(Total address bits)} - (\text{Offset bits} + \text{Index bits})
     $
   - The tag is compared to the stored tag in the cache line to determine if the line contains the requested address (i.e., hit or miss).

---

## 4. MESI Protocol

**Purpose**  
The MESI protocol maintains cache coherence across multiple cores. Each cache line can be in one of four states:

1. **M (Modified)**  
   - The line is valid and has been modified in this cache (dirty).  
   - Only one cache can hold a line in the Modified state.  
   - If evicted, it must be written back to memory.

2. **E (Exclusive)**  
   - The line is valid and clean; only this cache has it.  
   - If written to, it moves to Modified (M).  
   - No other copies exist in any other cache.

3. **S (Shared)**  
   - The line is valid and clean; it may be replicated in multiple caches.  
   - If written to, a transition to M may require invalidating or updating other caches.

4. **I (Invalid)**  
   - The line is not valid in this cache.  
   - Any access to this line triggers a miss.

---

## 5. LRU (Least Recently Used) Replacement

When a cache is **set-associative**, each set can hold multiple lines (ways). Upon a miss, the cache must choose which line to evict if all ways are full.

- **Least Recently Used (LRU)**: The line that was used (accessed) farthest in the past is replaced first.  
- **Other Policies**: FIFO (First In First Out), Random, etc., but LRU is very common in practice.

### LRU in a 2-Way Set-Associative Cache

In a **2-way** set-associative cache, LRU can be managed with **just one bit** per set:
- This bit indicates which of the two ways was used **most recently**.
- Whenever one way is accessed, the LRU bit is toggled to mark that way as the most recently used.
- Upon a miss requiring eviction, the cache checks the LRU bit to see which way is the **least** recently used (i.e., the other way) and evicts that line.

---

## 6. Miss Penalty

**Definition**  
A miss penalty is the extra time (in cycles) required to fetch a block from the next level of the memory hierarchy when a cache miss occurs.

1. **Types of Misses**  
   - **Cold (Compulsory)**: The first access to a block that has never been in cache.  
   - **Conflict**: A miss that occurs because multiple addresses map to the same set and evict each other repeatedly.  
   - **Capacity**: A miss that occurs because the cache cannot contain all the data needed at once.

2. **Calculation**  
   - Typically, the miss penalty includes:  
     1. Time to send the address to memory.  
     2. Memory access latency.  
     3. Transfer time of the entire cache block from memory.

   - If the block was in the **Modified (M)** state in another core’s cache, additional steps (like write-back) might be necessary before the new core can have the line.

---

## 7. Cache Write Policies

There are two main categories of write policies in caches:

1. **Write-Through**  
   - Every write to the cache line is also immediately written to the next level of memory.  
   - Ensures data in lower-level memory is always up to date.  
   - Generally simpler coherence but higher bandwidth usage because every write goes to memory.

2. **Write-Back**  
   - Writes go only to the cache line in the local cache, and the line is marked “dirty.”  
   - The updated data is written back to memory **only when** the line is evicted.  
   - Can reduce bandwidth usage, but the memory may be stale until eviction.

---

## 8. Floating-Point Representation

Two primary floating-point formats are commonly used (IEEE 754 standard):

1. **Single Precision (float)**  
   - **Size**: 32 bits total  
   - **In Bytes**: 4 bytes  
   - Typically, 1 bit for sign, 8 bits for exponent, 23 bits for fraction (mantissa).

2. **Double Precision (double)**  
   - **Size**: 64 bits total  
   - **In Bytes**: 8 bytes  
   - Typically, 1 bit for sign, 11 bits for exponent, 52 bits for fraction (mantissa).

---
