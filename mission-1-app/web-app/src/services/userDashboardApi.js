import {
  fetchAnnonces,
  fetchEvenements,
  fetchFormations,
  fetchUsers
} from "./publicApi";

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
  const [users, annonces, evenements, formations] = await Promise.allSettled([
    fetchUsers(),
    fetchAnnonces(),
    fetchEvenements(),
    fetchFormations()
  ]);

  const usersPayload = users.status === "fulfilled" ? users.value : null;
  const annoncesPayload = annonces.status === "fulfilled" ? annonces.value : null;
  const evenementsPayload = evenements.status === "fulfilled" ? evenements.value : null;
  const formationsPayload = formations.status === "fulfilled" ? formations.value : null;

  const planning =
    Array.isArray(evenementsPayload) && evenementsPayload.length > 0
      ? evenementsPayload.slice(0, 7).map((item) => ({
          id: item.id,
          task: item.titre ?? null,
          date: item.date_evenement ?? null
        }))
      : null;

  const advice =
    Array.isArray(formationsPayload) && formationsPayload.length > 0
      ? {
          title: formationsPayload[0].titre ?? null,
          content: formationsPayload[0].description ?? null
        }
      : null;

  return {
    user: normalizeUser(usersPayload),
    annonces: normalizeArray(annoncesPayload),
    planning,
    advice,
    notification: null
  };
}
