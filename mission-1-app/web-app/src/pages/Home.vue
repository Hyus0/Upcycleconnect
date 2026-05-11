<template>
    <div class="app-container">
        <main>
            <SiteNavbar
              :is-authenticated="isLoggedIn"
              :user-name="userName"
              user-role="Particulier"
              :user-score="userScore"
            />
            <section class="page-background">
                <div class="hero-content">
                    <button class="active-plateform">
                        Plateforme active · {{ plateform_user }}+ membres
                    </button>

                    <h1 class="hero-title">
                        Upcycle<span class="text-accent">Connect</span>
                    </h1>
                    <p class="hero-subtitle">
                        UpcycleConnect réunit particuliers, artisans et
                        entreprises autour<br />
                        de l'upcycling. Déposez, échangez, transformez — et
                        mesurez<br />
                        votre impact environnemental.
                    </p>

                    <div class="hero-buttons">
                        <button
                            class="btn btn--primary"
                            @click="$router.push('/inscription')"
                        >
                            ♻ Commencer gratuitement
                        </button>
                        <button
                            class="btn btn--outline"
                            @click="$router.push('/profil')"
                        >
                            Je suis professionnel →
                        </button>
                    </div>

                    <div class="hero-divider"></div>

                    <div class="hero-stats">
                        <div class="stat-item">
                            <h2 class="stat-number">
                                {{ formatCO2(platformStats.co2_evite) }}<span class="text-accent">{{ uniteCO2 }}</span>
                            </h2>
                            <p class="stat-label">CO₂ évité au total</p>
                        </div>
                        
                        <div class="stat-item">
                            <h2 class="stat-number">
                                {{ platformStats.objets_upcycles >= 1000 ? (platformStats.objets_upcycles / 1000).toFixed(1) : platformStats.objets_upcycles }}
                                <span class="text-accent" v-if="platformStats.objets_upcycles >= 1000">k+</span>
                            </h2>
                            <p class="stat-label">Objets upcyclés</p>
                        </div>
                        
                        <div class="stat-item">
                            <h2 class="stat-number">
                                {{ platformStats.artisans_actifs }}
                            </h2>
                            <p class="stat-label">Artisans partenaires</p>
                        </div>
                        
                        <div class="stat-item">
                            <h2 class="stat-number">{{ platformStats.sites_actifs }}</h2>
                            <p class="stat-label">Sites à Paris & IDF</p>
                        </div>
                    </div>
                </div>
            </section>

            <section class="all-you-need-part">
                <div class="alyn-container">
                    <p class="subtitle">FONCTIONNALITÉS</p>
                    <h1>
                        Tout ce dont vous<br />
                        avez besoin
                    </h1>
                    <p class="hero-subtitle1">
                        Une plateforme complète pour donner une seconde vie à
                        vos objets, du<br />
                        dépôt jusqu'à la création finale.
                    </p>
                    <section class="features-section">
                        <div class="features-grid">
                            <div class="feature-card">
                                <div class="feature-icon-box">📢</div>
                                <h3 class="feature-title">
                                    Annonces intelligentes
                                </h3>
                                <p class="feature-description">
                                    Publiez vos dons ou ventes avec photos,
                                    catégorie et localisation. Les artisans
                                    reçoivent des alertes personnalisées selon
                                    leurs besoins.
                                </p>
                            </div>

                            <div class="feature-card">
                                <div class="feature-icon-box">📦</div>
                                <h3 class="feature-title">
                                    Conteneurs sécurisés
                                </h3>
                                <p class="feature-description">
                                    Demandez un code-barres pour déposer votre
                                    objet dans l'un de nos 30+ conteneurs à
                                    Paris et en IDF. Simple, rapide, sécurisé.
                                </p>
                            </div>

                            <div class="feature-card">
                                <div class="feature-icon-box">🌿</div>
                                <h3 class="feature-title">Upcycling Score</h3>
                                <p class="feature-description">
                                    Chaque action compte. Visualisez en temps
                                    réel le CO₂ évité, les ressources
                                    économisées et votre contribution à
                                    l'économie circulaire.
                                </p>
                            </div>

                            <div class="feature-card">
                                <div class="feature-icon-box">🎓</div>
                                <h3 class="feature-title">
                                    Formations & Ateliers
                                </h3>
                                <p class="feature-description">
                                    Apprenez les techniques d'upcycling avec nos
                                    formateurs experts. Cours en ligne et
                                    ateliers en présentiel dans nos 7 sites.
                                </p>
                            </div>

                            <div class="feature-card">
                                <div class="feature-icon-box">🗓️</div>
                                <h3 class="feature-title">
                                    Planning personnel
                                </h3>
                                <p class="feature-description">
                                    Gérez vos inscriptions, formations, dépôts
                                    et événements depuis un calendrier
                                    centralisé et synchronisable.
                                </p>
                            </div>

                            <div class="feature-card">
                                <div class="feature-icon-box">🤝</div>
                                <h3 class="feature-title">
                                    Communauté vivante
                                </h3>
                                <p class="feature-description">
                                    Forums, événements, conseils et projets
                                    partagés. Rejoignez une communauté de
                                    passionnés qui agissent pour
                                    l'environnement.
                                </p>
                            </div>
                        </div>
                    </section>
                </div>
            </section>

            <section id="processus" class="processus-part">
                <div class="processus-container">
                    <p class="subtitle">PROCESSUS</p>
                    <h1>Comment ça marche ?</h1>

                    <section class="how-it-works">
                        <div class="steps-container">
                            <div class="steps-line"></div>
                            <div class="step-item">
                                <div class="step-number">1</div>
                                <h3>Déposez une annonce</h3>
                                <p>
                                    Photographiez et décrivez l'objet à donner
                                    ou vendre sur la plateforme.
                                </p>
                            </div>

                            <div class="step-item">
                                <div class="step-number">2</div>
                                <h3>Validation</h3>
                                <p>
                                    Notre équipe valide votre objet et vous
                                    envoie un code-barres d'accès.
                                </p>
                            </div>

                            <div class="step-item">
                                <div class="step-number">3</div>
                                <h3>Dépôt en conteneur</h3>
                                <p>
                                    Déposez l'objet dans le conteneur le plus
                                    proche grâce à votre code-barres.
                                </p>
                            </div>

                            <div class="step-item">
                                <div class="step-number">4</div>
                                <h3>Seconde vie !</h3>
                                <p>
                                    Un artisan récupère l'objet et lui donne une
                                    nouvelle vie.
                                </p>
                            </div>
                        </div>
                    </section>
                </div>
            </section>

            <section class="ready-part">
                <div class="ready-container">
                    <h1 class="hero-title1">
                        Prêt à agir pour<br />
                        <span class="navbar__logo--accent">la planète ?</span>
                    </h1>
                    <p class="hero-subtitle1">
                        Rejoignez les {{ plateform_user }} membres qui donnent
                        une seconde vie à leurs objets chaque jour.
                    </p>

                    <div class="hero-buttons">
                        <button
                            class="btn btn--primary"
                            @click="$router.push('/inscription')"
                        >
                            Créer mon compte gratuit
                        </button>
                        <button
                            class="btn btn--outline"
                            @click="$router.push('/annonces')"
                        >
                            Voir les annonces
                        </button>
                    </div>
                </div>
            </section>
        </main>
    </div>
    <SiteFooter />
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import SiteNavbar from "../components/SiteNavbar.vue";

import SiteFooter from "../components/SiteFooter.vue";

const API_URL = "http://localhost:8081";

const isLoggedIn = computed(() => {
  return !!sessionStorage.getItem("userToken");
});

const userName = computed(() => {
  const prenom = sessionStorage.getItem("userPrenom") || "";
  const nom = sessionStorage.getItem("userNom") || "";
  
  return (prenom || nom) ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

// Les variables réactives pour stocker les vraies données de la BDD
const platformStats = ref({
    co2_evite: 0,
    objets_upcycles: 0,
    artisans_actifs: 0,
    sites_actifs: 7 // On garde 7 par défaut si on n'a pas la table SITE
});

// Fonction pour formater le CO2
const formatCO2 = (kg) => {
    if (kg >= 1000) {
        return (kg / 1000).toFixed(1); // Retourne juste le chiffre (ex: 2.4)
    }
    return kg; // Retourne juste le chiffre
};

// L'unité de mesure associée au CO2 (t ou kg)
const uniteCO2 = computed(() => {
    if (platformStats.value.co2_evite >= 1000) {
        return 't';
    }
    return 'kg';
});

// Fonction pour récupérer les vraies statistiques
const fetchPlatformStats = async () => {
    try {
        const res = await fetch(`${API_URL}/stats/platform`);
        if (res.ok) {
            platformStats.value = await res.json();
        }
    } catch (e) {
        console.error("Erreur chargement des statistiques :", e);
    }
};

onMounted(() => {
    fetchPlatformStats();
});
</script>

<style>
body,
html {
    margin: 0;
    padding: 0;
    font-family: "Syne", sans-serif;
}

.app-container main {
    position: relative;
}

.logo-text {
    display: inline-block;
    white-space: nowrap;
    color: #ffffff;
    font-size: 1.4rem;
    font-weight: 800;
    cursor: pointer;
}

.navbar__logo--accent {
    color: #4ade80;
}

.btn {
    padding: 0.6rem 1.4rem;
    border-radius: 10px;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: 0.3s ease;
    border: 1px solid transparent;
}

.btn--outline {
    background: transparent;
    color: #ffffff;
    border-color: rgba(255, 255, 255, 0.2);
}

.btn--outline:hover {
    background: rgba(255, 255, 255, 0.05);
    border-color: #ffffff;
}

.btn--primary {
    background-color: var(--brand-green);
    color: #ffffff;
}

.page-background {
    min-height: 100vh;
    width: 100%;
    background-image:
        linear-gradient(rgba(0, 0, 0, 0.4), rgba(0, 0, 0, 0.4)),
        url("../components/background.png");
    background-size: cover;
    background-position: center;
    color: white;
    display: flex;
    align-items: center;
    padding-left: 10px;
    margin-top: -108px;
}

.active-plateform {
    width: fit-content;
    background: transparent;
    padding: 0.6rem 1rem;
    color: #4bab73;
    border-radius: 20px;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: 0.3s ease;
    border: 1px solid #4bab73;
    display: inline-block;
}
.hero-content {
    max-width: 900px;
    display: flex;
    flex-direction: column;
    padding: 165px 30px;
}

.hero-title {
    font-size: 3.5rem;
    font-weight: 800;
    line-height: 1.1;
    margin-bottom: 1.5rem;
}

.hero-title1 {
    font-size: 2.8rem;
    font-weight: 800;
    line-height: 1.1;
    margin-bottom: 1.5rem;
}
.text-accent {
    color: #4ade80;
}

.hero-subtitle {
    font-size: 1.1rem;
    color: rgba(255, 255, 255, 0.8);
    margin-bottom: 2rem;
}

.hero-buttons {
    display: flex;
    gap: 1.5rem;
    margin-bottom: 3.5rem;
}

.hero-divider {
    width: 100%;
    height: 1px;
    background: linear-gradient(
        to right,
        transparent,
        rgba(255, 255, 255, 0.2),
        transparent
    );
    margin-bottom: 1rem;
}

.hero-stats {
    display: flex;
    justify-content: center;
    gap: 5rem;
    width: 100%;
}

.stat-item {
    display: flex;
    flex-direction: column;
    align-items: center;
}

.stat-number {
    font-size: 2.5rem;
    font-weight: 800;
    color: #ffffff;
    margin: 0;
}

.stat-label {
    font-size: 0.75rem;
    color: rgba(255, 255, 255, 0.5);
    text-transform: uppercase;
    letter-spacing: 1px;
    margin-top: 0.5rem;
}

.all-you-need-part {
    padding: 50px 5%;
    background-color: #ffffff;
    color: black;
}

.alyn-container {
    max-width: 1300px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    text-align: left;
    padding: 50px 50px;
}

.hero-subtitle1 {
    font-size: 1.1rem;
    color: rgba(90, 102, 96, 100);
    margin-bottom: 2rem;
}

.subtitle {
    font-size: 1.1rem;
    color: var(--brand-green);
    margin-bottom: 2rem;
}

.features-grid {
    display: grid;
    width: 100%;
    max-width: 1400px;
    grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
    gap: 30px;
    margin: 0 auto;
}

.feature-card {
    background: #ffffff;
    padding: 30px;
    border-radius: 30px;
    border: 1px solid var(--border);
    display: flex;
    flex-direction: column;
    align-items: flex-start;
    transition: all 0.3s ease;
}

.feature-icon-box {
    width: 45px;
    height: 45px;
    background-color: #f2f4f3;
    border-radius: 12px;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 1.4rem;
    margin-bottom: 25px;
}

.feature-title {
    font-family: "Syne", sans-serif;
    font-size: 1.25rem;
    font-weight: 700;
    color: #1a1a1a;
    margin: 0 0 15px 0;
}

.feature-description {
    font-size: 0.95rem;
    line-height: 1.6;
    color: #666666;
    margin: 0;
}

.processus-part {
    padding: 50px 5%;
    background-color: var(--brand-mist);
    color: black;
}

.processus-container {
    max-width: 1300px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 50px 50px;
}

.ready-part {
    background-image: url("../components/pret-a-agir.png");
    background-size: cover;
    background-position: center;
    background-attachment: scroll;
    padding: 50px 5%;
    color: white;
}

.ready-container {
    max-width: 1300px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
}

.footer-class {
    background-color: #0b0f0e;
    padding: 80px 5% 40px 5%;
    color: #ffffff;
}

.footer-container {
    max-width: 1300px;
    margin: 0 auto;
}

.footer-main-content {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    gap: 150px;
    flex-wrap: wrap;
}

.footer-brand {
    flex: 2;
    min-width: 250px;
}

.footer-description {
    margin-top: 20px;
    color: #888;
    line-height: 1.6;
    font-size: 0.9rem;
    max-width: 300px;
}

.footer-links-wrapper {
    flex: 3;
    display: flex;
    justify-content: space-between;
    gap: 30px;
}

.footer-column h4 {
    font-size: 0.85rem;
    letter-spacing: 1.5px;
    margin-top: 0;
    margin-bottom: 20px;
    padding-top: 0;
    color: #ffffff;
    font-weight: 600;
}

.footer-column ul {
    list-style: none;
    padding: 0;
}

.footer-column ul li {
    margin-bottom: 12px;
}

.footer-column ul li a {
    color: #888;
    text-decoration: none;
    font-size: 0.9rem;
    transition: color 0.3s;
}

.footer-column ul li a:hover {
    color: #4ade80;
}

.footer-divider {
    height: 1px;
    background: rgba(255, 255, 255, 0.05);
    margin: 60px 0 30px 0;
}

.footer-bottom {
    display: flex;
    justify-content: space-between;
    color: #555;
    font-size: 0.85rem;
}

.footer-legal {
    display: flex;
    gap: 20px;
}

.footer-legal a {
    color: #555;
    text-decoration: none;
}

.how-it-works {
    padding: 50px 2%;
}

.steps-container {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    max-width: 1200px;
    margin: 0 auto;
    position: relative;
}

.steps-line {
    position: absolute;
    top: 45px;
    left: 10%;
    right: 10%;
    height: 2px;
    background-color: #d1e2d9;
    z-index: 1;
}

.step-item {
    flex: 1;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    z-index: 2;
    padding: 0 20px;
}

.step-number {
    width: 90px;
    height: 90px;
    background-color: #2d6a4f;
    color: white;
    border-radius: 50%;
    display: flex;
    justify-content: center;
    align-items: center;
    font-size: 2rem;
    font-weight: 900;
    margin-bottom: 30px;
    border: 8px solid #f9fbf9;
    box-shadow: 0 0 0 8px #e8f3ee;
}

.step-item h3 {
    font-size: 1.2rem;
    margin-bottom: 15px;
    color: #1a1a1a;
}

.step-item p {
    font-size: 0.95rem;
    color: #666;
    line-height: 1.5;
}
</style>
