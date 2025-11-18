# Kafka Writer Configuration Guide

## Current Configuration Analysis

### Current Settings:
- `BatchTimeout: 50ms` - ✅ Good for low latency
- `BatchSize: 10` - ⚠️ Could be optimized
- `RequiredAcks: RequireOne` - ✅ Good balance
- `Async: true` - ✅ Good for performance
- `Compression: Snappy` (only notification-service) - ✅ Good choice

## Recommended Production Configuration

### For High-Throughput Events (Chat, Notifications):
```go
writer := &kafka.Writer{
    Addr:         kafka.TCP(brokers...),
    Topic:        topic,
    Balancer:     &kafka.LeastBytes{},
    
    // Batching Configuration
    BatchSize:    100,              // Batch up to 100 messages
    BatchBytes:   1048576,          // 1MB max batch size
    BatchTimeout: 10 * time.Millisecond, // Flush every 10ms if batch incomplete
    
    // Reliability
    RequiredAcks: kafka.RequireOne, // Wait for leader ack (good balance)
    MaxAttempts:  10,               // Retry up to 10 times
    
    // Timeouts
    ReadTimeout:  10 * time.Second,
    WriteTimeout: 10 * time.Second,
    
    // Backoff for retries
    WriteBackoffMin: 100 * time.Millisecond,
    WriteBackoffMax: 1 * time.Second,
    
    // Performance
    Async: true,                    // Non-blocking writes
    
    // Compression
    Compression: kafka.Snappy,      // Best balance of speed/ratio
}
```

### For Low-Latency Critical Events (User Created):
```go
writer := &kafka.Writer{
    Addr:         kafka.TCP(brokers...),
    Topic:        topic,
    Balancer:     &kafka.LeastBytes{},
    
    // Smaller batches for lower latency
    BatchSize:    1,                // Send immediately
    BatchBytes:   1048576,          // 1MB max
    BatchTimeout: 0,                // No batching delay
    
    // Reliability
    RequiredAcks: kafka.RequireOne,
    MaxAttempts:  10,
    
    // Timeouts
    ReadTimeout:  5 * time.Second,  // Shorter timeout
    WriteTimeout: 5 * time.Second,
    
    // Backoff
    WriteBackoffMin: 50 * time.Millisecond,
    WriteBackoffMax: 500 * time.Millisecond,
    
    // Performance
    Async: true,
    
    // Compression (optional for small messages)
    Compression: kafka.Snappy,
}
```

## Compression Comparison

### Available Options:
1. **None** (`kafka.CompressionNone`)
   - Fastest encoding/decoding
   - No compression overhead
   - Best for: Very small messages (< 100 bytes)

2. **GZIP** (`kafka.CompressionGZIP`)
   - Best compression ratio (60-80% reduction)
   - Slowest encoding/decoding
   - Higher CPU usage
   - Best for: Large messages, archival, bandwidth-limited scenarios

3. **Snappy** (`kafka.CompressionSnappy`)
   - ✅ **RECOMMENDED** - Best balance
   - Good compression ratio (40-60% reduction)
   - Fast encoding/decoding
   - Low CPU usage
   - Best for: Real-time events, most use cases

4. **LZ4** (`kafka.CompressionLZ4`)
   - Fast compression (faster than Snappy)
   - Good compression ratio (30-50% reduction)
   - Best for: High-throughput scenarios

5. **Zstd** (`kafka.CompressionZstd`)
   - Excellent compression ratio (50-70% reduction)
   - Fast decoding, slower encoding
   - Best for: Storage optimization, high compression needs

### Recommendation by Event Type:

**Chat Messages:**
- Use **Snappy** - Good balance for JSON payloads (typically 200-500 bytes)
- Reduces network traffic by ~40-50%
- Minimal CPU overhead

**Notifications:**
- Use **Snappy** - Similar payload size to chat
- Consistent compression across event types

**User Created:**
- Use **Snappy** or **None** - Small payloads, compression may not be worth it
- If payload < 100 bytes, consider no compression

## Configuration Tuning Guidelines

### BatchSize:
- **Small (1-10)**: Low latency, higher overhead
- **Medium (10-100)**: ✅ Good balance
- **Large (100-1000)**: High throughput, higher latency

### BatchTimeout:
- **0ms**: Send immediately (no batching)
- **10-50ms**: ✅ Good for real-time events
- **100ms+**: Better throughput, higher latency

### RequiredAcks:
- **RequireNone (0)**: Fire-and-forget, fastest but no guarantee
- **RequireOne (1)**: ✅ Recommended - Leader confirms, good balance
- **RequireAll (-1)**: Highest reliability, slowest (waits for all replicas)

### Async:
- **true**: ✅ Recommended - Non-blocking, better performance
- **false**: Blocking, synchronous writes

## Performance Impact

### Current Config (50ms timeout, BatchSize 10, Async true):
- ✅ Good for low latency
- ✅ Non-blocking writes
- ⚠️ Could batch more for better throughput

### Recommended Config (10ms timeout, BatchSize 100, Async true):
- ✅ Better throughput
- ✅ Still low latency
- ✅ Better network efficiency
- ✅ Lower CPU per message

## Monitoring Metrics to Watch

1. **Write latency**: Time from WriteMessages() to completion
2. **Batch fill rate**: How often batches reach BatchSize before timeout
3. **Compression ratio**: Size reduction achieved
4. **CPU usage**: Compression overhead
5. **Network bandwidth**: Bytes sent vs uncompressed

