// Node.js server mimicking the unoptimized Go server
// Same pattern: workers, job queue, collect results, return random

const http = require("http");

const NUM_WORKERS = 10;

// Simulate a "channel" with async queue
class Channel {
  constructor() {
    this.queue = [];
    this.resolvers = [];
  }

  send(value) {
    if (this.resolvers.length > 0) {
      const resolve = this.resolvers.shift();
      resolve({ value, done: false });
    } else {
      this.queue.push(value);
    }
  }

  receive() {
    return new Promise((resolve) => {
      if (this.queue.length > 0) {
        resolve({ value: this.queue.shift(), done: false });
      } else {
        this.resolvers.push(resolve);
      }
    });
  }

  close() {
    this.resolvers.forEach((resolve) => resolve({ value: null, done: true }));
  }
}

// Worker function - like a goroutine
async function worker(jobsChan, resultsChan) {
  while (true) {
    const { value: job, done } = await jobsChan.receive();
    if (done) break;

    resultsChan.send({
      number: job,
      square: job * job,
      workers_used: NUM_WORKERS,
    });
  }
}

async function handleRequest(req, res) {
  if (req.url === "/health") {
    res.writeHead(200);
    res.end("OK");
    return;
  }

  if (req.url !== "/square") {
    res.writeHead(404);
    res.end("Not found");
    return;
  }

  const start = process.hrtime.bigint();

  // Create "channels"
  const jobsChan = new Channel();
  const resultsChan = new Channel();

  // Start workers (like goroutines)
  const workerPromises = [];
  for (let i = 0; i < NUM_WORKERS; i++) {
    workerPromises.push(worker(jobsChan, resultsChan));
  }

  // Send jobs
  for (let i = 0; i <= 100; i++) {
    jobsChan.send(i);
  }

  // Collect results
  const results = [];
  for (let i = 0; i <= 100; i++) {
    const { value } = await resultsChan.receive();
    results.push(value);
  }

  // Close jobs channel to stop workers
  jobsChan.close();

  // Wait for workers to finish
  await Promise.all(workerPromises);

  // Pick random result
  const randomResult = results[Math.floor(Math.random() * results.length)];
  randomResult.time_ns = Number(process.hrtime.bigint() - start);

  res.writeHead(200, { "Content-Type": "application/json" });
  res.end(JSON.stringify(randomResult));
}

const server = http.createServer(handleRequest);

server.listen(8080, () => {
  console.log("🟡 Node.js server running on http://localhost:8080");
  console.log("   GET /square - returns random square (with workers)");
  console.log("   GET /health - health check");
  console.log("\nBenchmark with:");
  console.log("   ~/go/bin/hey -n 100000 -c 1000 http://localhost:8080/square");
});
