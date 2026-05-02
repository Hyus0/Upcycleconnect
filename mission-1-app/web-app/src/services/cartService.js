const CART_EVENT = "upcycle-cart-change";

function resolveCartKey() {
  const userId = localStorage.getItem("userId") || sessionStorage.getItem("userId");
  return userId ? `upcycle-cart:${userId}` : "upcycle-cart:guest";
}

function readCart() {
  try {
    const raw = localStorage.getItem(resolveCartKey());
    const parsed = raw ? JSON.parse(raw) : [];
    return Array.isArray(parsed) ? parsed : [];
  } catch {
    return [];
  }
}

function writeCart(items) {
  localStorage.setItem(resolveCartKey(), JSON.stringify(items));
  window.dispatchEvent(new Event(CART_EVENT));
}

function itemKey(item) {
  return `${item.itemType || "annonce"}:${item.id}`;
}

function normalizeAnnonce(annonce) {
  return {
    itemType: "annonce",
    id: annonce.id,
    titre: annonce.titre || "NULL",
    description: annonce.description || "NULL",
    prix: annonce.prix ?? 0,
    type: annonce.type || "don",
    statut: annonce.statut || "en ligne",
    ville: annonce.ville || "NULL",
    code_postal: annonce.code_postal || "NULL",
    date_creation: annonce.date_creation || null
  };
}

function normalizeFormation(formation) {
  return {
    itemType: "formation",
    id: formation.id,
    titre: formation.titre || "Formation",
    description: formation.description || "NULL",
    prix: formation.prix_unitaire ?? formation.prix ?? 0,
    type: "formation",
    statut: formation.statut || "Ouvert",
    ville: formation.ville || "NULL",
    code_postal: formation.code_postal || formation.CodePostal || "NULL",
    date_creation: formation.date_debut || null,
    meta: {
      capacite_max: formation.capacite_max || 0,
      nb_inscrit: formation.nb_inscrit || 0,
      adresse: formation.adresse || "NULL",
      formateur: [formation.prenom_formateur, formation.nom_formateur].filter(Boolean).join(" ") || "Organisateur"
    }
  };
}

function normalizeItem(item, itemType = "annonce") {
  return itemType === "formation" ? normalizeFormation(item) : normalizeAnnonce(item);
}

export function getCartItems() {
  return readCart();
}

export function getCartCount() {
  return readCart().length;
}

export function isInCart(annonceId) {
  return isItemInCart("annonce", annonceId);
}

export function isItemInCart(itemType, itemId) {
  const key = `${itemType}:${itemId}`;
  return readCart().some((item) => itemKey(item) === key);
}

export function addToCart(annonce) {
  return addItemToCart(annonce, "annonce");
}

export function addItemToCart(item, itemType = "annonce") {
  const items = readCart();
  const normalized = normalizeItem(item, itemType);
  if (items.some((cartItem) => itemKey(cartItem) === itemKey(normalized))) {
    return items;
  }

  const nextItems = [...items, normalized];
  writeCart(nextItems);
  return nextItems;
}

export function removeFromCart(annonceId) {
  return removeItemFromCart("annonce", annonceId);
}

export function removeItemFromCart(itemType, itemId) {
  const key = `${itemType}:${itemId}`;
  const nextItems = readCart().filter((item) => itemKey(item) !== key);
  writeCart(nextItems);
  return nextItems;
}

export function clearCart() {
  writeCart([]);
}

export function savePurchase(items) {
  const userId = localStorage.getItem("userId") || "guest";
  const key = `upcycle-purchases:${userId}`;
  const existing = JSON.parse(localStorage.getItem(key) || "[]");
  const purchase = {
    id: `CMD-${Date.now()}`,
    created_at: new Date().toISOString(),
    status: "A confirmer",
    items,
    total: items.reduce((sum, item) => sum + (String(item.type).toLowerCase() === "don" ? 0 : Number(item.prix) || 0), 0)
  };
  localStorage.setItem(key, JSON.stringify([purchase, ...existing]));
  return purchase;
}

export function getPurchases() {
  const userId = localStorage.getItem("userId") || "guest";
  return JSON.parse(localStorage.getItem(`upcycle-purchases:${userId}`) || "[]");
}

export function onCartChange(handler) {
  window.addEventListener(CART_EVENT, handler);
  window.addEventListener("storage", handler);
  return () => {
    window.removeEventListener(CART_EVENT, handler);
    window.removeEventListener("storage", handler);
  };
}
