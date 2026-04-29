function authHeaders() {
  const token = localStorage.getItem("userToken");
  return token ? { Authorization: token } : {};
}

async function parseResponse(response) {
  const contentType = response.headers.get("content-type") ?? "";
  if (contentType.includes("application/json")) {
    return response.json();
  }
  return response.text();
}

async function apiRequest(path, options = {}) {
  const response = await fetch(path, {
    headers: {
      "Content-Type": "application/json",
      ...authHeaders(),
      ...(options.headers ?? {})
    },
    ...options
  });

  const payload = await parseResponse(response);
  if (!response.ok) {
    const message =
      typeof payload === "string"
        ? payload
        : payload?.message ?? payload?.error ?? `Erreur HTTP ${response.status}`;
    throw new Error(message);
  }
  return payload;
}

export async function fetchPlatformOverview() {
  return apiRequest("/go/platform/overview");
}

export async function checkSession(userId) {
  return apiRequest(`/go/check-session?id=${userId}`);
}

export async function loginUser(payload) {
  return apiRequest("/go/login", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function registerUser(payload) {
  return apiRequest("/go/users", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function fetchUsers() {
  return apiRequest("/go/users");
}

export async function fetchUser(userId) {
  return apiRequest(`/go/users/${userId}`);
}

export async function fetchUserStats(userId) {
  return apiRequest(`/go/users/${userId}/stats`);
}

export async function fetchUserPlanning(userId) {
  return apiRequest(`/go/users/${userId}/planning`);
}

export async function updateUser(userId, payload) {
  return apiRequest(`/go/users/${userId}`, {
    method: "PUT",
    body: JSON.stringify(payload)
  });
}

export async function updatePassword(userId, payload) {
  return apiRequest(`/go/users/${userId}/password`, {
    method: "PUT",
    body: JSON.stringify(payload)
  });
}

export async function fetchAnnonces() {
  return apiRequest("/go/annonces");
}

export async function fetchAnnonce(id) {
  return apiRequest(`/go/annonces/${id}`);
}

export async function fetchUserAnnonces(userId) {
  return apiRequest(`/go/users/${userId}/annonces`);
}

export async function createAnnonce(payload) {
  return apiRequest("/go/annonces", {
    method: "POST",
    body: JSON.stringify(payload)
  });
}

export async function updateAnnonce(id, payload) {
  return apiRequest(`/go/annonces/${id}`, {
    method: "PUT",
    body: JSON.stringify(payload)
  });
}

export async function deleteAnnonce(id) {
  return apiRequest(`/go/annonces/${id}`, {
    method: "DELETE"
  });
}

export async function fetchCategories() {
  return apiRequest("/go/categories");
}

export async function fetchCategory(id) {
  return apiRequest(`/go/category/${id}`);
}

export async function fetchEvenements() {
  return apiRequest("/go/evenements");
}

export async function fetchEvenementInscriptionStatus(evenementId, userId) {
  return apiRequest(`/go/api/evenements/${evenementId}/inscription-status?user_id=${userId}`);
}

export async function joinEvenement(evenementId, userId) {
  return apiRequest(`/go/api/evenements/${evenementId}/join`, {
    method: "POST",
    body: JSON.stringify({ id_utilisateur: Number(userId) })
  });
}

export async function quitEvenement(evenementId, userId) {
  return apiRequest(`/go/api/evenements/${evenementId}/quit`, {
    method: "POST",
    body: JSON.stringify({ id_utilisateur: Number(userId) })
  });
}

export async function fetchFormations() {
  return apiRequest("/go/formations");
}

export async function fetchFormation(id, userId = 0) {
  return apiRequest(`/go/formations/${id}?user_id=${Number(userId) || 0}`);
}

export async function joinFormation(formationId, userId) {
  return apiRequest(`/go/api/formations/${formationId}/join`, {
    method: "POST",
    body: JSON.stringify({ id_utilisateur: Number(userId) })
  });
}

export async function quitFormation(formationId, userId) {
  return apiRequest(`/go/api/formations/${formationId}/quit`, {
    method: "POST",
    body: JSON.stringify({ id_utilisateur: Number(userId) })
  });
}

export async function fetchProjets() {
  return apiRequest("/go/projets");
}

export async function fetchProjet(id) {
  return apiRequest(`/go/projet/${id}`);
}

export async function fetchProjetLikeStatus(projetId, userId) {
  return apiRequest(`/go/projets/${projetId}/like-status/${userId}`);
}

export async function toggleProjetLike(projetId, userId) {
  return apiRequest(`/go/projets/${projetId}/like/${userId}`, {
    method: "POST"
  });
}

export async function fetchSites() {
  return apiRequest("/go/sites");
}

export async function fetchSite(id) {
  return apiRequest(`/go/site/${id}`);
}

export async function reserveCasier(annonceId, siteId) {
  return apiRequest(`/go/annonces/${annonceId}/reserver`, {
    method: "POST",
    body: JSON.stringify({ site_id: Number(siteId) })
  });
}

export async function retirerAnnonceDuCasier(annonceId) {
  return apiRequest(`/go/annonces/${annonceId}/retirer`, {
    method: "POST"
  });
}
