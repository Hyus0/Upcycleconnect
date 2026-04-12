<template>
  <div class="register-page">

    <!-- Colonne gauche : panneau vert -->
    <div class="register-left">
          <div class="register-left__content">

    <!-- Logo -->
    <div class="register-left__logo">
      🔄 Upcycle <span class="accent">Connect</span>
    </div>

    <!-- Titre principal -->
    <h1 class="register-left__title">
      Rejoignez la <span class="accent">communauté</span><br> qui agit.
    </h1>

    <!-- Description -->
    <p class="register-left__desc">
      Plus de 12 400 membres, 340 artisans et des <br>centaines de projets d'upcycling chaque mois.
      <br>Commencez gratuitement dès aujourd'hui.
    </p>

    <!-- Stats -->
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
    <!-- Témoignage -->
     <div class="register-testimonial">
        <p class="register-testimonial__text">
            "Grâce à UpcycleConnect, j'ai trouvé les matériaux parfaits pour mes créations. La plateforme a transformé mon activité d'artisan."
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

    <!-- Colonne droite : formulaire -->
<div class="register-right">
  <div class="register-right__content">
    <!-- Logo -->
  <div class="register-right__logo">
    🔄 Upcycle <span class="accent-green">Connect</span>
  </div>

    <!-- En-tête -->
    <h2 class="register-right__title">Créer un compte</h2>
    <p class="register-right__subtitle">
      Déjà membre ? <a href="#" class="register-right__link">Se connecter</a>
    </p>

    <!-- Choix du type de compte -->
    <p class="register-right__label">Je suis...</p>
    <div class="register-type">
<button 
  class="register-type__btn"
  :class="{ 'register-type__btn--active': accountType === 'particulier' }"
  @click="accountType = 'particulier'"
>
  🏠 <br>Particulier
</button>
<button 
  class="register-type__btn"
  :class="{ 'register-type__btn--active': accountType === 'pro' }"
  @click="accountType = 'pro'"
>
  🔨<br> Pro / Artisan
</button>
    </div>

    <!-- Champs prénom / nom -->
    <div class="register-row">
      <div class="register-field">
        <label>Prénom</label>
        <input type="text" placeholder="Marie" v-model="prenom" />
      </div>
      <div class="register-field">
        <label>Nom</label>
        <input type="text" placeholder="Lambert" v-model="nom" />
      </div>
    </div>

    <!-- Email -->
    <div class="register-field">
      <label>Adresse e-mail</label>
      <input type="email" placeholder="marie.lambert@exemple.fr" v-model="email" />
    </div>

    <!-- Mot de passe -->
    <div class="register-field">
      <label>Mot de passe</label>
      <input type="password" placeholder="••••••••••" v-model="motDePasse" />
    </div>

    <!-- Code postal -->
    <div class="register-field">
      <label>Code postal</label>
      <input type="text" placeholder="75011" v-model="codePostal" />
    </div>

    <!-- CGU -->
    <div class="register-cgu">
      <input type="checkbox" id="cgu" v-model="cguAccepte" />
      <label for="cgu">
        J'accepte les <a href="#" class="register-right__link">CGU</a> et la
        <a href="#" class="register-right__link">politique de confidentialité</a>
      </label>
    </div>

    <!-- Bouton principal -->
    <button class="register-submit" @click="handleSubmit">
      Créer mon compte gratuitement →
    </button>

    <!-- Séparateur -->
    <div class="register-separator">
      <span>ou</span>
    </div>

    <!-- Google -->
    <button class="register-google">
      Continuer avec Google
    </button>

  </div>
</div>

  </div>
</template>

<script setup>
import { ref } from "vue";

// Type de compte sélectionné
const accountType = ref("particulier");

// Champs du formulaire
const prenom = ref("");
const nom = ref("");
const email = ref("");
const motDePasse = ref("");
const codePostal = ref("");
const cguAccepte = ref(false);

function handleSubmit() {
  if (!cguAccepte.value) {
    alert("Veuillez accepter les CGU pour continuer.");
    return;
  }

  console.log("Formulaire envoyé :", {
    accountType: accountType.value,
    prenom: prenom.value,
    nom: nom.value,
    email: email.value,
    motDePasse: motDePasse.value,
    codePostal: codePostal.value,
  });

  alert(`Compte créé avec succès ! Bienvenue ${prenom.value} 🎉`);
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
  color: #ffffff;
  font-family: 'Syne', sans-serif;
  font-size: 1.1rem;
  font-weight: 700;
  margin-bottom: 0;
  padding-bottom: 60px;
}

.accent {
  color: #338454;
}

.register-left__title {
  font-family: 'Syne', sans-serif;
  font-size: 2.8rem;
  font-weight: 800;
  color: #ffffff;
  line-height: 1.15;
  margin: 0 0 24px 0;
}

.register-left__desc {
  color: rgba(255,255,255,0.6);
  font-size: 0.95rem;
  line-height: 1.6;
  margin: 0 0 48px 0;
  max-width: 340px;
}

.register-left__stats {
  display: flex;
  gap: 40px;
}

.stat strong {
  display: block;
  font-family: 'Syne', sans-serif;
  font-size: 2rem;
  font-weight: 800;
  color: #ffffff;
}

.stat span {
  display: block;
  font-size: 0.75rem;
  color: rgba(255,255,255,0.5);
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
  color: rgba(255,255,255,0.5);
  margin-top: 2px;
}

.register-right__logo {
  font-family: 'Syne', sans-serif;
  font-weight: 700;
  font-size: 1rem;
  color: #17201b;
  margin-bottom: 32px;
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
  font-family: 'Syne', sans-serif;
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

/* Boutons Particulier / Pro */
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

/* Rangée prénom / nom */
.register-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px;
}

/* Champs */
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

/* CGU */
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

/* Bouton principal */
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

/* Séparateur */
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
  content: '';
  flex: 1;
  height: 1px;
  background: #dfe9e1;
}

/* Bouton Google */
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
</style>