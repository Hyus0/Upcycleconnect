<template>
  <footer class="footer-class">
    <div class="footer-container">
      <div class="footer-main-content">
        <div class="footer-brand">
          <div class="logo-text" @click="$router.push('/')">
            Upcycle<span class="navbar__logo--accent">Connect</span>
          </div>
          <p class="footer-description">
            {{ t.FooterDesc || "La plateforme de référence pour l'upcycling en France et en Europe. Rejoignez une communauté engagée pour l'économie circulaire." }}
          </p>
        </div>

        <div class="footer-links-wrapper">
          <div class="footer-column">
            <h4>{{ t.FooterCol1 || 'PLATEFORME' }}</h4>
            <ul>
              <li><router-link to="/#processus">{{ t.FooterHowItWorks || 'Comment ça marche' }}</router-link></li>
              <li><router-link to="/annonces">{{ t.FooterAds || 'Annonces' }}</router-link></li>
              <li><router-link to="/profil">{{ t.FooterDashboard || 'Tableau de bord' }}</router-link></li>
              <li><router-link to="/profil/annonces">{{ t.FooterMyAds || 'Mes annonces' }}</router-link></li>
            </ul>
          </div>

          <div class="footer-column">
            <h4>{{ t.FooterCol2 || 'PROFESSIONNELS' }}</h4>
            <ul>
              <li><router-link to="/inscription">{{ t.FooterSubOffers || "Offres d'abonnement" }}</router-link></li>
              <li><router-link to="/profil/depots">{{ t.FooterContainers || 'Conteneurs' }}</router-link></li>
              <li><router-link to="/profil/annonces">{{ t.FooterProjects || 'Projets upcycling' }}</router-link></li>
              <li><router-link to="/inscription">{{ t.FooterAdvertising || 'Publicité' }}</router-link></li>
            </ul>
          </div>

          <div class="footer-column">
            <h4>{{ t.FooterCol3 || 'ENTREPRISE' }}</h4>
            <ul>
              <li><router-link to="/">{{ t.FooterAbout || 'À propos' }}</router-link></li>
              <li><router-link to="/">{{ t.FooterLocations || 'Nos sites' }}</router-link></li>
              <li><router-link to="/">{{ t.FooterBlog || 'Blog' }}</router-link></li>
              <li><router-link to="/connexion">{{ t.FooterContact || 'Contact' }}</router-link></li>
            </ul>
          </div>
        </div>
      </div>

      <div class="footer-divider"></div>

      <div class="footer-bottom">
        <p>© 2026 UpcycleConnect — AlterNext · ESGI 2A</p>
        <div class="footer-legal">
          <router-link to="/">{{ t.FooterLegal || 'Mentions légales' }}</router-link>
          <router-link to="/">{{ t.FooterTOS || 'CGU' }}</router-link>
          <router-link to="/">{{ t.FooterGDPR || 'RGPD' }}</router-link>
        </div>
      </div>
    </div>
  </footer>
</template>

<script setup>
import { useTraduction } from "../composables/useTraduction";
const { t } = useTraduction();
</script>

<style scoped>

.footer-class {
  background-color: #0b0f0e;
  padding: 80px 5% 40px 5%;
  color: #ffffff;
  font-family: "Syne", sans-serif;
}

.footer-container {
  max-width: 1300px;
  margin: 0 auto;
}

.footer-main-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: 60px;
  flex-wrap: wrap;
}

.footer-brand {
  flex: 2;
  min-width: 250px;
}

.logo-text {
  display: inline-block;
  white-space: nowrap;
  color: #ffffff;
  font-size: 1.6rem;
  font-weight: 800;
  cursor: pointer;
  margin-bottom: 20px;
}

.navbar__logo--accent {
  color: #4ade80;
}

.footer-description {
  color: #888;
  line-height: 1.6;
  font-size: 0.95rem;
  max-width: 350px;
}

.footer-links-wrapper {
  flex: 3;
  display: flex;
  justify-content: space-between;
  gap: 30px;
  min-width: 300px;
}

.footer-column h4 {
  font-size: 0.85rem;
  letter-spacing: 1.5px;
  margin-bottom: 25px;
  color: #ffffff;
  font-weight: 700;
  text-transform: uppercase;
}

.footer-column ul {
  list-style: none;
  padding: 0;
  margin: 0;
}

.footer-column ul li {
  margin-bottom: 12px;
}

.footer-column ul li a {
  color: #888;
  text-decoration: none;
  font-size: 0.9rem;
  transition: all 0.3s ease;
}

.footer-column ul li a:hover {
  color: #4ade80;
  padding-left: 5px;
}

.footer-divider {
  height: 1px;
  background: rgba(255, 255, 255, 0.05);
  margin: 60px 0 30px 0;
}

.footer-bottom {
  display: flex;
  justify-content: space-between;
  align-items: center;
  color: #555;
  font-size: 0.85rem;
  flex-wrap: wrap;
  gap: 20px;
}

.footer-legal {
  display: flex;
  gap: 25px;
}

.footer-legal a {
  color: #555;
  text-decoration: none;
  transition: color 0.3s;
}

.footer-legal a:hover {
  color: #ffffff;
}

@media (max-width: 768px) {
  .footer-links-wrapper {
    flex-direction: column;
    gap: 40px;
  }
  
  .footer-bottom {
    flex-direction: column;
    text-align: center;
  }
}
</style>