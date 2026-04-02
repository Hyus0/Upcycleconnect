<template>
  <section class="auth-page">
    <aside class="auth-showcase">
      <div class="auth-showcase__glow"></div>

      <header class="auth-brand">
        <span class="auth-brand__text">Upcycle <strong>Connect</strong></span>
      </header>

      <div class="auth-copy">
        <h1>
          Rejoignez la
          <span>communaute</span>
          qui agit.
        </h1>
        <p>
          Les informations de cette zone seront recuperees depuis la base de donnees.
          Tant que la BDD n'est pas initialisee, les valeurs restent a NULL.
        </p>
      </div>

      <div class="auth-stats">
        <article v-for="stat in stats" :key="stat.label">
          <strong>{{ displayValue(stat.value) }}</strong>
          <span>{{ stat.label }}</span>
        </article>
      </div>

      <article class="auth-testimonial">
        <p>"{{ displayValue(showcase.quote) }}"</p>
        <div class="auth-testimonial__author">
          <div class="auth-avatar">{{ displayInitials(showcase.author) }}</div>
          <div>
            <strong>{{ displayValue(showcase.author) }}</strong>
            <span>{{ displayValue(showcase.meta) }}</span>
          </div>
        </div>
      </article>
    </aside>

    <main class="auth-panel">
      <div class="auth-panel__inner">
        <header class="auth-panel__header">
          <h2>Creer un compte</h2>
          <p>
            Deja membre ?
            <a href="#">Se connecter</a>
          </p>
        </header>

        <div class="auth-role">
          <span>Je suis...</span>
          <div class="auth-role__grid">
            <button
              v-for="option in roleOptions"
              :key="option.value"
              type="button"
              class="role-card"
              :class="{ active: form.role === option.value }"
              @click="form.role = option.value"
            >
              <span class="role-card__icon">{{ option.icon }}</span>
              <strong>{{ option.label }}</strong>
            </button>
          </div>
        </div>

        <form class="auth-form" @submit.prevent>
          <div class="auth-form__row">
            <label>
              <span>Prenom</span>
              <input v-model="form.firstName" type="text" placeholder="Marie" />
            </label>
            <label>
              <span>Nom</span>
              <input v-model="form.lastName" type="text" placeholder="Lambert" />
            </label>
          </div>

          <label>
            <span>Adresse e-mail</span>
            <input v-model="form.email" type="email" placeholder="marie.lambert@exemple.fr" />
          </label>

          <label>
            <span>Mot de passe</span>
            <input v-model="form.password" type="password" placeholder="********" />
          </label>

          <label>
            <span>Code postal</span>
            <input v-model="form.postalCode" type="text" placeholder="75011" />
          </label>

          <label class="auth-checkbox">
            <input v-model="form.accepted" type="checkbox" />
            <span>J'accepte les <a href="#">CGU</a> et la <a href="#">politique de confidentialite</a>.</span>
          </label>

          <button type="submit" class="auth-submit">Creer mon compte gratuitement -></button>

          <div class="auth-divider">
            <span></span>
            <small>ou</small>
            <span></span>
          </div>

          <button type="button" class="auth-google">Continuer avec Google</button>
        </form>
      </div>
    </main>
  </section>
</template>

<script setup>
import { reactive } from "vue";

const stats = [
  { label: "CO2 evite / mois", value: null },
  { label: "Objets upcycles", value: null },
  { label: "Artisans actifs", value: null }
];

const showcase = {
  quote: null,
  author: null,
  meta: null
};

const roleOptions = [
  { value: "particulier", label: "Particulier", icon: "P" },
  { value: "artisan", label: "Pro / Artisan", icon: "A" }
];

const form = reactive({
  role: "particulier",
  firstName: "",
  lastName: "",
  email: "",
  password: "",
  postalCode: "",
  accepted: false
});

function displayValue(value) {
  return value === null || value === undefined || value === "" ? "NULL" : value;
}

function displayInitials(value) {
  if (!value) return "NU";

  const initials = value
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? "")
    .join("");

  return initials || "NU";
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 1.08fr 1fr;
  background: #f6f4ef;
}

.auth-showcase {
  position: relative;
  overflow: hidden;
  padding: 42px 24px 30px;
  background: #171e1a;
  color: rgba(255, 255, 255, 0.94);
  display: grid;
  grid-template-rows: auto auto auto 1fr;
  gap: 42px;
}

.auth-showcase__glow {
  position: absolute;
  inset: auto auto -18% -12%;
  width: 72%;
  aspect-ratio: 1;
  background: radial-gradient(circle, rgba(50, 143, 84, 0.32), transparent 68%);
  pointer-events: none;
}

.auth-brand,
.auth-copy,
.auth-stats,
.auth-testimonial {
  position: relative;
  z-index: 1;
}

.auth-brand__text {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  color: #f1f6f2;
  font-family: "Syne", sans-serif;
  font-size: 1.15rem;
  letter-spacing: -0.04em;
}

.auth-brand__text strong {
  color: #57c17d;
  font-weight: 800;
}

.auth-copy {
  max-width: 520px;
}

.auth-copy h1 {
  margin: 0;
  font-family: "Syne", sans-serif;
  font-size: clamp(3rem, 7vw, 5rem);
  line-height: 0.94;
  letter-spacing: -0.07em;
}

.auth-copy h1 span {
  display: block;
  color: #57c17d;
}

.auth-copy p {
  margin: 22px 0 0;
  max-width: 36ch;
  color: rgba(255, 255, 255, 0.5);
  font-size: 1.05rem;
  line-height: 1.65;
}

.auth-stats {
  display: flex;
  gap: 18px;
  flex-wrap: wrap;
}

.auth-stats article {
  min-width: 92px;
  padding-right: 14px;
  border-right: 1px solid rgba(255, 255, 255, 0.08);
}

.auth-stats article:last-child {
  border-right: 0;
}

.auth-stats strong {
  display: block;
  font-family: "Syne", sans-serif;
  font-size: 2.2rem;
  letter-spacing: -0.06em;
}

.auth-stats span {
  display: block;
  margin-top: 6px;
  color: rgba(255, 255, 255, 0.38);
  font-size: 0.82rem;
}

.auth-testimonial {
  align-self: end;
  max-width: 610px;
  padding: 20px 18px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(10px);
}

.auth-testimonial p {
  margin: 0;
  color: rgba(255, 255, 255, 0.64);
  font-style: italic;
  line-height: 1.55;
}

.auth-testimonial__author {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-top: 18px;
}

.auth-avatar {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  display: grid;
  place-items: center;
  background: rgba(87, 193, 125, 0.22);
  color: #8be4ab;
  font-weight: 700;
}

.auth-testimonial__author span {
  display: block;
  margin-top: 2px;
  color: rgba(255, 255, 255, 0.36);
  font-size: 0.82rem;
}

.auth-panel {
  display: grid;
  place-items: center;
  padding: 40px 24px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.92), rgba(244, 245, 241, 0.94)),
    radial-gradient(circle at top, rgba(87, 193, 125, 0.08), transparent 28%);
}

.auth-panel__inner {
  width: min(100%, 430px);
}

.auth-panel__header {
  text-align: center;
}

.auth-panel__header h2 {
  margin: 0;
  color: #161c18;
  font-family: "Syne", sans-serif;
  font-size: clamp(2.2rem, 5vw, 3.4rem);
  letter-spacing: -0.06em;
}

.auth-panel__header p {
  margin: 10px 0 0;
  color: #79817d;
  font-size: 0.96rem;
}

.auth-panel__header a,
.auth-checkbox a {
  color: #318554;
}

.auth-role {
  margin-top: 28px;
}

.auth-role > span {
  display: block;
  margin-bottom: 12px;
  color: #47504b;
  font-size: 0.92rem;
}

.auth-role__grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.role-card {
  border: 1px solid rgba(33, 79, 54, 0.12);
  border-radius: 12px;
  padding: 16px 12px;
  background: rgba(255, 255, 255, 0.72);
  color: #202723;
  display: grid;
  gap: 8px;
  justify-items: center;
}

.role-card.active {
  background: rgba(87, 193, 125, 0.1);
  border-color: rgba(49, 133, 84, 0.56);
}

.role-card__icon {
  width: 28px;
  height: 28px;
  display: grid;
  place-items: center;
  border-radius: 999px;
  background: rgba(87, 193, 125, 0.12);
  color: #2f754c;
  font-size: 0.82rem;
  font-weight: 700;
}

.auth-form {
  display: grid;
  gap: 14px;
  margin-top: 16px;
}

.auth-form label {
  display: grid;
  gap: 8px;
}

.auth-form label > span {
  color: #313833;
  font-size: 0.9rem;
}

.auth-form__row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.auth-form input {
  border-radius: 12px;
  border: 1px solid #d4ded6;
  background: rgba(255, 255, 255, 0.92);
  color: #16201a;
  padding: 14px 15px;
}

.auth-form input:focus {
  border-color: rgba(49, 133, 84, 0.5);
  outline: none;
  box-shadow: 0 0 0 3px rgba(87, 193, 125, 0.1);
}

.auth-checkbox {
  grid-template-columns: auto 1fr;
  align-items: start;
  gap: 10px;
  margin-top: 2px;
}

.auth-checkbox input {
  width: 14px;
  height: 14px;
  margin-top: 4px;
  accent-color: #318554;
}

.auth-checkbox span {
  color: #76807b;
  font-size: 0.82rem;
  line-height: 1.5;
}

.auth-submit,
.auth-google {
  border: 0;
  border-radius: 12px;
  min-height: 52px;
  font-weight: 700;
  font-size: 0.98rem;
}

.auth-submit {
  margin-top: 2px;
  background: linear-gradient(180deg, #368858, #2f754c);
  color: white;
}

.auth-google {
  background: rgba(255, 255, 255, 0.84);
  color: #1b221e;
  border: 1px solid #d7dfd8;
}

.auth-divider {
  display: grid;
  grid-template-columns: 1fr auto 1fr;
  gap: 12px;
  align-items: center;
  color: #9ca59f;
}

.auth-divider span {
  display: block;
  height: 1px;
  background: #d9dfda;
}

@media (max-width: 1100px) {
  .auth-page {
    grid-template-columns: 1fr;
  }

  .auth-showcase {
    min-height: auto;
  }
}

@media (max-width: 640px) {
  .auth-showcase,
  .auth-panel {
    padding: 24px 18px;
  }

  .auth-form__row,
  .auth-role__grid {
    grid-template-columns: 1fr;
  }

  .auth-stats {
    display: grid;
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }

  .auth-stats article {
    border-right: 0;
    padding-right: 0;
  }
}
</style>
