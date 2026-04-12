<template>
    <div class="app-container">
        <header class="navbar">
            <div class="navbar__container">
                <div class="navbar__logo">
                    <img
                        src="../components/logo.png"
                        class="navbar__logo-img"
                    />

                    <div class="logo-text">
                        Upcycle<span class="navbar__logo--accent">Connect</span>
                    </div>
                </div>
                <div class="navbar__right-side">
                    <nav class="navbar__menu">
                        <ul class="navbar__list">
                            <li>
                                <a href="#" class="navbar__link"
                                    >Comment ça marche</a
                                >
                            </li>
                            <li>
                                <a href="#" class="navbar__link">Annonces</a>
                            </li>
                            <li>
                                <a href="#" class="navbar__link">Formations</a>
                            </li>
                            <li>
                                <a href="#" class="navbar__link">Communauté</a>
                            </li>
                        </ul>
                    </nav>

                    <div class="navbar__actions">
                        <button class="btn btn--outline">Se connecter</button>
                        <button class="btn btn--primary" @click="$router.push('/inscription')">
                            Rejoindre gratuitement →
                        </button>
                    </div>
                </div>
            </div>
        </header>
        <main>
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
                        <button class="btn btn--primary" @click="$router.push('/inscription')">
                            ♻ Commencer gratuitement
                        </button>
                        <button class="btn btn--outline" @click="$router.push('/inscription')">
                            Je suis professionnel →
                        </button>
                    </div>

                    <div class="hero-divider"></div>

                    <div class="hero-stats">
                        <div class="stat-item">
                            <h2 class="stat-number">
                                {{ objets_evites
                                }}<span class="text-accent">{{
                                    objet_grandeur
                                }}</span>
                            </h2>
                            <p class="stat-label">CO₂ évité ce mois</p>
                        </div>
                        <div class="stat-item">
                            <h2 class="stat-number">
                                {{ objets_upcycle
                                }}<span
                                    class="text-accent"
                                    v-if="objets_upcycle > 1"
                                    >k+</span
                                >
                            </h2>
                            <p class="stat-label">Objets upcyclés</p>
                        </div>
                        <div class="stat-item">
                            <h2 class="stat-number">
                                {{ artisans_partenaire }}
                            </h2>
                            <p class="stat-label">Artisans partenaires</p>
                        </div>
                        <div class="stat-item">
                            <h2 class="stat-number">{{ sites }}</h2>
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

            <section class="processus-part">
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
                        <button class="btn btn--primary">
                            Créer mon compte gratuit
                        </button>
                        <button class="btn btn--outline">
                            Voir les annonces
                        </button>
                    </div>
                </div>
            </section>
        </main>
    </div>
    <footer class="footer-class">
        <div class="footer-container">
            <div class="footer-main-content">
                <div class="footer-brand">
                    <div class="logo-text">
                        Upcycle<span class="navbar__logo--accent">Connect</span>
                    </div>
                    <p class="footer-description">
                        La plateforme de référence pour l'upcycling en France et
                        en Europe. Rejoignez une communauté engagée pour
                        l'économie circulaire.
                    </p>
                </div>

                <div class="footer-links-wrapper">
                    <div class="footer-column">
                        <h4>PLATEFORME</h4>
                        <ul>
                            <li><a href="#">Comment ça marche</a></li>
                            <li><a href="#">Annonces</a></li>
                            <li><a href="#">Formations</a></li>
                            <li><a href="#">Communauté</a></li>
                        </ul>
                    </div>
                    <div class="footer-column">
                        <h4>PROFESSIONNELS</h4>
                        <ul>
                            <li><a href="#">Offres d'abonnement</a></li>
                            <li><a href="#">Conteneurs</a></li>
                            <li><a href="#">Projets upcycling</a></li>
                            <li><a href="#">Publicité</a></li>
                        </ul>
                    </div>
                    <div class="footer-column">
                        <h4>ENTREPRISE</h4>
                        <ul>
                            <li><a href="#">À propos</a></li>
                            <li><a href="#">Nos sites</a></li>
                            <li><a href="#">Blog</a></li>
                            <li><a href="#">Contact</a></li>
                        </ul>
                    </div>
                </div>
            </div>

            <div class="footer-divider"></div>

            <div class="footer-bottom">
                <p>© 2026 UpcycleConnect — AlterNext · ESGI 2A</p>
                <div class="footer-legal">
                    <a href="#">Mentions légales</a>
                    <a href="#">CGU</a>
                    <a href="#">RGPD</a>
                </div>
            </div>
        </div>
    </footer>
</template>

<script>
export default {
    data() {
        return {
            objets_evites: 2.4,
            objet_grandeur: "t",
            objets_upcycle: 8,
            artisans_partenaire: 340,
            sites: 7,
            plateform_user: 12400,
        };
    },
};
</script>

<style>
body,
html {
    margin: 0;
    padding: 0;
    font-family: "Syne", sans-serif;
}

.navbar {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    z-index: 2000;
    padding: 1.25rem 0;
    background-color: rgba(0, 0, 0, 0.3);
    backdrop-filter: blur(10px);
}

.navbar__container {
    max-width: 1500px;
    margin: 0 auto;
    padding: 0 1.5rem;
    display: flex;
    align-items: center;
    justify-content: space-between;
}

.navbar__right-side {
    display: flex;
    align-items: center;
    gap: 3rem;
}

.navbar__logo {
    display: flex;
    align-items: center;
    gap: 12px;
    cursor: pointer;
}

.logo-text {
    display: inline-block; 
    white-space: nowrap; 
    color: #FFFFFF;
    font-size: 1.4rem;
    font-weight: 800;
    cursor: pointer;
}

.navbar__logo--accent {
    color: #4ade80;
}

.navbar__logo-img {
    height: 35px;
    width: auto;
    object-fit: contain;
}

.navbar__list {
    display: flex;
    gap: 2.5rem;
    list-style: none;
    margin: 0;
    padding: 0;
}

.navbar__link {
    color: rgba(255, 255, 255, 0.6);
    text-decoration: none;
    font-size: 0.95rem;
    font-weight: 500;
    transition: 0.3s ease;
}

.navbar__link:hover {
    color: #FFFFFF;
}

.navbar__actions {
    display: flex;
    gap: 1rem;
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
    color: #FFFFFF;
    border-color: rgba(255, 255, 255, 0.2);
}

.btn--outline:hover {
    background: rgba(255, 255, 255, 0.05);
    border-color: #FFFFFF;
}

.btn--primary {
    background-color: var(--brand-green);
    color: #FFFFFF;
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
    padding-top: 30px;
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
    color: #FFFFFF;
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
    background-color: #FFFFFF;
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
  background: #FFFFFF;
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
  font-family: 'Syne', sans-serif;
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
    background-image: url('../components/pret-a-agir.png');
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