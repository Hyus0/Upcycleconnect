<template>
  <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

  <main class="auth-shell">
    <section class="auth-layout auth-layout--register">
      <aside class="auth-panel auth-panel--brand">
        <div class="auth-panel__eyebrow">Nouveau compte</div>
        <img src="../components/logo_texte.png" alt="UpcycleConnect" class="auth-panel__logo" />
        <h1>Rejoignez la communauté qui agit.</h1>
        <p>
          Créez votre compte pour publier des annonces, suivre vos dépôts,
          rejoindre les événements et accéder à votre profil connecté.
        </p>

        <div class="auth-stats">
          <article class="auth-stat">
            <strong>8k+</strong>
            <span>objets upcyclés</span>
          </article>
          <article class="auth-stat">
            <strong>340</strong>
            <span>artisans actifs</span>
          </article>
          <article class="auth-stat">
            <strong>2.4t</strong>
            <span>CO₂ évité / mois</span>
          </article>
        </div>
      </aside>

      <section class="auth-panel auth-panel--form">
        <div class="auth-card auth-card--wide">
          <div class="auth-card__header">
            <div class="auth-card__eyebrow">Inscription</div>
            <h2>Créer un compte</h2>
            <p>
              Déjà membre ?
              <router-link to="/connexion">Se connecter</router-link>
            </p>
          </div>

          <form class="auth-form" @submit.prevent="handleSubmit">
            <div class="auth-grid">
              <div class="auth-field">
                <label for="register-firstname">Prénom</label>
                <input id="register-firstname" v-model="prenom" type="text" placeholder="Marie" />
              </div>
              <div class="auth-field">
                <label for="register-lastname">Nom</label>
                <input id="register-lastname" v-model="nom" type="text" placeholder="Lambert" />
              </div>
            </div>

            <div class="auth-field">
              <label for="register-email">Adresse e-mail</label>
              <input
                id="register-email"
                v-model="email"
                type="email"
                placeholder="marie.lambert@exemple.fr"
                autocomplete="email"
              />
            </div>

            <div class="auth-grid">
              <div class="auth-field">
                <label for="register-password">Mot de passe</label>
                <input
                  id="register-password"
                  v-model="motDePasse"
                  type="password"
                  placeholder="Minimum 8 caractères"
                  autocomplete="new-password"
                />
              </div>
              <div class="auth-field">
                <label for="register-postal">Code postal</label>
                <input
                  id="register-postal"
                  v-model="codePostal"
                  type="text"
                  placeholder="75011"
                  inputmode="numeric"
                />
              </div>
            </div>

            <label class="auth-checkbox">
              <input v-model="cguAccepte" type="checkbox" />
              <span>J'accepte les CGU et la politique de confidentialité.</span>
            </label>

            <div v-if="errorMessages.length" class="auth-error">
              <ul>
                <li v-for="(err, index) in errorMessages" :key="index">{{ err }}</li>
              </ul>
            </div>

            <button class="auth-submit" type="submit">Créer mon compte</button>
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
import { registerUser } from "../services/publicApi";

const router = useRouter();
const isLoggedIn = computed(() => Boolean(localStorage.getItem("userToken")));
const userName = computed(() => {
  const prenom = localStorage.getItem("userPrenom") || "";
  const nom = localStorage.getItem("userNom") || "";
  return prenom || nom ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const errorMessages = ref([]);
const prenom = ref("");
const nom = ref("");
const email = ref("");
const motDePasse = ref("");
const codePostal = ref("");
const cguAccepte = ref(false);

async function handleSubmit() {
  errorMessages.value = [];

  if (
    !prenom.value.trim() ||
    !nom.value.trim() ||
    !email.value.trim() ||
    !motDePasse.value.trim() ||
    !codePostal.value.trim()
  ) {
    errorMessages.value = ["Il manque des informations."];
    return;
  }
  if (!cguAccepte.value) {
    errorMessages.value = ["Veuillez accepter les CGU pour continuer."];
    return;
  }

  try {
    await registerUser({
      prenom: prenom.value,
      nom: nom.value,
      mail: email.value,
      password: motDePasse.value,
      code_postal: codePostal.value,
      role: "Particulier",
      id_langue: 1
    });
    router.push("/connexion");
  } catch (error) {
    console.error("Erreur inscription :", error);
    errorMessages.value = [error.message || "Inscription impossible."];
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
  grid-template-columns: minmax(320px, 1fr) minmax(460px, 620px);
  border-radius: 32px;
  overflow: hidden;
  background: #ffffff;
  box-shadow: 0 30px 70px rgba(15, 33, 22, 0.12);
}

.auth-panel {
  min-height: 760px;
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
  font-size: clamp(2.2rem, 4vw, 4rem);
  line-height: 1;
  color: #ffffff;
}

.auth-panel--brand p {
  max-width: 500px;
  margin: 0;
  font-size: 1.02rem;
  line-height: 1.7;
  color: rgba(244, 248, 244, 0.82);
}

.auth-stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 14px;
  margin-top: 8px;
}

.auth-stat {
  padding: 18px 16px;
  border-radius: 18px;
  border: 1px solid rgba(126, 182, 137, 0.16);
  background: rgba(255, 255, 255, 0.05);
}

.auth-stat strong,
.auth-stat span {
  display: block;
}

.auth-stat strong {
  font-size: 1.45rem;
  color: #ffffff;
}

.auth-stat span {
  margin-top: 6px;
  color: rgba(244, 248, 244, 0.72);
  line-height: 1.45;
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
  width: min(100%, 500px);
}

.auth-card--wide {
  width: min(100%, 540px);
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
  gap: 16px;
  margin-top: 32px;
}

.auth-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 14px;
}

.auth-field {
  display: grid;
  gap: 10px;
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

.auth-checkbox {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  font-weight: 500;
  color: #445248;
}

.auth-checkbox input {
  width: 18px;
  height: 18px;
  margin-top: 2px;
  padding: 0;
  border-radius: 6px;
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

@media (max-width: 1080px) {
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
}

@media (max-width: 720px) {
  .auth-shell {
    padding: 12px;
  }

  .auth-layout {
    border-radius: 24px;
  }

  .auth-grid,
  .auth-stats {
    grid-template-columns: 1fr;
  }

  .auth-panel__logo {
    width: 180px;
  }

  .auth-card__header h2 {
    font-size: 1.8rem;
  }
}
</style>
