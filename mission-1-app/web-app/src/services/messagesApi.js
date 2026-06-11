const API_BASE = "/go";

function authHeaders() {
  const token = sessionStorage.getItem("userToken");
  return {
    "Content-Type": "application/json",
    ...(token ? { Authorization: token } : {})
  };
}

export function currentUserId() {
  return Number(sessionStorage.getItem("userId") || sessionStorage.getItem("id")) || 0;
}

export async function fetchSubscriptionStatus(userId = currentUserId()) {
  const response = await fetch(`${API_BASE}/users/${userId}/subscription`, {
    headers: authHeaders()
  });
  if (!response.ok) throw new Error(await response.text());
  return response.json();
}

export async function startConversation({ targetUserId, annonceId = null }) {
  const userId = currentUserId();
  const response = await fetch(`${API_BASE}/users/${userId}/messages/start`, {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify({
      target_user_id: targetUserId,
      annonce_id: annonceId
    })
  });

  const payload = await response.json().catch(() => ({}));
  if (!response.ok) {
    const error = new Error(payload.message || "Conversation impossible");
    error.status = response.status;
    error.payload = payload;
    throw error;
  }
  return payload;
}

export async function fetchConversations(userId = currentUserId()) {
  const response = await fetch(`${API_BASE}/users/${userId}/messages`, {
    headers: authHeaders()
  });
  if (!response.ok) throw new Error(await response.text());
  return response.json();
}

export async function fetchMessages(conversationId, userId = currentUserId()) {
  const response = await fetch(`${API_BASE}/users/${userId}/messages/${conversationId}`, {
    headers: authHeaders()
  });
  if (!response.ok) throw new Error(await response.text());
  return response.json();
}

export async function fetchConversationState(conversationId, userId = currentUserId()) {
  const response = await fetch(`${API_BASE}/users/${userId}/messages/${conversationId}/state`, {
    headers: authHeaders()
  });
  if (!response.ok) throw new Error(await response.text());
  return response.json();
}

export async function sendMessage(conversationId, content, userId = currentUserId()) {
  const response = await fetch(`${API_BASE}/users/${userId}/messages/${conversationId}`, {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify({ content })
  });
  if (!response.ok) throw new Error(await response.text());
  return response.json();
}

export async function createOffer(conversationId, amount, userId = currentUserId()) {
  const response = await fetch(`${API_BASE}/users/${userId}/messages/${conversationId}/offers`, {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify({ amount: Number(amount) })
  });
  if (!response.ok) throw new Error(await response.text());
  return response.json();
}

export async function respondOffer(offerId, action, userId = currentUserId()) {
  const response = await fetch(`${API_BASE}/users/${userId}/messages/offers/${offerId}`, {
    method: "PATCH",
    headers: authHeaders(),
    body: JSON.stringify({ action })
  });
  if (!response.ok) throw new Error(await response.text());
  return response.json();
}

export async function confirmSaleReception(saleId, userId = currentUserId()) {
  const response = await fetch(`${API_BASE}/users/${userId}/messages/sales/${saleId}/reception`, {
    method: "POST",
    headers: authHeaders()
  });
  if (!response.ok) throw new Error(await response.text());
  return response.json();
}

export async function reviewSale(saleId, { note, commentaire }, userId = currentUserId()) {
  const response = await fetch(`${API_BASE}/users/${userId}/messages/sales/${saleId}/review`, {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify({ note: Number(note), commentaire })
  });
  if (!response.ok) throw new Error(await response.text());
  return response.json();
}

export async function addSubscriptionToCart(userId = currentUserId()) {
  const response = await fetch(`${API_BASE}/users/${userId}/panier`, {
    method: "POST",
    headers: authHeaders(),
    body: JSON.stringify({
      type_item: "Abonnement",
      reference_id: 0,
      prix_unitaire: 2.99
    })
  });
  if (!response.ok) throw new Error(await response.text());
  return response.json();
}
