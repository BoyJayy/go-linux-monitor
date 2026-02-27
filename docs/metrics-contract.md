
json exmple

---
```json
{
  "device_id": "string",
  "ts": "RFC3339 string",
  "cpu": {
    "usage_percent": 0.0,
    "iowait_percent": 0.0,
    "load1": 0.0
  },
  "mem": {
    "total_bytes": 0,
    "available_bytes": 0,
    "used_percent": 0.0
  },
  "disk": {
    "path": "/",
    "total_bytes": 0,
    "free_bytes": 0,
    "used_percent": 0.0
  },
  "net": {
    "rx_bps": 0.0,
    "tx_bps": 0.0
  },
  "system": {
    "uptime_sec": 0,
    "process_count": 0,
    "cpu_temp_c": 0.0
  }
}
```
---

desc

## Root fields

### device_id
Type: string  
Required: yes  
Description: Unique and stable identifier of the device.

### ts
Type: RFC3339 string  
Required: yes  
Description: Timestamp when metrics snapshot was taken on the device.

---

## CPU

### cpu.usage_percent
Type: float (0..100)  
Unit: percent  
Description: CPU usage calculated over sampling interval.

### cpu.iowait_percent
Type: float (0..100)  
Unit: percent  
Description: Percentage of time CPU spent waiting for I/O during sampling interval.

### cpu.load1
Type: float  
Description: 1-minute load average.

---

## Memory

### mem.total_bytes
Type: uint64  
Unit: bytes  
Description: Total physical memory.

### mem.available_bytes
Type: uint64  
Unit: bytes  
Description: Available memory (MemAvailable).

### mem.used_percent
Type: float (0..100)  
Unit: percent  
Description: Percentage of used memory.

---

## Disk

### disk.path
Type: string  
Description: Mount point (V1 uses "/").

### disk.total_bytes
Type: uint64  
Unit: bytes  
Description: Total filesystem size.

### disk.free_bytes
Type: uint64  
Unit: bytes  
Description: Free space.

### disk.used_percent
Type: float (0..100)  
Unit: percent  
Description: Percentage of used disk space.

---

## Network

### net.rx_bps
Type: float  
Unit: bytes_per_sec  
Description: Incoming network traffic rate.

### net.tx_bps
Type: float  
Unit: bytes_per_sec  
Description: Outgoing network traffic rate.

---

## System

### system.uptime_sec
Type: uint64  
Unit: seconds  
Description: System uptime.

### system.process_count
Type: uint64  
Description: Number of running processes.

### system.cpu_temp_c
Type: float  
Unit: celsius  
Description: CPU temperature (if available).
