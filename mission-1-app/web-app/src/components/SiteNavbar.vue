<template>
  <header class="site-navbar">
    <div class="site-navbar__inner">
      <div class="site-navbar__left">
        <RouterLink class="site-navbar__brand" to="/front" aria-label="Accueil UpcycleConnect">
          <img :src="logoSrc" alt="UpcycleConnect" class="site-navbar__logo" />
        </RouterLink>

        <div v-if="isAuthenticated" class="site-navbar__account nav-menu">
          <button class="site-navbar__account-button" type="button" aria-haspopup="true">
            <span class="site-navbar__avatar">{{ userInitials }}</span>
            <span class="site-navbar__account-text">
              <strong>{{ userName }} <span v-if="isPremium" class="premium-check-mini">✓</span></strong>
              <small>{{ activeUserRole }} · Score {{ userScore }} pts</small>
            </span>
          </button>
          
          <div class="nav-menu__panel nav-menu__panel--account">
            <RouterLink to="/profil/informations" class="nav-menu__item">
              <span>Compte</span>
              <small>Informations personnelles</small>
            </RouterLink>
            <RouterLink to="/profil/notifications" class="nav-menu__item">
              <span>Notifications</span>
              <small>Alertes casiers, messages et rappels</small>
            </RouterLink>
            <RouterLink to="/messages" class="nav-menu__item">
              <span>Messagerie</span>
              <small>Discussions avec vendeurs et membres</small>
            </RouterLink>
            <RouterLink to="/profil/factures" class="nav-menu__item">
              <span>Factures</span>
              <small>Historique d'achats et reçus</small>
            </RouterLink>
            
            <RouterLink v-if="activeUserRole === 'Prestataire'" to="/abonnement" class="nav-menu__item">
              <span>Abonnement</span>
              <small>Gestion du forfait pro et alertes</small>
            </RouterLink>

            <button class="nav-menu__item nav-menu__item--danger" type="button" @click="handleLogout">
              <span>Se deconnecter</span>
              <small>Fermer la session locale</small>
            </button>
          </div>
        </div>
      </div>

      <nav class="site-navbar__links" aria-label="Navigation principale">
        <div
          v-for="group in dynamicNavGroups"
          :key="group.label"
          class="nav-menu"
          :class="{ 'is-active': isGroupActive(group) }"
        >
          <RouterLink
            v-if="group.to"
            :to="group.to"
            class="site-navbar__link"
            :class="{ 'is-active': isGroupActive(group) }"
          >
            {{ group.label }}
          </RouterLink>
          <button
            v-else
            class="site-navbar__link site-navbar__link--button"
            :class="{ 'is-active': isGroupActive(group) }"
            type="button"
            aria-haspopup="true"
          >
            {{ group.label }}
            <span class="site-navbar__chevron">⌄</span>
          </button>

          <div v-if="group.children?.length" class="nav-menu__panel">
            <RouterLink
              v-for="item in group.children"
              :key="item.label"
              :to="item.to"
              class="nav-menu__item"
              :class="{ 'is-active': isActive(item) }"
            >
              <span>{{ item.label }}</span>
              <small>{{ item.description }}</small>
            </RouterLink>
          </div>
        </div>
      </nav>

      <div class="site-navbar__actions" v-if="isAuthenticated">
        <RouterLink class="site-navbar__button site-navbar__button--ghost" to="/messages">
        Messages
        </RouterLink>
        <RouterLink class="site-navbar__button site-navbar__button--primary" to="/panier">
        Panier
        </RouterLink>
        <RouterLink v-if="activeUserRole !== 'Prestataire'" class="site-navbar__button site-navbar__button--primary" to="/profil/annonces">
          + Deposer
        </RouterLink>
      </div>

      <div class="site-navbar__actions" v-else>
        <RouterLink class="site-navbar__button site-navbar__button--ghost" to="/connexion">
          Connexion
        </RouterLink>
        <RouterLink class="site-navbar__button site-navbar__button--primary" to="/inscription">
          S'inscrire
        </RouterLink>
      </div>
      <div class="nav-menu site-navbar__langue">
        <button class="site-navbar__link site-navbar__link--button" type="button" aria-haspopup="true">
            {{ currentLangCode.toUpperCase() }}
            <span class="site-navbar__chevron">⌄</span>
        </button>
        
        <div class="nav-menu__panel nav-menu__panel--lang">
            <button 
            v-for="langue in langues" 
            :key="langue.id" 
            class="nav-menu__item nav-menu__item--lang"
            @click="changerLangue(langue)"
            type="button"
            >
            <span>{{ langue.code.toUpperCase() }}</span>
            </button>
        </div>
        </div>
    </div>
  </header>
</template>

<script setup>
import { computed, ref, onMounted } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import logoSrc from "./logo_texte.png";

const route = useRoute();
const router = useRouter();

const API_URL = "/go";

const langues = ref([]);
const currentLangCode = ref(localStorage.getItem("langCode") || "fr");
const t = ref({});

const fetchLangues = async () => {
  try {
    const res = await fetch(`${API_URL}/langues`);
    if (res.ok) langues.value = await res.json();
  } catch (e) { console.error(e); }
};

const fetchTraductions = async (code) => {
  try {
    const res = await fetch(`${API_URL}/traductions/${code}`);
    if (res.ok) t.value = await res.json();
  } catch (e) { console.error(e); }
};

const changerLangue = async (langue) => {
  if (!langue) return;
  
  if (document.activeElement) {
    document.activeElement.blur();
  }

  currentLangCode.value = langue.code;
  localStorage.setItem("langCode", langue.code);
  localStorage.setItem("langId", langue.id);
  await fetchTraductions(langue.code);

  window.dispatchEvent(new Event("lang-changed"));

  if (props.isAuthenticated) {
    const userId = sessionStorage.getItem("userId");
    if (userId) {
      fetch(`${API_URL}/users/${userId}/langue`, {
        method: "PUT",
        headers: { 
          "Content-Type": "application/json",
          "Authorization": `Bearer ${sessionStorage.getItem("userToken")}`
        },
        body: JSON.stringify({ id_langue: langue.id })
      }).catch(e => console.error(e));
    }
  }
};

onMounted(() => {
  fetchLangues();
  fetchTraductions(currentLangCode.value);
});

const props = defineProps({
  variant: {
    type: String,
    default: "public"
  },
  userName: {
    type: String,
    default: "Marie Lambert"
  },
  userRole: {
    type: String,
    default: "Particulier"
  },
  userScore: {
    type: [String, Number],
    default: () => sessionStorage.getItem("userScore") || 0
  },
  isAuthenticated: {
    type: Boolean,
    default: false
  }
});

const activeUserRole = computed(() => {
  return sessionStorage.getItem("userRole") || props.userRole;
});

const isPremium = computed(() => sessionStorage.getItem("isPremium") === "true");

const dashboardChildrenParticulier = [
  {
    label: "Vue d'ensemble",
    to: "/profil",
    description: "Score, planning et activite"
  },
  {
    label: "Mes annonces",
    to: "/profil/annonces",
    description: "Objets publies et brouillons"
  },
  {
    label: "Mes depots conteneurs",
    to: "/profil/depots",
    description: "Suivi des depots et collectes"
  },
  {
    label: "Mes projets likés",
    to: "/profil/projet-favoris",
    description: "Projets likés et sauvegardés"
  },
];

const dashboardChildrenPrestataire = [
  {
    label: "Vue d'ensemble",
    to: "/profil",
    description: "Récupération objets en conteneurs"
  },
  {
    label: "Récupérer des objets",
    to: "/profil/recuperations",
    description: "Objets et matériaux en attente de récupération"
  },
  {
    label: "Mes annonces likées",
    to: "/profil/favoris",
    description: "Objets et matériaux sauvegardés"
  },
  {
    label: "Mes projets",
    to: "/profil/projets",
    description: "Créations et vitrine d'upcycling"
  },
];

const dashboardChildrenSalarie = [
  {
    label: "Vue d'ensemble",
    to: "/profil",
    description: "Score, planning et activite"
  },
  {
      label: "Mes Tips",
      to: "/profil/tips",
      description: "Gérez vos astuces partagées et tutoriels d'upcycling"
    },
    {
      label: "Mes Formations",
      to: "/profil/formations",
      description: "Consultez vos cours suivis et parcours d'apprentissage"
    },
    {
      label: "Mes Evenements",
      to: "/profil/evenements",
      description: "Retrouvez vos inscriptions et vos ateliers à venir"
    },
    {
      label: "Modération des forums",
      to: "/profil/forum",
      description: "Gérez les signalements et veillez aux règles de la communauté"
    }
];
const serviceChildren = [
  {
    label: "Upcycling Score",
    to: "/profil",
    description: "Impact et points utilisateur"
  },
  {
    label: "Formations & Ateliers",
    to: "/formations",
    description: "Sessions et apprentissage"
  },
  {
    label: "Mon planning",
    to: "/profil/planning",
    description: "Evenements et rendez-vous"
  },
  {
    label: "Espace Conseils",
    to: "/conseils",
    description: "Guides et astuces"
  }
];

const communityChildren = [
  {
    label: "Catalogue offres",
    to: "/catalogue",
    description: "Toutes les annonces publiques"
  },
  {
    label: "Forums",
    to: "/forums",
    description: "Discussions, entraide et projets"
  },
  {
    label: "Messagerie",
    to: "/messages",
    description: "DM vendeurs, artisans et membres"
  },
  {
    label: "DM Plus",
    to: "/abonnement",
    description: "Messagerie illimitee"
  },
  {
    label: "Projets Upcycling",
    to: "/projets",
    description: "Rejoignez des initiatives de création"
  },
  {
    label: "Evenements",
    to: "/evenements",
    description: "Participez à des événements"
  }
];

const objetsChildren = [
  {
    label: "Depot Objet",
    to: "/deposer",
    description: "Déposez vos objets"
  },
  {
    label: "Récupérer Objet",
    to: "/claim",
    description: "Récupérez vos objets"
  },
];

const dataChildren = [
  {
    label: "Modèle logique",
    to: "/module/science/1",
    description: "Modèle de ML de prédiction"
  },
];

const dynamicNavGroups = computed(() => {
  let dashboardChildren = dashboardChildrenParticulier;

  if (activeUserRole.value === "Prestataire") {
    dashboardChildren = dashboardChildrenPrestataire;
  } else if (activeUserRole.value === "Salarie") {
    dashboardChildren = dashboardChildrenSalarie;
  }

  const groups = [
    {
      label: "Tableau de bord",
      children: dashboardChildren
    },
    {
      label: "Services",
      children: serviceChildren
    },
    {
      label: "Communaute",
      children: communityChildren
    }
  ];

  if (activeUserRole.value === "Salarie") {
    groups.push({
      label: "Réclamer ou Retirer un objet", 
      children: objetsChildren
    });
  }
  if (activeUserRole.value === "Salarie" || activeUserRole.value === "Admin") {
    groups.push({
      label: "Data Mining", 
      children: dataChildren
    });
  }

  return groups;
});

const userInitials = computed(() =>
  props.userName
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? "")
    .join("")
);

function itemPath(item) {
  return typeof item.to === "string" ? item.to : item.to?.path;
}

function isActive(item) {
  const path = itemPath(item);
  if (!path) return false;
  if (path === "/front" && item.to?.hash) {
    return route.path === "/front" && route.hash === item.to.hash;
  }
  if (path === "/profil") {
    return route.path === "/profil";
  }
  return route.path === path || route.path.startsWith(`${path}/`);
}

function isGroupActive(group) {
  if (group.children?.length) {
    return group.children.some((item) => isActive(item));
  }
  return isActive(group);
}

function handleLogout() {
  sessionStorage.removeItem("userToken");
  sessionStorage.removeItem("userId");
  sessionStorage.removeItem("userRole"); 
  sessionStorage.removeItem("isPremium");
  router.push("/connexion");
}
</script>

<style>
.site-navbar__right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.site-navbar__langue {
  position: relative;
  margin-left: 10px;
}

.nav-menu__panel--lang {
  min-width: 60px; 
  right: 0 !important; 
  left: auto !important; 
  transform: none !important; 
  padding: 5px;
}

.nav-menu__panel--lang {
  width: 60px !important;      
  min-width: 0 !important;    
  right: 0 !important;
  left: auto !important;
  transform: none !important;
  padding: 5px;
}

.nav-menu__item--lang:hover {
  background: #f0f7f3;
  border-radius: 6px;
}

.premium-check-mini {
  display: inline-grid;
  width: 18px;
  height: 18px;
  place-items: center;
  margin-left: 4px;
  border-radius: 999px;
  color: #102018;
  background: #8ef0a8;
  font-size: 0.72rem;
  font-weight: 900;
}
</style>
