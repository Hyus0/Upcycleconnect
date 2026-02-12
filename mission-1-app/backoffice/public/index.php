<?php

declare(strict_types=1);

header('Content-Type: application/json');
header('Access-Control-Allow-Origin: *');
header('Access-Control-Allow-Methods: GET, POST, OPTIONS');
header('Access-Control-Allow-Headers: Content-Type, Authorization');

if ($_SERVER['REQUEST_METHOD'] === 'OPTIONS') {
    http_response_code(204);
    exit;
}

$storagePath = __DIR__ . '/../storage/data.json';

function loadData(string $path): array
{
    if (!file_exists($path)) {
        return [];
    }

    $raw = file_get_contents($path);
    if ($raw === false || trim($raw) === '') {
        return [];
    }

    $decoded = json_decode($raw, true);
    return is_array($decoded) ? $decoded : [];
}

function saveData(string $path, array $data): void
{
    file_put_contents($path, json_encode($data, JSON_PRETTY_PRINT | JSON_UNESCAPED_UNICODE));
}

function respond(array $payload, int $status = 200): void
{
    http_response_code($status);
    echo json_encode($payload, JSON_UNESCAPED_UNICODE);
    exit;
}

function body(): array
{
    $raw = file_get_contents('php://input');
    if ($raw === false || $raw === '') {
        return [];
    }

    $decoded = json_decode($raw, true);
    return is_array($decoded) ? $decoded : [];
}

function nextId(string $prefix): string
{
    return $prefix . '-' . bin2hex(random_bytes(4));
}

$method = $_SERVER['REQUEST_METHOD'];
$path = parse_url($_SERVER['REQUEST_URI'], PHP_URL_PATH) ?? '/';

if ($path === '/health') {
    respond(['status' => 'ok', 'service' => 'backoffice-php']);
}

if ($path === '/admin') {
    header('Content-Type: text/html; charset=utf-8');
    echo '<!doctype html><html><body><h1>Backoffice UpcycleConnect</h1><p>Utiliser les endpoints /api/*</p></body></html>';
    exit;
}

if (!str_starts_with($path, '/api/')) {
    respond(['error' => 'not_found', 'path' => $path], 404);
}

$route = substr($path, 5);
$data = loadData($storagePath);

if ($method === 'GET' && $route === 'conseils') {
    respond(['items' => $data['conseils'] ?? []]);
}

if ($method === 'GET' && $route === 'formations') {
    respond(['items' => $data['formations'] ?? []]);
}

if ($method === 'GET' && $route === 'news') {
    respond(['items' => $data['news'] ?? []]);
}

if ($route === 'annonces') {
    if ($method === 'GET') {
        $items = $data['annonces'] ?? [];
        $type = $_GET['type'] ?? null;
        if ($type !== null) {
            $items = array_values(array_filter($items, fn ($row) => ($row['type'] ?? '') === $type));
        }
        respond(['items' => $items]);
    }

    if ($method === 'POST') {
        $payload = body();
        $item = [
            'id' => nextId('a'),
            'title' => (string)($payload['title'] ?? ''),
            'description' => (string)($payload['description'] ?? ''),
            'type' => (string)($payload['type'] ?? 'don'),
            'price' => (float)($payload['price'] ?? 0)
        ];
        $data['annonces'][] = $item;
        saveData($storagePath, $data);
        respond(['created' => $item], 201);
    }
}

if ($route === 'planning') {
    if ($method === 'GET') {
        $userId = $_GET['userId'] ?? '';
        $items = $data['planning'] ?? [];
        if ($userId !== '') {
            $items = array_values(array_filter($items, fn ($row) => ($row['userId'] ?? '') === $userId));
        }
        respond(['items' => $items]);
    }

    if ($method === 'POST') {
        $payload = body();
        $item = [
            'id' => nextId('p'),
            'userId' => (string)($payload['userId'] ?? ''),
            'date' => (string)($payload['date'] ?? ''),
            'task' => (string)($payload['task'] ?? '')
        ];
        $data['planning'][] = $item;
        saveData($storagePath, $data);
        respond(['created' => $item], 201);
    }
}

if ($method === 'POST' && $route === 'container-deposits') {
    $payload = body();
    $item = [
        'id' => nextId('d'),
        'containerId' => (string)($payload['containerId'] ?? ''),
        'barcode' => (string)($payload['barcode'] ?? ''),
        'status' => (string)($payload['status'] ?? 'pending')
    ];
    $data['containerDeposits'][] = $item;
    saveData($storagePath, $data);
    respond(['created' => $item], 201);
}

if ($method === 'POST' && $route === 'pro/subscriptions') {
    $payload = body();
    $item = [
        'id' => nextId('sub'),
        'plan' => (string)($payload['plan'] ?? 'starter'),
        'billingState' => 'pending_stripe_integration'
    ];
    $data['proSubscriptions'][] = $item;
    saveData($storagePath, $data);
    respond($item, 201);
}

if ($method === 'POST' && $route === 'pro/projects') {
    $payload = body();
    $item = [
        'id' => nextId('proj'),
        'name' => (string)($payload['name'] ?? ''),
        'summary' => (string)($payload['summary'] ?? ''),
        'featured' => true
    ];
    $data['proProjects'][] = $item;
    saveData($storagePath, $data);
    respond(['created' => $item], 201);
}

if ($method === 'POST' && $route === 'employee/trainings') {
    $payload = body();
    $item = [
        'id' => nextId('et'),
        'title' => (string)($payload['title'] ?? ''),
        'date' => (string)($payload['date'] ?? '')
    ];
    $data['employeeTrainings'][] = $item;
    saveData($storagePath, $data);
    respond(['created' => $item], 201);
}

if ($method === 'POST' && $route === 'news') {
    $payload = body();
    $item = [
        'id' => nextId('n'),
        'title' => (string)($payload['title'] ?? ''),
        'content' => (string)($payload['content'] ?? '')
    ];
    $data['news'][] = $item;
    saveData($storagePath, $data);
    respond(['created' => $item], 201);
}

if ($method === 'POST' && $route === 'employee/moderation') {
    $payload = body();
    $item = [
        'id' => nextId('mod'),
        'threadId' => (string)($payload['threadId'] ?? ''),
        'action' => (string)($payload['action'] ?? 'approve')
    ];
    $data['moderationActions'][] = $item;
    saveData($storagePath, $data);
    respond(['created' => $item], 201);
}

if ($method === 'GET' && $route === 'admin/metrics') {
    $metrics = [
        'users' => count($data['users'] ?? []),
        'annonces' => count($data['annonces'] ?? []),
        'formations' => count($data['formations'] ?? []),
        'deposits' => count($data['containerDeposits'] ?? []),
        'subscriptions' => count($data['proSubscriptions'] ?? [])
    ];
    respond(['metrics' => $metrics]);
}

if ($method === 'GET' && $route === 'admin/users') {
    respond(['items' => $data['users'] ?? []]);
}

if ($method === 'POST' && $route === 'admin/notifications') {
    $payload = body();
    $item = [
        'id' => nextId('notif'),
        'title' => (string)($payload['title'] ?? ''),
        'channel' => 'onesignal_pending_integration'
    ];
    $data['notifications'][] = $item;
    saveData($storagePath, $data);
    respond(['created' => $item], 201);
}

respond(['error' => 'route_not_found', 'route' => $route], 404);
