function authHeaders() {
  const token = localStorage.getItem("userToken");
  return token ? { Authorization: token } : {};
}

const API_BASE_URL = "http://localhost:8081";

async function parseResponse(response) {
  const contentType = response.headers.get("content-type") ?? "";
  if (contentType.includes("application/json")) {
    return response.json();
  }
  return response.text();
}

async function apiRequest(path, options = {}) {
  const response = await fetch(`${API_BASE_URL}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...authHeaders(),
      ...(options.headers ?? {})
    },
    ...options
  });

  const payload = await parseResponse(response);
  if (!response.ok) {
    const message = Array.isArray(payload)
      ? payload.join(" ")
      : payload?.errors && Array.isArray(payload.errors)
        ? payload.errors.join(" ")
        : typeof payload === "string"
          ? payload
          : payload?.message ?? payload?.error ?? `Erreur HTTP ${response.status}`;
    throw new Error(message);
  }
  return payload;
}

export async function fetchPlatformOverview() {
  return apiRequest("/platform/overview");
}

export async function checkSession(userId) {
  return apiRequest(`/check-session?id=${userId}`);
}

export async function loginUser(payload) {
  return apiRequest("/login", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function registerUser(payload) {
  return apiRequest("/users", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function fetchUsers() {
  return apiRequest("/users");
}

export async function fetchUser(userId) {
  return apiRequest(`/users/${userId}`);
}

export async function fetchUserStats(userId) {
  return apiRequest(`/users/${userId}/stats`);
}

export async function fetchUserPlanning(userId) {
  return apiRequest(`/users/${userId}/planning`);
}

export async function updateUser(userId, payload) {
  return apiRequest(`/users/${userId}`, {
    method: "PUT",
    body: JSON.stringify(payload)
  });
}

export async function updatePassword(userId, payload) {
  return apiRequest(`/users/${userId}/password`, {
    method: "PUT",
    body: JSON.stringify(payload)
  });
}

export async function fetchAnnonces() {
  return apiRequest("/annonces");
}

export async function fetchAnnonce(id) {
  return apiRequest(`/annonces/${id}`);
}

export async function fetchUserAnnonces(userId) {
  return apiRequest(`/users/${userId}/annonces`);
}

export async function createAnnonce(payload) {
  return apiRequest("/annonces", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function updateAnnonce(id, payload) {
  return apiRequest(`/annonces/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload)
  });
}

export async function deleteAnnonce(id) {
  return apiRequest(`/annonces/${id}`, {
    method: "DELETE"
  });
}

export async function fetchCategories() {
  return apiRequest("/categories");
}

export async function fetchCategory(id) {
  return apiRequest(`/category/${id}`);
}

export async function fetchEvenements() {
  return apiRequest("/evenements");
}

export async function fetchEvenement(id) {
  return apiRequest(`/evenements/${id}`);
}

export async function fetchEvenementInscriptionStatus(evenementId, userId) {
  return apiRequest(`/api/evenements/${evenementId}/inscription-status?user_id=${userId}`);
}

export async function joinEvenement(evenementId, userId) {
  return apiRequest(`/api/evenements/${evenementId}/join`, {
    method: "POST",
    body: JSON.stringify({ id_utilisateur: Number(userId) })
  });
}

export async function quitEvenement(evenementId, userId) {
  return apiRequest(`/api/evenements/${evenementId}/quit`, {
    method: "POST",
    body: JSON.stringify({ id_utilisateur: Number(userId) })
  });
}

export async function fetchFormations() {
  return apiRequest("/formations");
}

export async function fetchFormation(id, userId = 0) {
  return apiRequest(`/formations/${id}?user_id=${Number(userId) || 0}`);
}

export async function joinFormation(formationId, userId) {
  return apiRequest(`/api/formations/${formationId}/join`, {
    method: "POST",
    body: JSON.stringify({ id_utilisateur: Number(userId) })
  });
}

export async function quitFormation(formationId, userId) {
  return apiRequest(`/api/formations/${formationId}/quit`, {
    method: "POST",
    body: JSON.stringify({ id_utilisateur: Number(userId) })
  });
}

export async function fetchProjets() {
  return apiRequest("/projets");
}

export async function fetchProjet(id) {
  return apiRequest(`/projet/${id}`);
}

export async function fetchProjetLikeStatus(projetId, userId) {
  return apiRequest(`/projets/${projetId}/like-status/${userId}`);
}

export async function toggleProjetLike(projetId, userId) {
  return apiRequest(`/projets/${projetId}/like/${userId}`, {
    method: "POST"
  });
}

export async function fetchSites() {
  return apiRequest("/sites");
}

export async function fetchSite(id) {
  return apiRequest(`/site/${id}`);
}

export async function reserveCasier(annonceId, siteId) {
  return apiRequest(`/annonces/${annonceId}/reserver`, {
    method: "POST",
    body: JSON.stringify({ site_id: Number(siteId) })
  });
}

export async function retirerAnnonceDuCasier(annonceId) {
  return apiRequest(`/annonces/${annonceId}/retirer`, {
    method: "POST"
  });
}
