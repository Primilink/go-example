<?php
// PHP server mimicking the unoptimized Go server
// Same pattern: workers, job queue, collect results, return random

define('NUM_WORKERS', 10);

// Simulate a "channel" with a queue
class Channel {
    private array $queue = [];
    private bool $closed = false;

    public function send(mixed $value): void {
        $this->queue[] = $value;
    }

    public function receive(): ?array {
        if (empty($this->queue) && $this->closed) {
            return ['value' => null, 'done' => true];
        }
        if (empty($this->queue)) {
            return ['value' => null, 'done' => true];
        }
        return ['value' => array_shift($this->queue), 'done' => false];
    }

    public function close(): void {
        $this->closed = true;
    }

    public function hasData(): bool {
        return !empty($this->queue);
    }
}

// Worker function - simulating goroutine work
function runWorkers(Channel $jobsChan, Channel $resultsChan): void {
    // Process all jobs (simulating 10 workers)
    for ($w = 0; $w < NUM_WORKERS; $w++) {
        while ($jobsChan->hasData()) {
            $result = $jobsChan->receive();
            if ($result['done']) break;

            $job = $result['value'];
            $resultsChan->send([
                'number' => $job,
                'square' => $job * $job,
                'workers_used' => NUM_WORKERS,
            ]);
        }
    }
}

function handleSquare(): string {
    $start = hrtime(true);

    // Create "channels"
    $jobsChan = new Channel();
    $resultsChan = new Channel();

    // Send all jobs first
    for ($i = 0; $i <= 100; $i++) {
        $jobsChan->send($i);
    }

    // Run workers (sequentially in PHP - no real concurrency)
    while ($jobsChan->hasData()) {
        $result = $jobsChan->receive();
        if ($result['done']) break;

        $job = $result['value'];
        $resultsChan->send([
            'number' => $job,
            'square' => $job * $job,
            'workers_used' => NUM_WORKERS,
        ]);
    }

    // Collect results
    $results = [];
    while ($resultsChan->hasData()) {
        $result = $resultsChan->receive();
        if (!$result['done']) {
            $results[] = $result['value'];
        }
    }

    // Pick random result
    $randomResult = $results[array_rand($results)];
    $randomResult['time_ns'] = hrtime(true) - $start;

    return json_encode($randomResult);
}

// Simple router
$uri = $_SERVER['REQUEST_URI'] ?? '/';
$path = parse_url($uri, PHP_URL_PATH);

header('Content-Type: application/json');

switch ($path) {
    case '/square':
        echo handleSquare();
        break;
    case '/health':
        echo 'OK';
        break;
    default:
        http_response_code(404);
        echo json_encode(['error' => 'Not found']);
}
