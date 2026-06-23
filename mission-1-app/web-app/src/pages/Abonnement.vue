<template>
  <main class="public-dashboard">
    <SiteNavbar :is-authenticated="isLoggedIn" :user-name="userName" variant="public" />

    <header class="content-header">
      <div class="header-left">
        <p class="sidebar__category2">ACCUEIL > MON COMPTE > ESPACE PRO</p>
        <h1 class="hero-title1">Mon Espace Pro</h1>
        <p class="classic-text">
          Gérez votre accès aux outils avancés, vos collectes et vos campagnes.
        </p>
      </div>
      <button class="btn-secondary" type="button" @click="$router.push('/profil')">
        🠔 Retour au profil
      </button>
    </header>
    
    <div v-if="loading" class="state-card" style="margin-top: 24px; text-align: center;">
      Chargement de votre espace...
    </div>
    
    <section v-else class="section-container" style="margin-top: 24px;">
      
      <div v-if="!subscription.is_premium" class="pricing-container text-center">
        <div class="section-header justify-center" style="margin-bottom: 2rem;">
          <div>
            <h2 style="font-size: 1.8rem; margin-bottom: 8px;">Passez à la vitesse supérieure</h2>
            <p class="classic-text">Débloquez des outils puissants pour votre activité (28 € / mois).</p>
          </div>
        </div>
    
        <div class="pricing-grid">
          <div class="card card--white pricing-card">
            <h3 class="pricing-title">Artisan Débutant (Freemium)</h3>
            <p class="pricing-price">0 € <span>/mois</span></p>
            <ul class="pricing-features">
                <li><Check :size="18" class="text-green-500" /> Publication de créations & recherche</li>
                <li><Check :size="18" class="text-green-500" /> Réservation standard de matériaux</li>
                <li class="locked">🔒 Statistiques & Impact écologique</li>
                <li class="locked">🔒 Alertes de collecte priorisées</li>
                <li class="locked">🔒 Campagnes de Sponsoring</li>
            </ul>
            <button class="btn-secondary w-full" disabled style="width: 100%; margin-top: auto;">
              Plan actuel
            </button>
          </div>
    
          <div class="card card--score pricing-card premium-card">
            <div class="tag-score popular-tag">RECOMMANDÉ</div>
            <h3 class="pricing-title text-white">Pro & Entreprises</h3>
            <p class="pricing-price text-white">25.00 € <span class="text-white-50">/mois</span></p>
            <ul class="pricing-features text-white">
              <li><Check :size="18" class="text-white" /> Tout le plan gratuit</li>
              <li><Check :size="18" class="text-white" /> Tableaux de bord avancés</li>
              <li><Check :size="18" class="text-white" /> Analyse d'impact écologique détaillée</li>
              <li><Check :size="18" class="text-white" /> Alertes priorisées pour la collecte</li>
              <li><Check :size="18" class="text-white" /> Création de campagnes publicitaires</li>
            </ul>
            <button 
              class="btn-main-action w-full justify-center btn-premium" 
              style="width: 100%; justify-content: center; margin-top: auto;"
              @click="goToPaiement"
            >
              S'abonner maintenant
            </button>
          </div>
        </div>
      </div>
    
      <div v-else>
    
        <div class="pricing-container text-center" style="margin-bottom: 4rem;">
          <div class="section-header justify-center" style="margin-bottom: 2rem;">
            <div>
              <h2 style="font-size: 1.8rem; margin-bottom: 8px;">Gestion de votre abonnement</h2>
              <p class="classic-text" style="font-size: 1.1rem;">
                Votre abonnement est actuellement <strong>{{ subscription.statut.toLowerCase() }}</strong>
                Il est valable jusqu'au <strong>{{ formatDate(subscription.date_fin) || 'Facturation continue' }}</strong>.
              </p>
            </div>
          </div>
    
          <div class="pricing-grid">
            <div class="card card--white pricing-card">
              <h3 class="pricing-title">Artisan Débutant (Freemium)</h3>
              <p class="pricing-price">0 € <span>/mois</span></p>
              <ul class="pricing-features">
                <li><Check :size="18" class="text-green-500" /> Accès gratuit pour publier et chercher</li>
                <li><Check :size="18" class="text-green-500" /> Dépôt dans les conteneurs</li>
                <li class="locked">🔒 Statistiques & Impact écologique</li>
                <li class="locked">🔒 Alertes de collecte priorisées</li>
                <li class="locked">🔒 Campagnes de Sponsoring</li>
              </ul>
              <button 
                v-if="subscription.statut === 'Actif'" 
                class="btn-danger w-full" 
                style="width: 100%; margin-top: auto; justify-content: center;"
                :disabled="processing"
                @click="resilierAbonnement"
              >
                {{ processing ? "Traitement..." : "Résilier et repasser au gratuit" }}
              </button>
            </div>
    
            <div class="card card--score pricing-card premium-card">
              <div class="tag-score popular-tag" style="background: white; color: #2d7a4f;">ACTIF</div>
              <h3 class="pricing-title text-white">Pro & Entreprises</h3>
              <p class="pricing-price text-white">25.00 € <span class="text-white-50">/mois</span></p>
              <ul class="pricing-features text-white">
                <li><Check :size="18" class="text-white" /> Tout le plan gratuit</li>
                <li><Check :size="18" class="text-white" /> Tableaux de bord avancés</li>
                <li><Check :size="18" class="text-white" /> Analyse d'impact écologique détaillée</li>
                <li><Check :size="18" class="text-white" /> Alertes priorisées pour la collecte</li>
                <li><Check :size="18" class="text-white" /> Création de campagnes publicitaires</li>
              </ul>
              <button class="btn-secondary w-full" disabled style="width: 100%; margin-top: auto; color: #2d7a4f; background: white; opacity: 0.9;">
                Plan actuel
              </button>
            </div>
          </div>
        </div>
    
        <hr style="border: 0; border-top: 1px solid #eee; margin: 3rem 0;">
    
        <div class="section-container" style="border: none; padding: 0;">
          <div class="section-header" style="display: flex; justify-content: space-between; align-items: center;">
            <div>
              <h2 style="font-size: 1.8rem; margin-bottom: 8px;">Vos campagnes de sponsoring</h2>
              <p class="classic-text">Mettez en avant vos projets d'upcycling (100€ - 500€ / campagne).</p>
            </div>
            <button class="btn-main-action" @click="createCampaign">
              + Nouvelle campagne
            </button> 
          </div>
          <table class="data-table">
            <thead>
              <tr>
                <th>Projet Sponsorisé</th>
                <th>Budget</th>
                <th>Vues générées</th>
                <th>Statut</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="campagne in campaigns" :key="campagne.id">
                <td><strong>{{ campagne.titre }}</strong></td>
                <td>{{ campagne.budget }} €</td>
                <td>{{ campagne.vues }}</td>
                <td>
                  <span :class="campagne.statut === 'Active' ? 'status-valid' : 'status-pending'">
                    {{ campagne.statut.toUpperCase() }}
                  </span>
                </td>
              </tr>
              <tr v-if="campaigns.length === 0">
                <td colspan="4" class="text-center text-grey py-4">Vous n'avez aucune campagne en cours.</td>
              </tr>
            </tbody>
          </table>
        </div>
    
      </div>
    </section>
  </main>
  <SiteFooter />
</template>

<script setup>
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import SiteNavbar from "../components/SiteNavbar.vue";
import SiteFooter from "../components/SiteFooter.vue";
import { Check } from 'lucide-vue-next';

const router = useRouter();
const API_URL = "http://localhost:8081";

const loading = ref(true);
const processing = ref(false);

const activeTab = ref('dashboard');
const stats = ref({});
const priorityAlerts = ref([]);
const campaigns = ref([]);

const subscription = ref({ is_premium: false, statut: "" }); 

const currentUserId = computed(() => Number(sessionStorage.getItem("userId")) || 0);
const isLoggedIn = computed(() => !!sessionStorage.getItem("userToken"));
const userName = computed(() => {
  const prenom = sessionStorage.getItem("userPrenom") || "";
  const nom = sessionStorage.getItem("userNom") || "";
  return (prenom || nom) ? `${prenom} ${nom}`.trim() : "Utilisateur";
});

const getHeaders = () => ({
  "Content-Type": "application/json",
  Authorization: sessionStorage.getItem("userToken") || "",
});

const formatDate = (dateStr) => {
  if (!dateStr) return "";
  const date = new Date(dateStr.replace(/Z$/, ""));
  if (Number.isNaN(date.getTime())) return "";
  return new Intl.DateTimeFormat("fr-FR", { day: '2-digit', month: 'long', year: 'numeric' }).format(date);
};

const loadProData = async () => {
  stats.value = {
    materiaux_sauves_kg: 1250,
    co2_evite_kg: 450,
    top_materiau: "Palettes Bois",
    projets_actifs: 4
  };

  priorityAlerts.value = [
    { id: 1, materiau: "Bois Massif (50kg)", localisation: "Déchèterie Nord", temps: "Il y a 2h" },
    { id: 2, materiau: "Chutes de cuir", localisation: "Atelier Textile Centre", temps: "Il y a 5h" }
  ];

  campaigns.value = [
    { id: 1, titre: "Table basse industrielle", budget: 150, vues: 4230, statut: "Active" }
  ];
};

const fetchAbonnement = async () => {
  if (!currentUserId.value) {
    loading.value = false;
    return;
  }
  loading.value = true;
  try {
    const res = await fetch(`${API_URL}/users/${currentUserId.value}/abonnement`, {
      headers: getHeaders(),
    });
    if (res.ok) {
      subscription.value = await res.json();
      if (subscription.value.is_premium) {
        await loadProData(); 
      }
    }
  } catch (error) {
    console.error("Erreur de chargement de l'abonnement:", error);
  } finally {
    loading.value = false;
  }
};

const goToPaiement = () => {
  if (!currentUserId.value) {
    alert("Veuillez vous connecter pour vous abonner.");
    router.push("/connexion");
    return;
  }

  router.push({
    path: '/paiement',
    query: {
      plan_id: 2,
      prix: 28,
      nom_plan: 'Pro & Entreprises'
    }
  });
};

const resilierAbonnement = async () => {
  if (!confirm("Êtes-vous sûr de vouloir résilier ? Vous garderez l'accès jusqu'à la fin du mois en cours.")) return;

  processing.value = true;
  try {
    const res = await fetch(`${API_URL}/users/${currentUserId.value}/abonnement/resilier`, {
      method: "POST",
      headers: getHeaders(),
    });

    if (res.ok) {
      alert("Votre abonnement a été résilié.");
      await fetchAbonnement();
    } else {
      alert("Erreur lors de la résiliation.");
    }
  } catch (error) {
    console.error("Erreur de résiliation:", error);
  } finally {
    processing.value = false;
  }
};

const createCampaign = () => {
  alert("Ouverture du module de création de campagne publicitaire !");
};

onMounted(fetchAbonnement);
</script>

<style scoped>

.public-dashboard {
  min-height: 100vh;
  padding: 20px;
  background: var(--bg-light, #f7f9f7);
}

.content-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-end;
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #e5ede7;
}

.sidebar__category2 { 
  margin: 0;
  color: #a0ada7;
  font-family: "Space Mono", monospace;
  font-size: 0.65rem;
  letter-spacing: 1px;
  text-transform: uppercase;
}

.hero-title1 { 
  font-size: 2.2rem; 
  font-weight: 900; 
  color: #122018; 
  margin: 0; 
}

.classic-text {
  font-size: 0.95rem;
  color: #63746a;
  margin-top: 8px;
  line-height: 1.5;
}

.btn-secondary {
  padding: 8px 16px;
  border-radius: 10px;
  border: 1px solid #ddd;
  background: white;
  cursor: pointer;
  font-weight: bold;
}

.btn-secondary:disabled {
  background: #f5f5f5;
  color: #999;
  border-color: #eee;
  cursor: not-allowed;
}

.section-container {
  background: white;
  padding: 2rem;
  border-radius: 16px;
  border: 1px solid #eee;
  margin-bottom: 2rem;
}

.state-card {
  border: 1px dashed #cfe0d4;
  border-radius: 14px;
  padding: 26px;
  color: var(--text-grey, #666);
  background: #fbfdfb;
}


.pricing-container {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.pricing-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  max-width: 900px;
  width: 100%;
}

@media (max-width: 768px) {
  .pricing-grid { grid-template-columns: 1fr; }
}

.pricing-card {
  padding: 2.5rem;
  border-radius: 16px;
  border: 1px solid #eee;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.02);
  display: flex;
  flex-direction: column;
  position: relative;
  text-align: left;
}

.premium-card {
  background: #2d7a4f;
  color: white;
  border: none;
  transform: scale(1.02);
  box-shadow: 0 24px 60px rgba(16, 32, 24, 0.25);
}

.pricing-title {
  font-size: 1.3rem;
  font-weight: 800;
  margin: 0;
  font-family: 'Syne', sans-serif;
}

.pricing-price {
  font-size: 2.8rem;
  font-weight: 900;
  margin: 1rem 0;
}

.pricing-price span {
  font-size: 1.1rem;
  font-weight: normal;
}

.pricing-features {
  list-style: none;
  padding: 0;
  margin: 2rem 0;
}

.pricing-features li {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 1rem;
  font-size: 0.95rem;
  font-weight: 600;
}

.pricing-features .locked {
  color: #999;
  text-decoration: line-through;
}

.popular-tag {
  position: absolute;
  top: 0;
  right: 0;
  border-radius: 0 16px 0 12px;
  background: #ffb300;
  color: #333;
  font-weight: bold;
  font-size: 0.75rem;
  padding: 6px 12px;
}

.text-white { color: white !important; }
.text-white-50 { color: rgba(255,255,255,0.7) !important; }
.text-green-500 { color: #2d7a4f; }

.tabs-container {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
  border-bottom: 2px solid #eee;
  padding-bottom: 10px;
}

.tab-btn {
  background: transparent;
  border: none;
  font-size: 1.1rem;
  font-weight: bold;
  color: #666;
  padding: 8px 16px;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s;
}

.tab-btn:hover { background: #eaf4ed; color: #2d7a4f; }
.tab-btn.active { background: #2d7a4f; color: white; }

.stats-grid {
  display: grid;
  gap: 1.5rem;
  grid-template-columns: 1.4fr 0.8fr 0.8fr;
}

@media (max-width: 920px) {
  .stats-grid { grid-template-columns: 1fr; }
}

.card {
  background: white;
  padding: 1.5rem;
  border-radius: 12px;
  border: 1px solid #eee;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.02);
}

.card--score { background: #2d7a4f; color: white; border: none; }
.tag-score {
  background: rgba(255, 255, 255, 0.2);
  color: white;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: 800;
  display: inline-block;
  margin-bottom: 1rem;
}

.score-value {
  font-size: 2.8rem;
  font-weight: 900;
  line-height: 1;
  margin-bottom: 0.5rem;
}

.score-value span { font-size: 1.2rem; font-weight: normal; opacity: 0.8; }
.card-num, .card-num2 { font-size: 2.5rem; font-weight: bold; line-height: 1; margin-bottom: 0.5rem; }

.score-footer {
  display: flex;
  justify-content: space-between;
  border-top: 1px solid rgba(255, 255, 255, 0.2);
  padding-top: 1rem;
  margin-top: 1rem;
}

.badge {
  padding: 4px 10px;
  border-radius: 20px;
  font-size: 0.8rem;
  font-weight: bold;
  display: inline-block;
}

.badge--green { background: #eaf4ed; color: #2d7a4f; }
.badge--orange { background: #fff3cd; color: #856404; }

.data-table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th, .data-table td {
  padding: 12px;
  text-align: left;
  border-bottom: 1px solid #eee;
}

.status-logistique {
  background: #f0f0f0;
  color: #555;
  padding: 4px 10px;
  border-radius: 6px;
  font-size: 0.75rem;
  font-weight: bold;
  text-transform: uppercase;
}

.status-valid { color: #2d7a4f; font-weight: bold; }
.status-pending { color: #856404; font-weight: bold; }

.active-sub-card {
  background: #fff;
  border: 1px solid #e5ede7;
  border-radius: 16px;
  padding: 40px;
  max-width: 700px;
  margin: 0 auto;
  box-shadow: 0 12px 24px rgba(44, 126, 79, 0.08);
}

.sub-header { display: flex; align-items: center; gap: 16px; margin-bottom: 20px; }
.sub-perks {
  background: #fbfdfb;
  padding: 24px;
  border-radius: 12px;
  border: 1px solid #f0f4f1;
  margin-bottom: 30px;
}
.sub-perks ul { list-style: none; padding: 0; margin: 0; }
.sub-perks li { display: flex; align-items: center; gap: 10px; margin-bottom: 12px; color: #1a1a1a; font-weight: 600;}
.check-icon-active { color: #2c7e4f; }

.sub-actions { display: flex; gap: 16px; }

.btn-main-action {
  display: inline-flex;
  align-items: center;
  background: #2d7a4f;
  color: white;
  padding: 10px 20px;
  border-radius: 10px;
  text-decoration: none;
  font-weight: bold;
  border: none;
  cursor: pointer;
}
.btn-main-action:hover:not(:disabled) { background: #1b4d31; }

.btn-premium { background: white; color: #2d7a4f; }
.btn-premium:hover:not(:disabled) { background: #f0f0f0; }

.btn-danger {
  background: transparent;
  color: #c94b4b;
  border: 1px solid #f7d7d7;
  border-radius: 12px;
  padding: 10px 18px;
  font-size: 0.95rem;
  font-family: "Syne", sans-serif;
  font-weight: 700;
  cursor: pointer;
  transition: 0.2s;
}
.btn-danger:hover:not(:disabled) { background: #fceaea; }
.btn-danger:disabled, .btn-main-action:disabled { opacity: 0.6; cursor: not-allowed; }
</style>