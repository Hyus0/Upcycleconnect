<template>
    <SiteNavbar
        :is-authenticated="isLoggedIn"
        :user-name="userName"
        variant="public"
    />
    <div class="login-page">
        <div class="login-left">
            <div class="login-left__content">
                <div
                    class="login-left__logo"
                    @click="$router.push('/')"
                    style="cursor: pointer"
                >
                    <img
                        src="../components/logo.png"
                        class="navbar__logo-img"
                    />
                    Upcycle <span class="accent">Connect</span>
                </div>
                <h1 class="login-left__title">
                    Heureux de vous <span class="accent">revoir</span>.
                </h1>
                <p class="login-left__desc">
                    Connectez-vous pour retrouver vos annonces, vos projets en
                    cours et continuer à faire grandir la communauté de
                    l'upcycling.
                </p>

                <div class="login-left__stats">
                    <div class="stat">
                        <strong>2.4t</strong>
                        <span>CO₂ évité / mois</span>
                    </div>
                    <div class="stat">
                        <strong>8k+</strong>
                        <span>Objets upcyclés</span>
                    </div>
                    <div class="stat">
                        <strong>340</strong>
                        <span>Artisans actifs</span>
                    </div>
                </div>

                <div class="login-testimonial">
                    <p class="login-testimonial__text">
                        "La simplicité de gestion de mes dépôts sur
                        UpcycleConnect m'a permis de me concentrer sur ce que
                        j'aime : créer."
                    </p>
                    <div class="login-testimonial__author">
                        <div class="login-testimonial__avatar">ML</div>
                        <div>
                            <strong>Marie L.</strong>
                            <span>Artisane ébéniste — Paris 11e</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="login-right">
            <div class="login-right__content">
                <div class="login-right__logo">
                    <img
                        src="../components/logo.png"
                        class="navbar__logo-img"
                    />
                    Upcycle <span class="accent-green">Connect</span>
                </div>

                <h2 class="login-right__title">Bon retour parmi nous</h2>
                <p class="login-right__subtitle">
                    Pas encore membre ?
                    <router-link to="/inscription" class="login-right__link">
                        Créer un compte gratuitement
                    </router-link>
                </p>

                <form @submit.prevent="handleLogin">
                    <div class="login-field">
                        <label>Adresse e-mail</label>
                        <input
                            type="email"
                            placeholder="marie.lambert@exemple.fr"
                            v-model="email"
                            :class="{ 'input--error': isEmailInvalid }"
                        />
                    </div>

                    <div class="login-field">
                        <div
                            style="
                                display: flex;
                                justify-content: space-between;
                                align-items: center;
                            "
                        >
                            <label>Mot de passe</label>
                            <router-link to="/connexion" class="login-right__link" style="font-size: 0.8rem">Oublié ?</router-link>
                        </div>
                        <input
                            type="password"
                            placeholder="••••••••••"
                            v-model="motDePasse"
                            autocomplete="current-password"
                        />
                    </div>

                    <div class="login-cgu">
                        <input
                            type="checkbox"
                            id="remember"
                            v-model="rememberMe"
                        />
                        <label for="remember">Se souvenir de moi</label>
                    </div>

                    <div
                        v-if="errorMessages.length > 0"
                        style="
                            color: #e74c3c;
                            background: #fdeaea;
                            padding: 10px;
                            border-radius: 4px;
                            margin-bottom: 15px;
                            font-size: 0.9rem;
                        "
                    >
                        <ul style="margin: 0; padding-left: 20px">
                            <li
                                v-for="(error, index) in errorMessages"
                                :key="index"
                            >
                                {{ error }}
                            </li>
                        </ul>
                    </div>

                    <button type="submit" class="login-submit">
                        Se connecter →
                    </button>
                </form>

                <div class="login-separator">
                    <span>ou continuer avec</span>
                </div>

                <button class="login-google">Google</button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from "vue";
import {useRouter} from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";

const router = useRouter();

const email = ref("");
const motDePasse = ref("");
const rememberMe = ref(false);
const errorMessages = ref([]);
const isLoggedIn = ref(Boolean(sessionStorage.getItem("userToken") || localStorage.getItem("userToken")));
const userName = ref("Marie Lambert");

const isEmailInvalid = computed(() => {
    return email.value.length > 0 && !email.value.includes("@");
});

async function handleLogin() {
    errorMessages.value = [];

    if (!email.value.trim() || !motDePasse.value.trim()) {
        errorMessages.value = ["Veuillez remplir tous les champs."];
        return;
    }

    const userData = {
        email: email.value,
        password: motDePasse.value,
    };

    try {
        const response = await fetch("http://localhost:8081/login", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(userData),
        });

        const data = await response.json();

        if (response.ok) {
            if (data.token) {
                localStorage.setItem("userToken", data.token);
            }
            localStorage.setItem("userId", data.userId);
            localStorage.setItem("userPrenom", data.prenom);
            localStorage.setItem("userNom", data.nom);
            
            alert("Connecté avec succès !");
            router.push("/profil");
            return;
        } else {
            errorMessages.value = Array.isArray(data) ? data : [data.message || "Erreur de connexion"];
        }
    } catch (error) {
        console.error("Détail :", error);
        errorMessages.value = ["Le serveur est injoignable."];
    }
}
</script>

<style scoped>
.login-page {
    display: flex;
    min-height: calc(100vh - 86px);
    margin-top: -108px;
}

.login-left {
    width: 50%;
    background-color: #17201b;
}

.login-right {
    width: 50%;
    background: #ffffff;
}

.login-left__content {
    padding: 60px 50px;
    height: 100vh;
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
}

.login-left__logo {
    display: flex;
    align-items: center;
    gap: 12px;
    color: #ffffff;
    font-family: "Syne", sans-serif;
    font-size: 1.1rem;
    font-weight: 700;
    margin-bottom: 0;
    padding-bottom: 60px;
}

.accent {
    color: #338454;
}

.login-left__title {
    font-family: "Syne", sans-serif;
    font-size: 2.8rem;
    font-weight: 800;
    color: #ffffff;
    line-height: 1.15;
    margin: 0 0 24px 0;
}

.login-left__desc {
    color: rgba(255, 255, 255, 0.6);
    font-size: 0.95rem;
    line-height: 1.6;
    margin: 0 0 48px 0;
    max-width: 340px;
}

.login-left__stats {
    display: flex;
    gap: 40px;
}

.navbar__logo {
    display: flex;
    align-items: center;
    gap: 12px;
    cursor: pointer;
}

.navbar__logo-img {
    height: 35px;
    width: auto;
    object-fit: contain;
}

.stat strong {
    display: block;
    font-family: "Syne", sans-serif;
    font-size: 2rem;
    font-weight: 800;
    color: #ffffff;
}

.stat span {
    display: block;
    font-size: 0.75rem;
    color: rgba(255, 255, 255, 0.5);
    margin-top: 4px;
}

.login-testimonial {
    margin-top: auto;
    padding: 20px 24px;
    border-radius: 14px;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.1);
}

.login-testimonial__text {
    font-size: 0.9rem;
    color: rgba(255, 255, 255, 0.7);
    line-height: 1.6;
    margin: 0;
    font-style: italic;
}

.login-testimonial__author {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-top: 16px;
}

.login-testimonial__avatar {
    width: 36px;
    height: 36px;
    border-radius: 10px;
    background: #338454;
    color: white;
    font-size: 0.75rem;
    font-weight: 700;
    display: flex;
    align-items: center;
    justify-content: center;
}

.login-testimonial__author strong {
    display: block;
    font-size: 0.85rem;
    color: #ffffff;
}

.login-testimonial__author span {
    display: block;
    font-size: 0.75rem;
    color: rgba(255, 255, 255, 0.5);
    margin-top: 2px;
}

.login-right__logo {
    display: flex;
    align-items: center;
    gap: 12px;
    color: #17201b;
    font-family: "Syne", sans-serif;
    font-size: 1.1rem;
    font-weight: 700;
    margin-bottom: 0;
    padding-bottom: 30px;
}

.accent-green {
    color: #338454;
}
.login-right__content {
    padding: 60px 50px;
    max-width: 520px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    justify-content: center;
    min-height: 100vh;
}

.login-right__title {
    font-family: "Syne", sans-serif;
    font-size: 2rem;
    font-weight: 800;
    color: #17201b;
    margin: 0 0 8px 0;
    white-space: nowrap;
}

.login-right__subtitle {
    font-size: 0.9rem;
    color: #7b857f;
    margin: 0 0 28px 0;
}

.login-right__link {
    color: #338454;
    text-decoration: none;
    font-weight: 500;
}

.login-right__label {
    font-size: 0.9rem;
    color: #17201b;
    margin: 0 0 10px 0;
    font-weight: 500;
}

.login-type {
    display: flex;
    gap: 12px;
    margin-bottom: 24px;
}

.login-type__btn {
    flex: 1;
    padding: 14px;
    border-radius: 12px;
    border: 1.5px solid #dfe9e1;
    background: #ffffff;
    font-size: 0.9rem;
    color: #17201b;
    font-weight: 500;
    transition: all 0.2s ease;
}

.login-type__btn--active {
    border-color: #338454;
    background: #ecf8f0;
    color: #338454;
    font-weight: 700;
}

.login-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 14px;
}

.login-field {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 16px;
}

.login-field label {
    font-size: 0.85rem;
    font-weight: 500;
    color: #17201b;
}

.login-field input {
    padding: 12px 14px;
    border-radius: 10px;
    border: 1.5px solid #dfe9e1;
    background: #ffffff;
    font-size: 0.9rem;
    color: #17201b;
    outline: none;
    transition: border-color 0.2s ease;
}

.login-field input:focus {
    border-color: #338454;
}

.login-cgu {
    display: flex;
    align-items: flex-start;
    gap: 10px;
    margin-bottom: 20px;
    font-size: 0.85rem;
    color: #7b857f;
}

.login-cgu input[type="checkbox"] {
    margin-top: 2px;
    accent-color: #338454;
}

.login-submit {
    width: 100%;
    padding: 16px;
    border-radius: 12px;
    border: none;
    background: #338454;
    color: #ffffff;
    font-size: 1rem;
    font-weight: 700;
    margin-bottom: 16px;
    transition: background 0.2s ease;
}

.login-submit:hover {
    background: #1f6d43;
}

.login-separator {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
    color: #7b857f;
    font-size: 0.85rem;
}

.login-separator::before,
.login-separator::after {
    content: "";
    flex: 1;
    height: 1px;
    background: #dfe9e1;
}

.login-google {
    width: 100%;
    padding: 14px;
    border-radius: 12px;
    border: 1.5px solid #dfe9e1;
    background: #ffffff;
    font-size: 0.95rem;
    font-weight: 500;
    color: #17201b;
    transition: background 0.2s ease;
}

.login-google:hover {
    background: #f5f8f5;
}

.input--error {
    border-color: #ff4d4f !important;
}

.error-text {
    color: #ff4d4f;
    font-size: 0.75rem;
    margin-top: -10px;
    margin-bottom: 10px;
}

.login-field div label {
    margin-bottom: 0;
}

.input--error {
    border-color: #ff4d4f !important;
    background-color: #fff2f0 !important;
}
</style>
