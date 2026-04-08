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
  const items = Array.isArray(payload?.items) ? payload.items : [];
  const first =
    items.find((item) => (item.role ?? "").toLowerCase() === "particulier") ??
    items.find((item) => (item.fullName ?? "").toLowerCase().includes("marie")) ??
    items[0];
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
  const [users, annonces, evenements, formations] = await Promise.all([
    getJson("/api-go/admin/users"),
    getJson("/annonces"),
    getJson("/evenements"),
    getJson("/formations")
  ]);

  const planning =
    Array.isArray(evenements) && evenements.length > 0
      ? evenements.slice(0, 7).map((item) => ({
          id: item.id,
          task: item.titre ?? null,
          date: item.date_evenement ?? null
        }))
      : null;

  const advice =
    Array.isArray(formations) && formations.length > 0
      ? {
          title: formations[0].titre ?? null,
          content: formations[0].description ?? null
        }
      : null;

  return {
    user: normalizeUser(users),
    annonces: normalizeArray(annonces),
    planning,
    advice,
    notification: null
  };
}
