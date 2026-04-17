<template>
    <div class="register-page">
        <div class="register-left">
            <div class="register-left__content">
                <div
                    class="register-left__logo"
                    @click="$router.push('/')"
                    style="cursor: pointer"
                >
                    <img
                        src="../components/logo.png"
                        class="navbar__logo-img"
                    />
                    Upcycle <span class="accent">Connect</span>
                </div>
                <h1 class="register-left__title">
                    Rejoignez la <span class="accent">communauté</span><br />
                    qui agit.
                </h1>
                <p class="register-left__desc">
                    Plus de 12 400 membres, 340 artisans et des <br />centaines
                    de projets d'upcycling chaque mois. <br />Commencez
                    gratuitement dès aujourd'hui.
                </p>

                <div class="register-left__stats">
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

                <div class="register-testimonial">
                    <p class="register-testimonial__text">
                        "Grâce à UpcycleConnect, j'ai trouvé les matériaux
                        parfaits pour mes créations. La plateforme a transformé
                        mon activité d'artisan."
                    </p>
                    <div class="register-testimonial__author">
                        <div class="register-testimonial__avatar">ML</div>
                        <div>
                            <strong>Marie L.</strong>
                            <span>Artisane ébéniste — Paris 11e</span>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <div class="register-right">
            <div class="register-right__content">
                <div class="register-right__logo">
                    <img
                        src="../components/logo.png"
                        class="navbar__logo-img"
                    />
                    Upcycle <span class="accent-green">Connect</span>
                </div>

                <h2 class="register-right__title">Créer un compte</h2>
                <p class="register-right__subtitle">
                    Déjà membre ?
                    <router-link to="/connexion" class="register-right__link">
                        Se connecter
                    </router-link>
                </p>

                <p class="register-right__label">Je suis...</p>
                <div class="register-type">
                    <button
                        class="register-type__btn"
                        :class="{
                            'register-type__btn--active':
                                accountType === 'particulier',
                        }"
                        @click="accountType = 'particulier'"
                    >
                        🏠 <br />Particulier
                    </button>
                    <button
                        class="register-type__btn"
                        :class="{
                            'register-type__btn--active': accountType === 'pro',
                        }"
                        @click="accountType = 'pro'"
                    >
                        🔨<br />
                        Pro / Artisan
                    </button>
                </div>

                <div class="register-row">
                    <div class="register-field">
                        <label>Prénom</label>
                        <input
                            type="text"
                            placeholder="Marie"
                            v-model="prenom"
                        />
                    </div>
                    <div class="register-field">
                        <label>Nom</label>
                        <input
                            type="text"
                            placeholder="Lambert"
                            v-model="nom"
                        />
                    </div>
                </div>

                <div class="register-field">
                    <label>Adresse e-mail</label>
                    <input
                        type="email"
                        placeholder="marie.lambert@exemple.fr"
                        v-model="email"
                        :class="{ 'input--error': isEmailInvalid }"
                    />
                </div>

                <div class="register-field">
                    <label>Mot de passe</label>
                    <input
                        type="password"
                        placeholder="••••••••••"
                        v-model="motDePasse"
                    />
                </div>
                
                <div v-if="errorMessages.length > 0" 
                style="background-color: #fee2e2; border: 1px solid #ef4444; color: #b91c1c; 
                padding: 10px; border-radius: 8px; margin-bottom: 15px;">
                    <ul style="margin: 0; padding-left: 20px;">
                        <li v-for="(err, index) in errorMessages" :key="index">
                            {{ err }}
                        </li>
                    </ul>
                </div>
                
                <div class="register-field">
                    <label>Code postal</label>
                    <input
                        type="text"
                        placeholder="75011"
                        v-model="codePostal"
                        :class="{ 'input--error': isPostalCodeInvalid }"
                    />
                </div>

                <div class="register-cgu">
                    <input type="checkbox" id="cgu" v-model="cguAccepte" />
                    <label for="cgu">
                        J'accepte les
                        <a href="#" class="register-right__link">CGU</a> et la
                        <a href="#" class="register-right__link"
                            >politique de confidentialité</a
                        >
                    </label>
                </div>

                <button class="register-submit" @click="handleSubmit">
                    Créer mon compte gratuitement →
                </button>

                <div class="register-separator">
                    <span>ou</span>
                </div>

                <button class="register-google">Continuer avec Google</button>
            </div>
        </div>
    </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const accountType = ref("particulier");

const errorMessages = ref([]);
const successMessage = ref('');
const prenom = ref("");
const nom = ref("");
const email = ref("");
const motDePasse = ref("");
const codePostal = ref("");
const cguAccepte = ref(false);

const isEmailInvalid = computed(() => {
    return email.value.length > 0 && !email.value.includes("@");
});

const isPostalCodeInvalid = computed(() => {
    const onlyDigits = /^\d+$/.test(codePostal.value);
    return (
        codePostal.value.length > 0 &&
        (codePostal.value.length !== 5 || !onlyDigits)
    );
});

async function handleSubmit() {
    errorMessages.value = [];
    
    if (!prenom.value.trim() || !nom.value.trim() || !email.value.trim() || 
        !motDePasse.value.trim() || !codePostal.value.trim()) {
        errorMessages.value = ["Il manque des informations."];
        return;
    }

    const userData = {
        prenom: prenom.value,
        nom: nom.value,
        mail: email.value,
        password: motDePasse.value,
        code_postal: codePostal.value,
        role: "Particulier", 
        id_langue: 1     
    };
    
    try {
        const response = await fetch("http://localhost:8081/users", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(userData),
        });

        if (response.ok) {
            alert("Compte créé avec succès ! 🎉");
            router.push("/connexion");
            return;
        } else {
            const data = await response.json();
            if (Array.isArray(data)) {
                errorMessages.value = data;
            } else {
                errorMessages.value = [data.message || "Erreur de validation"];
            }
        }

    } catch (error) {
        console.error("Détail :", error);
        
        if (error instanceof SyntaxError) {
            errorMessages.value = ["Le serveur a réussi mais a renvoyé une réponse illisible."];
        } else {
            errorMessages.value = ["Le serveur est injoignable (éteint ou problème réseau)."];
        }
    }
}
</script>

<style scoped>
.register-page {
    display: flex;
    min-height: 100vh;
}

.register-left {
    width: 50%;
    background-color: #17201b;
}

.register-right {
    width: 50%;
    background: #ffffff;
}

.register-left__content {
    padding: 60px 50px;
    height: 100vh;
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
}

.register-left__logo {
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

.register-left__title {
    font-family: "Syne", sans-serif;
    font-size: 2.8rem;
    font-weight: 800;
    color: #ffffff;
    line-height: 1.15;
    margin: 0 0 24px 0;
}

.register-left__desc {
    color: rgba(255, 255, 255, 0.6);
    font-size: 0.95rem;
    line-height: 1.6;
    margin: 0 0 48px 0;
    max-width: 340px;
}

.register-left__stats {
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

.register-testimonial {
    margin-top: auto;
    padding: 20px 24px;
    border-radius: 14px;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.1);
}

.register-testimonial__text {
    font-size: 0.9rem;
    color: rgba(255, 255, 255, 0.7);
    line-height: 1.6;
    margin: 0;
    font-style: italic;
}

.register-testimonial__author {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-top: 16px;
}

.register-testimonial__avatar {
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

.register-testimonial__author strong {
    display: block;
    font-size: 0.85rem;
    color: #ffffff;
}

.register-testimonial__author span {
    display: block;
    font-size: 0.75rem;
    color: rgba(255, 255, 255, 0.5);
    margin-top: 2px;
}

.register-right__logo {
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
.register-right__content {
    padding: 60px 50px;
    max-width: 520px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    justify-content: center;
    min-height: 100vh;
}

.register-right__title {
    font-family: "Syne", sans-serif;
    font-size: 2rem;
    font-weight: 800;
    color: #17201b;
    margin: 0 0 8px 0;
    white-space: nowrap;
}

.register-right__subtitle {
    font-size: 0.9rem;
    color: #7b857f;
    margin: 0 0 28px 0;
}

.register-right__link {
    color: #338454;
    text-decoration: none;
    font-weight: 500;
}

.register-right__label {
    font-size: 0.9rem;
    color: #17201b;
    margin: 0 0 10px 0;
    font-weight: 500;
}

.register-type {
    display: flex;
    gap: 12px;
    margin-bottom: 24px;
}

.register-type__btn {
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

.register-type__btn--active {
    border-color: #338454;
    background: #ecf8f0;
    color: #338454;
    font-weight: 700;
}

.register-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 14px;
}

.register-field {
    display: flex;
    flex-direction: column;
    gap: 6px;
    margin-bottom: 16px;
}

.register-field label {
    font-size: 0.85rem;
    font-weight: 500;
    color: #17201b;
}

.register-field input {
    padding: 12px 14px;
    border-radius: 10px;
    border: 1.5px solid #dfe9e1;
    background: #ffffff;
    font-size: 0.9rem;
    color: #17201b;
    outline: none;
    transition: border-color 0.2s ease;
}

.register-field input:focus {
    border-color: #338454;
}

.register-cgu {
    display: flex;
    align-items: flex-start;
    gap: 10px;
    margin-bottom: 20px;
    font-size: 0.85rem;
    color: #7b857f;
}

.register-cgu input[type="checkbox"] {
    margin-top: 2px;
    accent-color: #338454;
}

.register-submit {
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

.register-submit:hover {
    background: #1f6d43;
}

.register-separator {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 16px;
    color: #7b857f;
    font-size: 0.85rem;
}

.register-separator::before,
.register-separator::after {
    content: "";
    flex: 1;
    height: 1px;
    background: #dfe9e1;
}

.register-google {
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

.register-google:hover {
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
</style>
