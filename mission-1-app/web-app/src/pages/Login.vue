<template>
  <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

  <main class="auth-shell">
    <section class="auth-layout">
      <aside class="auth-panel auth-panel--brand">
        <div class="auth-panel__eyebrow">Espace membre</div>
        <img src="../components/logo_texte.png" alt="UpcycleConnect" class="auth-panel__logo" />
        <h1>Heureux de vous revoir.</h1>
        <p>
          Connectez-vous pour retrouver vos annonces, vos inscriptions, votre planning
          et votre espace personnel connecté au back Go.
        </p>

        <div class="auth-highlights">
          <article class="auth-highlight">
            <strong>Suivi centralisé</strong>
            <span>Annonces, dépôts, événements et formations dans un seul espace.</span>
          </article>
          <article class="auth-highlight">
            <strong>Compte persistant</strong>
            <span>Votre session pilote le front et les données stockées côté API.</span>
          </article>
        </div>
      </aside>

      <section class="auth-panel auth-panel--form">
        <div class="auth-card">
          <div class="auth-card__header">
            <div class="auth-card__eyebrow">Connexion</div>
            <h2>Accéder à mon compte</h2>
            <p>
              Pas encore membre ?
              <router-link to="/inscription">Créer un compte</router-link>
            </p>
          </div>

          <form @submit.prevent="handleLogin" class="auth-form">
            <label for="login-email">Adresse e-mail</label>
            <input
              id="login-email"
              v-model="email"
              type="email"
              placeholder="marie.lambert@exemple.fr"
              autocomplete="email"
            />

            <label for="login-password">Mot de passe</label>
            <input
              id="login-password"
              v-model="motDePasse"
              type="password"
              placeholder="Votre mot de passe"
              autocomplete="current-password"
            />

            <div v-if="errorMessages.length" class="auth-error">
              <ul>
                <li v-for="(error, index) in errorMessages" :key="index">{{ error }}</li>
              </ul>
            </div>

            <button type="submit" class="auth-submit">Se connecter</button>
          </form>
        </div>
      </section>
    </section>
  </main>
</template>

<script setup>
import { computed, ref } from "vue";
import { useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import { loginUser } from "../services/publicApi";

const router = useRouter();
const email = ref("");
const motDePasse = ref("");
const errorMessages = ref([]);

const isLoggedIn = computed(() => Boolean(localStorage.getItem("userToken")));
const userName = computed(() => {
  const prenom = localStorage.getItem("userPrenom") || "";
  const nom = localStorage.getItem("userNom") || "";
  return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

async function handleLogin() {
  errorMessages.value = [];
  if (!email.value.trim() || !motDePasse.value.trim()) {
    errorMessages.value = ["Veuillez remplir tous les champs."];
    return;
  }

  try {
    const data = await loginUser({
      email: email.value,
      password: motDePasse.value
    });

    if (data.token) {
      localStorage.setItem("userToken", data.token);
    }
    localStorage.setItem("userId", data.userId);
    localStorage.setItem("userPrenom", data.prenom || "");
    localStorage.setItem("userNom", data.nom || "");
    window.dispatchEvent(new Event("auth-change"));
    router.push("/profil");
  } catch (error) {
    console.error("Erreur connexion :", error);
    errorMessages.value = [error.message || "Connexion impossible."];
  }
}
</script>

<style scoped>
.auth-shell {
  min-height: calc(100vh - 92px);
  padding: 24px;
  background:
    radial-gradient(circle at top left, rgba(61, 145, 90, 0.18), transparent 38%),
    linear-gradient(180deg, #eef3ef 0%, #f8fbf8 100%);
}

.auth-layout {
  max-width: 1380px;
  margin: 0 auto;
  display: grid;
  grid-template-columns: minmax(320px, 1fr) minmax(420px, 560px);
  border-radius: 32px;
  overflow: hidden;
  background: #ffffff;
  box-shadow: 0 30px 70px rgba(15, 33, 22, 0.12);
}

.auth-panel {
  min-height: 720px;
}

.auth-panel--brand {
  padding: 56px;
  color: #f4f8f4;
  background:
    radial-gradient(circle at 20% 20%, rgba(93, 205, 108, 0.22), transparent 32%),
    radial-gradient(circle at 80% 75%, rgba(51, 132, 84, 0.18), transparent 28%),
    linear-gradient(160deg, #152019 0%, #101713 100%);
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 24px;
}

.auth-panel__eyebrow,
.auth-card__eyebrow {
  font-size: 0.78rem;
  letter-spacing: 0.24em;
  text-transform: uppercase;
  color: #76cb88;
}

.auth-panel__logo {
  width: 220px;
  max-width: 100%;
}

.auth-panel--brand h1 {
  margin: 0;
  font-size: clamp(2.4rem, 4vw, 4.3rem);
  line-height: 0.98;
  color: #ffffff;
}

.auth-panel--brand p {
  max-width: 480px;
  margin: 0;
  font-size: 1.02rem;
  line-height: 1.7;
  color: rgba(244, 248, 244, 0.82);
}

.auth-highlights {
  display: grid;
  gap: 14px;
  margin-top: 8px;
}

.auth-highlight {
  padding: 18px 20px;
  border: 1px solid rgba(126, 182, 137, 0.16);
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(8px);
}

.auth-highlight strong,
.auth-highlight span {
  display: block;
}

.auth-highlight strong {
  margin-bottom: 8px;
  font-size: 0.98rem;
}

.auth-highlight span {
  color: rgba(244, 248, 244, 0.78);
  line-height: 1.55;
}

.auth-panel--form {
  padding: 56px 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  background:
    linear-gradient(180deg, rgba(243, 247, 243, 0.96), rgba(255, 255, 255, 1));
}

.auth-card {
  width: min(100%, 430px);
}

.auth-card__header h2 {
  margin: 10px 0 12px;
  font-size: 2.2rem;
  color: #17201b;
}

.auth-card__header p {
  margin: 0;
  color: #5d6c62;
}

.auth-card__header a {
  color: #338454;
  font-weight: 700;
  text-decoration: none;
}

.auth-form {
  display: grid;
  gap: 14px;
  margin-top: 32px;
}

.auth-form label {
  font-size: 0.98rem;
  font-weight: 700;
  color: #17201b;
}

.auth-form input {
  width: 100%;
  padding: 16px 18px;
  border: 1px solid #d9e2db;
  border-radius: 16px;
  background: #f5f8f5;
  color: #17201b;
  font-size: 1rem;
  transition: border-color 0.2s ease, box-shadow 0.2s ease, background 0.2s ease;
}

.auth-form input:focus {
  outline: none;
  border-color: #4ca466;
  box-shadow: 0 0 0 4px rgba(76, 164, 102, 0.12);
  background: #ffffff;
}

.auth-error {
  padding: 14px 16px;
  border-radius: 16px;
  border: 1px solid #ef4444;
  background: #fef0f0;
  color: #b91c1c;
}

.auth-error ul {
  margin: 0;
  padding-left: 18px;
}

.auth-submit {
  margin-top: 8px;
  border: none;
  border-radius: 16px;
  padding: 16px 20px;
  font-size: 1rem;
  font-weight: 800;
  color: #ffffff;
  background: linear-gradient(135deg, #338454, #429f63);
  box-shadow: 0 14px 24px rgba(51, 132, 84, 0.18);
  cursor: pointer;
}

.auth-submit:hover {
  transform: translateY(-1px);
}

@media (max-width: 980px) {
  .auth-layout {
    grid-template-columns: 1fr;
  }

  .auth-panel {
    min-height: auto;
  }

  .auth-panel--brand,
  .auth-panel--form {
    padding: 36px 24px;
  }

  .auth-panel__logo {
    width: 180px;
  }
}

@media (max-width: 640px) {
  .auth-shell {
    padding: 12px;
  }

  .auth-layout {
    border-radius: 24px;
  }

  .auth-card__header h2 {
    font-size: 1.8rem;
  }
}
</style>
