// 🔥 Concurrent HTTP benchmark script
// Run with: node benchmark.js

const TOTAL_REQUESTS = 100000;
const CONCURRENCY = 1000; // simultaneous requests

async function makeRequest() {
  const start = performance.now();
  try {
    const res = await fetch("http://localhost:8080/square");
    const data = await res.json();
    return {
      success: true,
      latencyMs: performance.now() - start,
      data,
    };
  } catch (error) {
    return {
      success: false,
      latencyMs: performance.now() - start,
      error: error.message,
    };
  }
}

async function runBatch(batchSize) {
  const promises = [];
  for (let i = 0; i < batchSize; i++) {
    promises.push(makeRequest());
  }
  return Promise.all(promises);
}

async function benchmark() {
  console.log("🔥 Starting benchmark...");
  console.log(`   Total requests: ${TOTAL_REQUESTS}`);
  console.log(`   Concurrency: ${CONCURRENCY}`);
  console.log("");

  const allResults = [];
  const startTime = performance.now();

  const batches = Math.ceil(TOTAL_REQUESTS / CONCURRENCY);

  for (let i = 0; i < batches; i++) {
    const batchSize = Math.min(CONCURRENCY, TOTAL_REQUESTS - i * CONCURRENCY);
    const results = await runBatch(batchSize);
    allResults.push(...results);

    // Progress update every 10%
    const progress = Math.floor(((i + 1) / batches) * 100);
    if (progress % 10 === 0) {
      process.stdout.write(`\r   Progress: ${progress}%`);
    }
  }

  const totalTime = performance.now() - startTime;

  // Calculate stats
  const successful = allResults.filter((r) => r.success);
  const failed = allResults.filter((r) => !r.success);
  const latencies = successful.map((r) => r.latencyMs).sort((a, b) => a - b);

  const avgLatency = latencies.reduce((a, b) => a + b, 0) / latencies.length;
  const p50 = latencies[Math.floor(latencies.length * 0.5)];
  const p95 = latencies[Math.floor(latencies.length * 0.95)];
  const p99 = latencies[Math.floor(latencies.length * 0.99)];
  const minLatency = latencies[0];
  const maxLatency = latencies[latencies.length - 1];

  const requestsPerSecond = (TOTAL_REQUESTS / (totalTime / 1000)).toFixed(2);

  console.log("\n");
  console.log("📊 RESULTS");
  console.log("═══════════════════════════════════════");
  console.log(`   Total time:        ${(totalTime / 1000).toFixed(2)}s`);
  console.log(`   Requests/sec:      ${requestsPerSecond} 🚀`);
  console.log(`   Successful:        ${successful.length}`);
  console.log(`   Failed:            ${failed.length}`);
  console.log("");
  console.log("⏱️  LATENCY");
  console.log("═══════════════════════════════════════");
  console.log(`   Min:               ${minLatency.toFixed(2)}ms`);
  console.log(`   Avg:               ${avgLatency.toFixed(2)}ms`);
  console.log(`   P50 (median):      ${p50.toFixed(2)}ms`);
  console.log(`   P95:               ${p95.toFixed(2)}ms`);
  console.log(`   P99:               ${p99.toFixed(2)}ms`);
  console.log(`   Max:               ${maxLatency.toFixed(2)}ms`);
  console.log("");

  if (parseFloat(requestsPerSecond) >= 20000) {
    console.log("👑 YASSS QUEEN! You hit 20k+ req/s! 💅🔥");
  } else if (parseFloat(requestsPerSecond) >= 10000) {
    console.log("🔥 Fire! 10k+ req/s - getting close to the goal!");
  } else if (parseFloat(requestsPerSecond) >= 5000) {
    console.log("💪 Solid! 5k+ req/s - room to optimize!");
  } else {
    console.log("🚧 Warming up! Try increasing concurrency or check the server.");
  }

  // Show a sample response
  if (successful.length > 0) {
    console.log("");
    console.log("📦 Sample response:");
    console.log(JSON.stringify(successful[0].data, null, 2));
  }
}

// Health check first
async function healthCheck() {
  try {
    const res = await fetch("http://localhost:8080/health");
    if (res.ok) {
      console.log("✅ Server is healthy!\n");
      return true;
    }
  } catch (e) {
    console.log("❌ Server not responding. Start it with:");
    console.log("   cd benchmark && go run server.go\n");
    return false;
  }
}

async function main() {
  const healthy = await healthCheck();
  if (healthy) {
    await benchmark();
  }
}

main();
