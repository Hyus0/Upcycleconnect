<template>
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

    <div class="auth-page">
        <section class="auth-side auth-side--dark">
            <div class="auth-side__content">
                <div class="auth-brand" @click="$router.push('/')">
                    <img src="../components/logo.png" class="navbar__logo-img" alt="Upcycle Connect" />
                    <span>Upcycle <strong>Connect</strong></span>
                </div>
                <h1>Rejoignez la communaute qui agit.</h1>
                <p>Le compte cree ici sera disponible partout sur le front et administrable depuis le back.</p>
            </div>
        </section>

        <section class="auth-side auth-side--light">
            <div class="auth-card">
                <h2>Creer un compte</h2>
                <p>Deja membre ? <router-link to="/connexion">Se connecter</router-link></p>

                <form class="auth-form" @submit.prevent="handleSubmit">
                    <div class="auth-grid">
                        <div><label>Prenom</label><input v-model="prenom" type="text" placeholder="Marie" /></div>
                        <div><label>Nom</label><input v-model="nom" type="text" placeholder="Lambert" /></div>
                    </div>
                    <label>Adresse e-mail</label>
                    <input v-model="email" type="email" placeholder="marie.lambert@exemple.fr" />
                    <label>Mot de passe</label>
                    <input v-model="motDePasse" type="password" placeholder="Minimum 8 caracteres" />
                    <label>Code postal</label>
                    <input v-model="codePostal" type="text" placeholder="75011" />

                    <label class="cgu-line"><input v-model="cguAccepte" type="checkbox" /> J'accepte les CGU.</label>

                    <div v-if="errorMessages.length" class="auth-error">
                        <ul><li v-for="(err, index) in errorMessages" :key="index">{{ err }}</li></ul>
                    </div>

                    <button class="login-submit" type="submit">Creer mon compte</button>
                </form>
            </div>
        </section>
    </div>
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

    if (!prenom.value.trim() || !nom.value.trim() || !email.value.trim() || !motDePasse.value.trim() || !codePostal.value.trim()) {
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
.auth-page { display: grid; grid-template-columns: 1fr 1fr; min-height: calc(100vh - 86px); margin-top: -108px; }
.auth-side { display: flex; align-items: center; justify-content: center; padding: 48px; }
.auth-side--dark { background: #17201b; color: white; }
.auth-side--light { background: #ffffff; }
.auth-brand { display: flex; align-items: center; gap: 12px; font-size: 1.15rem; font-weight: 700; cursor: pointer; margin-bottom: 2rem; }
.auth-brand strong { color: #338454; }
.auth-card { width: min(480px, 100%); }
.auth-form { display: grid; gap: 0.85rem; margin-top: 1.5rem; }
.auth-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 0.85rem; }
.auth-form input { padding: 0.95rem 1rem; border: 1px solid #d9dfdb; border-radius: 12px; width: 100%; }
.auth-error { color: #b91c1c; background: #fee2e2; border: 1px solid #ef4444; padding: 10px; border-radius: 10px; }
.login-submit { margin-top: 0.5rem; background: #338454; color: white; border: none; border-radius: 12px; padding: 0.95rem 1rem; font-weight: 700; }
.cgu-line { display: flex; align-items: center; gap: 8px; }
</style>
