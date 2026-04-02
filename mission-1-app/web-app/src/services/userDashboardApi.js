async function getJson(url) {
  try {
    const response = await fetch(url);
    if (!response.ok) {
      return null;
    }
    return await response.json();
  } catch {
    return null;
  }
}

function normalizeUser(payload) {
  const first = payload?.items?.[0];
  if (!first) return null;

  const firstName = first.prenom ?? first.firstName ?? null;
  const lastName = first.nom ?? first.lastName ?? null;

  return {
    firstName,
    fullName: [firstName, lastName].filter(Boolean).join(" ") || null,
    role: first.role ?? null
  };
}

function normalizeArray(payload) {
  return Array.isArray(payload?.items) ? payload.items : Array.isArray(payload) ? payload : null;
}

function normalizeSingle(payload) {
  return Array.isArray(payload?.items) ? payload.items[0] ?? null : null;
}

export async function fetchUserDashboard() {
  const [users, annonces, planning, conseils, news] = await Promise.all([
    getJson("/api/admin/users"),
    getJson("/annonces"),
    getJson("/api/planning?userId=u1"),
    getJson("/api/conseils"),
    getJson("/api/news")
  ]);

  return {
    user: normalizeUser(users),
    annonces: normalizeArray(annonces),
    planning: normalizeArray(planning),
    advice: normalizeSingle(conseils),
    notification: normalizeSingle(news)
  };
}
