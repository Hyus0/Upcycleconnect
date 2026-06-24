<template>
  <section class="auth-page">
    <div class="auth-card">
      <header class="auth-head">
        <span class="auth-brand">Upcycle <strong>Connect</strong></span>
        <h1>Back-office</h1>
        <p>Acces reserve aux administrateurs.</p>
      </header>

      <form class="auth-form" @submit.prevent="onSubmit">
        <label>
          <span>Adresse e-mail</span>
          <input
            v-model.trim="email"
            type="email"
            autocomplete="username"
            placeholder="admin@upcycleconnect.local"
            required
          />
        </label>

        <label>
          <span>Mot de passe</span>
          <input
            v-model="password"
            type="password"
            autocomplete="current-password"
            placeholder="********"
            required
          />
        </label>

        <p v-if="error" class="auth-error">{{ error }}</p>

        <button type="submit" class="auth-submit" :disabled="loading">
          {{ loading ? "Connexion..." : "Se connecter" }}
        </button>
      </form>
    </div>
  </section>
</template>

<script setup>
import { ref } from "vue";
import { useRouter, useRoute } from "vue-router";
import { setSession } from "../services/auth";

const router = useRouter();
const route = useRoute();

const email = ref("");
const password = ref("");
const error = ref("");
const loading = ref(false);

async function onSubmit() {
  error.value = "";
  loading.value = true;
  try {
    const res = await fetch("/go/login", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ email: email.value, password: password.value })
    });
    const data = await res.json().catch(() => ({}));

    if (!res.ok) {
      error.value = data.message || "Email ou mot de passe incorrect.";
      return;
    }

    const role = data.role || "";
    if (role !== "Admin") {
      error.value = "Acces refuse : compte non administrateur.";
      return;
    }

    setSession({
      token: data.token,
      role,
      name: `${data.prenom ?? ""} ${data.nom ?? ""}`.trim()
    });

    const redirect = typeof route.query.redirect === "string" ? route.query.redirect : null;
    if (redirect) {
      router.replace(redirect);
    } else {
      router.replace({ name: "dashboard" });
    }
  } catch (e) {
    error.value = "Le serveur est injoignable.";
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.auth-page {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
  background:
    linear-gradient(180deg, rgba(255, 255, 255, 0.94), rgba(244, 245, 241, 0.96)),
    radial-gradient(circle at top, rgba(87, 193, 125, 0.1), transparent 30%);
}

.auth-card {
  width: min(100%, 400px);
  background: #fff;
  border: 1px solid #e2e8e3;
  border-radius: 18px;
  padding: 36px 30px;
  box-shadow: 0 18px 48px rgba(23, 30, 26, 0.08);
}

.auth-head {
  text-align: center;
  margin-bottom: 26px;
}

.auth-brand {
  font-family: "Syne", sans-serif;
  font-size: 1.1rem;
  letter-spacing: -0.04em;
  color: #171e1a;
}

.auth-brand strong {
  color: #318554;
}

.auth-head h1 {
  margin: 14px 0 6px;
  font-family: "Syne", sans-serif;
  font-size: 2rem;
  letter-spacing: -0.05em;
  color: #161c18;
}

.auth-head p {
  margin: 0;
  color: #79817d;
  font-size: 0.92rem;
}

.auth-form {
  display: grid;
  gap: 16px;
}

.auth-form label {
  display: grid;
  gap: 8px;
}

.auth-form label > span {
  color: #313833;
  font-size: 0.9rem;
}

.auth-form input {
  border-radius: 12px;
  border: 1px solid #d4ded6;
  background: #fff;
  color: #16201a;
  padding: 14px 15px;
  font-size: 0.95rem;
}

.auth-form input:focus {
  border-color: rgba(49, 133, 84, 0.5);
  outline: none;
  box-shadow: 0 0 0 3px rgba(87, 193, 125, 0.12);
}

.auth-error {
  margin: 0;
  padding: 10px 12px;
  border-radius: 10px;
  background: #fcebec;
  color: #b3261e;
  font-size: 0.85rem;
}

.auth-submit {
  margin-top: 4px;
  border: 0;
  border-radius: 12px;
  min-height: 50px;
  font-weight: 700;
  font-size: 0.98rem;
  background: linear-gradient(180deg, #368858, #2f754c);
  color: #fff;
  cursor: pointer;
}

.auth-submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
