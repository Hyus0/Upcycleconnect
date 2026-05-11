<template>
  <footer class="footer-class">
    <div class="footer-container">
      <div class="footer-main-content">
        <div class="footer-brand">
          <div class="logo-text" @click="$router.push('/')">
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
              <li><router-link to="/#processus">Comment ça marche</router-link></li>
              <li><router-link to="/annonces">Annonces</router-link></li>
              <li><router-link to="/profil">Tableau de bord</router-link></li>
              <li><router-link to="/profil/annonces">Mes annonces</router-link></li>
            </ul>
          </div>

          <div class="footer-column">
            <h4>PROFESSIONNELS</h4>
            <ul>
              <li><router-link to="/inscription">Offres d'abonnement</router-link></li>
              <li><router-link to="/profil/depots">Conteneurs</router-link></li>
              <li><router-link to="/profil/annonces">Projets upcycling</router-link></li>
              <li><router-link to="/inscription">Publicité</router-link></li>
            </ul>
          </div>

          <div class="footer-column">
            <h4>ENTREPRISE</h4>
            <ul>
              <li><router-link to="/">À propos</router-link></li>
              <li><router-link to="/">Nos sites</router-link></li>
              <li><router-link to="/">Blog</router-link></li>
              <li><router-link to="/connexion">Contact</router-link></li>
            </ul>
          </div>
        </div>
      </div>

      <div class="footer-divider"></div>

      <div class="footer-bottom">
        <p>© 2026 UpcycleConnect — AlterNext · ESGI 2A</p>
        <div class="footer-legal">
          <router-link to="/">Mentions légales</router-link>
          <router-link to="/">CGU</router-link>
          <router-link to="/">RGPD</router-link>
        </div>
      </div>
    </div>
  </footer>
</template>

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