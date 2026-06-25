// Gestion de session du back-office admin.
// Le token vient de /go/login (colonne UTILISATEUR.token cote API Go).

const TOKEN_KEY = "adminToken";
const ROLE_KEY = "adminRole";
const NAME_KEY = "adminName";

export function setSession({ token, role, name }) {
  localStorage.setItem(TOKEN_KEY, token ?? "");
  localStorage.setItem(ROLE_KEY, role ?? "");
  localStorage.setItem(NAME_KEY, name ?? "");
}

export function clearSession() {
  localStorage.removeItem(TOKEN_KEY);
  localStorage.removeItem(ROLE_KEY);
  localStorage.removeItem(NAME_KEY);
}

export function getToken() {
  return localStorage.getItem(TOKEN_KEY) || "";
}

export function getRole() {
  return localStorage.getItem(ROLE_KEY) || "";
}

export function getName() {
  return localStorage.getItem(NAME_KEY) || "";
}

export function isAuthenticated() {
  return Boolean(getToken());
}
