const STORAGE_KEY = "upcycleconnect-admin-db";

async function loadSeed() {
  const response = await fetch("/mock-db.json");
  return response.json();
}

async function ensureDb() {
  const existing = window.localStorage.getItem(STORAGE_KEY);
  if (existing) {
    return JSON.parse(existing);
  }

  const seed = await loadSeed();
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(seed));
  return seed;
}

async function saveDb(db) {
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(db));
  return db;
}

function createId(prefix) {
  return `${prefix}${Math.random().toString(36).slice(2, 8)}`;
}

export async function readCollection(name) {
  const db = await ensureDb();
  return db[name] ?? [];
}

export async function insertCollectionItem(name, item, prefix) {
  const db = await ensureDb();
  const record = { ...item, id: item.id ?? createId(prefix) };
  db[name] = [record, ...(db[name] ?? [])];
  await saveDb(db);
  return record;
}

export async function updateCollectionItem(name, id, patch) {
  const db = await ensureDb();
  db[name] = (db[name] ?? []).map((item) => (item.id === id ? { ...item, ...patch } : item));
  await saveDb(db);
  return (db[name] ?? []).find((item) => item.id === id) ?? null;
}

export async function deleteCollectionItem(name, id) {
  const db = await ensureDb();
  db[name] = (db[name] ?? []).filter((item) => item.id !== id);
  await saveDb(db);
  return true;
}

export async function getDashboardSnapshot() {
  const db = await ensureDb();
  return db;
}
