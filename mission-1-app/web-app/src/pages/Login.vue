<template>
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

    <div class="auth-page">
        <section class="auth-side auth-side--dark">
            <div class="auth-side__content">
                <div class="auth-brand" @click="$router.push('/')">
                    <img src="../components/logo.png" class="navbar__logo-img" alt="Upcycle Connect" />
                    <span>Upcycle <strong>Connect</strong></span>
                </div>
                <h1>Heureux de vous revoir.</h1>
                <p>Connectez-vous pour retrouver vos annonces, vos inscriptions et votre espace personnel pilote par le back.</p>
            </div>
        </section>

        <section class="auth-side auth-side--light">
            <div class="auth-card">
                <h2>Connexion</h2>
                <p>Pas encore membre ? <router-link to="/inscription">Creer un compte</router-link></p>

                <form @submit.prevent="handleLogin" class="auth-form">
                    <label>Adresse e-mail</label>
                    <input v-model="email" type="email" placeholder="marie.lambert@exemple.fr" />

                    <label>Mot de passe</label>
                    <input v-model="motDePasse" type="password" placeholder="Votre mot de passe" autocomplete="current-password" />

                    <div v-if="errorMessages.length" class="auth-error">
                        <ul><li v-for="(error, index) in errorMessages" :key="index">{{ error }}</li></ul>
                    </div>

                    <button type="submit" class="login-submit">Se connecter</button>
                </form>
            </div>
        </section>
    </div>
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

        if (data.token) localStorage.setItem("userToken", data.token);
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
.auth-page { display: grid; grid-template-columns: 1fr 1fr; min-height: calc(100vh - 86px); margin-top: -108px; }
.auth-side { display: flex; align-items: center; justify-content: center; padding: 48px; }
.auth-side--dark { background: #17201b; color: white; }
.auth-side--light { background: #ffffff; }
.auth-brand { display: flex; align-items: center; gap: 12px; font-size: 1.15rem; font-weight: 700; cursor: pointer; margin-bottom: 2rem; }
.auth-brand strong { color: #338454; }
.auth-card { width: min(420px, 100%); }
.auth-form { display: grid; gap: 0.85rem; margin-top: 1.5rem; }
.auth-form input { padding: 0.95rem 1rem; border: 1px solid #d9dfdb; border-radius: 12px; }
.auth-error { color: #b91c1c; background: #fee2e2; border: 1px solid #ef4444; padding: 10px; border-radius: 10px; }
.login-submit { margin-top: 0.5rem; background: #338454; color: white; border: none; border-radius: 12px; padding: 0.95rem 1rem; font-weight: 700; }
</style>
